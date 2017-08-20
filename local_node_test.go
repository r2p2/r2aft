package r2aft_test

import "testing"
import "github.com/r2p2/r2aft"

func TestCreateLocalNode(t *testing.T) {
	nodeId := uint64(1)
	node := r2aft.New(nodeId)
	if node.Id() != nodeId {
		t.Fail()
	}

	if node.State() != r2aft.Follower {
		t.Fail()
	}
}

func TestBecomingCandidate(t *testing.T) {
	nodeId := uint64(1)
	node := r2aft.New(nodeId)

	if node.State() != r2aft.Follower {
		t.Fail()
	}

	node.Timeout()

	if node.State() != r2aft.Candidate {
		t.Fail()
	}
}
