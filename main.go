package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"

	"github.com/funkymcb/guillocut/internal/export"
	"github.com/funkymcb/guillocut/internal/packer"
)

//go:embed web
var webFS embed.FS

type optimizeRequest struct {
	Stocks []packer.Stock `json:"stocks"`
	Pieces []packer.Piece `json:"pieces"`
	packer.Options
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	static, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServerFS(static))
	mux.HandleFunc("POST /api/optimize", handleOptimize)
	mux.HandleFunc("POST /api/export/csv", handleExportCSV)
	mux.HandleFunc("POST /api/export/dxf", handleExportDXF)

	log.Printf("guillocut listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func handleOptimize(w http.ResponseWriter, r *http.Request) {
	res, ok := optimize(w, r)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("encoding response: %v", err)
	}
}

func handleExportCSV(w http.ResponseWriter, r *http.Request) {
	res, ok := optimize(w, r)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="guillocut-cutlist.csv"`)
	if err := export.CSV(w, res); err != nil {
		log.Printf("csv export: %v", err)
	}
}

func handleExportDXF(w http.ResponseWriter, r *http.Request) {
	res, ok := optimize(w, r)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="guillocut-dxf.zip"`)
	if err := export.DXFZip(w, res); err != nil {
		log.Printf("dxf export: %v", err)
	}
}

// optimize decodes the request body and runs the packer, writing an error
// response and returning ok=false on failure.
func optimize(w http.ResponseWriter, r *http.Request) (res *packer.Result, ok bool) {
	var req optimizeRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return nil, false
	}
	res, err := packer.Optimize(req.Stocks, req.Pieces, req.Options)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return nil, false
	}
	return res, true
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
