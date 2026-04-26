package experimental

import (
	"fmt"
)

func Step1() {
	N := 1000
	M := 20
	K := 11
	E := 30

	nodes := GenerateNodes(N, 0.1)
	candidates := SelectCandidates(nodes, M)
	delagates := SelectDelegates(candidates, K)
	blocks := GenerateBlocks(delagates, E)

	stat := ComputeMetrics(blocks)
	fmt.Println(stat)
}
