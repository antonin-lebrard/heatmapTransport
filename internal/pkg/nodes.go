package pkg

import (
	"bufio"
	"bytes"
	"io"
	"math"
	"os"
	"strconv"
	"time"
)

func readNumberChar(c uint8) int64 {
	switch c {
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	case '3':
		return 3
	case '4':
		return 4
	case '5':
		return 5
	case '6':
		return 6
	case '7':
		return 7
	case '8':
		return 8
	case '9':
		return 9
	}
	return 0
}

var hours, minutes, seconds int64
// 11:08:00
func ParseStopTime(stopTime string) (int64, int64, int64) {
	hours = readNumberChar(stopTime[0]) * 10 + readNumberChar(stopTime[1])
	minutes = readNumberChar(stopTime[3]) * 10 + readNumberChar(stopTime[4])
	seconds = readNumberChar(stopTime[6]) * 10 + readNumberChar(stopTime[7])
	return hours, minutes, seconds
}

var fromHours, fromMinutes, fromSeconds, toHours, toMinutes, toSeconds, fromTotal, toTotal int64
func DiffInSecondsBetweenTwoStopTimes(fromStopTime, toStopTime string) int64 {
	fromHours, fromMinutes, fromSeconds = ParseStopTime(fromStopTime)
	toHours, toMinutes, toSeconds = ParseStopTime(toStopTime)

	fromTotal = fromHours * 3600 + fromMinutes * 60 + fromSeconds
	if fromHours == 23 && toHours == 0 {
		toHours = 24
	}
	toTotal = toHours * 3600 + toMinutes * 60 + toSeconds

	return toTotal - fromTotal
}

type Node struct {
	ID uint64
	Stop *Stop
	DistanceNext []uint16
	Next []*Node
	DepartureTime string
}
type Nodes []*Node

func (n Nodes) Len() int      { return len(n) }
func (n Nodes) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return DiffInSecondsBetweenTwoStopTimes(n[i].DepartureTime, n[j].DepartureTime) > 0 }

