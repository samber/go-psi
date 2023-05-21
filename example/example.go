package main

import (
	"fmt"

	"github.com/samber/go-psi"
)

func main() {
	stats, err := psi.AllPSIStats()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("PSI stats:\n------\n\n%s\n", stats)

	stream, done, err := psi.NotifyStarvation(psi.CPU, psi.Avg10, 70, 90)
	for {
		last, _ := <-stream
		fmt.Printf("\nALERT %t\nCPU: %f%%\n", last.Starved, last.Current)
	}

	<-done // never called 🤡
}
