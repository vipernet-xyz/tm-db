package db

import (
	"math"

	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/syndtr/goleveldb/leveldb/errors"
)

const (
	GoLevelMemDBBackend BackendType = "golevelmemdb"
)

func init() {
	dbCreator := func(name string, dir string) (DB, error) {
		return NewGoLevelMemDB(), nil
	}
	registerDBCreator(GoLevelMemDBBackend, dbCreator, false)
}

type GoLevelMemDB struct {
	db *memdb.DB
}

var _ DB = (*GoLevelMemDB)(nil)

func NewGoLevelMemDB() *GoLevelMemDB {
	return &GoLevelMemDB{memdb.New(comparer.DefaultComparer, math.MaxInt32)}
}

func NewGoLevelMemDBWithCapacity(capacity int) *GoLevelMemDB {
	return &GoLevelMemDB{memdb.New(comparer.DefaultComparer, capacity)}
}

// Get implements DB.
func (db *GoLevelMemDB) Get(key []byte) ([]byte, error) {
	key = nonNilBytes(key)
	res, err := db.db.Get(key)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

// Has implements DB.
func (db *GoLevelMemDB) Has(key []byte) (bool, error) {
	bytes, err := db.Get(key)
	if err != nil {
		return false, err
	}
	return bytes != nil, nil
}

// Set implements DB.
func (db *GoLevelMemDB) Set(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	if err := db.db.Put(key, value); err != nil {
		return err
	}
	return nil
}

// SetSync implements DB.
func (db *GoLevelMemDB) SetSync(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	if err := db.db.Put(key, value); err != nil {
		return err
	}
	return nil
}

// Delete implements DB.
func (db *GoLevelMemDB) Delete(key []byte) error {
	key = nonNilBytes(key)
	if err := db.db.Delete(key); err != nil {
		return err
	}
	return nil
}

// DeleteSync implements DB.
func (db *GoLevelMemDB) DeleteSync(key []byte) error {
	key = nonNilBytes(key)
	err := db.db.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (db *GoLevelMemDB) DB() *memdb.DB {
	return db.db
}

// Close implements DB.
func (db *GoLevelMemDB) Close() error {
	db.db.Free()
	return nil
}

// Print implements DB.
func (db *GoLevelMemDB) Print() error {
	return nil
}

// Stats implements DB.
func (db *GoLevelMemDB) Stats() map[string]string {
	return nil
}

// Manual compaction of the leveldb
func (db *GoLevelMemDB) Compact(r util.Range) error {
	return nil
}

// NewBatch implements DB.
func (db *GoLevelMemDB) NewBatch() Batch {
	return newgoMemDBBatch(db)
}

// Iterator implements DB.
func (db *GoLevelMemDB) Iterator(start, end []byte) (Iterator, error) {
	itr := db.db.NewIterator(&util.Range{})
	return newGoMemDBIterator(itr, start, end, false), nil
}

// ReverseIterator implements DB.
func (db *GoLevelMemDB) ReverseIterator(start, end []byte) (Iterator, error) {
	itr := db.db.NewIterator(&util.Range{})
	return newGoMemDBIterator(itr, start, end, true), nil
}
