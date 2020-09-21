package pkg

type StopIds []uint32

func (arr StopIds) Includes(element uint32) bool {
	var el uint32
	for _, el = range arr {
		if el == element {
			return true
		}
	}
	return false
}

func (arr StopIds) Last() uint32 {
	return arr[len(arr) - 1]
}

type Traveler struct {
	TotalTime        uint64
	StopIdsVisited   StopIds
	MaxTimeAvailable *uint64
	Swarm            *SwarmOfTravelers
	Node             *Node
}

func (t *Traveler) Destroy() {
	t.MaxTimeAvailable = nil
	t.StopIdsVisited = nil
	t.Swarm = nil
	t.Node = nil
}

func (t *Traveler) Advance() {
	var next *Node
	var other *Traveler
	var newTraveler *Traveler
	var copyOfStopIds StopIds
	var present bool
	var totalTimeNext uint64
	var i int
	for i, next = range t.Node.Next {
		totalTimeNext = t.TotalTime + uint64(t.Node.DistanceNext[i])
		if other, present = t.Swarm.BestTravelersByStopId[next.Stop.StopId]; present == true && other.TotalTime <= totalTimeNext {
			continue
		}
		if totalTimeNext > *t.MaxTimeAvailable {
			continue
		}
		if t.StopIdsVisited.Includes(next.Stop.StopId) {
			continue
		}
		copyOfStopIds = make(StopIds, len(t.StopIdsVisited), len(t.StopIdsVisited) + 1)
		copy(copyOfStopIds, t.StopIdsVisited)
		newTraveler = &Traveler{
			TotalTime:        t.TotalTime + uint64(t.Node.DistanceNext[i]),
			StopIdsVisited:   append(copyOfStopIds, next.Stop.StopId),
			MaxTimeAvailable: t.MaxTimeAvailable,
			Swarm:            t.Swarm,
			Node:             next,
		}
		newTraveler.Report()
		newTraveler.Advance()
	}
}

func (t *Traveler) Report() {
	lastStopId := t.StopIdsVisited.Last()
	other, present := t.Swarm.BestTravelersByStopId[lastStopId]
	if !present {
		t.Swarm.BestTravelersByStopId[lastStopId] = t
	} else {
		other.Destroy()
		other = nil
		t.Swarm.BestTravelersByStopId[lastStopId] = t
	}
}

type SwarmOfTravelers struct {
	BestTravelersByStopId map[uint32]*Traveler
	MaxTimeAvailable *uint64
}

func GenerateHeatmapData(fromNodes Nodes, maxTime uint64) *SwarmOfTravelers {
	swarm := SwarmOfTravelers{
		BestTravelersByStopId: make(map[uint32]*Traveler),
		MaxTimeAvailable: &maxTime,
	}
	var travelerFromOrigin Traveler
	for _, node := range fromNodes {
		travelerFromOrigin = Traveler{
			TotalTime:        0,
			StopIdsVisited:   []uint32{node.Stop.StopId},
			MaxTimeAvailable: &maxTime,
			Swarm:            &swarm,
			Node:             node,
		}
		travelerFromOrigin.Report()
		travelerFromOrigin.Advance()
	}
	return &swarm
}

func GenTravelerPath(t *Traveler, stopIdToStop map[uint32]*Stop) string {
	var s = "Path: "
	for _, stopId := range t.StopIdsVisited {
		s += stopIdToStop[stopId].StopName + " => "
	}
	return s[:len(s)-4]
}