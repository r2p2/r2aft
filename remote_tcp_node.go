package r2aft

type RemoteTcpNode struct {
	id uint64
	nextIndex []uint64
	matchIndex []uint64
	localNode LocalNode
}

func (self *RemoteTcpNode) Id() uint64 {
	return self.id
}

func (self *RemoteTcpNode) RequestVote(
	term uint64,
	candidateId uint64,
	lastLogIndex uint64,
	lastLogTerm uint64,
) {
}

func (self *RemoteTcpNode) AppendEntries(
	term uint64,
	leaderId uint64,
	prevLogIndex uint64,
	prevLogTerm uint64,
	entries []Entry,
	leaderCommit uint64,
) {
}
