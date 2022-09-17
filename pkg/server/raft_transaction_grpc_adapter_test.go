/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package server

import (
	"context"
	"net"
	"testing"
	"time"

	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3"
	dclient "github.com/lni/dragonboat/v3/client"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestRaftTransactionGrpcAdapter(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping")
	}
	suite.Run(t, new(RaftTransactionGrpcAdapterTestSuite))
}

type RaftTransactionGrpcAdapterTestSuite struct {
	suite.Suite
	logger         zerolog.Logger
	conn           *grpc.ClientConn
	srv            *grpc.Server
	nh             *dragonboat.NodeHost
	rtm            *raftTransactionManager
	rtmAdapter     *raftTransactionGrpcAdapter
	defaultTimeout time.Duration
}

// SetupTest represents a remote Pleiades host
func (r *RaftTransactionGrpcAdapterTestSuite) SetupTest() {
	r.logger = utils.NewTestLogger(r.T())
	r.defaultTimeout = 300 * time.Millisecond

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	ctx := context.Background()
	r.srv = grpc.NewServer()

	r.nh = buildTestNodeHost(r.T())
	r.rtm = &raftTransactionManager{
		logger:       r.logger,
		nh:           r.nh,
		sessionCache: make(map[uint64]*dclient.Session),
	}
	r.rtmAdapter = &raftTransactionGrpcAdapter{
		logger:             r.logger,
		transactionManager: r.rtm,
	}

	kvstorev1.RegisterTransactionsServiceServer(r.srv, r.rtmAdapter)

	go func() {
		if err := r.srv.Serve(listener); err != nil {
			panic(err)
		}
	}()

	r.conn, _ = grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
}

func (r *RaftTransactionGrpcAdapterTestSuite) TearDownTest() {
	// safely close things.
	r.conn.Close()
	r.srv.Stop()

	// clear out the values
	r.srv = nil
	r.conn = nil
	r.rtm = nil
	r.nh = nil
}

func (r *RaftTransactionGrpcAdapterTestSuite) TestTransactionLifecycle() {

	shardConfig := buildTestShardConfig(r.T())
	shardConfig.SnapshotEntries = 5
	shardConfig.ClusterID = 1000
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = r.nh.RaftAddress()

	err := r.nh.StartCluster(members, false, newTestStateMachine, shardConfig)
	r.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(3000 * time.Millisecond)

	// create a transaction

	ctx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	client := kvstorev1.NewTransactionsServiceClient(r.conn)
	resp, err := client.NewTransaction(ctx, &kvstorev1.NewTransactionRequest{
		ShardId: shardConfig.ClusterID,
	})
	r.Require().NoError(err, "there must not be an error when requesting a new transaction")
	r.Require().NotNil(resp, "the reponse must not be nil")
	r.Require().NotNil(resp.GetTransaction(), "the transaction must not be nil")

	// verify it's in the cache and configured properly

	initialTransaction := resp.GetTransaction()
	cs, ok := r.rtm.sessionCache[initialTransaction.GetClientId()]
	r.Require().True(ok, "the session must exist in the session cache")
	r.Require().Equal(cs.ClusterID, initialTransaction.GetShardId(), "the shard ids must")
	r.Require().Equal(cs.ClientID, initialTransaction.GetClientId(), "the client ids must")
	r.Require().Equal(cs.SeriesID, initialTransaction.GetTransactionId(), "the transaction ids must")
	r.Require().Equal(cs.RespondedTo, initialTransaction.GetRespondedTo(), "the responded values must")

	// add something to the shard out-of-band

	proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	_, err = r.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	r.Require().NoError(err, "there must not be an error when proposing a new message")

	// commit the change

	postCommitResp, err := client.Commit(context.Background(), &kvstorev1.CommitRequest{Transaction: initialTransaction})
	r.Require().NoError(err, "there must not be an error when committing a transaction")
	r.Require().NotNil(postCommitResp, "the response but not be nil")
	r.Require().NotNil(postCommitResp.GetTransaction(), "the transaction must not be nil")

	postCommitTransaction := postCommitResp.GetTransaction()
	r.Require().NotNil(postCommitTransaction, "the post commit transaction must not be nil")

	cs, ok = r.rtm.sessionCache[initialTransaction.GetClientId()]
	r.Require().True(ok, "the session must exist in the session cache")
	r.Require().Equal(cs.ClusterID, postCommitTransaction.GetShardId(), "the shard ids must match")
	r.Require().Equal(cs.ClientID, postCommitTransaction.GetClientId(), "the client ids must match")
	r.Require().Equal(cs.SeriesID, postCommitTransaction.GetTransactionId(), "the transaction ids must match")
	r.Require().Equal(cs.RespondedTo, postCommitTransaction.GetRespondedTo(), "the responded values must match")

	// continue with the same transaction chain out of band

	proposeContext, _ = context.WithTimeout(context.Background(), 3000*time.Millisecond)
	_, err = r.nh.SyncPropose(proposeContext, cs, []byte("test-message-2"))
	r.Require().NoError(err, "there must not be an error when proposing a new message")

	// verify a second commit on the same transaction chain is possible

	postSecondCommitResponse, err := client.Commit(context.Background(), &kvstorev1.CommitRequest{Transaction: postCommitTransaction})
	r.Require().NoError(err, "there must not be an error when committing a transaction")
	r.Require().NotNil(postCommitResp, "the response but not be nil")
	r.Require().NotNil(postCommitResp.GetTransaction(), "the transaction must not be nil")

	postSecondCommitTransaction := postSecondCommitResponse.GetTransaction()
	r.Require().NotNil(postSecondCommitTransaction, "the post commit transaction must not be nil")

	cs, ok = r.rtm.sessionCache[initialTransaction.GetClientId()]
	r.Require().True(ok, "the session must exist in the session cache")
	r.Require().Equal(cs.ClusterID, postSecondCommitTransaction.GetShardId(), "the shard ids must match")
	r.Require().Equal(cs.ClientID, postSecondCommitTransaction.GetClientId(), "the client ids must match")
	r.Require().Equal(cs.SeriesID, postSecondCommitTransaction.GetTransactionId(), "the transaction ids must match")
	r.Require().Equal(cs.RespondedTo, postSecondCommitTransaction.GetRespondedTo(), "the responded values must match")

	// close the session

	closeCtx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	_, err = client.CloseTransaction(closeCtx, &kvstorev1.CloseTransactionRequest{
		Transaction: postSecondCommitTransaction,
		Timeout:     int64(3000 * time.Millisecond),
	})
	r.Require().NoError(err, "there must not be an error when closing the transaction")

	cs, ok = r.rtm.sessionCache[initialTransaction.GetClientId()]
	r.Require().False(ok, "the session must not exist in the session cache")
	r.Require().Nil(cs, "the value must be nil")
}
