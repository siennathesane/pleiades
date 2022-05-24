package fsm

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/lni/dragonboat/v3/statemachine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.uber.org/fx/fxtest"
	"r3t.io/pleiades/pkg/conf"
)

type TestBBoltFsm struct {
	suite.Suite
	lifecycle *fxtest.Lifecycle
	client    *api.Client
	env       *conf.EnvironmentConfig
}

func TestBBoltStateMachine(t *testing.T) {
	suite.Run(t, new(TestBBoltFsm))
}

func (bfsm *TestBBoltFsm) SetupSuite() {
	var err error
	bfsm.lifecycle = fxtest.NewLifecycle(bfsm.T())
	bfsm.client, err = conf.NewConsulClient(bfsm.lifecycle)
	require.Nil(bfsm.T(), err, "failed to connect to consul")
	require.NotNil(bfsm.T(), bfsm.client, "the consul client can't be empty")

	bfsm.env, err = conf.NewEnvironmentConfig(bfsm.client)
	require.Nil(bfsm.T(), err, "the environment config is needed")
	require.NotNil(bfsm.T(), bfsm.env, "the environment config must be rendered")
}

func (bfsm *TestBBoltFsm) BeforeTest(suiteName, testName string) {
	if err := os.RemoveAll(bfsm.env.BaseDir); err != nil {
		bfsm.T().Errorf("error cleaning test directory: %s", err)
	}
}

func (bfsm *TestBBoltFsm) TestNewBBoltStateMachine() {
	opts := &bbolt.Options{
		Timeout:         0,
		NoGrowSync:      false,
		NoFreelistSync:  false,
		FreelistType:    bbolt.FreelistMapType,
		ReadOnly:        false,
		InitialMmapSize: 0,
		PageSize:        0,
		NoSync:          false,
		OpenFile:        nil,
		Mlock:           false,
	}

	require.NotPanics(bfsm.T(), func() {
		NewBBoltStateMachine(1, 1, bfsm.env.BaseDir, opts)
	}, "creating a new bbolt fsm must not throw an error")
}

func (bfsm *TestBBoltFsm) TestBBoltStateMachineOpen() {
	testOpts := &bbolt.Options{
		Timeout:         0,
		NoGrowSync:      false,
		NoFreelistSync:  false,
		FreelistType:    bbolt.FreelistMapType,
		ReadOnly:        false,
		InitialMmapSize: 0,
		PageSize:        0,
		NoSync:          false,
		OpenFile:        nil,
		Mlock:           false,
	}

	fsm := NewBBoltStateMachine(1, 1, bfsm.env.BaseDir, testOpts)

	// verify the core path constructs for humans are available
	dbPath := fsm.dbPath(true)
	assert.Contains(bfsm.T(), dbPath, "cluster-")
	assert.Contains(bfsm.T(), dbPath, "node-")
	assert.Contains(bfsm.T(), dbPath, ".db")

	// verify there's no ".db"
	dbRootPath := fsm.dbPath(false)
	assert.Contains(bfsm.T(), dbRootPath, "cluster-")
	assert.Contains(bfsm.T(), dbRootPath, "node-")
	assert.NotContains(bfsm.T(), dbRootPath, ".db")

	// open a blank database
	var index uint64
	var err error
	require.NotPanics(bfsm.T(), func() {
		index, err = fsm.Open(make(<-chan struct{}, 1))
	})
	require.Equal(bfsm.T(), uint64(0), index, "opening a brand new fsm requires a 0 index")
	require.NoError(bfsm.T(), err, "there must not be an error opening a brand new fsm")
	require.NoError(bfsm.T(), fsm.db.Close(), "there must not be an error when closing the the brand new database")

	fi, err := os.Lstat(bfsm.env.BaseDir)
	val := fi.Mode().Perm()
	assert.NotEmpty(bfsm.T(), val)
	err = os.RemoveAll(bfsm.env.BaseDir)
	require.NoError(bfsm.T(), err, "there must not be an error when deleting the test directory")

	err = os.MkdirAll(fsm.dbPath(false), os.FileMode(dbDirModeVal))
	require.NoError(bfsm.T(), err, "there must not be an error when creating the database path")

	// create a database with an existing index
	db, err := bbolt.Open(fsm.dbPath(true), os.FileMode(dbFileModeVal), testOpts)
	require.NoError(bfsm.T(), err, "there must not be an error when opening the test database")
	require.NotNil(bfsm.T(), db, "the test database must be opened")

	// set the monotonic key, monotonicLogKey, to be an arbitrary value
	testIndexVal := uint64(55)
	err = db.Update(func(tx *bbolt.Tx) error {
		internalBucket, err := tx.CreateBucketIfNotExists([]byte(monotonicLogBucket))
		if err != nil {
			return err
		}

		val := make([]byte, 8)
		binary.LittleEndian.PutUint64(val, testIndexVal)
		return internalBucket.Put([]byte(monotonicLogKey), val)
	})
	require.NoError(bfsm.T(), err, "there must not be an error when prepping the monotonic log")

	err = db.Close()
	require.NoError(bfsm.T(), err, "there must not be an error when closing the test database")

	require.NotPanics(bfsm.T(), func() {
		index, err = fsm.Open(make(<-chan struct{}, 1))
	})
	require.Equal(bfsm.T(), testIndexVal, index, fmt.Sprintf("the existing index must be %d", testIndexVal))
	require.NoError(bfsm.T(), err, "there must not be an error opening an existing fsm")
	require.NoError(bfsm.T(), fsm.db.Close(), "there must not be an error when closing the database manually")
}