func saveNodesGraph(nodes Nodes, dataDir string) {
	f := PanicIfErrOrReturn(os.OpenFile(dataDir + "/graph.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)).(*os.File)
	defer f.Close()

	var s, distS, nextS string
	var distance uint16
	var next *Node
	for i, node := range nodes {
		L("saved", i, "/", len(nodes), "nodes")
		s = strconv.FormatUint(node.ID, 10) + ";" +
			strconv.FormatUint(uint64(node.Stop.StopId), 10) + ";" +
			node.DepartureTime + ";"
		distS = ""
		for _, distance = range node.DistanceNext {
			distS += "," + strconv.FormatUint(uint64(distance), 10)
		}
		if len(distS) > 0 { distS = distS[1:] }
		nextS = ""
		for _, next = range node.Next {
			nextS += "," + strconv.FormatUint(next.ID, 10)
		}
		if len(nextS) > 0 { nextS = nextS[1:] }
		s += distS + ";" + nextS + "\n"
		PanicIfErrOrReturn(f.WriteString(s))
		//runtime.GC()
	}
}

func readNodesGraph(stopIdToStop map[uint32]*Stop, dataDir string) map[uint32]Nodes {
	f := PanicIfErrOrReturn(os.OpenFile(dataDir + "/graph.txt", os.O_RDONLY|os.O_CREATE, 0644)).(*os.File)
	defer f.Close()

	mappingStopIdToNodes := make(map[uint32]Nodes, len(stopIdToStop))
	var nodes = make(Nodes, 0, 14305101)
	var node *Node
	var nextsIdsForNodes = make([][]uint64, 0, 14305101)
	var lineB []byte

	var bNodeId, bStopId, bDepartureTime, bDistances, bNextStops, bDist, bNext []byte
	var nodeId, dist, next, stopId uint64
	var idxSep1, idxSep2, idxSep3, idxSep4 int

	var dists = make([]uint16, 0, 58)

	var nexts = make([]uint64, 0, 58)

	var idxLastSep, idxNextSep int

	var startupTime = time.Now()
	var i = 0
	var err error = nil
	var scanner = bufio.NewReader(f)
	lineB, err = scanner.ReadBytes('\n')
	for ;err == nil; i++ {
		if i != 0 && i % 500000 == 0 {
			L(i, "/", cap(nodes))
			printRemainingTime(startupTime, i, cap(nodes))
		}
		idxSep1 = bytes.IndexByte(lineB, ';')
		idxSep2 = bytes.IndexByte(lineB[idxSep1 + 1:], ';') + idxSep1 + 1
		idxSep3 = bytes.IndexByte(lineB[idxSep2 + 1:], ';') + idxSep2 + 1
		idxSep4 = bytes.IndexByte(lineB[idxSep3 + 1:], ';') + idxSep3 + 1
		bNodeId = lineB[:idxSep1]
		bStopId = lineB[idxSep1+1:idxSep2]
		bDepartureTime = lineB[idxSep2+1:idxSep3]
		bDistances = lineB[idxSep3+1:idxSep4]
		bNextStops = lineB[idxSep4+1:]

		idxLastSep = 0
		for ; idxLastSep != -1 ; {
			idxNextSep = bytes.IndexByte(bDistances[idxLastSep:], ',') + idxLastSep
			if idxNextSep < idxLastSep {
				bDist = bDistances[idxLastSep:]
			} else {
				bDist = bDistances[idxLastSep:idxNextSep]
			}
			if len(bDist) == 0 {
				break
			}
			dist, err = strconv.ParseUint(string(bDist), 10, 16)
			if err != nil { panic(err) }
			dists = append(dists, uint16(dist))
			if idxNextSep < idxLastSep {
				break
			}
			idxLastSep = idxNextSep + 1
		}

		idxLastSep = 0
		for ; idxLastSep != -1 ; {
			idxNextSep = bytes.IndexByte(bNextStops[idxLastSep:], ',') + idxLastSep
			if idxNextSep < idxLastSep {
				bNext = bNextStops[idxLastSep:len(bNextStops) - 1]
			} else {
				bNext = bNextStops[idxLastSep:idxNextSep]
			}
			if len(bNext) == 0 {
				break
			}
			next, err = strconv.ParseUint(string(bNext), 10, 64)
			if err != nil { panic(err) }
			nexts = append(nexts, next)
			if idxNextSep < idxLastSep {
				break
			}
			idxLastSep = idxNextSep + 1
		}

		nodeId, err = strconv.ParseUint(string(bNodeId), 10, 64)
		if err != nil { panic(err) }
		stopId, err = strconv.ParseUint(string(bStopId), 10, 32)
		if err != nil { panic(err) }

		distsCopy := make([]uint16, len(dists), len(dists))
		copy(distsCopy, dists)
		dists = dists[:0]

		nextsCopy := make([]uint64, len(nexts), len(nexts))
		copy(nextsCopy, nexts)
		nexts = nexts[:0]

		nextsIdsForNodes = append(nextsIdsForNodes, nextsCopy)
		node = &Node{
			ID:            nodeId,
			Stop:          stopIdToStop[uint32(stopId)],
			DistanceNext:  distsCopy,
			Next:          make(Nodes, 0, len(nextsCopy)),
			DepartureTime: string(bDepartureTime),
		}
		mappingStopIdToNodes[uint32(stopId)] = append(mappingStopIdToNodes[uint32(stopId)], node)
		nodes = append(nodes, node)
		lineB, err = scanner.ReadBytes('\n')
	}
	if err != io.EOF {
		panic(err)
	}

	L("read all the graph file, and construct all nodes")

	for i, node := range nodes {
		for _, nodeId := range nextsIdsForNodes[i] {
			node.Next = append(node.Next, nodes[nodeId])
		}
	}

	L("linked all the nodes")

	return mappingStopIdToNodes
}

func buildNodesGraph(stopTimes []*StopTime, stops []*Stop, transfers []*Transfer, stopIdToStop map[uint32]*Stop, dataDir string) map[uint32]Nodes {
	mappingStopIdToNodes := make(map[uint32]Nodes, len(stops))
	var nodes Nodes

	var id uint64
	id = 0
	var precedentNode *Node
	for i, stopTime := range stopTimes {
		node := Node{
			ID: id,
			Stop: stopIdToStop[stopTime.StopId],
			DistanceNext: make([]uint16, 0),
			Next: make(Nodes, 0),
			DepartureTime: stopTime.ArrivalTime,
		}
		id++
		mappingStopIdToNodes[stopTime.StopId] = append(mappingStopIdToNodes[stopTime.StopId], &node)
		nodes = append(nodes, &node)
		if precedentNode != nil && stopTimes[i-1].TripId == stopTime.TripId {
			precedentNode.DistanceNext = append(precedentNode.DistanceNext, uint16(DiffInSecondsBetweenTwoStopTimes(precedentNode.DepartureTime, node.DepartureTime)))
			precedentNode.Next = append(precedentNode.Next, &node)
		}
		precedentNode = &node
	}

	for _, stopTime := range stopTimes {
		stopTime.Next = nil
		stopTime.Precedent = nil
	}
	stopTimes = nil

    L("len(nodes)", len(nodes))
	PrintMemUsage()

	startupTime := time.Now()

	lenTransfers := len(transfers)
	var nodeFrom, nodeTo, nodeFromToLink, nodeToFromLink *Node
	var minPossibleTimeFromTo, minPossibleTimeToFrom int64
	var timeFromTo, timeToFrom int64
	for i, transfer := range transfers {
		if i % 20 == 0 {
			L("done", i, "/", lenTransfers, "transfers linking")
			printRemainingTime(startupTime, i, lenTransfers)
		}
		for _, nodeFrom = range mappingStopIdToNodes[transfer.FromStopId] {
			nodeFromToLink = nil
			minPossibleTimeFromTo = math.MaxInt64
			for _, nodeTo = range mappingStopIdToNodes[transfer.ToStopId] {
				timeFromTo = DiffInSecondsBetweenTwoStopTimes(nodeFrom.DepartureTime, nodeTo.DepartureTime)
				if timeFromTo > int64(transfer.TransferTime) && timeFromTo < minPossibleTimeFromTo {
					minPossibleTimeFromTo = timeFromTo
					nodeFromToLink = nodeTo
				}
			}
			if nodeFromToLink != nil {
				nodeFrom.Next = append(nodeFrom.Next, nodeFromToLink)
				nodeFrom.DistanceNext = append(nodeFrom.DistanceNext, uint16(minPossibleTimeFromTo))
			}
		}
		for _, nodeTo = range mappingStopIdToNodes[transfer.ToStopId] {
			nodeToFromLink = nil
			minPossibleTimeToFrom = math.MaxInt64
			for _, nodeFrom = range mappingStopIdToNodes[transfer.FromStopId] {
				timeToFrom = DiffInSecondsBetweenTwoStopTimes(nodeTo.DepartureTime, nodeFrom.DepartureTime)
				if timeToFrom > int64(transfer.TransferTime) && timeToFrom < minPossibleTimeToFrom {
					minPossibleTimeToFrom = timeToFrom
					nodeToFromLink = nodeFrom
				}
			}
			if nodeToFromLink != nil {
				nodeTo.Next = append(nodeTo.Next, nodeToFromLink)
				nodeTo.DistanceNext = append(nodeTo.DistanceNext, uint16(minPossibleTimeToFrom))
			}
		}
	}

	saveNodesGraph(nodes, dataDir)

	return mappingStopIdToNodes
}

func LoadNodesGraphAndStops(dataDir string) (map[uint32]*Stop, map[uint32]Nodes) {
	var stopIdToStop map[uint32]*Stop
	var mappingStopIdToNodes map[uint32]Nodes

	if _, err := os.Stat(dataDir + "/graph.txt"); os.IsNotExist(err) {
		stopTimes, stops, transfers := LoadDataFromDisk(dataDir)
		stopIdToStop = make(map[uint32]*Stop, len(stops))
		for _, stop := range stops {
			stopIdToStop[stop.StopId] = stop
		}
		mappingStopIdToNodes = buildNodesGraph(stopTimes, stops, transfers, stopIdToStop, dataDir)
	} else {
		stops := LoadStopsFromDisk(dataDir)
		stopIdToStop = make(map[uint32]*Stop, len(stops))
		for _, stop := range stops {
			stopIdToStop[stop.StopId] = stop
		}
		mappingStopIdToNodes = readNodesGraph(stopIdToStop, dataDir)
	}

	for stopId, stop := range stopIdToStop {
		if mappingStopIdToNodes[stopId].Len() == 0 {
			L("deleting stop", stop.StopName, "because no node pass through it")
			delete(stopIdToStop, stopId)
			delete(mappingStopIdToNodes, stopId)
		}
	}

	return stopIdToStop, mappingStopIdToNodes
}