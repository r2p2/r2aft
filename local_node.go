package r2aft

type State int
const (
	Follower State = iota
	Candidate
	Leader
)

type LocalNode struct {
	id uint64
	state State

	// persistent on all servers
	currentTerm uint64
	votesReceived uint64
	votedFor Node
	log []Entry

	// volatile on all servers
	commitIndex uint64
	lastApplied uint64

	// volatile on leaders; reinit after election
	remoteNodes []RemoteNode
}

func New(id uint64) LocalNode {
	return LocalNode {
		id: id,
		state: Follower,
		currentTerm: 0,
		votesReceived: 0,
		votedFor: nil,
		log: make([]Entry, 0, 100),
		commitIndex: 0,
		lastApplied: 0,
		remoteNodes: make([]RemoteNode, 0, 10),
	}
}

func (self *LocalNode) State() State {
	return self.state
}

func (self *LocalNode) Id() uint64 {
	return self.id
}

func (self *LocalNode) RequestVote(
	term uint64,
	candidateId uint64,
	lastLogIndex uint64,
	lastLogTerm uint64,
) (uint64, error) {
	return 0, nil
}

func (self *LocalNode) VoteReply(
	term uint64,
	err error,
) {
	if term != self.currentTerm {
		return
	} else if err != nil {
		return
	}

	self.votesReceived += 1
	if self.votesReceived <= uint64((len(self.remoteNodes)+1)/2) {
		return
	}

	self.becomeLeader()
}

func (self *LocalNode) AppendEntries(
	term uint64,
	leaderId uint64,
	prevLogIndex uint64,
	prevLogTerm uint64,
	entries []Entry,
	leaderCommit uint64,
) (uint64, error) {
	return 0, nil
}

func (self *LocalNode) AppendEntriesReply(
	term uint64,
	err error,
) {

}

func (self *LocalNode) Timeout() {
	if self.state == Leader {
		self.emitHeartbeat()
		return
	}

	self.startElection()
}

func (self *LocalNode) startElection() {
	self.state = Candidate
	self.currentTerm += 1
	self.votesReceived = 1
	self.votedFor = self
	for _, rn := range self.remoteNodes {
		rn.RequestVote(
			self.currentTerm,
			self.id,
			uint64(len(self.log)),
			self.log[len(self.log)-1].Term(),
		)
	}
}

func (self *LocalNode) becomeLeader() {
	self.state = Leader
}

func (self *LocalNode) emitHeartbeat() {
	for _, rn := range self.remoteNodes {
		rn.AppendEntries(
			self.currentTerm,
			self.id,
			uint64(len(self.log)),
			self.log[len(self.log)-1].Term(),
			self.log,
			self.commitIndex,
		)
	}
}
