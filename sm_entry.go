package r2aft

type Entry interface {
	Term() uint64
	Apply() uint64
}
