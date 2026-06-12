package packer

import "testing"

func TestPerfectFit(t *testing.T) {
	stocks := []Stock{{Name: "board", W: 100, H: 100, Count: 1}}
	pieces := []Piece{{Name: "half", W: 100, H: 50, Count: 2}}

	res, err := Optimize(stocks, pieces, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Sheets) != 1 {
		t.Fatalf("expected 1 sheet, got %d", len(res.Sheets))
	}
	if len(res.Unplaced) != 0 {
		t.Fatalf("expected no unplaced pieces, got %v", res.Unplaced)
	}
	if res.WastePct != 0 {
		t.Fatalf("expected 0%% waste, got %f%%", res.WastePct)
	}
}

func TestRotationRequired(t *testing.T) {
	stocks := []Stock{{Name: "board", W: 100, H: 50, Count: 1}}
	pieces := []Piece{{Name: "p", W: 40, H: 90, Count: 1}}

	res, err := Optimize(stocks, pieces, Options{AllowRotation: false})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Unplaced) != 1 {
		t.Fatalf("expected piece to be unplaceable without rotation, got %v", res.Unplaced)
	}

	res, err = Optimize(stocks, pieces, Options{AllowRotation: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Unplaced) != 0 {
		t.Fatalf("expected piece to fit with rotation, got unplaced %v", res.Unplaced)
	}
	if !res.Sheets[0].Placements[0].Rotated {
		t.Fatal("expected placement to be rotated")
	}
}

func TestMultipleSheets(t *testing.T) {
	stocks := []Stock{{Name: "board", W: 100, H: 100, Count: 5}}
	pieces := []Piece{{Name: "p", W: 60, H: 60, Count: 4}}

	res, err := Optimize(stocks, pieces, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Unplaced) != 0 {
		t.Fatalf("expected all pieces placed, got unplaced %v", res.Unplaced)
	}
	// Only one 60x60 piece fits per 100x100 sheet.
	if len(res.Sheets) != 4 {
		t.Fatalf("expected 4 sheets, got %d", len(res.Sheets))
	}
	checkInvariants(t, res)
}

func TestKerfConsumesSpace(t *testing.T) {
	stocks := []Stock{{Name: "board", W: 100, H: 50, Count: 2}}
	// Two 50x50 pieces fit on one sheet without kerf, but not with kerf.
	pieces := []Piece{{Name: "p", W: 50, H: 50, Count: 2}}

	res, err := Optimize(stocks, pieces, Options{Kerf: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Unplaced) != 0 {
		t.Fatalf("expected all pieces placed, got unplaced %v", res.Unplaced)
	}
	if len(res.Sheets) != 2 {
		t.Fatalf("expected kerf to force 2 sheets, got %d", len(res.Sheets))
	}
	checkInvariants(t, res)
}

func TestCabinetExample(t *testing.T) {
	stocks := []Stock{{Name: "MDF 19mm", W: 2800, H: 2070, Count: 3}}
	pieces := []Piece{
		{Name: "side", W: 600, H: 1800, Count: 2},
		{Name: "top/bottom", W: 762, H: 600, Count: 2},
		{Name: "shelf", W: 762, H: 580, Count: 4},
		{Name: "back", W: 800, H: 1800, Count: 1},
		{Name: "door", W: 398, H: 1780, Count: 2},
	}

	res, err := Optimize(stocks, pieces, Options{Kerf: 3, AllowRotation: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Unplaced) != 0 {
		t.Fatalf("expected all pieces placed, got unplaced %v", res.Unplaced)
	}
	// Total piece area is ~7.7m², one sheet is ~5.8m², so 2 sheets is the minimum.
	if len(res.Sheets) != 2 {
		t.Fatalf("expected cabinet to fit on 2 sheets, got %d", len(res.Sheets))
	}
	checkInvariants(t, res)
}

func TestValidation(t *testing.T) {
	cases := []struct {
		name   string
		stocks []Stock
		pieces []Piece
		opts   Options
	}{
		{"no stocks", nil, []Piece{{W: 1, H: 1, Count: 1}}, Options{}},
		{"no pieces", []Stock{{W: 1, H: 1, Count: 1}}, nil, Options{}},
		{"negative kerf", []Stock{{W: 1, H: 1, Count: 1}}, []Piece{{W: 1, H: 1, Count: 1}}, Options{Kerf: -1}},
		{"zero stock dim", []Stock{{W: 0, H: 1, Count: 1}}, []Piece{{W: 1, H: 1, Count: 1}}, Options{}},
		{"zero piece count", []Stock{{W: 1, H: 1, Count: 1}}, []Piece{{W: 1, H: 1, Count: 0}}, Options{}},
	}
	for _, c := range cases {
		if _, err := Optimize(c.stocks, c.pieces, c.opts); err == nil {
			t.Errorf("%s: expected error", c.name)
		}
	}
}

// checkInvariants verifies that no placement overlaps another or leaves
// the sheet bounds.
func checkInvariants(t *testing.T, res *Result) {
	t.Helper()
	for si, sheet := range res.Sheets {
		for i, a := range sheet.Placements {
			if a.X < 0 || a.Y < 0 || a.X+a.W > sheet.W || a.Y+a.H > sheet.H {
				t.Errorf("sheet %d: placement %v out of bounds", si, a)
			}
			for j, b := range sheet.Placements {
				if i >= j {
					continue
				}
				if a.X < b.X+b.W && b.X < a.X+a.W && a.Y < b.Y+b.H && b.Y < a.Y+a.H {
					t.Errorf("sheet %d: placements %v and %v overlap", si, a, b)
				}
			}
		}
	}
}
