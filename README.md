# guillocut

A small web app that computes a **guillotine cutting pattern** for panel
cutting (e.g. building a cabinet from MDF sheets): you enter the base plates
you have and the pieces you need, and it produces a layout that

- uses **guillotine cuts only** (every cut goes edge to edge), and
- keeps the **waste of the base plates minimal**.

## Run

```sh
go run .            # listens on :8080
go run . -addr :9000
```

Then open <http://localhost:8080>, fill in your stock plates and pieces
(dimensions in mm), set the saw kerf, and hit *Optimize*. The result is shown
as an SVG pattern per sheet, with waste statistics.

## API

`POST /api/optimize`

```json
{
  "stocks": [{"name": "MDF 19mm", "w": 2800, "h": 2070, "count": 3}],
  "pieces": [{"name": "side", "w": 600, "h": 1800, "count": 2}],
  "kerf": 3,
  "allowRotation": true
}
```

Returns the placed sheets with piece coordinates, per-sheet and total waste,
and any pieces that could not be placed.

## How it works

The optimizer (`internal/packer`) uses a free-rectangle guillotine heuristic:
each placed piece splits the remaining free space with one straight
edge-to-edge cut, so the layout is always producible on a panel saw. Several
heuristic variants (piece sort orders × split rules) are run and the result
consuming the least stock material wins. The kerf is accounted for on the
cuts right of and below each piece.

Exact guillotine cutting-stock is NP-hard; the heuristic gives good, fast
results but not provably optimal ones.

## Development

```sh
go test ./...
```

The frontend is a single static page in `web/`, embedded into the binary via
`go:embed` — `go build` produces one self-contained executable.
