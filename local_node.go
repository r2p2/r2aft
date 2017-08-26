package r2aft

import "errors"

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
	votedFor uint64
	log []Entry

	// volatile on all servers
	commitIndex uint64
	lastApplied uint64

	// volatile on leaders; reinit after election
	remoteNodes []RNode
}

func New(id uint64) LocalNode {
	return LocalNode {
		id: id,
		state: Follower,
		currentTerm: 0,
		votesReceived: 0,
		votedFor: 0,
		log: make([]Entry, 0, 100),
		commitIndex: 0,
		lastApplied: 0,
		remoteNodes: make([]RNode, 0, 10),
	}
}

func (self *LocalNode) Add(remoteNode RNode) {
	self.remoteNodes = append(self.remoteNodes, remoteNode)
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
	if term < self.currentTerm {
		return self.currentTerm, errors.New("term < current term")
	} else if term > self.currentTerm {
		self.currentTerm = term
		self.votedFor = 0
	}

	if self.votedFor != 0 && self.votedFor != candidateId {
		return self.currentTerm, errors.New("did vote already")
	}

	// TODO LocalNode::RequestVote: Add log checks $5.2, $5.4

	self.votedFor = candidateId
	return self.currentTerm, nil
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
	if self.votesReceived <= uint64((len(self.remoteNodes))/2) {
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
	if term < self.currentTerm {
		return self.currentTerm, errors.New("term < current term")
	}

	// TODO LocalNode::AppendEntries impl check for $5.3
	// TODO LocalNode::AppendEntries append entries not already in log
	// TODO LocalNode::AppendEntries update leader commit


	if term > self.currentTerm {
		self.currentTerm = term;
		self.votedFor = 0
	}

	return self.currentTerm, nil
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
	self.votesReceived = 0
	self.votedFor = self.Id();

	var lastLogTerm uint64 = 0
	if len(self.log) > 0 {
		lastLogTerm = self.log[len(self.log)-1].Term()
	}
	for _, rn := range self.remoteNodes {
		rn.RequestVote(
			self.currentTerm,
			self.id,
			uint64(len(self.log)),
			lastLogTerm,
		)
	}
}

func (self *LocalNode) becomeLeader() {
	self.state = Leader
	self.emitHeartbeat()
}

func (self *LocalNode) emitHeartbeat() {
	var lastLogTerm uint64 = 0
	if len(self.log) > 0 {
		lastLogTerm = self.log[len(self.log)-1].Term()
	}
	for _, rn := range self.remoteNodes {
		rn.AppendEntries(
			self.currentTerm,
			self.id,
			uint64(len(self.log)),
			lastLogTerm,
			self.log,
			self.commitIndex,
		)
	}
}