func (bfsm *TestBBoltFsm) TestBBoltStateMachineClose() {
	testOpts := &bbolt.Options{
		Timeout:         0,
		NoGrowSync:      false,
		NoFreelistSync:  false,
		FreelistType:    bbolt.FreelistMapType,
		ReadOnly:        false,
		InitialMmapSize: 0,
		PageSize:        0,
		NoSync:          false,
		OpenFile:        nil,
		Mlock:           false,
	}

	fsm := NewBBoltStateMachine(1, 1, bfsm.env.BaseDir, testOpts)

	index, err := fsm.Open(make(<-chan struct{}))
	require.NoError(bfsm.T(), err, "there must not be an error when opening the database")
	require.Equal(bfsm.T(), uint64(0), index, "the index must equal as there are no records")

	err = fsm.Close()
	require.NoError(bfsm.T(), err, "there must not be an error while closing the database")
	assert.Nil(bfsm.T(), fsm.db, "the database should be dereferenced")
	require.Panics(bfsm.T(), func() {
		_ = fsm.db.View(func(tx *bbolt.Tx) error {
			return nil
		})
	}, "there must be a nil reference error when trying to access the database after it's been closed")
}

func (bfsm *TestBBoltFsm) TestBBoltStateMachineUpdate() {
	testOpts := &bbolt.Options{
		Timeout:         0,
		NoGrowSync:      false,
		NoFreelistSync:  false,
		FreelistType:    bbolt.FreelistMapType,
		ReadOnly:        false,
		InitialMmapSize: 0,
		PageSize:        0,
		NoSync:          false,
		OpenFile:        nil,
		Mlock:           false,
	}

	fsm := NewBBoltStateMachine(1, 1, bfsm.env.BaseDir, testOpts)

	index, err := fsm.Open(make(<-chan struct{}))
	require.NoError(bfsm.T(), err, "there must not be an error when opening the database")
	require.Equal(bfsm.T(), uint64(0), index, "the index must equal as there are no records")

	rootPrn := &PleiadesResourceName{
		Partition:    GlobalPartition,
		Service:      Pleiades,
		Region:       GlobalRegion,
		AccountId:    testAccountKey,
		ResourceType: Bucket,
		ResourceId:   "test-bucket",
	}

	testKvps := []mvccpb.KeyValue{
		{
			Key:            []byte(rootPrn.ToFsmRootPath("test-bucket") + "/" + "test-key-0"),
			CreateRevision: 0,
			ModRevision:    0,
			Version:        1,
			Value:          []byte("test-value-0"),
			Lease:          0,
		},
		{
			Key:            []byte(rootPrn.ToFsmRootPath("test-bucket") + "/" + "test-key-1"),
			CreateRevision: 0,
			ModRevision:    0,
			Version:        1,
			Value:          []byte("test-value-1"),
			Lease:          0,
		},
		{
			Key:            []byte(rootPrn.ToFsmRootPath("test-bucket") + "/" + "test-key-2"),
			CreateRevision: 0,
			ModRevision:    0,
			Version:        1,
			Value:          []byte("test-value-2"),
			Lease:          0,
		},
	}

	testUpdates := make([]statemachine.Entry, 0)
	for idx := range testKvps {
		val, err := testKvps[0].Marshal()
		if err != nil {
			bfsm.T().Error(err)
		}
		testUpdates = append(testUpdates, statemachine.Entry{
			Index:  uint64(idx),
			Cmd:    val,
			Result: statemachine.Result{},
		})
	}

	var endingIndex []statemachine.Entry
	require.NotPanics(bfsm.T(), func() {
		endingIndex, err = fsm.Update(testUpdates)
	})
	require.NoError(bfsm.T(), err, "there must not be an error delivering updates")
	require.Equal(bfsm.T(),
		testUpdates[len(testUpdates)-1].Index,
		endingIndex[len(endingIndex)-1].Index,
		fmt.Sprintf("the ending index must be %d", testUpdates[len(testUpdates)-1].Index))
}

func (bfsm *TestBBoltFsm) TestPrepareSnapshot() {
	testOpts := &bbolt.Options{
		Timeout:         0,
		NoGrowSync:      false,
		NoFreelistSync:  false,
		FreelistType:    bbolt.FreelistMapType,
		ReadOnly:        false,
		InitialMmapSize: 0,
		PageSize:        0,
		NoSync:          false,
		OpenFile:        nil,
		Mlock:           false,
	}

	fsm := NewBBoltStateMachine(1, 1, bfsm.env.BaseDir, testOpts)
	empty, err := fsm.PrepareSnapshot()
	require.Empty(bfsm.T(), empty, "this must be a noop")
	require.NoError(bfsm.T(), err, "there must be no error")
}

