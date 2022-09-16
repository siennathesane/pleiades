import * as grpcWeb from 'grpc-web';

import * as api_v1_database_transactions_pb from '../../api/v1/database/transactions_pb';
import * as api_v1_database_kv_pb from '../../api/v1/database/kv_pb';
import * as api_v1_raft_raft_shard_pb from '../../api/v1/raft/raft_shard_pb';
import * as api_v1_raft_raft_host_pb from '../../api/v1/raft/raft_host_pb';


export class ShardManagerClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  addReplica(
    request: api_v1_raft_raft_shard_pb.AddReplicaRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.AddReplicaReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.AddReplicaReply>;

  addReplicaObserver(
    request: api_v1_raft_raft_shard_pb.AddReplicaObserverRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.AddReplicaObserverReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.AddReplicaObserverReply>;

  addReplicaWitness(
    request: api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.AddReplicaWitnessReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.AddReplicaWitnessReply>;

  getLeaderId(
    request: api_v1_raft_raft_shard_pb.GetLeaderIdRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.GetLeaderIdReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.GetLeaderIdReply>;

  getShardMembers(
    request: api_v1_raft_raft_shard_pb.GetShardMembersRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.GetShardMembersReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.GetShardMembersReply>;

  newShard(
    request: api_v1_raft_raft_shard_pb.NewShardRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.NewShardReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.NewShardReply>;

  removeData(
    request: api_v1_raft_raft_shard_pb.RemoveDataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.RemoveDataReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.RemoveDataReply>;

  removeReplica(
    request: api_v1_raft_raft_shard_pb.DeleteReplicaRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.DeleteReplicaReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.DeleteReplicaReply>;

  startReplica(
    request: api_v1_raft_raft_shard_pb.StartReplicaRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.StartReplicaReply>;

  startReplicaObserver(
    request: api_v1_raft_raft_shard_pb.StartReplicaRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.StartReplicaReply>;

  stopReplica(
    request: api_v1_raft_raft_shard_pb.StopReplicaRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_shard_pb.StopReplicaReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_shard_pb.StopReplicaReply>;

}

export class RaftHostClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  compact(
    request: api_v1_raft_raft_host_pb.CompactRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_host_pb.CompactReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_host_pb.CompactReply>;

  getHostConfig(
    request: api_v1_raft_raft_host_pb.GetHostConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_host_pb.GetHostConfigReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_host_pb.GetHostConfigReply>;

  snapshot(
    request: api_v1_raft_raft_host_pb.SnapshotRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_host_pb.SnapshotReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_host_pb.SnapshotReply>;

  stop(
    request: api_v1_raft_raft_host_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_raft_raft_host_pb.StopReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_raft_raft_host_pb.StopReply>;

}

export class TransactionsClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  newTransaction(
    request: api_v1_database_transactions_pb.NewTransactionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_transactions_pb.NewTransactionReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_transactions_pb.NewTransactionReply>;

  closeTransaction(
    request: api_v1_database_transactions_pb.CloseTransactionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_transactions_pb.CloseTransactionReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_transactions_pb.CloseTransactionReply>;

  commit(
    request: api_v1_database_transactions_pb.CommitRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_transactions_pb.CommitReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_transactions_pb.CommitReply>;

}

export class KVStoreServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createAccount(
    request: api_v1_database_kv_pb.CreateAccountRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_kv_pb.CreateAccountReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_kv_pb.CreateAccountReply>;

  deleteAccount(
    request: api_v1_database_kv_pb.DeleteAccountRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_kv_pb.DeleteAccountReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_kv_pb.DeleteAccountReply>;

  createBucket(
    request: api_v1_database_kv_pb.CreateBucketRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_kv_pb.CreateBucketReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_kv_pb.CreateBucketReply>;

  deleteBucket(
    request: api_v1_database_kv_pb.DeleteBucketRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_kv_pb.DeleteBucketReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_kv_pb.DeleteBucketReply>;

  getKey(
    request: api_v1_database_kv_pb.GetKeyRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_kv_pb.GetKeyReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_kv_pb.GetKeyReply>;

  putKey(
    request: api_v1_database_kv_pb.PutKeyRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_kv_pb.PutKeyReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_kv_pb.PutKeyReply>;

  deleteKey(
    request: api_v1_database_kv_pb.DeleteKeyRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: api_v1_database_kv_pb.DeleteKeyReply) => void
  ): grpcWeb.ClientReadableStream<api_v1_database_kv_pb.DeleteKeyReply>;

}

export class ShardManagerPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  addReplica(
    request: api_v1_raft_raft_shard_pb.AddReplicaRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.AddReplicaReply>;

  addReplicaObserver(
    request: api_v1_raft_raft_shard_pb.AddReplicaObserverRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.AddReplicaObserverReply>;

  addReplicaWitness(
    request: api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.AddReplicaWitnessReply>;

  getLeaderId(
    request: api_v1_raft_raft_shard_pb.GetLeaderIdRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.GetLeaderIdReply>;

  getShardMembers(
    request: api_v1_raft_raft_shard_pb.GetShardMembersRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.GetShardMembersReply>;

  newShard(
    request: api_v1_raft_raft_shard_pb.NewShardRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.NewShardReply>;

  removeData(
    request: api_v1_raft_raft_shard_pb.RemoveDataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.RemoveDataReply>;

  removeReplica(
    request: api_v1_raft_raft_shard_pb.DeleteReplicaRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.DeleteReplicaReply>;

  startReplica(
    request: api_v1_raft_raft_shard_pb.StartReplicaRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.StartReplicaReply>;

  startReplicaObserver(
    request: api_v1_raft_raft_shard_pb.StartReplicaRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.StartReplicaReply>;

  stopReplica(
    request: api_v1_raft_raft_shard_pb.StopReplicaRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_shard_pb.StopReplicaReply>;

}

export class RaftHostPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  compact(
    request: api_v1_raft_raft_host_pb.CompactRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_host_pb.CompactReply>;

  getHostConfig(
    request: api_v1_raft_raft_host_pb.GetHostConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_host_pb.GetHostConfigReply>;

  snapshot(
    request: api_v1_raft_raft_host_pb.SnapshotRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_host_pb.SnapshotReply>;

  stop(
    request: api_v1_raft_raft_host_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_raft_raft_host_pb.StopReply>;

}

export class TransactionsPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  newTransaction(
    request: api_v1_database_transactions_pb.NewTransactionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_transactions_pb.NewTransactionReply>;

  closeTransaction(
    request: api_v1_database_transactions_pb.CloseTransactionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_transactions_pb.CloseTransactionReply>;

  commit(
    request: api_v1_database_transactions_pb.CommitRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_transactions_pb.CommitReply>;

}

export class KVStoreServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createAccount(
    request: api_v1_database_kv_pb.CreateAccountRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_kv_pb.CreateAccountReply>;

  deleteAccount(
    request: api_v1_database_kv_pb.DeleteAccountRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_kv_pb.DeleteAccountReply>;

  createBucket(
    request: api_v1_database_kv_pb.CreateBucketRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_kv_pb.CreateBucketReply>;

  deleteBucket(
    request: api_v1_database_kv_pb.DeleteBucketRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_kv_pb.DeleteBucketReply>;

  getKey(
    request: api_v1_database_kv_pb.GetKeyRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_kv_pb.GetKeyReply>;

  putKey(
    request: api_v1_database_kv_pb.PutKeyRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_kv_pb.PutKeyReply>;

  deleteKey(
    request: api_v1_database_kv_pb.DeleteKeyRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_v1_database_kv_pb.DeleteKeyReply>;

}

