package r2aft

import "errors"

type State int
const (
	Follower State = iota
	Candidate
	Leader
)

type Node struct {
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
	remoteNodes []RemoteNode
}

func NewNode(id uint64) Node {
	return Node {
		id: id,
		state: Follower,
		currentTerm: 0,
		votesReceived: 0,
		votedFor: 0,
		log: make([]Entry, 0, 100),
		commitIndex: 0,
		lastApplied: 0,
		remoteNodes: make([]RemoteNode, 0, 10),
	}
}

func (self *Node) Add(remoteNode RemoteNode) {
	self.remoteNodes = append(self.remoteNodes, remoteNode)
}

func (self *Node) State() State {
	return self.state
}

func (self *Node) Id() uint64 {
	return self.id
}

func (self *Node) RequestVote(
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

	// TODO Node::RequestVote: Add log checks $5.2, $5.4

	self.votedFor = candidateId
	return self.currentTerm, nil
}

func (self *Node) VoteReply(
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

func (self *Node) AppendEntries(
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

	// TODO Node::AppendEntries impl check for $5.3
	// TODO Node::AppendEntries append entries not already in log
	// TODO Node::AppendEntries update leader commit


	if term > self.currentTerm {
		self.currentTerm = term;
		self.votedFor = 0
	}

	return self.currentTerm, nil
}

func (self *Node) AppendEntriesReply(
	term uint64,
	err error,
) {
	if term > self.currentTerm {
		self.currentTerm = term
		self.state = Follower
		self.votedFor = 0
	}
}

func (self *Node) Timeout() {
	if self.state == Leader {
		self.emitHeartbeat()
		return
	}

	self.startElection()
}

func (self *Node) startElection() {
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

func (self *Node) becomeLeader() {
	self.state = Leader
	self.emitHeartbeat()
}

func (self *Node) emitHeartbeat() {
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
