package db

import "github.com/pkg/errors"

// goMemDBBatch operations

// goMemDBBatch handles in-memory batching.
type goMemDBBatch struct {
	db  *GoLevelMemDB
	ops []operation
}

var _ Batch = (*goMemDBBatch)(nil)

// newgoMemDBBatch creates a new goMemDBBatch
func newgoMemDBBatch(db *GoLevelMemDB) *goMemDBBatch {
	return &goMemDBBatch{
		db:  db,
		ops: []operation{},
	}
}

func (b *goMemDBBatch) assertOpen() {
	if b.ops == nil {
		panic("batch has been written or closed")
	}
}

// Set implements Batch.
func (b *goMemDBBatch) Set(key, value []byte) {
	b.assertOpen()
	b.ops = append(b.ops, operation{opTypeSet, key, value})
}

// Delete implements Batch.
func (b *goMemDBBatch) Delete(key []byte) {
	b.assertOpen()
	b.ops = append(b.ops, operation{opTypeDelete, key, nil})
}

// Write implements Batch.
func (b *goMemDBBatch) Write() error {
	b.assertOpen()

	for _, op := range b.ops {
		switch op.opType {
		case opTypeSet:
			b.db.Set(op.key, op.value)
		case opTypeDelete:
			b.db.Delete(op.key)
		default:
			return errors.Errorf("unknown operation type %v (%v)", op.opType, op)
		}
	}

	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	b.Close()
	return nil
}

// WriteSync implements Batch.
func (b *goMemDBBatch) WriteSync() error {
	return b.Write()
}

// Close implements Batch.
func (b *goMemDBBatch) Close() {
	b.ops = nil
}
