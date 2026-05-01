package experimental

import (
	"math/rand"
)

func VoteRandom(nodes []Node, candidates []Candidate) {
	for i := range candidates {
		candidates[i].Votes = 0
	}

	for _, n := range nodes {
		choice := rand.Intn(len(candidates))
		candidates[choice].Votes += n.Stake
	}
}

func VoteWhaleForOneCandidate(
	nodes []Node,
	candidates []Candidate,
	whaleIDs []int,
	K int) []Candidate {
	for i := range nodes {
		nodes[i].VotedFor = -1
	}

	for i := range candidates {
		candidates[i].Votes = 0
	}

	whales := []Node{}
	others := []Node{}

	for i := range nodes {
		if nodes[i].IsWhale {
			whales = append(whales, nodes[i])
		} else {
			others = append(others, nodes[i])
		}
	}

	// Each whale votes for a random candidate
	targetCount := K
	if len(whales) < K {
		targetCount = len(whales)
	}

	// Randomly select targetCount candidates for whales to vote for
	perm := rand.Perm(len(candidates))
	whaleChoices := perm[:targetCount]

	// Each whale votes for their chosen candidate
	for i := range whales {
		candidateIndex := whaleChoices[i%targetCount]

		whales[i].VotedFor = candidateIndex
		candidates[candidateIndex].Votes += whales[i].Stake
	}
	//Others vote randomly
	for i := range others {
		choice := rand.Intn(len(candidates))

		others[i].VotedFor = choice
		candidates[choice].Votes += others[i].Stake
	}

	return candidates
}
