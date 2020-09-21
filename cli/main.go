package main

import (
	"../internal/pkg"
	"flag"
	"fmt"
	"image/png"
	"os"
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

	imgHeatmap := pkg.GenerateHeatMapPNG(mappingStopIdToNodes, mappingStopIdToNodes[uint32(5121998)][0], 3600)

	imgF := pkg.PanicIfErrOrReturn(os.Create("heatmapPernety.png")).(*os.File)

	pkg.PanicIfErr(png.Encode(imgF, imgHeatmap))

}