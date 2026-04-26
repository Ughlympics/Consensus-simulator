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

func GenerateBlocks(delegates []Delegate, E int) []Block {

	blocks := make([]Block, 0)

	for r := 0; r < E; r++ {
		d := delegates[r%len(delegates)]

		valid := true

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
