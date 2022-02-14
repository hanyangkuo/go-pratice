//go:build windows
// +build windows

package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"
)

func main() {
	//v, _ := mem.VirtualMemory()
	//
	//// almost every return value is a struct
	//fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	//// convert to JSON. String() is also implemented
	//fmt.Println(v)
	//var mod = syscall.NewLazyDLL("kernel32.dll")
	//var proc = mod.NewProc("GetPhysicallyInstalledSystemMemory")
	//var mem uint64
	//
	//ret, _, err := proc.Call(uintptr(unsafe.Pointer(&mem)))
	//fmt.Printf("Ret: %d, err: %v, Physical memory: %f\n", ret, err, bToMb(mem))


	//Print our starting memory usage (should be around 0mb)
	PrintMemUsage()

	var overall [][]int
	for i := 0; i<10; i++ {

		// Allocate memory using make() and append to overall (so it doesn't get
		// garbage collected). This is to create an ever increasing memory usage
		// which we can track. We're just using []int as an example.
		a := make([]int, 0, 999999)
		overall = append(overall, a)

		// Print our memory usage at each interval
		PrintMemUsage()
		time.Sleep(time.Second)
	}
	for {
		PrintMemUsage()
		time.Sleep(time.Second)
	}
	// Clear our memory and print usage, unless the GC has run 'Alloc' will remain the same
	overall = nil
	PrintMemUsage()

	// Force GC to clear up, should see a memory drop
	runtime.GC()
	PrintMemUsage()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	// Method 1.
	p := &process.Process{Pid: int32(os.Getpid())}
	m, err := p.MemoryInfo()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("RSS: %f\tHWM: %f\tVMS: %d\n", bToMb(m.RSS), bToMb(m.HWM), m.VMS)

	// Method 2.
	//v, _ := mem.VirtualMemory()
	//
	//// almost every return value is a struct
	//fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)


	// Method 3.
	//var m runtime.MemStats
	//runtime.ReadMemStats(&m)
	//// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	//
	//fmt.Printf("\tStackSys = %v MiB", bToMb(m.StackSys))
	//fmt.Printf("\tHeapSys = %v MiB", bToMb(m.HeapSys))
	//fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	//
	//fmt.Printf("\tHeapAlloc = %v MiB", bToMb(m.HeapAlloc))
	//fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	//
	//fmt.Printf("\tNumGC = %v\n", m.NumGC)

	// Method 4. Physical memory
	//// convert to JSON. String() is also implemented
	//fmt.Println(v)
	//var mod = syscall.NewLazyDLL("kernel32.dll")
	//var proc = mod.NewProc("GetPhysicallyInstalledSystemMemory")
	//var mem uint64
	//
	//ret, _, err := proc.Call(uintptr(unsafe.Pointer(&mem)))
	//fmt.Printf("Ret: %d, err: %v, Physical memory: %f\n", ret, err, bToMb(mem))
}

func bToMb(b uint64) float32 {
	return float32(b) / 1024 / 1024
}


func GetCurrentProcessTimes() {
	handle, _ := syscall.GetCurrentProcess()
	var u syscall.Rusage
	e := syscall.GetProcessTimes(handle, &u.CreationTime, &u.ExitTime, &u.KernelTime, &u.UserTime)
	if e != nil {
		log.Println(e)
	}
}
