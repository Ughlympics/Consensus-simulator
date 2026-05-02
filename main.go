package main

import (
	"DPos/experimental"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	//experimental.Step1()

	stat := experimental.RunSimulation(2500, 50, 21, 2500, 50, 0.4, 0.2, 0.15, 0.15, 0.1, true)
	experimental.PrintMetrics(stat)
}
