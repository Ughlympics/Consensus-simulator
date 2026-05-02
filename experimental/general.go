package experimental

import (
	"math/rand"
	"sort"
)

type Node struct {
	ID       int
	Stake    float64
	VotedFor int
	IsWhale  bool
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

type GlobalMetrics struct {
	TotalSlots     int
	ProducedBlocks int
	ValidBlocks    int
	InvalidBlocks  int

	Availability float64
	InvalidShare float64
	ForkRate     float64
	Latency      float64

	Capture       bool
	Concentration float64
}

func GenerateNodes(N int, gamma float64) []Node {
	nodes := make([]Node, 0, N)

	whaleCount := int(0.01 * float64(N))
	if whaleCount < 1 {
		whaleCount = 1
	}

	totalStake := 10000.0
	otherCount := N - whaleCount

	whaleTotal := totalStake * gamma
	otherTotal := totalStake * (1 - gamma)

	whaleStake := whaleTotal / float64(whaleCount)
	otherStake := otherTotal / float64(otherCount)

	for i := 0; i < whaleCount; i++ {
		nodes = append(nodes, Node{
			ID:      i,
			Stake:   whaleStake,
			IsWhale: true,
		})
	}

	for i := 0; i < otherCount; i++ {
		nodes = append(nodes, Node{
			ID:      whaleCount + i,
			Stake:   otherStake,
			IsWhale: false,
		})
	}

	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	return nodes
}

func SelectCandidates(nodes []Node, M int) []Candidate {
	candidates := make([]Candidate, 0, M)

	shuffled := make([]Node, len(nodes))
	copy(shuffled, nodes)

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	for i := 0; i < M; i++ {
		candidates = append(candidates, Candidate{
			ID: shuffled[i].ID,
		})
	}

	return candidates
}

func SelectDelegates(candidates []Candidate, K int) []Delegate {
	if K > len(candidates) {
		K = len(candidates)
	}

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

func GenerateBlocks(
	delegates []Delegate,
	E int,
	alpha float64,
	pOff float64,
	pInv float64,
	tau float64) []Block {

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

		// if delegate is malicious, block is invalid with probability pInv
		if d.IsMalicious && rand.Float64() < pInv {
			valid = false
		}

		forked := false
		if rand.Float64() < tau {
			forked = true
		}

		block := Block{
			Round:    r,
			Producer: d.ID,
			Valid:    valid,
			Forked:   forked,
		}

		blocks = append(blocks, block)
	}

	return blocks
}

func ComputeMetrics(blocks []Block, totalSlots int, delegates []Delegate) GlobalMetrics {
	produced := len(blocks)
	valid := 0
	invalid := 0
	forks := 0

	lastRound := -1
	totalGap := 0
	gapCount := 0

	for _, b := range blocks {
		if b.Valid {
			valid++
		} else {
			invalid++
		}

		if b.Forked {
			forks++
		}

		// latency (gap between produced blocks)
		if lastRound != -1 {
			gap := b.Round - lastRound
			totalGap += gap
			gapCount++
		}
		lastRound = b.Round
	}

	availability := 0.0
	if totalSlots > 0 {
		availability = float64(valid) / float64(totalSlots)
	}

	invalidShare := 0.0
	forkRate := 0.0

	if produced > 0 {
		invalidShare = float64(invalid) / float64(produced)
		forkRate = float64(forks) / float64(produced)
	}

	latency := 0.0
	if gapCount > 0 {
		latency = float64(totalGap) / float64(gapCount)
	}

	//Capture
	malicious := 0
	for _, d := range delegates {
		if d.IsMalicious {
			malicious++
		}
	}
	capture := float64(malicious)/float64(len(delegates)) > 0.5

	//Concentration (top-3 delegates by votes / total votes)
	sort.Slice(delegates, func(i, j int) bool {
		return delegates[i].Votes > delegates[j].Votes
	})

	top := 3
	if len(delegates) < top {
		top = len(delegates)
	}

	topVotes := 0.0
	totalVotes := 0.0

	for _, d := range delegates {
		totalVotes += d.Votes
	}

	for i := 0; i < top; i++ {
		topVotes += delegates[i].Votes
	}

	concentration := 0.0
	if totalVotes > 0 {
		concentration = topVotes / totalVotes
	}

	return GlobalMetrics{
		TotalSlots:     totalSlots,
		ProducedBlocks: produced,

		ValidBlocks:   valid,
		InvalidBlocks: invalid,

		Availability: availability,
		InvalidShare: invalidShare,
		ForkRate:     forkRate,
		Latency:      latency,

		Capture:       capture,
		Concentration: concentration,
	}
}

func AssignMalicious(delegates []Delegate, alpha float64) {
	for i := range delegates {
		if rand.Float64() < alpha {
			delegates[i].IsMalicious = true
		}
	}
}
