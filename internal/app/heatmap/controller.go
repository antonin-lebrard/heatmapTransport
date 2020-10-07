package heatmap

import (
	"errors"
	"math"
	"net/http"
	"heatmapTransport/internal/pkg"
	"regexp"
	"sort"
	"strconv"
)

var mappingStopIdToStop map[uint32]*pkg.Stop
var mappingStopIdToNodes map[uint32]pkg.Nodes

func AllStops(w http.ResponseWriter, r *http.Request) {
	var stops = make([]*pkg.Stop, 0, len(mappingStopIdToStop))
	for stopId, stop := range mappingStopIdToStop {
		if stop == nil {
			pkg.L(stopId, "stop == nil")
		}
		stops = append(stops, stop)
	}
	pkg.WriteResponse(stops, w)
}

func AllDepartureForStop(w http.ResponseWriter, r *http.Request) {
	stopIdstr := r.FormValue("stop_id")
	stopId64, err := strconv.ParseUint(stopIdstr, 10, 32)
	if pkg.HandleErrorInHTTPRequest(400, err, w) {
		return
	}

	var departureTimes []string
	for _, node := range mappingStopIdToNodes[uint32(stopId64)] {
		departureTimes = append(departureTimes, node.DepartureTime)
	}

	sort.Strings(departureTimes)
	pkg.WriteResponse(departureTimes, w)
}

func getNodesForGeneratingHeatmapFromParams(w http.ResponseWriter, r *http.Request) (nodes pkg.Nodes, sentAnError bool) {
	stopIdstr := r.FormValue("stop_id")
	departureTimestr := r.FormValue("from_time")

	stopId64, err := strconv.ParseUint(stopIdstr, 10, 32)
	if pkg.HandleErrorInHTTPRequest(400, err, w) {
		return nil, true
	}
	if len(departureTimestr) != 8 || !departureTimeRegexp.MatchString(departureTimestr) {
		pkg.HandleErrorInHTTPRequest(400, errors.New(`from_time should be in this form "17:50:01"`), w)
		return nil, true
	}

	var nodesFrom pkg.Nodes
	var minPositiveDiff int64 = math.MaxInt64
	var diff int64
	var nextStopIdsAlreadyReachableFromNodesFrom = make(map[uint32]bool)
	for _, node := range mappingStopIdToNodes[uint32(stopId64)] {
		diff = pkg.DiffInSecondsBetweenTwoStopTimes(departureTimestr, node.DepartureTime)
		if diff > 0 && diff < minPositiveDiff {
			minPositiveDiff = diff
			nodesFrom = pkg.Nodes{node}
		}
	}

	if nodesFrom == nil {
		pkg.L("404 Not Found, no nodes found for request stop and hour")
		pkg.L(mappingStopIdToNodes[uint32(stopId64)].Len())
		pkg.Return404(w, r)
		return nil, true
	}

	for _, nextNode := range nodesFrom[0].Next {
		nextStopIdsAlreadyReachableFromNodesFrom[nextNode.Stop.StopId] = true
	}
	var present = false
	var added = false
	for _, node := range mappingStopIdToNodes[uint32(stopId64)] {
		diff = pkg.DiffInSecondsBetweenTwoStopTimes(departureTimestr, node.DepartureTime)
		if diff > 0 && diff < 3600 {
			added = false
			for _, nextNode := range node.Next {
				if present, _ = nextStopIdsAlreadyReachableFromNodesFrom[nextNode.Stop.StopId]; present == false {
					nextStopIdsAlreadyReachableFromNodesFrom[nextNode.Stop.StopId] = true
					if added == false {
						added = true
						nodesFrom = append(nodesFrom, node)
					}
				}
			}
		}
	}
	return nodesFrom, false
}

var departureTimeRegexp, _ = regexp.Compile("^([0-1][0-9])|(2[0-3]):[0-5][0-9]:[0-5][0-9]$")

type StopFromNode struct {
	StopId uint32 `json:"stop_id"`
	TotalTime uint64 `json:"total_time"`
}

func GenHeatmapForNode(w http.ResponseWriter, r *http.Request) {
	maxTimeInSecondsstr := r.FormValue("max_time_seconds")
	maxTimeInSeconds, err := strconv.ParseUint(maxTimeInSecondsstr, 10, 64)
	if pkg.HandleErrorInHTTPRequest(400, err, w) {
		return
	}

	nodesFrom, sentAnError := getNodesForGeneratingHeatmapFromParams(w, r)
	if sentAnError { return }

	swarm := pkg.GenerateHeatmapData(nodesFrom, maxTimeInSeconds)

	var stops = make([]*StopFromNode, 0, len(swarm.BestTravelersByStopId))
	for stopId, traveler := range swarm.BestTravelersByStopId {
		stops = append(stops, &StopFromNode{
			StopId:    stopId,
			TotalTime: traveler.TotalTime,
		})
	}

	pkg.WriteResponse(stops, w)
}

type PathFromNode struct {
	StopIdFrom uint32 `json:"from"`
	StopIdTo uint32 `json:"to"`
	TotalTime uint64 `json:"time"`
}

func GenPathsFromNode(w http.ResponseWriter, r *http.Request) {
	maxTimeInSecondsstr := r.FormValue("max_time_seconds")
	maxTimeInSeconds, err := strconv.ParseUint(maxTimeInSecondsstr, 10, 64)
	if pkg.HandleErrorInHTTPRequest(400, err, w) {
		return
	}

	nodesFrom, sentAnError := getNodesForGeneratingHeatmapFromParams(w, r)
	if sentAnError { return }

	swarm := pkg.GenerateHeatmapData(nodesFrom, maxTimeInSeconds)

	var paths = make([]*PathFromNode, 0, len(swarm.BestTravelersByStopId))
	for _, traveler := range swarm.BestTravelersByStopId {
		if len(traveler.StopIdsVisited) > 1 {
			paths = append(paths, &PathFromNode{
				StopIdFrom: traveler.StopIdsVisited[len(traveler.StopIdsVisited)-2],
				StopIdTo:   traveler.StopIdsVisited[len(traveler.StopIdsVisited)-1],
				TotalTime:  traveler.TotalTime,
			})
		}
	}

	pkg.WriteResponse(paths, w)
}

func SetMappingStopsAndNodes(stopIdToStop map[uint32]*pkg.Stop, stopIdToNodes map[uint32]pkg.Nodes) {
	mappingStopIdToStop = stopIdToStop
	mappingStopIdToNodes = stopIdToNodes
}
