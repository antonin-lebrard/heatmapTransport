package pkg

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
)

func L(smt ...interface{}) {
	log.Println(smt...)
}

func PanicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func PanicIfErrOrReturn(smt interface{}, err error) interface{} {
	PanicIfErr(err)
	return smt
}

func secondsToReadableDurationString(seconds float64) string {
	var hoursRemaining = int(seconds / 60.0 / 60.0)
	var minRemaining = (int(seconds) % (60 * 60)) / 60
	var secondsRemaining = int(seconds) % 60
	return strconv.Itoa(hoursRemaining) + "h" + strconv.Itoa(minRemaining) + "m" + strconv.Itoa(secondsRemaining) + "s"
}

func printRemainingTime(startupTime time.Time, nbDone, nbTotal int) {
	var duration = time.Since(startupTime)
	var dividendForNbDone = float64(nbTotal) / float64(nbDone)
	var totalExpectedTime = duration.Seconds() * dividendForNbDone
	var remainingTime = totalExpectedTime - duration.Seconds()

	L(secondsToReadableDurationString(duration.Seconds()), "spent, remaining:", secondsToReadableDurationString(remainingTime))
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
var m runtime.MemStats
func PrintMemUsage() {
	L("")
	L("Before GC")
	fmt.Printf("Alloc = %v MiB\n", m.Alloc / 1024 / 1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc / 1024 / 1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys / 1024 / 1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)
	runtime.GC()
	runtime.ReadMemStats(&m)
	L("")
	L("After GC")
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB\n", m.Alloc / 1024 / 1024)
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc / 1024 / 1024)
	fmt.Printf("Sys = %v MiB\n", m.Sys / 1024 / 1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)
	L("")
}
