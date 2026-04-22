package experimental

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
