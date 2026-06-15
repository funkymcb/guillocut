// Package export renders optimization results into formats that CNC and
// CAM software can consume: a flat CSV cutting list and per-sheet DXF
// drawings (bundled in a zip).
package export

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/funkymcb/guillocut/internal/packer"
)

// CSV writes a cutting list with one row per placed piece. Coordinates are
// in millimetres, origin at the top-left corner of each sheet (matching the
// on-screen layout).
func CSV(w io.Writer, res *packer.Result) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{
		"sheet", "stock", "piece", "x_mm", "y_mm", "width_mm", "height_mm", "rotated",
	}); err != nil {
		return err
	}
	for i, sheet := range res.Sheets {
		for _, p := range sheet.Placements {
			if err := cw.Write([]string{
				strconv.Itoa(i + 1),
				sheet.Stock,
				p.Name,
				num(p.X), num(p.Y), num(p.W), num(p.H),
				strconv.FormatBool(p.Rotated),
			}); err != nil {
				return err
			}
		}
	}
	cw.Flush()
	return cw.Error()
}

// DXFZip writes a zip archive containing one DXF drawing per sheet.
func DXFZip(w io.Writer, res *packer.Result) error {
	zw := zip.NewWriter(w)
	for i, sheet := range res.Sheets {
		f, err := zw.Create(fmt.Sprintf("sheet-%d.dxf", i+1))
		if err != nil {
			return err
		}
		if _, err := io.WriteString(f, dxf(sheet)); err != nil {
			return err
		}
	}
	return zw.Close()
}

func num(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// dxf renders a single sheet as an AutoCAD R12 (AC1009) ASCII DXF. R12 is
// the most widely importable DXF flavour. The stock outline, cut pieces and
// labels go on separate layers so an operator can toggle them. DXF uses a
// Y-up coordinate system, so piece positions (Y-down, top-left origin) are
// flipped vertically.
func dxf(sheet packer.Sheet) string {
	b := &builder{}
	b.header()
	b.tables()
	b.pair(0, "SECTION")
	b.pair(2, "ENTITIES")

	b.rect("STOCK", 0, 0, sheet.W, sheet.H)
	for _, p := range sheet.Placements {
		y := sheet.H - p.Y - p.H // flip to Y-up
		b.rect("PIECES", p.X, y, p.W, p.H)

		label := p.Name
		if p.Rotated {
			label += " (rot)"
		}
		h := max(5, min(50, min(p.W, p.H)*0.2))
		b.text("LABELS", p.X+p.W/2, y+p.H/2, h, label)
	}

	b.pair(0, "ENDSEC")
	b.pair(0, "EOF")
	return b.sb.String()
}

type builder struct{ sb strings.Builder }

func (b *builder) header() {
	b.pair(0, "SECTION")
	b.pair(2, "HEADER")
	b.pair(9, "$ACADVER")
	b.pair(1, "AC1009")
	b.pair(9, "$INSUNITS")
	b.pair(70, "4") // millimetres
	b.pair(0, "ENDSEC")
}

func (b *builder) tables() {
	b.pair(0, "SECTION")
	b.pair(2, "TABLES")
	b.pair(0, "TABLE")
	b.pair(2, "LAYER")
	b.pair(70, "3")
	b.layer("STOCK", 8)  // dark grey
	b.layer("PIECES", 5) // blue
	b.layer("LABELS", 3) // green
	b.pair(0, "ENDTAB")
	b.pair(0, "ENDSEC")
}

func (b *builder) layer(name string, color int) {
	b.pair(0, "LAYER")
	b.pair(2, name)
	b.pair(70, "0")
	b.pair(62, strconv.Itoa(color))
	b.pair(6, "CONTINUOUS")
}

// rect emits a closed polyline rectangle with lower-left corner at (x, y).
func (b *builder) rect(layer string, x, y, w, h float64) {
	b.pair(0, "POLYLINE")
	b.pair(8, layer)
	b.pair(66, "1") // vertices follow
	b.pair(70, "1") // closed
	b.f(10, 0)
	b.f(20, 0)
	b.f(30, 0)
	for _, p := range [4][2]float64{{x, y}, {x + w, y}, {x + w, y + h}, {x, y + h}} {
		b.pair(0, "VERTEX")
		b.pair(8, layer)
		b.f(10, p[0])
		b.f(20, p[1])
	}
	b.pair(0, "SEQEND")
	b.pair(8, layer)
}

// text emits a string centred on (cx, cy).
func (b *builder) text(layer string, cx, cy, h float64, s string) {
	b.pair(0, "TEXT")
	b.pair(8, layer)
	b.f(10, cx)
	b.f(20, cy)
	b.f(40, h)
	b.pair(1, s)
	b.pair(72, "1") // horizontally centred
	b.f(11, cx)
	b.f(21, cy)
	b.pair(73, "2") // vertically middle
}

func (b *builder) pair(code int, val string) {
	b.sb.WriteString(strconv.Itoa(code))
	b.sb.WriteByte('\n')
	b.sb.WriteString(val)
	b.sb.WriteByte('\n')
}

func (b *builder) f(code int, v float64) {
	b.pair(code, strconv.FormatFloat(v, 'f', 4, 64))
}
