package fsm

import (
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lni/dragonboat/v3/statemachine"
	"go.etcd.io/bbolt"
)

type BBoltStateMachine struct {
	ClusterId uint64
	NodeId    uint64
	BasePath  string
	Options   *bbolt.Options

	db *bbolt.DB
}

func NewBBoltStateMachine(clusterId uint64, nodeId uint64, basePath string, options *bbolt.Options) *BBoltStateMachine {
	return &BBoltStateMachine{ClusterId: clusterId, NodeId: nodeId, BasePath: basePath, Options: options}
}

func (b *BBoltStateMachine) dbPath() string {
	return filepath.Join(b.BasePath, strconv.FormatUint(b.ClusterId, 10), strconv.FormatUint(b.NodeId, 10))
}

// Open opens the bbolt backend.
// todo (sienna): leverage stopc at some point on bbolt.Open
func (b *BBoltStateMachine) Open(stopc <-chan struct{}) (uint64, error) {
	var err error
	b.db, err = bbolt.Open(b.dbPath(), 0600, b.Options)
	if err != nil {
		return 0, err
	}

	var val uint64
	err = b.db.Update(func(tx *bbolt.Tx) error {
		// todo (sienna): implement db stats on open
		//tx.Stats()

		internalBucket, err := tx.CreateBucketIfNotExists([]byte("monotonic-log"))
		if err != nil {
			return err
		}

		key, _ := internalBucket.Cursor().Last()
		val = binary.LittleEndian.Uint64(key)
		return nil
	})
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (b BBoltStateMachine) Update(entries []statemachine.Entry) ([]statemachine.Entry, error) {
	//TODO implement me
	panic("implement me")
}

func (b BBoltStateMachine) Lookup(i interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (b BBoltStateMachine) Sync() error {
	return b.db.Sync()
}

func (b BBoltStateMachine) PrepareSnapshot() (interface{}, error) {
	return nil, nil
}

func (b *BBoltStateMachine) SaveSnapshot(ctx interface{}, writer io.Writer, done <-chan struct{}) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.WriteTo(writer)
		return err
	})
}

func (b BBoltStateMachine) RecoverFromSnapshot(reader io.Reader, i <-chan struct{}) error {
	fn := func(r io.Reader) error {
		target, err := os.Create(b.dbPath())
		if err != nil {
			return err
		}
		_, err = io.Copy(target, reader)
		if err != nil {
			return err
		}
		return nil
	}

	// verify the existing database is closed
	err := b.db.Close()
	if err != nil {
		return err
	}

	_, err = os.Stat(b.dbPath())
	if err != nil {
		if os.IsNotExist(err) {
			return fn(reader)
		}
		return err
	}

	err = os.Remove(b.dbPath())
	if err != nil {
		return err
	}

	return fn(reader)
}

func (b *BBoltStateMachine) Close() error {
	return b.db.Close()
}
