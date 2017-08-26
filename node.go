package r2aft

type LNode interface {
	Id() uint64

	RequestVote(
		term uint64,
		candidateId uint64,
		lastLogIndex uint64,
		lastLogTerm uint64,
	) (uint64, error)

	VoteReply(
		term uint64,
		err error,
	)

	AppendEntries(
		term uint64,
		leaderId uint64,
		prevLogIndex uint64,
		prevLogTerm uint64,
		entries []Entry,
		leaderCommit uint64,
	) (uint64, error)

	AppendEntriesReply(
		term uint64,
		err error,
	)
}

type RNode interface {
	Id() uint64

	RequestVote(
		term uint64,
		candidateId uint64,
		lastLogIndex uint64,
		lastLogTerm uint64,
	)

	AppendEntries(
		term uint64,
		leaderId uint64,
		prevLogIndex uint64,
		prevLogTerm uint64,
		entries []Entry,
		leaderCommit uint64,
	)
}
