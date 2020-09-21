package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"personnel/heatmapTransport/internal/app"
	"personnel/heatmapTransport/internal/app/heatmap"
	"personnel/heatmapTransport/internal/pkg"
)

var dataDir = ""

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&dataDir, "data-dir", "./ratp", "path to the folder containing the stops, transfers, stop_times, and after graph building graph.txt")
}

func main() {
	//f, err := os.Create("cpu.prof")
	//if err != nil {
	//	log.Fatal("could not create CPU profile: ", err)
	//}
	//defer f.Close() // error handling omitted for example
	//if err := pprof.StartCPUProfile(f); err != nil {
	//	log.Fatal("could not start CPU profile: ", err)
	//}
	//defer pprof.StopCPUProfile()

	flag.Usage = usage
	flag.Parse()

	stopIdToStop, mappingStopIdToNodes := pkg.LoadNodesGraphAndStops(dataDir)

	pkg.L(len(stopIdToStop), len(mappingStopIdToNodes))

	heatmap.SetMappingStopsAndNodes(stopIdToStop, mappingStopIdToNodes)

	var mux = http.NewServeMux()
	app.RegisterHandlers(mux)

	serverAddr := "127.0.0.1:4000"

	server := &http.Server{
		Addr: serverAddr,
		Handler: mux,
	}

	log.Println("Serving heatmap http://" + serverAddr)
	pkg.PanicIfErr(server.ListenAndServe())
}