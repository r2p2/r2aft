package r2aft

type RemoteNode struct {
	id uint64
	nextIndex []uint64
	matchIndex []uint64
	localNode LocalNode
}

func (self *RemoteNode) Id() uint64 {
	return self.id
}

func (self *RemoteNode) RequestVote(
	term uint64,
	candidateId uint64,
	lastLogIndex uint64,
	lastLogTerm uint64,
) {
}

func (self *RemoteNode) AppendEntries(
	term uint64,
	leaderId uint64,
	prevLogIndex uint64,
	prevLogTerm uint64,
	entries []Entry,
	leaderCommit uint64,
) {
}
