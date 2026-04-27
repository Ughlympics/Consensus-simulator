package experimental

import (
	"math/rand"
	"sort"
)

type Node struct {
	ID       int
	Stake    float64
	VotedFor int
}

type Candidate struct {
	ID          int
	Votes       float64
	IsMalicious bool
}

type Delegate struct {
	Candidate
}

type Block struct {
	Producer int
	Valid    bool
	Round    int
	Forked   bool
}

type Metrics struct {
	TotalBlocks   int
	InvalidBlocks int
	InvalidRate   float64
}

func GenerateNodes(N int, gamma float64) []Node {
	nodes := make([]Node, 0, N)

	for i := 0; i < N; i++ {
		stake := gamma * float64(100) // Example stake calculation, can be modified as needed
		node := Node{ID: i, Stake: stake}
		nodes = append(nodes, node)
	}
	return nodes
}

func SelectCandidates(nodes []Node, M int) []Candidate {
	candidates := make([]Candidate, 0, M)

	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	for i := 0; i < M; i++ {
		candidates = append(candidates, Candidate{
			ID: nodes[i].ID,
		})
	}

	return candidates
}

func SelectDelegates(candidates []Candidate, K int) []Delegate {
	delegates := make([]Delegate, 0, K)

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Votes > candidates[j].Votes
	})

	for i := 0; i < K; i++ {
		delegates = append(delegates, Delegate{
			Candidate: candidates[i],
		})
	}

	return delegates
}

func GenerateBlocks(delegates []Delegate, E int, alpha float64, pOff float64) []Block {

	blocks := make([]Block, 0)

	AssignMalicious(delegates, alpha)

	for r := 0; r < E; r++ {
		d := delegates[r%len(delegates)]

		randVal := rand.Float64()

		// delegate is offline with probability pOff
		if randVal < pOff {
			continue
		}

		valid := true

		// if delegate is malicious, block is invalid
		if d.IsMalicious == true {
			valid = false
		}

		block := Block{
			Round:    r,
			Producer: d.ID,
			Valid:    valid,
		}

		blocks = append(blocks, block)
	}

	return blocks
}

func ComputeMetrics(blocks []Block) Metrics {
	total := len(blocks)
	invalid := 0

	for _, b := range blocks {
		if !b.Valid {
			invalid++
		}
	}

	invalidRate := 0.0
	if total > 0 {
		invalidRate = float64(invalid) / float64(total)
	}

	return Metrics{
		TotalBlocks:   total,
		InvalidBlocks: invalid,
		InvalidRate:   invalidRate,
	}
}

func AssignMalicious(delegates []Delegate, alpha float64) {
	for i := range delegates {
		if rand.Float64() < alpha {
			delegates[i].IsMalicious = true
		}
	}
}
