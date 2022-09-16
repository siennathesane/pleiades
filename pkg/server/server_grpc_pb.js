// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var api_v1_database_kv_pb = require('../../api/v1/database/kv_pb.js');
var api_v1_database_transactions_pb = require('../../api/v1/database/transactions_pb.js');
var api_v1_raft_raft_shard_pb = require('../../api/v1/raft/raft_shard_pb.js');
var api_v1_raft_raft_host_pb = require('../../api/v1/raft/raft_host_pb.js');

function serialize_database_CloseTransactionReply(arg) {
  if (!(arg instanceof api_v1_database_transactions_pb.CloseTransactionReply)) {
    throw new Error('Expected argument of type database.CloseTransactionReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_CloseTransactionReply(buffer_arg) {
  return api_v1_database_transactions_pb.CloseTransactionReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_CloseTransactionRequest(arg) {
  if (!(arg instanceof api_v1_database_transactions_pb.CloseTransactionRequest)) {
    throw new Error('Expected argument of type database.CloseTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_CloseTransactionRequest(buffer_arg) {
  return api_v1_database_transactions_pb.CloseTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_CommitReply(arg) {
  if (!(arg instanceof api_v1_database_transactions_pb.CommitReply)) {
    throw new Error('Expected argument of type database.CommitReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_CommitReply(buffer_arg) {
  return api_v1_database_transactions_pb.CommitReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_CommitRequest(arg) {
  if (!(arg instanceof api_v1_database_transactions_pb.CommitRequest)) {
    throw new Error('Expected argument of type database.CommitRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_CommitRequest(buffer_arg) {
  return api_v1_database_transactions_pb.CommitRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_CreateAccountReply(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.CreateAccountReply)) {
    throw new Error('Expected argument of type database.CreateAccountReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_CreateAccountReply(buffer_arg) {
  return api_v1_database_kv_pb.CreateAccountReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_CreateAccountRequest(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.CreateAccountRequest)) {
    throw new Error('Expected argument of type database.CreateAccountRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_CreateAccountRequest(buffer_arg) {
  return api_v1_database_kv_pb.CreateAccountRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_CreateBucketReply(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.CreateBucketReply)) {
    throw new Error('Expected argument of type database.CreateBucketReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_CreateBucketReply(buffer_arg) {
  return api_v1_database_kv_pb.CreateBucketReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_CreateBucketRequest(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.CreateBucketRequest)) {
    throw new Error('Expected argument of type database.CreateBucketRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_CreateBucketRequest(buffer_arg) {
  return api_v1_database_kv_pb.CreateBucketRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_DeleteAccountReply(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.DeleteAccountReply)) {
    throw new Error('Expected argument of type database.DeleteAccountReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_DeleteAccountReply(buffer_arg) {
  return api_v1_database_kv_pb.DeleteAccountReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_DeleteAccountRequest(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.DeleteAccountRequest)) {
    throw new Error('Expected argument of type database.DeleteAccountRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_DeleteAccountRequest(buffer_arg) {
  return api_v1_database_kv_pb.DeleteAccountRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_DeleteBucketReply(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.DeleteBucketReply)) {
    throw new Error('Expected argument of type database.DeleteBucketReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_DeleteBucketReply(buffer_arg) {
  return api_v1_database_kv_pb.DeleteBucketReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_DeleteBucketRequest(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.DeleteBucketRequest)) {
    throw new Error('Expected argument of type database.DeleteBucketRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_DeleteBucketRequest(buffer_arg) {
  return api_v1_database_kv_pb.DeleteBucketRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_DeleteKeyReply(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.DeleteKeyReply)) {
    throw new Error('Expected argument of type database.DeleteKeyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_DeleteKeyReply(buffer_arg) {
  return api_v1_database_kv_pb.DeleteKeyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_DeleteKeyRequest(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.DeleteKeyRequest)) {
    throw new Error('Expected argument of type database.DeleteKeyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_DeleteKeyRequest(buffer_arg) {
  return api_v1_database_kv_pb.DeleteKeyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_GetKeyReply(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.GetKeyReply)) {
    throw new Error('Expected argument of type database.GetKeyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_GetKeyReply(buffer_arg) {
  return api_v1_database_kv_pb.GetKeyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_GetKeyRequest(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.GetKeyRequest)) {
    throw new Error('Expected argument of type database.GetKeyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_GetKeyRequest(buffer_arg) {
  return api_v1_database_kv_pb.GetKeyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_NewTransactionReply(arg) {
  if (!(arg instanceof api_v1_database_transactions_pb.NewTransactionReply)) {
    throw new Error('Expected argument of type database.NewTransactionReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_NewTransactionReply(buffer_arg) {
  return api_v1_database_transactions_pb.NewTransactionReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_NewTransactionRequest(arg) {
  if (!(arg instanceof api_v1_database_transactions_pb.NewTransactionRequest)) {
    throw new Error('Expected argument of type database.NewTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_NewTransactionRequest(buffer_arg) {
  return api_v1_database_transactions_pb.NewTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_PutKeyReply(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.PutKeyReply)) {
    throw new Error('Expected argument of type database.PutKeyReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_PutKeyReply(buffer_arg) {
  return api_v1_database_kv_pb.PutKeyReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_database_PutKeyRequest(arg) {
  if (!(arg instanceof api_v1_database_kv_pb.PutKeyRequest)) {
    throw new Error('Expected argument of type database.PutKeyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_database_PutKeyRequest(buffer_arg) {
  return api_v1_database_kv_pb.PutKeyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_AddReplicaObserverReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.AddReplicaObserverReply)) {
    throw new Error('Expected argument of type raft.AddReplicaObserverReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_AddReplicaObserverReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.AddReplicaObserverReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_AddReplicaObserverRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.AddReplicaObserverRequest)) {
    throw new Error('Expected argument of type raft.AddReplicaObserverRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_AddReplicaObserverRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.AddReplicaObserverRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_AddReplicaReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.AddReplicaReply)) {
    throw new Error('Expected argument of type raft.AddReplicaReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_AddReplicaReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.AddReplicaReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_AddReplicaRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.AddReplicaRequest)) {
    throw new Error('Expected argument of type raft.AddReplicaRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_AddReplicaRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.AddReplicaRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_AddReplicaWitnessReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.AddReplicaWitnessReply)) {
    throw new Error('Expected argument of type raft.AddReplicaWitnessReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_AddReplicaWitnessReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.AddReplicaWitnessReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_AddReplicaWitnessRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest)) {
    throw new Error('Expected argument of type raft.AddReplicaWitnessRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_AddReplicaWitnessRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_CompactReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_host_pb.CompactReply)) {
    throw new Error('Expected argument of type raft.CompactReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_CompactReply(buffer_arg) {
  return api_v1_raft_raft_host_pb.CompactReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_CompactRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_host_pb.CompactRequest)) {
    throw new Error('Expected argument of type raft.CompactRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_CompactRequest(buffer_arg) {
  return api_v1_raft_raft_host_pb.CompactRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_DeleteReplicaReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.DeleteReplicaReply)) {
    throw new Error('Expected argument of type raft.DeleteReplicaReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_DeleteReplicaReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.DeleteReplicaReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_DeleteReplicaRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.DeleteReplicaRequest)) {
    throw new Error('Expected argument of type raft.DeleteReplicaRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_DeleteReplicaRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.DeleteReplicaRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_GetHostConfigReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_host_pb.GetHostConfigReply)) {
    throw new Error('Expected argument of type raft.GetHostConfigReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_GetHostConfigReply(buffer_arg) {
  return api_v1_raft_raft_host_pb.GetHostConfigReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_GetHostConfigRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_host_pb.GetHostConfigRequest)) {
    throw new Error('Expected argument of type raft.GetHostConfigRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_GetHostConfigRequest(buffer_arg) {
  return api_v1_raft_raft_host_pb.GetHostConfigRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_GetLeaderIdReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.GetLeaderIdReply)) {
    throw new Error('Expected argument of type raft.GetLeaderIdReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_GetLeaderIdReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.GetLeaderIdReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_GetLeaderIdRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.GetLeaderIdRequest)) {
    throw new Error('Expected argument of type raft.GetLeaderIdRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_GetLeaderIdRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.GetLeaderIdRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_GetShardMembersReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.GetShardMembersReply)) {
    throw new Error('Expected argument of type raft.GetShardMembersReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_GetShardMembersReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.GetShardMembersReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_GetShardMembersRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.GetShardMembersRequest)) {
    throw new Error('Expected argument of type raft.GetShardMembersRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_GetShardMembersRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.GetShardMembersRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_NewShardReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.NewShardReply)) {
    throw new Error('Expected argument of type raft.NewShardReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_NewShardReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.NewShardReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_NewShardRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.NewShardRequest)) {
    throw new Error('Expected argument of type raft.NewShardRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_NewShardRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.NewShardRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_RemoveDataReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.RemoveDataReply)) {
    throw new Error('Expected argument of type raft.RemoveDataReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_RemoveDataReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.RemoveDataReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_RemoveDataRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.RemoveDataRequest)) {
    throw new Error('Expected argument of type raft.RemoveDataRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_RemoveDataRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.RemoveDataRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_SnapshotReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_host_pb.SnapshotReply)) {
    throw new Error('Expected argument of type raft.SnapshotReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_SnapshotReply(buffer_arg) {
  return api_v1_raft_raft_host_pb.SnapshotReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_SnapshotRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_host_pb.SnapshotRequest)) {
    throw new Error('Expected argument of type raft.SnapshotRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_SnapshotRequest(buffer_arg) {
  return api_v1_raft_raft_host_pb.SnapshotRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_StartReplicaReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.StartReplicaReply)) {
    throw new Error('Expected argument of type raft.StartReplicaReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_StartReplicaReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.StartReplicaReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_StartReplicaRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.StartReplicaRequest)) {
    throw new Error('Expected argument of type raft.StartReplicaRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_StartReplicaRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.StartReplicaRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_StopReplicaReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.StopReplicaReply)) {
    throw new Error('Expected argument of type raft.StopReplicaReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_StopReplicaReply(buffer_arg) {
  return api_v1_raft_raft_shard_pb.StopReplicaReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_StopReplicaRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_shard_pb.StopReplicaRequest)) {
    throw new Error('Expected argument of type raft.StopReplicaRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_StopReplicaRequest(buffer_arg) {
  return api_v1_raft_raft_shard_pb.StopReplicaRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_StopReply(arg) {
  if (!(arg instanceof api_v1_raft_raft_host_pb.StopReply)) {
    throw new Error('Expected argument of type raft.StopReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_StopReply(buffer_arg) {
  return api_v1_raft_raft_host_pb.StopReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_raft_StopRequest(arg) {
  if (!(arg instanceof api_v1_raft_raft_host_pb.StopRequest)) {
    throw new Error('Expected argument of type raft.StopRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_raft_StopRequest(buffer_arg) {
  return api_v1_raft_raft_host_pb.StopRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


var ShardManagerService = exports.ShardManagerService = {
  addReplica: {
    path: '/server.ShardManager/AddReplica',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.AddReplicaRequest,
    responseType: api_v1_raft_raft_shard_pb.AddReplicaReply,
    requestSerialize: serialize_raft_AddReplicaRequest,
    requestDeserialize: deserialize_raft_AddReplicaRequest,
    responseSerialize: serialize_raft_AddReplicaReply,
    responseDeserialize: deserialize_raft_AddReplicaReply,
  },
  addReplicaObserver: {
    path: '/server.ShardManager/AddReplicaObserver',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.AddReplicaObserverRequest,
    responseType: api_v1_raft_raft_shard_pb.AddReplicaObserverReply,
    requestSerialize: serialize_raft_AddReplicaObserverRequest,
    requestDeserialize: deserialize_raft_AddReplicaObserverRequest,
    responseSerialize: serialize_raft_AddReplicaObserverReply,
    responseDeserialize: deserialize_raft_AddReplicaObserverReply,
  },
  addReplicaWitness: {
    path: '/server.ShardManager/AddReplicaWitness',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest,
    responseType: api_v1_raft_raft_shard_pb.AddReplicaWitnessReply,
    requestSerialize: serialize_raft_AddReplicaWitnessRequest,
    requestDeserialize: deserialize_raft_AddReplicaWitnessRequest,
    responseSerialize: serialize_raft_AddReplicaWitnessReply,
    responseDeserialize: deserialize_raft_AddReplicaWitnessReply,
  },
  getLeaderId: {
    path: '/server.ShardManager/GetLeaderId',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.GetLeaderIdRequest,
    responseType: api_v1_raft_raft_shard_pb.GetLeaderIdReply,
    requestSerialize: serialize_raft_GetLeaderIdRequest,
    requestDeserialize: deserialize_raft_GetLeaderIdRequest,
    responseSerialize: serialize_raft_GetLeaderIdReply,
    responseDeserialize: deserialize_raft_GetLeaderIdReply,
  },
  getShardMembers: {
    path: '/server.ShardManager/GetShardMembers',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.GetShardMembersRequest,
    responseType: api_v1_raft_raft_shard_pb.GetShardMembersReply,
    requestSerialize: serialize_raft_GetShardMembersRequest,
    requestDeserialize: deserialize_raft_GetShardMembersRequest,
    responseSerialize: serialize_raft_GetShardMembersReply,
    responseDeserialize: deserialize_raft_GetShardMembersReply,
  },
  newShard: {
    path: '/server.ShardManager/NewShard',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.NewShardRequest,
    responseType: api_v1_raft_raft_shard_pb.NewShardReply,
    requestSerialize: serialize_raft_NewShardRequest,
    requestDeserialize: deserialize_raft_NewShardRequest,
    responseSerialize: serialize_raft_NewShardReply,
    responseDeserialize: deserialize_raft_NewShardReply,
  },
  removeData: {
    path: '/server.ShardManager/RemoveData',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.RemoveDataRequest,
    responseType: api_v1_raft_raft_shard_pb.RemoveDataReply,
    requestSerialize: serialize_raft_RemoveDataRequest,
    requestDeserialize: deserialize_raft_RemoveDataRequest,
    responseSerialize: serialize_raft_RemoveDataReply,
    responseDeserialize: deserialize_raft_RemoveDataReply,
  },
  removeReplica: {
    path: '/server.ShardManager/RemoveReplica',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.DeleteReplicaRequest,
    responseType: api_v1_raft_raft_shard_pb.DeleteReplicaReply,
    requestSerialize: serialize_raft_DeleteReplicaRequest,
    requestDeserialize: deserialize_raft_DeleteReplicaRequest,
    responseSerialize: serialize_raft_DeleteReplicaReply,
    responseDeserialize: deserialize_raft_DeleteReplicaReply,
  },
  startReplica: {
    path: '/server.ShardManager/StartReplica',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.StartReplicaRequest,
    responseType: api_v1_raft_raft_shard_pb.StartReplicaReply,
    requestSerialize: serialize_raft_StartReplicaRequest,
    requestDeserialize: deserialize_raft_StartReplicaRequest,
    responseSerialize: serialize_raft_StartReplicaReply,
    responseDeserialize: deserialize_raft_StartReplicaReply,
  },
  startReplicaObserver: {
    path: '/server.ShardManager/StartReplicaObserver',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.StartReplicaRequest,
    responseType: api_v1_raft_raft_shard_pb.StartReplicaReply,
    requestSerialize: serialize_raft_StartReplicaRequest,
    requestDeserialize: deserialize_raft_StartReplicaRequest,
    responseSerialize: serialize_raft_StartReplicaReply,
    responseDeserialize: deserialize_raft_StartReplicaReply,
  },
  stopReplica: {
    path: '/server.ShardManager/StopReplica',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_shard_pb.StopReplicaRequest,
    responseType: api_v1_raft_raft_shard_pb.StopReplicaReply,
    requestSerialize: serialize_raft_StopReplicaRequest,
    requestDeserialize: deserialize_raft_StopReplicaRequest,
    responseSerialize: serialize_raft_StopReplicaReply,
    responseDeserialize: deserialize_raft_StopReplicaReply,
  },
};

exports.ShardManagerClient = grpc.makeGenericClientConstructor(ShardManagerService);
var RaftHostService = exports.RaftHostService = {
  compact: {
    path: '/server.RaftHost/Compact',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_host_pb.CompactRequest,
    responseType: api_v1_raft_raft_host_pb.CompactReply,
    requestSerialize: serialize_raft_CompactRequest,
    requestDeserialize: deserialize_raft_CompactRequest,
    responseSerialize: serialize_raft_CompactReply,
    responseDeserialize: deserialize_raft_CompactReply,
  },
  getHostConfig: {
    path: '/server.RaftHost/GetHostConfig',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_host_pb.GetHostConfigRequest,
    responseType: api_v1_raft_raft_host_pb.GetHostConfigReply,
    requestSerialize: serialize_raft_GetHostConfigRequest,
    requestDeserialize: deserialize_raft_GetHostConfigRequest,
    responseSerialize: serialize_raft_GetHostConfigReply,
    responseDeserialize: deserialize_raft_GetHostConfigReply,
  },
  //  rpc LeaderTransfer(raft.LeaderTransferRequest) returns (raft.LeaderTransferReply);
snapshot: {
    path: '/server.RaftHost/Snapshot',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_host_pb.SnapshotRequest,
    responseType: api_v1_raft_raft_host_pb.SnapshotReply,
    requestSerialize: serialize_raft_SnapshotRequest,
    requestDeserialize: deserialize_raft_SnapshotRequest,
    responseSerialize: serialize_raft_SnapshotReply,
    responseDeserialize: deserialize_raft_SnapshotReply,
  },
  stop: {
    path: '/server.RaftHost/Stop',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_raft_raft_host_pb.StopRequest,
    responseType: api_v1_raft_raft_host_pb.StopReply,
    requestSerialize: serialize_raft_StopRequest,
    requestDeserialize: deserialize_raft_StopRequest,
    responseSerialize: serialize_raft_StopReply,
    responseDeserialize: deserialize_raft_StopReply,
  },
};

exports.RaftHostClient = grpc.makeGenericClientConstructor(RaftHostService);
var TransactionsService = exports.TransactionsService = {
  newTransaction: {
    path: '/server.Transactions/NewTransaction',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_transactions_pb.NewTransactionRequest,
    responseType: api_v1_database_transactions_pb.NewTransactionReply,
    requestSerialize: serialize_database_NewTransactionRequest,
    requestDeserialize: deserialize_database_NewTransactionRequest,
    responseSerialize: serialize_database_NewTransactionReply,
    responseDeserialize: deserialize_database_NewTransactionReply,
  },
  closeTransaction: {
    path: '/server.Transactions/CloseTransaction',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_transactions_pb.CloseTransactionRequest,
    responseType: api_v1_database_transactions_pb.CloseTransactionReply,
    requestSerialize: serialize_database_CloseTransactionRequest,
    requestDeserialize: deserialize_database_CloseTransactionRequest,
    responseSerialize: serialize_database_CloseTransactionReply,
    responseDeserialize: deserialize_database_CloseTransactionReply,
  },
  commit: {
    path: '/server.Transactions/Commit',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_transactions_pb.CommitRequest,
    responseType: api_v1_database_transactions_pb.CommitReply,
    requestSerialize: serialize_database_CommitRequest,
    requestDeserialize: deserialize_database_CommitRequest,
    responseSerialize: serialize_database_CommitReply,
    responseDeserialize: deserialize_database_CommitReply,
  },
};

exports.TransactionsClient = grpc.makeGenericClientConstructor(TransactionsService);
var KVStoreServiceService = exports.KVStoreServiceService = {
  createAccount: {
    path: '/server.KVStoreService/CreateAccount',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_kv_pb.CreateAccountRequest,
    responseType: api_v1_database_kv_pb.CreateAccountReply,
    requestSerialize: serialize_database_CreateAccountRequest,
    requestDeserialize: deserialize_database_CreateAccountRequest,
    responseSerialize: serialize_database_CreateAccountReply,
    responseDeserialize: deserialize_database_CreateAccountReply,
  },
  deleteAccount: {
    path: '/server.KVStoreService/DeleteAccount',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_kv_pb.DeleteAccountRequest,
    responseType: api_v1_database_kv_pb.DeleteAccountReply,
    requestSerialize: serialize_database_DeleteAccountRequest,
    requestDeserialize: deserialize_database_DeleteAccountRequest,
    responseSerialize: serialize_database_DeleteAccountReply,
    responseDeserialize: deserialize_database_DeleteAccountReply,
  },
  createBucket: {
    path: '/server.KVStoreService/CreateBucket',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_kv_pb.CreateBucketRequest,
    responseType: api_v1_database_kv_pb.CreateBucketReply,
    requestSerialize: serialize_database_CreateBucketRequest,
    requestDeserialize: deserialize_database_CreateBucketRequest,
    responseSerialize: serialize_database_CreateBucketReply,
    responseDeserialize: deserialize_database_CreateBucketReply,
  },
  deleteBucket: {
    path: '/server.KVStoreService/DeleteBucket',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_kv_pb.DeleteBucketRequest,
    responseType: api_v1_database_kv_pb.DeleteBucketReply,
    requestSerialize: serialize_database_DeleteBucketRequest,
    requestDeserialize: deserialize_database_DeleteBucketRequest,
    responseSerialize: serialize_database_DeleteBucketReply,
    responseDeserialize: deserialize_database_DeleteBucketReply,
  },
  getKey: {
    path: '/server.KVStoreService/GetKey',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_kv_pb.GetKeyRequest,
    responseType: api_v1_database_kv_pb.GetKeyReply,
    requestSerialize: serialize_database_GetKeyRequest,
    requestDeserialize: deserialize_database_GetKeyRequest,
    responseSerialize: serialize_database_GetKeyReply,
    responseDeserialize: deserialize_database_GetKeyReply,
  },
  putKey: {
    path: '/server.KVStoreService/PutKey',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_kv_pb.PutKeyRequest,
    responseType: api_v1_database_kv_pb.PutKeyReply,
    requestSerialize: serialize_database_PutKeyRequest,
    requestDeserialize: deserialize_database_PutKeyRequest,
    responseSerialize: serialize_database_PutKeyReply,
    responseDeserialize: deserialize_database_PutKeyReply,
  },
  deleteKey: {
    path: '/server.KVStoreService/DeleteKey',
    requestStream: false,
    responseStream: false,
    requestType: api_v1_database_kv_pb.DeleteKeyRequest,
    responseType: api_v1_database_kv_pb.DeleteKeyReply,
    requestSerialize: serialize_database_DeleteKeyRequest,
    requestDeserialize: deserialize_database_DeleteKeyRequest,
    responseSerialize: serialize_database_DeleteKeyReply,
    responseDeserialize: deserialize_database_DeleteKeyReply,
  },
};

exports.KVStoreServiceClient = grpc.makeGenericClientConstructor(KVStoreServiceService);
