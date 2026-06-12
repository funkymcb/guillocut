package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"

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

	log.Printf("guillocut listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func handleOptimize(w http.ResponseWriter, r *http.Request) {
	var req optimizeRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	res, err := packer.Optimize(req.Stocks, req.Pieces, req.Options)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("encoding response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
