package fsm

import (
	"encoding/binary"
	"io"
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

// Open opens the bbolt backend.
// todo (sienna): leverage stopc at some point on bbolt.Open
func (b *BBoltStateMachine) Open(stopc <-chan struct{}) (uint64, error) {
	var err error
	b.db, err = bbolt.Open(filepath.Join(b.BasePath, strconv.FormatUint(b.ClusterId, 10), strconv.FormatUint(b.NodeId, 10)), 0600, b.Options)
	if err != nil {
		panic(err)
	}

	var val uint64
	err = b.db.View(func(tx *bbolt.Tx) error {
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
		return -1, err
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

}

func (b *BBoltStateMachine) Close() error {
	return b.db.Close()
}
