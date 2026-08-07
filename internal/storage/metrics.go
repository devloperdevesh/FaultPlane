package storage

type Metrics struct {
	Saves   uint64
	Loads   uint64
	Deletes uint64

	Items int
}
