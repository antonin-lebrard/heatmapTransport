package app

import (
	"net/http"
	"personnel/heatmapTransport/internal/app/heatmap"
	"personnel/heatmapTransport/internal/pkg"
)

func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", pkg.Return404)
	mux.HandleFunc("/stops", pkg.HandleGet(heatmap.AllStops))
	mux.HandleFunc("/departures", pkg.HandleGet(heatmap.AllDepartureForStop))
	mux.HandleFunc("/heatmap", pkg.HandleGet(heatmap.GenHeatmapForNode))
	mux.HandleFunc("/paths", pkg.HandleGet(heatmap.GenPathsFromNode))
}
