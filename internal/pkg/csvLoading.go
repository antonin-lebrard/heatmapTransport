package pkg

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	RATP_DATA_TYPE = 0
	SNCF_DATA_TYPE = 1
	UNKNOWN_DATA_TYPE = 2
)

func readCsvLine(line string) []string {
	var arr []string
	workingS := line[:]
	idxApo := strings.IndexRune(workingS, '"')
	idxSep := strings.IndexRune(workingS, ',')
	for ; idxSep != -1 ; {
		if idxApo == 0 {
			workingS = workingS[1:]
			idxApo = strings.IndexRune(workingS, '"')
			arr = append(arr, workingS[:idxApo])
			workingS = workingS[idxApo + 2:]
		} else {
			arr = append(arr, workingS[:idxSep])
			workingS = workingS[idxSep + 1:]
		}
		idxApo = strings.IndexRune(workingS, '"')
		idxSep = strings.IndexRune(workingS, ',')
	}
	arr = append(arr, workingS)
	return arr
}

type StopTime struct {
	TripId uint64
	ArrivalTime string
	StopId uint32
	StopSequence uint16
	Precedent *StopTime
	Next *StopTime
}

func stopTimeForRatp(s []string) *StopTime {
	return &StopTime{
		TripId:       PanicIfErrOrReturn(strconv.ParseUint(s[0], 10, 64)).(uint64),
		ArrivalTime:  s[1],
		StopId:       uint32(PanicIfErrOrReturn(strconv.ParseUint(s[3], 10, 32)).(uint64)),
		StopSequence: uint16(PanicIfErrOrReturn(strconv.ParseUint(s[4], 10, 16)).(uint64)),
	}
}

func stopTimeForSncf(s []string) *StopTime {
	var tripId = s[0][strings.Index(s[0], "-1_") + 3 : ]
	var stopId = s[3][strings.Index(s[3], ":DUA") + 4 : ]
	return &StopTime{
		TripId:       PanicIfErrOrReturn(strconv.ParseUint(tripId, 10, 64)).(uint64),
		ArrivalTime:  s[1],
		StopId:       uint32(PanicIfErrOrReturn(strconv.ParseUint(stopId, 10, 32)).(uint64)),
		StopSequence: uint16(PanicIfErrOrReturn(strconv.ParseUint(s[4], 10, 16)).(uint64)),
	}
}

func loadStopTimes(filename string) []*StopTime {
	file := PanicIfErrOrReturn(os.Open(filename)).(*os.File)
	defer file.Close()

	var typeOfStopTimesToLoad = UNKNOWN_DATA_TYPE
	if strings.Contains(filename, "ratp") {
		typeOfStopTimesToLoad = RATP_DATA_TYPE
	} else if strings.Contains(filename, "sncf") {
		typeOfStopTimesToLoad = SNCF_DATA_TYPE
	}

	reader := bufio.NewReader(file)
	// ignore first line
	PanicIfErrOrReturn(reader.ReadString(byte('\n')))

	var data []*StopTime
	for ;; {
		line, err := reader.ReadString(byte('\n'))
		if line != "" {
			line = line[:len(line)-1]
			s := strings.Split(line, ",")
			var current *StopTime
			if typeOfStopTimesToLoad == UNKNOWN_DATA_TYPE {
				L("unkown data type, support only ratp and sncf")
				os.Exit(1)
				return nil
			} else if typeOfStopTimesToLoad == RATP_DATA_TYPE {
				current = stopTimeForRatp(s)
			} else if typeOfStopTimesToLoad == SNCF_DATA_TYPE {
				current = stopTimeForSncf(s)
			}
			data = append(data, current)
			if len(data) > 1 && data[len(data) - 2].TripId == data[len(data) - 1].TripId {
				data[len(data) - 2].Next = current
				current.Precedent = data[len(data) - 2]
			}
		}
		if err != nil && err == io.EOF {
			break
		}
		PanicIfErr(err)
	}
	return data
}

type Stop struct {
	StopId uint32 `json:"id"`
	StopName string `json:"name"`
	StopLat float64 `json:"lat"`
	StopLon float64 `json:"lon"`
}

func stopForRatp(s []string) *Stop {
	return &Stop{
		StopId:   uint32(PanicIfErrOrReturn(strconv.ParseUint(s[0], 10, 32)).(uint64)),
		StopName: s[2],
		StopLat:  PanicIfErrOrReturn(strconv.ParseFloat(s[4], 64)).(float64),
		StopLon:  PanicIfErrOrReturn(strconv.ParseFloat(s[5], 64)).(float64),
	}
}

func stopForSncf(s []string) *Stop {
	var stopId = s[0][strings.Index(s[0], ":DUA")+4:]
	return &Stop{
		StopId:   uint32(PanicIfErrOrReturn(strconv.ParseUint(stopId, 10, 32)).(uint64)),
		StopName: s[1],
		StopLat:  PanicIfErrOrReturn(strconv.ParseFloat(s[3], 64)).(float64),
		StopLon:  PanicIfErrOrReturn(strconv.ParseFloat(s[4], 64)).(float64),
	}
}

