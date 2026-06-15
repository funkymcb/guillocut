package export

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/funkymcb/guillocut/internal/packer"
)

func sampleResult(t *testing.T) *packer.Result {
	t.Helper()
	res, err := packer.Optimize(
		[]packer.Stock{{Name: "MDF 19mm", W: 2800, H: 2070, Count: 3}},
		[]packer.Piece{
			{Name: "side", W: 600, H: 1800, Count: 2},
			{Name: "shelf", W: 762, H: 580, Count: 4},
		},
		packer.Options{Kerf: 3, AllowRotation: true},
	)
	if err != nil {
		t.Fatalf("optimize: %v", err)
	}
	return res
}

func TestCSV(t *testing.T) {
	var buf bytes.Buffer
	if err := CSV(&buf, sampleResult(t)); err != nil {
		t.Fatalf("CSV: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "sheet,stock,piece,x_mm,y_mm,width_mm,height_mm,rotated" {
		t.Errorf("unexpected header: %q", lines[0])
	}
	// 2 sides + 4 shelves = 6 placements + header row.
	if len(lines) != 7 {
		t.Errorf("got %d lines, want 7", len(lines))
	}
}

func TestDXFZip(t *testing.T) {
	res := sampleResult(t)
	var buf bytes.Buffer
	if err := DXFZip(&buf, res); err != nil {
		t.Fatalf("DXFZip: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	if len(zr.File) != len(res.Sheets) {
		t.Fatalf("got %d files, want %d sheets", len(zr.File), len(res.Sheets))
	}

	rc, err := zr.File[0].Open()
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(rc)
	rc.Close()
	dxf := string(data)

	for _, want := range []string{"AC1009", "ENTITIES", "STOCK", "PIECES", "LABELS", "POLYLINE"} {
		if !strings.Contains(dxf, want) {
			t.Errorf("DXF missing %q", want)
		}
	}
	if !strings.HasSuffix(strings.TrimSpace(dxf), "EOF") {
		t.Errorf("DXF does not end with EOF")
	}
	// stock outline + one polyline per placement on sheet 1.
	wantPoly := 1 + len(res.Sheets[0].Placements)
	if got := strings.Count(dxf, "\nPOLYLINE\n"); got != wantPoly {
		t.Errorf("got %d polylines, want %d", got, wantPoly)
	}
}
