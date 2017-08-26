package r2aft_test

import "testing"
import "github.com/r2p2/r2aft"

type MockRemoteNode struct {
	id uint64

	term uint64
	candidateId uint64
	leaderId uint64
	lastLogIndex uint64
	prevLogIndex uint64
	lastLogTerm uint64
	prevLogTerm uint64
	entries []r2aft.Entry
	leaderCommit uint64
}

func NewMockRemoteNode(id uint64) r2aft.RemoteNode {
	return &MockRemoteNode {
		id: id,
	}
}

func (self *MockRemoteNode) Id() uint64 {
	return self.id;
}

func (self *MockRemoteNode) RequestVote(
	term uint64,
	candidateId uint64,
	lastLogIndex uint64,
	lastLogTerm uint64,
) {
	self.term = term
	self.candidateId = candidateId
	self.lastLogTerm = lastLogTerm
	self.lastLogTerm = lastLogTerm
}

func (self *MockRemoteNode) AppendEntries(
	term uint64,
	leaderId uint64,
	prevLogIndex uint64,
	prevLogTerm uint64,
	entries []r2aft.Entry,
	leaderCommit uint64,
) {
	self.term = term
	self.prevLogIndex = prevLogIndex
	self.prevLogTerm = prevLogTerm
	self.entries = entries
	self.leaderCommit = leaderCommit
}

func TestCreateLocalNode(t *testing.T) {
	nodeId := uint64(1)
	node := r2aft.NewNode(nodeId)
	if node.Id() != nodeId {
		t.Fail()
	}

	if node.State() != r2aft.Follower {
		t.Fail()
	}
}

func TestBecomingCandidate(t *testing.T) {
	node        := r2aft.NewNode(3)
	remoteNode1 := NewMockRemoteNode(1)
	remoteNode2 := NewMockRemoteNode(2)

	node.Add(remoteNode1)
	node.Add(remoteNode2)

	if node.State() != r2aft.Follower {
		t.Fail()
	}

	node.Timeout()

	if node.State() != r2aft.Candidate {
		t.Fail()
	}

	if remoteNode1.(*MockRemoteNode).term != 1 {
		t.Fail()
	}

	if remoteNode1.(*MockRemoteNode).candidateId != 3 {
		t.Fail()
	}

	if remoteNode2.(*MockRemoteNode).term != 1 {
		t.Fail()
	}

	if remoteNode2.(*MockRemoteNode).candidateId != 3 {
		t.Fail()
	}
}

func TestBecomingFollower(t *testing.T) {
	node        := r2aft.NewNode(3)
	remoteNode1 := NewMockRemoteNode(1)
	remoteNode2 := NewMockRemoteNode(2)

	node.Add(remoteNode1)
	node.Add(remoteNode2)

	if node.State() != r2aft.Follower {
		t.Fail()
	}

	node.Timeout()

	if node.State() != r2aft.Candidate {
		t.Fail()
	}

	if remoteNode1.(*MockRemoteNode).term != 1 {
		t.Fail()
	}

	if remoteNode1.(*MockRemoteNode).candidateId != 3 {
		t.Fail()
	}

	if remoteNode2.(*MockRemoteNode).term != 1 {
		t.Fail()
	}

	if remoteNode2.(*MockRemoteNode).candidateId != 3 {
		t.Fail()
	}

	term, err := node.AppendEntries(2, 3, 0, 0, make([]r2aft.Entry, 0), 0)
	if err != nil || term != 2 {
		t.Fail()
	}
}

func TestBecomingLeader(t *testing.T) {
	node        := r2aft.NewNode(3)
	remoteNode1 := NewMockRemoteNode(1)
	remoteNode2 := NewMockRemoteNode(2)

	node.Add(remoteNode1)
	node.Add(remoteNode2)

	if node.State() != r2aft.Follower {
		t.Fail()
	}

	node.Timeout()

	if node.State() != r2aft.Candidate {
		t.Fail()
	}

	if remoteNode1.(*MockRemoteNode).term != 1 {
		t.Fail()
	}

	if remoteNode1.(*MockRemoteNode).candidateId != 3 {
		t.Fail()
	}

	if remoteNode2.(*MockRemoteNode).term != 1 {
		t.Fail()
	}

	if remoteNode2.(*MockRemoteNode).candidateId != 3 {
		t.Fail()
	}

	node.VoteReply(1, nil) // grant vote from remote node 1

	if node.State() != r2aft.Candidate {
		t.Log("not in state candidate after first granted vote");
		t.Fail()
	}

	node.VoteReply(1, nil) // grant vote from remote node 2

	if node.State() != r2aft.Leader {
		t.Log("not in state follower after second granted vote");
		t.Fail()
	}
}

func TestBecomingGrantOneVote(t *testing.T) {
	node        := r2aft.NewNode(3)
	remoteNode1 := NewMockRemoteNode(1)
	remoteNode2 := NewMockRemoteNode(2)

	node.Add(remoteNode1)
	node.Add(remoteNode2)

	if node.State() != r2aft.Follower {
		t.Fail()
	}

	term, err := node.RequestVote(1, 1, 0, 0)
	if err != nil || term != 1 {
		t.Fail()
	}

	term, err = node.RequestVote(1, 2, 0, 0)
	if err == nil {
		t.Log("only one node can be voted for in a cycle")
		t.Fail()
	}

	term, err = node.RequestVote(2, 1, 0, 0)
	if err != nil || term != 2 {
		t.Log("another vote on new term is possible")
		t.Fail()
	}
}