func loadStops(filename string) []*Stop {
	file := PanicIfErrOrReturn(os.Open(filename)).(*os.File)
	defer file.Close()

	var typeOfStopTimesToLoad = UNKNOWN_DATA_TYPE
	if strings.Contains(filename, "ratp") {
		typeOfStopTimesToLoad = RATP_DATA_TYPE
	} else if strings.Contains(filename, "sncf") {
		typeOfStopTimesToLoad = SNCF_DATA_TYPE
	}

	reader := bufio.NewReader(file)
	// ignore first line
	PanicIfErrOrReturn(reader.ReadString(byte('\n')))

	var data []*Stop
	for ;; {
		line, err := reader.ReadString(byte('\n'))
		if line != "" {
			line = line[:len(line)-1]
			s := readCsvLine(line)
			if typeOfStopTimesToLoad == UNKNOWN_DATA_TYPE {
				L("unkown data type, support only ratp and sncf")
				os.Exit(1)
				return nil
			} else if typeOfStopTimesToLoad == RATP_DATA_TYPE {
				data = append(data, stopForRatp(s))
			} else if typeOfStopTimesToLoad == SNCF_DATA_TYPE {
				if strings.HasPrefix(s[0], "StopArea") {
					continue
				}
				data = append(data, stopForSncf(s))
			}
		}
		if err != nil && err == io.EOF {
			break
		}
		PanicIfErr(err)
	}
	return data
}

type Transfer struct {
	FromStopId uint32
	ToStopId uint32
	TransferTime uint16
}

func transferForRatp(s []string) *Transfer {
	var fromStopId = s[0][strings.Index(s[0], ":DUA")+4:]
	var toStopId = s[0][strings.Index(s[1], ":DUA")+4:]
	return &Transfer{
		FromStopId:    uint32(PanicIfErrOrReturn(strconv.ParseUint(fromStopId, 10, 32)).(uint64)),
		ToStopId:      uint32(PanicIfErrOrReturn(strconv.ParseUint(toStopId, 10, 32)).(uint64)),
		TransferTime:  uint16(PanicIfErrOrReturn(strconv.ParseUint(s[3], 10, 16)).(uint64)),
	}
}

func transferForSncf(s []string) *Transfer {
	var fromStopId = s[0][strings.Index(s[0], ":DUA")+4:]
	var toStopId = s[0][strings.Index(s[1], ":DUA")+4:]
	return &Transfer{
		FromStopId:    uint32(PanicIfErrOrReturn(strconv.ParseUint(fromStopId, 10, 32)).(uint64)),
		ToStopId:      uint32(PanicIfErrOrReturn(strconv.ParseUint(toStopId, 10, 32)).(uint64)),
		TransferTime:  uint16(PanicIfErrOrReturn(strconv.ParseUint(s[3][:len(s[3])-1], 10, 16)).(uint64)),
	}
}

func loadTransfers(filename string) []*Transfer {
	file := PanicIfErrOrReturn(os.Open(filename)).(*os.File)
	defer file.Close()

	var typeOfStopTimesToLoad = UNKNOWN_DATA_TYPE
	if strings.Contains(filename, "ratp") {
		typeOfStopTimesToLoad = RATP_DATA_TYPE
	} else if strings.Contains(filename, "sncf") {
		typeOfStopTimesToLoad = SNCF_DATA_TYPE
	}

	reader := bufio.NewReader(file)
	// ignore first line
	PanicIfErrOrReturn(reader.ReadString(byte('\n')))

	var data []*Transfer
	for ;; {
		line, err := reader.ReadString(byte('\n'))
		if line != "" {
			line = line[:len(line)-1]
			s := strings.Split(line, ",")
			if typeOfStopTimesToLoad == UNKNOWN_DATA_TYPE {
				L("unkown data type, support only ratp and sncf")
				os.Exit(1)
				return nil
			} else if typeOfStopTimesToLoad == RATP_DATA_TYPE {
				data = append(data, transferForRatp(s))
			} else if typeOfStopTimesToLoad == SNCF_DATA_TYPE {
				if strings.HasPrefix(s[0], "StopArea") {
					continue
				}
				data = append(data, transferForSncf(s))
			}
		}
		if err != nil && err == io.EOF {
			break
		}
		PanicIfErr(err)
	}
	return data
}

func LoadStopsFromDisk(dataDir string) []*Stop {
	stops := loadStops(dataDir + "/stops.txt")
	L("done loading stops", len(stops))
	return stops
}

func LoadDataFromDisk(dataDir string) ([]*StopTime, []*Stop, []*Transfer) {
	stopTimes := loadStopTimes(dataDir + "/stop_times.txt")
	L("done loading trips", len(stopTimes))
	stops := loadStops(dataDir + "/stops.txt")
	L("done loading stops", len(stops))
	transfers := loadTransfers(dataDir + "/transfers.txt")
	L("done loading transfer", len(transfers))
	return stopTimes, stops, transfers
}