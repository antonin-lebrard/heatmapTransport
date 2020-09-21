package pkg

import (
	"image"
	"image/color"
)

func ConvertLatLngToPngPoint(lat, long float64, maxX, maxY int) (x int, y int) {
	var minLat float64 = 48.470726
	var maxLat float64 = 49.126904
	var minLng float64 = 1.820545
	var maxLng float64 = 2.851565

	var diffLat = maxLat - minLat
	var diffLng = maxLng - minLng

	// The Y is reversed, as in the image generation going lower means growing Y, but the latitude is growing by going up
	var latToY = int(((maxLat - lat) / diffLat) * float64(maxY))
	var lngToX = int(((long - minLng) / diffLng) * float64(maxX))

	return lngToX, latToY
}

func ColorForTotalTime(t *Traveler) color.RGBA {
	var green, red uint8
	if t.TotalTime < 40 * 60 {
		green = 255
	}
	if t.TotalTime > 20 * 60 {
		red = 255
	}
	return color.RGBA{
		R: red,
		G: green,
		B: 0,
		A: 255,
	}
}

func GenerateHeatMapPNG(mappingStopIdToNodes map[uint32]Nodes, node *Node, maxTime uint64) *image.RGBA {
	var imgHeatmap = image.NewRGBA(image.Rectangle{ Min: image.Point{}, Max: image.Point{X: 1000, Y: 1000} })

	startingNode := mappingStopIdToNodes[uint32(5121998)][0]

	swarm := GenerateHeatmapData(Nodes{startingNode}, 3600)

	L("Started from", startingNode.Stop.StopName, "at:", startingNode.DepartureTime, "Lat:", startingNode.Stop.StopLat, "Long:", startingNode.Stop.StopLon)

	var c color.RGBA
	var x, y int
	for _, traveler := range swarm.BestTravelersByStopId {
		c = ColorForTotalTime(traveler)
		x, y = ConvertLatLngToPngPoint(traveler.Node.Stop.StopLat, traveler.Node.Stop.StopLon, 1000, 1000)
		imgHeatmap.Set(x, y, c)
	}

	return imgHeatmap
}
