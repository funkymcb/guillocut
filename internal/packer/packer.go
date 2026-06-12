// Package packer implements a guillotine cutting-stock optimizer.
//
// Pieces are placed into stock sheets using a free-rectangle guillotine
// heuristic: every placement splits the remaining free space with a single
// edge-to-edge cut, so the resulting layout is always producible with
// guillotine cuts only. Several heuristic variants (sort orders x split
// rules) are run and the best result by consumed stock area is returned.
package packer

import (
	"fmt"
	"sort"
)

// Stock describes one type of base plate available for cutting.
type Stock struct {
	Name  string  `json:"name"`
	W     float64 `json:"w"`
	H     float64 `json:"h"`
	Count int     `json:"count"`
}

// Piece describes one type of plate that should be cut out.
type Piece struct {
	Name  string  `json:"name"`
	W     float64 `json:"w"`
	H     float64 `json:"h"`
	Count int     `json:"count"`
}

// Options control the optimization run.
type Options struct {
	// Kerf is the saw blade width consumed by each cut.
	Kerf float64 `json:"kerf"`
	// AllowRotation permits placing pieces rotated by 90 degrees.
	AllowRotation bool `json:"allowRotation"`
}

// Placement is one piece positioned on a sheet. X/Y is the top-left corner.
type Placement struct {
	Name    string  `json:"name"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	W       float64 `json:"w"`
	H       float64 `json:"h"`
	Rotated bool    `json:"rotated"`
}

// Sheet is one used stock plate together with the pieces placed on it.
type Sheet struct {
	Stock      string      `json:"stock"`
	W          float64     `json:"w"`
	H          float64     `json:"h"`
	Placements []Placement `json:"placements"`
	UsedArea   float64     `json:"usedArea"`
	WastePct   float64     `json:"wastePct"`
}

// Result is the outcome of an optimization run.
type Result struct {
	Sheets    []Sheet `json:"sheets"`
	Unplaced  []Piece `json:"unplaced"`
	UsedArea  float64 `json:"usedArea"`
	SheetArea float64 `json:"sheetArea"`
	WastePct  float64 `json:"wastePct"`
}

const maxItems = 10000

// Optimize packs the requested pieces onto the available stock plates.
func Optimize(stocks []Stock, pieces []Piece, opts Options) (*Result, error) {
	if err := validate(stocks, pieces, opts); err != nil {
		return nil, err
	}

	items := expand(pieces)
	if len(items) > maxItems {
		return nil, fmt.Errorf("too many pieces requested (%d, max %d)", len(items), maxItems)
	}

	var best *Result
	for _, sortItems := range sortOrders {
		for _, rule := range []splitRule{splitShortAxis, splitLongAxis} {
			run := append([]item(nil), items...)
			sortItems(run)
			res := pack(stocks, run, opts, rule)
			if better(res, best) {
				best = res
			}
		}
	}
	return best, nil
}

func validate(stocks []Stock, pieces []Piece, opts Options) error {
	if len(stocks) == 0 {
		return fmt.Errorf("at least one stock plate is required")
	}
	if len(pieces) == 0 {
		return fmt.Errorf("at least one piece is required")
	}
	if opts.Kerf < 0 {
		return fmt.Errorf("kerf must not be negative")
	}
	for i, s := range stocks {
		if s.W <= 0 || s.H <= 0 {
			return fmt.Errorf("stock %q (#%d): dimensions must be positive", s.Name, i+1)
		}
		if s.Count <= 0 {
			return fmt.Errorf("stock %q (#%d): count must be positive", s.Name, i+1)
		}
	}
	for i, p := range pieces {
		if p.W <= 0 || p.H <= 0 {
			return fmt.Errorf("piece %q (#%d): dimensions must be positive", p.Name, i+1)
		}
		if p.Count <= 0 {
			return fmt.Errorf("piece %q (#%d): count must be positive", p.Name, i+1)
		}
	}
	return nil
}

// item is a single piece instance to place.
type item struct {
	name string
	w, h float64
}

func expand(pieces []Piece) []item {
	var items []item
	for _, p := range pieces {
		for i := 0; i < p.Count; i++ {
			items = append(items, item{name: p.Name, w: p.W, h: p.H})
		}
	}
	return items
}

var sortOrders = []func([]item){
	// Largest area first.
	func(s []item) {
		sort.SliceStable(s, func(i, j int) bool { return s[i].w*s[i].h > s[j].w*s[j].h })
	},
	// Longest side first.
	func(s []item) {
		sort.SliceStable(s, func(i, j int) bool {
			li, lj := max(s[i].w, s[i].h), max(s[j].w, s[j].h)
			if li != lj {
				return li > lj
			}
			return s[i].w*s[i].h > s[j].w*s[j].h
		})
	},
	// Largest perimeter first.
	func(s []item) {
		sort.SliceStable(s, func(i, j int) bool { return s[i].w+s[i].h > s[j].w+s[j].h })
	},
}

type splitRule int

const (
	splitShortAxis splitRule = iota
	splitLongAxis
)

type freeRect struct {
	x, y, w, h float64
}

type openSheet struct {
	stock      *Stock
	free       []freeRect
	placements []Placement
	usedArea   float64
}

// pack runs a single heuristic pass over the items.
func pack(stocks []Stock, items []item, opts Options, rule splitRule) *Result {
	remaining := make([]int, len(stocks))
	for i, s := range stocks {
		remaining[i] = s.Count
	}

	var sheets []*openSheet
	var unplaced []item

	for _, it := range items {
		sheet, rectIdx, rotated, ok := findBestFit(sheets, it, opts)
		if !ok {
			// No open sheet fits: open the smallest stock plate that can
			// hold the piece.
			idx := -1
			for i := range stocks {
				if remaining[i] == 0 || !fitsStock(stocks[i], it, opts) {
					continue
				}
				if idx == -1 || stocks[i].W*stocks[i].H < stocks[idx].W*stocks[idx].H {
					idx = i
				}
			}
			if idx == -1 {
				unplaced = append(unplaced, it)
				continue
			}
			remaining[idx]--
			sheet = &openSheet{
				stock: &stocks[idx],
				free:  []freeRect{{0, 0, stocks[idx].W, stocks[idx].H}},
			}
			sheets = append(sheets, sheet)
			_, rectIdx, rotated, ok = findBestFit([]*openSheet{sheet}, it, opts)
			if !ok {
				// Cannot happen after fitsStock, but never drop a piece silently.
				unplaced = append(unplaced, it)
				continue
			}
		}
		place(sheet, rectIdx, it, rotated, opts, rule)
	}

	return buildResult(sheets, unplaced)
}

func fitsStock(s Stock, it item, opts Options) bool {
	if it.w <= s.W && it.h <= s.H {
		return true
	}
	return opts.AllowRotation && it.h <= s.W && it.w <= s.H
}

// findBestFit returns the open sheet and free rectangle with the best
// (smallest leftover area) fit for the item.
func findBestFit(sheets []*openSheet, it item, opts Options) (sheet *openSheet, rectIdx int, rotated bool, ok bool) {
	bestArea := -1.0
	bestShort := -1.0
	for _, s := range sheets {
		for i, fr := range s.free {
			orientations := []bool{false}
			if opts.AllowRotation && it.w != it.h {
				orientations = append(orientations, true)
			}
			for _, rot := range orientations {
				w, h := it.w, it.h
				if rot {
					w, h = h, w
				}
				if w > fr.w || h > fr.h {
					continue
				}
				leftover := fr.w*fr.h - w*h
				short := min(fr.w-w, fr.h-h)
				if bestArea < 0 || leftover < bestArea || (leftover == bestArea && short < bestShort) {
					bestArea, bestShort = leftover, short
					sheet, rectIdx, rotated, ok = s, i, rot, true
				}
			}
		}
	}
	return sheet, rectIdx, rotated, ok
}

// place puts the item into the given free rectangle and performs the
// guillotine split of the remaining space.
func place(s *openSheet, rectIdx int, it item, rotated bool, opts Options, rule splitRule) {
	fr := s.free[rectIdx]
	w, h := it.w, it.h
	if rotated {
		w, h = h, w
	}

	s.placements = append(s.placements, Placement{
		Name: it.name, X: fr.x, Y: fr.y, W: w, H: h, Rotated: rotated,
	})
	s.usedArea += w * h

	// The kerf is consumed by the cuts right of and below the piece.
	usedW := w + opts.Kerf
	usedH := h + opts.Kerf
	remW := fr.w - usedW
	remH := fr.h - usedH

	horizontal := remW <= remH // bottom strip spans the full width
	if rule == splitLongAxis {
		horizontal = !horizontal
	}

	var right, bottom freeRect
	if horizontal {
		right = freeRect{fr.x + usedW, fr.y, remW, h}
		bottom = freeRect{fr.x, fr.y + usedH, fr.w, remH}
	} else {
		right = freeRect{fr.x + usedW, fr.y, remW, fr.h}
		bottom = freeRect{fr.x, fr.y + usedH, w, remH}
	}

	s.free = append(s.free[:rectIdx], s.free[rectIdx+1:]...)
	for _, r := range []freeRect{right, bottom} {
		if r.w > 0 && r.h > 0 {
			s.free = append(s.free, r)
		}
	}
}

func buildResult(sheets []*openSheet, unplaced []item) *Result {
	res := &Result{Sheets: []Sheet{}, Unplaced: []Piece{}}
	for _, s := range sheets {
		sheetArea := s.stock.W * s.stock.H
		res.Sheets = append(res.Sheets, Sheet{
			Stock:      s.stock.Name,
			W:          s.stock.W,
			H:          s.stock.H,
			Placements: s.placements,
			UsedArea:   s.usedArea,
			WastePct:   100 * (sheetArea - s.usedArea) / sheetArea,
		})
		res.UsedArea += s.usedArea
		res.SheetArea += sheetArea
	}
	if res.SheetArea > 0 {
		res.WastePct = 100 * (res.SheetArea - res.UsedArea) / res.SheetArea
	}

	// Group unplaced items back into piece counts.
	counts := map[item]int{}
	var order []item
	for _, it := range unplaced {
		if counts[it] == 0 {
			order = append(order, it)
		}
		counts[it]++
	}
	for _, it := range order {
		res.Unplaced = append(res.Unplaced, Piece{Name: it.name, W: it.w, H: it.h, Count: counts[it]})
	}
	return res
}

// better reports whether a should be preferred over b.
func better(a, b *Result) bool {
	if b == nil {
		return true
	}
	au, bu := totalCount(a.Unplaced), totalCount(b.Unplaced)
	if au != bu {
		return au < bu
	}
	if a.SheetArea != b.SheetArea {
		return a.SheetArea < b.SheetArea
	}
	return len(a.Sheets) < len(b.Sheets)
}

func totalCount(pieces []Piece) int {
	n := 0
	for _, p := range pieces {
		n += p.Count
	}
	return n
}