func (bfsm *TestBBoltFsm) TestSnapshotLifecycle() {
	testOpts := &bbolt.Options{
		Timeout:         0,
		NoGrowSync:      false,
		NoFreelistSync:  false,
		FreelistType:    bbolt.FreelistMapType,
		ReadOnly:        false,
		InitialMmapSize: 0,
		PageSize:        0,
		NoSync:          false,
		OpenFile:        nil,
		Mlock:           false,
	}

	fsm := NewBBoltStateMachine(1, 1, bfsm.env.BaseDir, testOpts)
	index, err := fsm.Open(make(<-chan struct{}))
	require.NoError(bfsm.T(), err, "there must not be an error when opening the database")
	require.Equal(bfsm.T(), uint64(0), index, "the index must equal as there are no records")

	rootPrn := &PleiadesResourceName{
		Partition:    GlobalPartition,
		Service:      Pleiades,
		Region:       GlobalRegion,
		AccountId:    testAccountKey,
		ResourceType: Bucket,
		ResourceId:   "test-bucket",
	}

	testKvps := []mvccpb.KeyValue{
		{
			Key:            []byte(rootPrn.ToFsmRootPath("test-bucket") + "/" + "test-key-0"),
			CreateRevision: 0,
			ModRevision:    0,
			Version:        1,
			Value:          []byte("test-value-0"),
			Lease:          0,
		},
		{
			Key:            []byte(rootPrn.ToFsmRootPath("test-bucket") + "/" + "test-key-1"),
			CreateRevision: 0,
			ModRevision:    0,
			Version:        1,
			Value:          []byte("test-value-1"),
			Lease:          0,
		},
		{
			Key:            []byte(rootPrn.ToFsmRootPath("test-bucket") + "/" + "test-key-2"),
			CreateRevision: 0,
			ModRevision:    0,
			Version:        1,
			Value:          []byte("test-value-2"),
			Lease:          0,
		},
	}

	testUpdates := make([]statemachine.Entry, 0)
	for idx := range testKvps {
		val, err := testKvps[idx].Marshal()
		if err != nil {
			bfsm.T().Error(err)
		}
		testUpdates = append(testUpdates, statemachine.Entry{
			Index:  uint64(idx),
			Cmd:    val,
			Result: statemachine.Result{},
		})
	}

	var endingIndex []statemachine.Entry
	require.NotPanics(bfsm.T(), func() {
		endingIndex, err = fsm.Update(testUpdates)
	})
	require.NoError(bfsm.T(), err, "there must not be an error delivering updates")
	require.Equal(bfsm.T(),
		testUpdates[len(testUpdates)-1].Index,
		endingIndex[len(endingIndex)-1].Index,
		fmt.Sprintf("the ending index must be %d", testUpdates[len(testUpdates)-1].Index))

	var buf bytes.Buffer
	require.NotPanics(bfsm.T(), func() {
		err = fsm.SaveSnapshot(context.Background(), &buf, make(<-chan struct{}))
	}, "saving the snapshot to the buffer must not panic")
	require.NoError(bfsm.T(), err, "there must not be an error when saving the snapshot")

	require.NotPanics(bfsm.T(), func() {
		err = fsm.RecoverFromSnapshot(&buf, make(<-chan struct{}))
	}, "recovering the snapshot must not panic")
	require.NoError(bfsm.T(), err, "there must not be an error when recovering from a snapshot")

	dbPath := fsm.dbPath(true)
	db, err := bbolt.Open(dbPath, os.FileMode(dbFileModeVal), nil)
	require.NoError(bfsm.T(), err)

	// todo (sienna): fix this to ensure the kvp is stored, not just the value
	var target []byte
	require.NoError(bfsm.T(), db.Update(func(tx *bbolt.Tx) error {
		// rootPrn.ToFsmRootPath("test-bucket") + "/" + "test-key-2"
		bucketHierarchy := strings.Split(rootPrn.ToFsmRootPath("test-bucket")+"/"+"test-key-2", "/")[1:]
		parentBucketName := bucketHierarchy[0]
		childBucketNames := bucketHierarchy[1 :]
		parentBucket := tx.Bucket([]byte(parentBucketName))
		resChan := make(chan []byte, 1)
		if err := keyOp(parentBucket, childBucketNames, make([]byte, 0), get, resChan); err != nil {
			return err
		}
		target = <-resChan
		return nil
	}))

	finalKvp := mvccpb.KeyValue{}
	if err := finalKvp.Unmarshal(target); err != nil {
		bfsm.T().Error(err)
	}

	require.Equal(bfsm.T(), testKvps[len(testKvps)-1], finalKvp, "the serialized result must match the initial value")
}
