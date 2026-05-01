package experimental

func RunSimulation(
	N int,
	M int,
	K int,
	R int,
	E int,
	gamma float64,
	alpha float64,
	pOff float64,
	pInv float64,
	tau float64,
	useWhales bool,
) GlobalMetrics {

	nodes := GenerateNodes(N, gamma)

	allBlocks := []Block{}

	roundsDone := 0

	for roundsDone < R {

		candidates := SelectCandidates(nodes, M)

		if useWhales {
			candidates = VoteWhaleForOneCandidate(nodes, candidates, K)
		} else {
			VoteRandom(nodes, candidates)
		}

		delegates := SelectDelegates(candidates, K)

		roundsLeft := R - roundsDone
		currentE := E
		if roundsLeft < E {
			currentE = roundsLeft
		}

		blocks := GenerateBlocks(delegates, currentE, alpha, pOff, pInv, tau)

		allBlocks = append(allBlocks, blocks...)

		roundsDone += currentE
	}

	return ComputeMetrics(allBlocks)
}
