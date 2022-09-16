// package: server
// file: pkg/server/server.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as pkg_server_server_pb from "../../pkg/server/server_pb";
import * as api_v1_database_kv_pb from "../../api/v1/database/kv_pb";
import * as api_v1_database_transactions_pb from "../../api/v1/database/transactions_pb";
import * as api_v1_raft_raft_shard_pb from "../../api/v1/raft/raft_shard_pb";
import * as api_v1_raft_raft_host_pb from "../../api/v1/raft/raft_host_pb";

interface IShardManagerService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    addReplica: IShardManagerService_IAddReplica;
    addReplicaObserver: IShardManagerService_IAddReplicaObserver;
    addReplicaWitness: IShardManagerService_IAddReplicaWitness;
    getLeaderId: IShardManagerService_IGetLeaderId;
    getShardMembers: IShardManagerService_IGetShardMembers;
    newShard: IShardManagerService_INewShard;
    removeData: IShardManagerService_IRemoveData;
    removeReplica: IShardManagerService_IRemoveReplica;
    startReplica: IShardManagerService_IStartReplica;
    startReplicaObserver: IShardManagerService_IStartReplicaObserver;
    stopReplica: IShardManagerService_IStopReplica;
}

interface IShardManagerService_IAddReplica extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.AddReplicaRequest, api_v1_raft_raft_shard_pb.AddReplicaReply> {
    path: "/server.ShardManager/AddReplica";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.AddReplicaRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.AddReplicaRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.AddReplicaReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.AddReplicaReply>;
}
interface IShardManagerService_IAddReplicaObserver extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.AddReplicaObserverRequest, api_v1_raft_raft_shard_pb.AddReplicaObserverReply> {
    path: "/server.ShardManager/AddReplicaObserver";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.AddReplicaObserverRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.AddReplicaObserverRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.AddReplicaObserverReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.AddReplicaObserverReply>;
}
interface IShardManagerService_IAddReplicaWitness extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest, api_v1_raft_raft_shard_pb.AddReplicaWitnessReply> {
    path: "/server.ShardManager/AddReplicaWitness";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.AddReplicaWitnessReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.AddReplicaWitnessReply>;
}
interface IShardManagerService_IGetLeaderId extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.GetLeaderIdRequest, api_v1_raft_raft_shard_pb.GetLeaderIdReply> {
    path: "/server.ShardManager/GetLeaderId";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.GetLeaderIdRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.GetLeaderIdRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.GetLeaderIdReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.GetLeaderIdReply>;
}
interface IShardManagerService_IGetShardMembers extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.GetShardMembersRequest, api_v1_raft_raft_shard_pb.GetShardMembersReply> {
    path: "/server.ShardManager/GetShardMembers";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.GetShardMembersRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.GetShardMembersRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.GetShardMembersReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.GetShardMembersReply>;
}
interface IShardManagerService_INewShard extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.NewShardRequest, api_v1_raft_raft_shard_pb.NewShardReply> {
    path: "/server.ShardManager/NewShard";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.NewShardRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.NewShardRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.NewShardReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.NewShardReply>;
}
interface IShardManagerService_IRemoveData extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.RemoveDataRequest, api_v1_raft_raft_shard_pb.RemoveDataReply> {
    path: "/server.ShardManager/RemoveData";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.RemoveDataRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.RemoveDataRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.RemoveDataReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.RemoveDataReply>;
}
interface IShardManagerService_IRemoveReplica extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.DeleteReplicaRequest, api_v1_raft_raft_shard_pb.DeleteReplicaReply> {
    path: "/server.ShardManager/RemoveReplica";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.DeleteReplicaRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.DeleteReplicaRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.DeleteReplicaReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.DeleteReplicaReply>;
}
interface IShardManagerService_IStartReplica extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.StartReplicaRequest, api_v1_raft_raft_shard_pb.StartReplicaReply> {
    path: "/server.ShardManager/StartReplica";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.StartReplicaRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.StartReplicaRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.StartReplicaReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.StartReplicaReply>;
}
interface IShardManagerService_IStartReplicaObserver extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.StartReplicaRequest, api_v1_raft_raft_shard_pb.StartReplicaReply> {
    path: "/server.ShardManager/StartReplicaObserver";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.StartReplicaRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.StartReplicaRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.StartReplicaReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.StartReplicaReply>;
}
interface IShardManagerService_IStopReplica extends grpc.MethodDefinition<api_v1_raft_raft_shard_pb.StopReplicaRequest, api_v1_raft_raft_shard_pb.StopReplicaReply> {
    path: "/server.ShardManager/StopReplica";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.StopReplicaRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.StopReplicaRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_shard_pb.StopReplicaReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_shard_pb.StopReplicaReply>;
}

export const ShardManagerService: IShardManagerService;

export interface IShardManagerServer {
    addReplica: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.AddReplicaRequest, api_v1_raft_raft_shard_pb.AddReplicaReply>;
    addReplicaObserver: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.AddReplicaObserverRequest, api_v1_raft_raft_shard_pb.AddReplicaObserverReply>;
    addReplicaWitness: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest, api_v1_raft_raft_shard_pb.AddReplicaWitnessReply>;
    getLeaderId: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.GetLeaderIdRequest, api_v1_raft_raft_shard_pb.GetLeaderIdReply>;
    getShardMembers: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.GetShardMembersRequest, api_v1_raft_raft_shard_pb.GetShardMembersReply>;
    newShard: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.NewShardRequest, api_v1_raft_raft_shard_pb.NewShardReply>;
    removeData: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.RemoveDataRequest, api_v1_raft_raft_shard_pb.RemoveDataReply>;
    removeReplica: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.DeleteReplicaRequest, api_v1_raft_raft_shard_pb.DeleteReplicaReply>;
    startReplica: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.StartReplicaRequest, api_v1_raft_raft_shard_pb.StartReplicaReply>;
    startReplicaObserver: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.StartReplicaRequest, api_v1_raft_raft_shard_pb.StartReplicaReply>;
    stopReplica: grpc.handleUnaryCall<api_v1_raft_raft_shard_pb.StopReplicaRequest, api_v1_raft_raft_shard_pb.StopReplicaReply>;
}

export interface IShardManagerClient {
    addReplica(request: api_v1_raft_raft_shard_pb.AddReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaReply) => void): grpc.ClientUnaryCall;
    addReplica(request: api_v1_raft_raft_shard_pb.AddReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaReply) => void): grpc.ClientUnaryCall;
    addReplica(request: api_v1_raft_raft_shard_pb.AddReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaReply) => void): grpc.ClientUnaryCall;
    addReplicaObserver(request: api_v1_raft_raft_shard_pb.AddReplicaObserverRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaObserverReply) => void): grpc.ClientUnaryCall;
    addReplicaObserver(request: api_v1_raft_raft_shard_pb.AddReplicaObserverRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaObserverReply) => void): grpc.ClientUnaryCall;
    addReplicaObserver(request: api_v1_raft_raft_shard_pb.AddReplicaObserverRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaObserverReply) => void): grpc.ClientUnaryCall;
    addReplicaWitness(request: api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaWitnessReply) => void): grpc.ClientUnaryCall;
    addReplicaWitness(request: api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaWitnessReply) => void): grpc.ClientUnaryCall;
    addReplicaWitness(request: api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaWitnessReply) => void): grpc.ClientUnaryCall;
    getLeaderId(request: api_v1_raft_raft_shard_pb.GetLeaderIdRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetLeaderIdReply) => void): grpc.ClientUnaryCall;
    getLeaderId(request: api_v1_raft_raft_shard_pb.GetLeaderIdRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetLeaderIdReply) => void): grpc.ClientUnaryCall;
    getLeaderId(request: api_v1_raft_raft_shard_pb.GetLeaderIdRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetLeaderIdReply) => void): grpc.ClientUnaryCall;
    getShardMembers(request: api_v1_raft_raft_shard_pb.GetShardMembersRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetShardMembersReply) => void): grpc.ClientUnaryCall;
    getShardMembers(request: api_v1_raft_raft_shard_pb.GetShardMembersRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetShardMembersReply) => void): grpc.ClientUnaryCall;
    getShardMembers(request: api_v1_raft_raft_shard_pb.GetShardMembersRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetShardMembersReply) => void): grpc.ClientUnaryCall;
    newShard(request: api_v1_raft_raft_shard_pb.NewShardRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.NewShardReply) => void): grpc.ClientUnaryCall;
    newShard(request: api_v1_raft_raft_shard_pb.NewShardRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.NewShardReply) => void): grpc.ClientUnaryCall;
    newShard(request: api_v1_raft_raft_shard_pb.NewShardRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.NewShardReply) => void): grpc.ClientUnaryCall;
    removeData(request: api_v1_raft_raft_shard_pb.RemoveDataRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.RemoveDataReply) => void): grpc.ClientUnaryCall;
    removeData(request: api_v1_raft_raft_shard_pb.RemoveDataRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.RemoveDataReply) => void): grpc.ClientUnaryCall;
    removeData(request: api_v1_raft_raft_shard_pb.RemoveDataRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.RemoveDataReply) => void): grpc.ClientUnaryCall;
    removeReplica(request: api_v1_raft_raft_shard_pb.DeleteReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.DeleteReplicaReply) => void): grpc.ClientUnaryCall;
    removeReplica(request: api_v1_raft_raft_shard_pb.DeleteReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.DeleteReplicaReply) => void): grpc.ClientUnaryCall;
    removeReplica(request: api_v1_raft_raft_shard_pb.DeleteReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.DeleteReplicaReply) => void): grpc.ClientUnaryCall;
    startReplica(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    startReplica(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    startReplica(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    startReplicaObserver(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    startReplicaObserver(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    startReplicaObserver(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    stopReplica(request: api_v1_raft_raft_shard_pb.StopReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StopReplicaReply) => void): grpc.ClientUnaryCall;
    stopReplica(request: api_v1_raft_raft_shard_pb.StopReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StopReplicaReply) => void): grpc.ClientUnaryCall;
    stopReplica(request: api_v1_raft_raft_shard_pb.StopReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StopReplicaReply) => void): grpc.ClientUnaryCall;
}

export class ShardManagerClient extends grpc.Client implements IShardManagerClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public addReplica(request: api_v1_raft_raft_shard_pb.AddReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaReply) => void): grpc.ClientUnaryCall;
    public addReplica(request: api_v1_raft_raft_shard_pb.AddReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaReply) => void): grpc.ClientUnaryCall;
    public addReplica(request: api_v1_raft_raft_shard_pb.AddReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaReply) => void): grpc.ClientUnaryCall;
    public addReplicaObserver(request: api_v1_raft_raft_shard_pb.AddReplicaObserverRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaObserverReply) => void): grpc.ClientUnaryCall;
    public addReplicaObserver(request: api_v1_raft_raft_shard_pb.AddReplicaObserverRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaObserverReply) => void): grpc.ClientUnaryCall;
    public addReplicaObserver(request: api_v1_raft_raft_shard_pb.AddReplicaObserverRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaObserverReply) => void): grpc.ClientUnaryCall;
    public addReplicaWitness(request: api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaWitnessReply) => void): grpc.ClientUnaryCall;
    public addReplicaWitness(request: api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaWitnessReply) => void): grpc.ClientUnaryCall;
    public addReplicaWitness(request: api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.AddReplicaWitnessReply) => void): grpc.ClientUnaryCall;
    public getLeaderId(request: api_v1_raft_raft_shard_pb.GetLeaderIdRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetLeaderIdReply) => void): grpc.ClientUnaryCall;
    public getLeaderId(request: api_v1_raft_raft_shard_pb.GetLeaderIdRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetLeaderIdReply) => void): grpc.ClientUnaryCall;
    public getLeaderId(request: api_v1_raft_raft_shard_pb.GetLeaderIdRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetLeaderIdReply) => void): grpc.ClientUnaryCall;
    public getShardMembers(request: api_v1_raft_raft_shard_pb.GetShardMembersRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetShardMembersReply) => void): grpc.ClientUnaryCall;
    public getShardMembers(request: api_v1_raft_raft_shard_pb.GetShardMembersRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetShardMembersReply) => void): grpc.ClientUnaryCall;
    public getShardMembers(request: api_v1_raft_raft_shard_pb.GetShardMembersRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.GetShardMembersReply) => void): grpc.ClientUnaryCall;
    public newShard(request: api_v1_raft_raft_shard_pb.NewShardRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.NewShardReply) => void): grpc.ClientUnaryCall;
    public newShard(request: api_v1_raft_raft_shard_pb.NewShardRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.NewShardReply) => void): grpc.ClientUnaryCall;
    public newShard(request: api_v1_raft_raft_shard_pb.NewShardRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.NewShardReply) => void): grpc.ClientUnaryCall;
    public removeData(request: api_v1_raft_raft_shard_pb.RemoveDataRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.RemoveDataReply) => void): grpc.ClientUnaryCall;
    public removeData(request: api_v1_raft_raft_shard_pb.RemoveDataRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.RemoveDataReply) => void): grpc.ClientUnaryCall;
    public removeData(request: api_v1_raft_raft_shard_pb.RemoveDataRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.RemoveDataReply) => void): grpc.ClientUnaryCall;
    public removeReplica(request: api_v1_raft_raft_shard_pb.DeleteReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.DeleteReplicaReply) => void): grpc.ClientUnaryCall;
    public removeReplica(request: api_v1_raft_raft_shard_pb.DeleteReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.DeleteReplicaReply) => void): grpc.ClientUnaryCall;
    public removeReplica(request: api_v1_raft_raft_shard_pb.DeleteReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.DeleteReplicaReply) => void): grpc.ClientUnaryCall;
    public startReplica(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    public startReplica(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    public startReplica(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    public startReplicaObserver(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    public startReplicaObserver(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    public startReplicaObserver(request: api_v1_raft_raft_shard_pb.StartReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StartReplicaReply) => void): grpc.ClientUnaryCall;
    public stopReplica(request: api_v1_raft_raft_shard_pb.StopReplicaRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StopReplicaReply) => void): grpc.ClientUnaryCall;
    public stopReplica(request: api_v1_raft_raft_shard_pb.StopReplicaRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StopReplicaReply) => void): grpc.ClientUnaryCall;
    public stopReplica(request: api_v1_raft_raft_shard_pb.StopReplicaRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_shard_pb.StopReplicaReply) => void): grpc.ClientUnaryCall;
}

interface IRaftHostService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    compact: IRaftHostService_ICompact;
    getHostConfig: IRaftHostService_IGetHostConfig;
    snapshot: IRaftHostService_ISnapshot;
    stop: IRaftHostService_IStop;
}

interface IRaftHostService_ICompact extends grpc.MethodDefinition<api_v1_raft_raft_host_pb.CompactRequest, api_v1_raft_raft_host_pb.CompactReply> {
    path: "/server.RaftHost/Compact";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_host_pb.CompactRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_host_pb.CompactRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_host_pb.CompactReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_host_pb.CompactReply>;
}
interface IRaftHostService_IGetHostConfig extends grpc.MethodDefinition<api_v1_raft_raft_host_pb.GetHostConfigRequest, api_v1_raft_raft_host_pb.GetHostConfigReply> {
    path: "/server.RaftHost/GetHostConfig";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_host_pb.GetHostConfigRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_host_pb.GetHostConfigRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_host_pb.GetHostConfigReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_host_pb.GetHostConfigReply>;
}
interface IRaftHostService_ISnapshot extends grpc.MethodDefinition<api_v1_raft_raft_host_pb.SnapshotRequest, api_v1_raft_raft_host_pb.SnapshotReply> {
    path: "/server.RaftHost/Snapshot";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_host_pb.SnapshotRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_host_pb.SnapshotRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_host_pb.SnapshotReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_host_pb.SnapshotReply>;
}
interface IRaftHostService_IStop extends grpc.MethodDefinition<api_v1_raft_raft_host_pb.StopRequest, api_v1_raft_raft_host_pb.StopReply> {
    path: "/server.RaftHost/Stop";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_raft_raft_host_pb.StopRequest>;
    requestDeserialize: grpc.deserialize<api_v1_raft_raft_host_pb.StopRequest>;
    responseSerialize: grpc.serialize<api_v1_raft_raft_host_pb.StopReply>;
    responseDeserialize: grpc.deserialize<api_v1_raft_raft_host_pb.StopReply>;
}

export const RaftHostService: IRaftHostService;

export interface IRaftHostServer {
    compact: grpc.handleUnaryCall<api_v1_raft_raft_host_pb.CompactRequest, api_v1_raft_raft_host_pb.CompactReply>;
    getHostConfig: grpc.handleUnaryCall<api_v1_raft_raft_host_pb.GetHostConfigRequest, api_v1_raft_raft_host_pb.GetHostConfigReply>;
    snapshot: grpc.handleUnaryCall<api_v1_raft_raft_host_pb.SnapshotRequest, api_v1_raft_raft_host_pb.SnapshotReply>;
    stop: grpc.handleUnaryCall<api_v1_raft_raft_host_pb.StopRequest, api_v1_raft_raft_host_pb.StopReply>;
}

export interface IRaftHostClient {
    compact(request: api_v1_raft_raft_host_pb.CompactRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.CompactReply) => void): grpc.ClientUnaryCall;
    compact(request: api_v1_raft_raft_host_pb.CompactRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.CompactReply) => void): grpc.ClientUnaryCall;
    compact(request: api_v1_raft_raft_host_pb.CompactRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.CompactReply) => void): grpc.ClientUnaryCall;
    getHostConfig(request: api_v1_raft_raft_host_pb.GetHostConfigRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.GetHostConfigReply) => void): grpc.ClientUnaryCall;
    getHostConfig(request: api_v1_raft_raft_host_pb.GetHostConfigRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.GetHostConfigReply) => void): grpc.ClientUnaryCall;
    getHostConfig(request: api_v1_raft_raft_host_pb.GetHostConfigRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.GetHostConfigReply) => void): grpc.ClientUnaryCall;
    snapshot(request: api_v1_raft_raft_host_pb.SnapshotRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.SnapshotReply) => void): grpc.ClientUnaryCall;
    snapshot(request: api_v1_raft_raft_host_pb.SnapshotRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.SnapshotReply) => void): grpc.ClientUnaryCall;
    snapshot(request: api_v1_raft_raft_host_pb.SnapshotRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.SnapshotReply) => void): grpc.ClientUnaryCall;
    stop(request: api_v1_raft_raft_host_pb.StopRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.StopReply) => void): grpc.ClientUnaryCall;
    stop(request: api_v1_raft_raft_host_pb.StopRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.StopReply) => void): grpc.ClientUnaryCall;
    stop(request: api_v1_raft_raft_host_pb.StopRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.StopReply) => void): grpc.ClientUnaryCall;
}

export class RaftHostClient extends grpc.Client implements IRaftHostClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public compact(request: api_v1_raft_raft_host_pb.CompactRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.CompactReply) => void): grpc.ClientUnaryCall;
    public compact(request: api_v1_raft_raft_host_pb.CompactRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.CompactReply) => void): grpc.ClientUnaryCall;
    public compact(request: api_v1_raft_raft_host_pb.CompactRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.CompactReply) => void): grpc.ClientUnaryCall;
    public getHostConfig(request: api_v1_raft_raft_host_pb.GetHostConfigRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.GetHostConfigReply) => void): grpc.ClientUnaryCall;
    public getHostConfig(request: api_v1_raft_raft_host_pb.GetHostConfigRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.GetHostConfigReply) => void): grpc.ClientUnaryCall;
    public getHostConfig(request: api_v1_raft_raft_host_pb.GetHostConfigRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.GetHostConfigReply) => void): grpc.ClientUnaryCall;
    public snapshot(request: api_v1_raft_raft_host_pb.SnapshotRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.SnapshotReply) => void): grpc.ClientUnaryCall;
    public snapshot(request: api_v1_raft_raft_host_pb.SnapshotRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.SnapshotReply) => void): grpc.ClientUnaryCall;
    public snapshot(request: api_v1_raft_raft_host_pb.SnapshotRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.SnapshotReply) => void): grpc.ClientUnaryCall;
    public stop(request: api_v1_raft_raft_host_pb.StopRequest, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.StopReply) => void): grpc.ClientUnaryCall;
    public stop(request: api_v1_raft_raft_host_pb.StopRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.StopReply) => void): grpc.ClientUnaryCall;
    public stop(request: api_v1_raft_raft_host_pb.StopRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_raft_raft_host_pb.StopReply) => void): grpc.ClientUnaryCall;
}

interface ITransactionsService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    newTransaction: ITransactionsService_INewTransaction;
    closeTransaction: ITransactionsService_ICloseTransaction;
    commit: ITransactionsService_ICommit;
}

interface ITransactionsService_INewTransaction extends grpc.MethodDefinition<api_v1_database_transactions_pb.NewTransactionRequest, api_v1_database_transactions_pb.NewTransactionReply> {
    path: "/server.Transactions/NewTransaction";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_transactions_pb.NewTransactionRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_transactions_pb.NewTransactionRequest>;
    responseSerialize: grpc.serialize<api_v1_database_transactions_pb.NewTransactionReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_transactions_pb.NewTransactionReply>;
}
interface ITransactionsService_ICloseTransaction extends grpc.MethodDefinition<api_v1_database_transactions_pb.CloseTransactionRequest, api_v1_database_transactions_pb.CloseTransactionReply> {
    path: "/server.Transactions/CloseTransaction";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_transactions_pb.CloseTransactionRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_transactions_pb.CloseTransactionRequest>;
    responseSerialize: grpc.serialize<api_v1_database_transactions_pb.CloseTransactionReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_transactions_pb.CloseTransactionReply>;
}
interface ITransactionsService_ICommit extends grpc.MethodDefinition<api_v1_database_transactions_pb.CommitRequest, api_v1_database_transactions_pb.CommitReply> {
    path: "/server.Transactions/Commit";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_transactions_pb.CommitRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_transactions_pb.CommitRequest>;
    responseSerialize: grpc.serialize<api_v1_database_transactions_pb.CommitReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_transactions_pb.CommitReply>;
}

export const TransactionsService: ITransactionsService;

export interface ITransactionsServer {
    newTransaction: grpc.handleUnaryCall<api_v1_database_transactions_pb.NewTransactionRequest, api_v1_database_transactions_pb.NewTransactionReply>;
    closeTransaction: grpc.handleUnaryCall<api_v1_database_transactions_pb.CloseTransactionRequest, api_v1_database_transactions_pb.CloseTransactionReply>;
    commit: grpc.handleUnaryCall<api_v1_database_transactions_pb.CommitRequest, api_v1_database_transactions_pb.CommitReply>;
}

export interface ITransactionsClient {
    newTransaction(request: api_v1_database_transactions_pb.NewTransactionRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.NewTransactionReply) => void): grpc.ClientUnaryCall;
    newTransaction(request: api_v1_database_transactions_pb.NewTransactionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.NewTransactionReply) => void): grpc.ClientUnaryCall;
    newTransaction(request: api_v1_database_transactions_pb.NewTransactionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.NewTransactionReply) => void): grpc.ClientUnaryCall;
    closeTransaction(request: api_v1_database_transactions_pb.CloseTransactionRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CloseTransactionReply) => void): grpc.ClientUnaryCall;
    closeTransaction(request: api_v1_database_transactions_pb.CloseTransactionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CloseTransactionReply) => void): grpc.ClientUnaryCall;
    closeTransaction(request: api_v1_database_transactions_pb.CloseTransactionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CloseTransactionReply) => void): grpc.ClientUnaryCall;
    commit(request: api_v1_database_transactions_pb.CommitRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CommitReply) => void): grpc.ClientUnaryCall;
    commit(request: api_v1_database_transactions_pb.CommitRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CommitReply) => void): grpc.ClientUnaryCall;
    commit(request: api_v1_database_transactions_pb.CommitRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CommitReply) => void): grpc.ClientUnaryCall;
}

export class TransactionsClient extends grpc.Client implements ITransactionsClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public newTransaction(request: api_v1_database_transactions_pb.NewTransactionRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.NewTransactionReply) => void): grpc.ClientUnaryCall;
    public newTransaction(request: api_v1_database_transactions_pb.NewTransactionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.NewTransactionReply) => void): grpc.ClientUnaryCall;
    public newTransaction(request: api_v1_database_transactions_pb.NewTransactionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.NewTransactionReply) => void): grpc.ClientUnaryCall;
    public closeTransaction(request: api_v1_database_transactions_pb.CloseTransactionRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CloseTransactionReply) => void): grpc.ClientUnaryCall;
    public closeTransaction(request: api_v1_database_transactions_pb.CloseTransactionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CloseTransactionReply) => void): grpc.ClientUnaryCall;
    public closeTransaction(request: api_v1_database_transactions_pb.CloseTransactionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CloseTransactionReply) => void): grpc.ClientUnaryCall;
    public commit(request: api_v1_database_transactions_pb.CommitRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CommitReply) => void): grpc.ClientUnaryCall;
    public commit(request: api_v1_database_transactions_pb.CommitRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CommitReply) => void): grpc.ClientUnaryCall;
    public commit(request: api_v1_database_transactions_pb.CommitRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_transactions_pb.CommitReply) => void): grpc.ClientUnaryCall;
}

interface IKVStoreServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    createAccount: IKVStoreServiceService_ICreateAccount;
    deleteAccount: IKVStoreServiceService_IDeleteAccount;
    createBucket: IKVStoreServiceService_ICreateBucket;
    deleteBucket: IKVStoreServiceService_IDeleteBucket;
    getKey: IKVStoreServiceService_IGetKey;
    putKey: IKVStoreServiceService_IPutKey;
    deleteKey: IKVStoreServiceService_IDeleteKey;
}

interface IKVStoreServiceService_ICreateAccount extends grpc.MethodDefinition<api_v1_database_kv_pb.CreateAccountRequest, api_v1_database_kv_pb.CreateAccountReply> {
    path: "/server.KVStoreService/CreateAccount";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_kv_pb.CreateAccountRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_kv_pb.CreateAccountRequest>;
    responseSerialize: grpc.serialize<api_v1_database_kv_pb.CreateAccountReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_kv_pb.CreateAccountReply>;
}
interface IKVStoreServiceService_IDeleteAccount extends grpc.MethodDefinition<api_v1_database_kv_pb.DeleteAccountRequest, api_v1_database_kv_pb.DeleteAccountReply> {
    path: "/server.KVStoreService/DeleteAccount";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_kv_pb.DeleteAccountRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_kv_pb.DeleteAccountRequest>;
    responseSerialize: grpc.serialize<api_v1_database_kv_pb.DeleteAccountReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_kv_pb.DeleteAccountReply>;
}
interface IKVStoreServiceService_ICreateBucket extends grpc.MethodDefinition<api_v1_database_kv_pb.CreateBucketRequest, api_v1_database_kv_pb.CreateBucketReply> {
    path: "/server.KVStoreService/CreateBucket";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_kv_pb.CreateBucketRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_kv_pb.CreateBucketRequest>;
    responseSerialize: grpc.serialize<api_v1_database_kv_pb.CreateBucketReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_kv_pb.CreateBucketReply>;
}
interface IKVStoreServiceService_IDeleteBucket extends grpc.MethodDefinition<api_v1_database_kv_pb.DeleteBucketRequest, api_v1_database_kv_pb.DeleteBucketReply> {
    path: "/server.KVStoreService/DeleteBucket";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_kv_pb.DeleteBucketRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_kv_pb.DeleteBucketRequest>;
    responseSerialize: grpc.serialize<api_v1_database_kv_pb.DeleteBucketReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_kv_pb.DeleteBucketReply>;
}
interface IKVStoreServiceService_IGetKey extends grpc.MethodDefinition<api_v1_database_kv_pb.GetKeyRequest, api_v1_database_kv_pb.GetKeyReply> {
    path: "/server.KVStoreService/GetKey";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_kv_pb.GetKeyRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_kv_pb.GetKeyRequest>;
    responseSerialize: grpc.serialize<api_v1_database_kv_pb.GetKeyReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_kv_pb.GetKeyReply>;
}
interface IKVStoreServiceService_IPutKey extends grpc.MethodDefinition<api_v1_database_kv_pb.PutKeyRequest, api_v1_database_kv_pb.PutKeyReply> {
    path: "/server.KVStoreService/PutKey";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_kv_pb.PutKeyRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_kv_pb.PutKeyRequest>;
    responseSerialize: grpc.serialize<api_v1_database_kv_pb.PutKeyReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_kv_pb.PutKeyReply>;
}
interface IKVStoreServiceService_IDeleteKey extends grpc.MethodDefinition<api_v1_database_kv_pb.DeleteKeyRequest, api_v1_database_kv_pb.DeleteKeyReply> {
    path: "/server.KVStoreService/DeleteKey";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_v1_database_kv_pb.DeleteKeyRequest>;
    requestDeserialize: grpc.deserialize<api_v1_database_kv_pb.DeleteKeyRequest>;
    responseSerialize: grpc.serialize<api_v1_database_kv_pb.DeleteKeyReply>;
    responseDeserialize: grpc.deserialize<api_v1_database_kv_pb.DeleteKeyReply>;
}

export const KVStoreServiceService: IKVStoreServiceService;

export interface IKVStoreServiceServer {
    createAccount: grpc.handleUnaryCall<api_v1_database_kv_pb.CreateAccountRequest, api_v1_database_kv_pb.CreateAccountReply>;
    deleteAccount: grpc.handleUnaryCall<api_v1_database_kv_pb.DeleteAccountRequest, api_v1_database_kv_pb.DeleteAccountReply>;
    createBucket: grpc.handleUnaryCall<api_v1_database_kv_pb.CreateBucketRequest, api_v1_database_kv_pb.CreateBucketReply>;
    deleteBucket: grpc.handleUnaryCall<api_v1_database_kv_pb.DeleteBucketRequest, api_v1_database_kv_pb.DeleteBucketReply>;
    getKey: grpc.handleUnaryCall<api_v1_database_kv_pb.GetKeyRequest, api_v1_database_kv_pb.GetKeyReply>;
    putKey: grpc.handleUnaryCall<api_v1_database_kv_pb.PutKeyRequest, api_v1_database_kv_pb.PutKeyReply>;
    deleteKey: grpc.handleUnaryCall<api_v1_database_kv_pb.DeleteKeyRequest, api_v1_database_kv_pb.DeleteKeyReply>;
}

export interface IKVStoreServiceClient {
    createAccount(request: api_v1_database_kv_pb.CreateAccountRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateAccountReply) => void): grpc.ClientUnaryCall;
    createAccount(request: api_v1_database_kv_pb.CreateAccountRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateAccountReply) => void): grpc.ClientUnaryCall;
    createAccount(request: api_v1_database_kv_pb.CreateAccountRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateAccountReply) => void): grpc.ClientUnaryCall;
    deleteAccount(request: api_v1_database_kv_pb.DeleteAccountRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteAccountReply) => void): grpc.ClientUnaryCall;
    deleteAccount(request: api_v1_database_kv_pb.DeleteAccountRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteAccountReply) => void): grpc.ClientUnaryCall;
    deleteAccount(request: api_v1_database_kv_pb.DeleteAccountRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteAccountReply) => void): grpc.ClientUnaryCall;
    createBucket(request: api_v1_database_kv_pb.CreateBucketRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateBucketReply) => void): grpc.ClientUnaryCall;
    createBucket(request: api_v1_database_kv_pb.CreateBucketRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateBucketReply) => void): grpc.ClientUnaryCall;
    createBucket(request: api_v1_database_kv_pb.CreateBucketRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateBucketReply) => void): grpc.ClientUnaryCall;
    deleteBucket(request: api_v1_database_kv_pb.DeleteBucketRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteBucketReply) => void): grpc.ClientUnaryCall;
    deleteBucket(request: api_v1_database_kv_pb.DeleteBucketRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteBucketReply) => void): grpc.ClientUnaryCall;
    deleteBucket(request: api_v1_database_kv_pb.DeleteBucketRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteBucketReply) => void): grpc.ClientUnaryCall;
    getKey(request: api_v1_database_kv_pb.GetKeyRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.GetKeyReply) => void): grpc.ClientUnaryCall;
    getKey(request: api_v1_database_kv_pb.GetKeyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.GetKeyReply) => void): grpc.ClientUnaryCall;
    getKey(request: api_v1_database_kv_pb.GetKeyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.GetKeyReply) => void): grpc.ClientUnaryCall;
    putKey(request: api_v1_database_kv_pb.PutKeyRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.PutKeyReply) => void): grpc.ClientUnaryCall;
    putKey(request: api_v1_database_kv_pb.PutKeyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.PutKeyReply) => void): grpc.ClientUnaryCall;
    putKey(request: api_v1_database_kv_pb.PutKeyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.PutKeyReply) => void): grpc.ClientUnaryCall;
    deleteKey(request: api_v1_database_kv_pb.DeleteKeyRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteKeyReply) => void): grpc.ClientUnaryCall;
    deleteKey(request: api_v1_database_kv_pb.DeleteKeyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteKeyReply) => void): grpc.ClientUnaryCall;
    deleteKey(request: api_v1_database_kv_pb.DeleteKeyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteKeyReply) => void): grpc.ClientUnaryCall;
}

export class KVStoreServiceClient extends grpc.Client implements IKVStoreServiceClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public createAccount(request: api_v1_database_kv_pb.CreateAccountRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateAccountReply) => void): grpc.ClientUnaryCall;
    public createAccount(request: api_v1_database_kv_pb.CreateAccountRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateAccountReply) => void): grpc.ClientUnaryCall;
    public createAccount(request: api_v1_database_kv_pb.CreateAccountRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateAccountReply) => void): grpc.ClientUnaryCall;
    public deleteAccount(request: api_v1_database_kv_pb.DeleteAccountRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteAccountReply) => void): grpc.ClientUnaryCall;
    public deleteAccount(request: api_v1_database_kv_pb.DeleteAccountRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteAccountReply) => void): grpc.ClientUnaryCall;
    public deleteAccount(request: api_v1_database_kv_pb.DeleteAccountRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteAccountReply) => void): grpc.ClientUnaryCall;
    public createBucket(request: api_v1_database_kv_pb.CreateBucketRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateBucketReply) => void): grpc.ClientUnaryCall;
    public createBucket(request: api_v1_database_kv_pb.CreateBucketRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateBucketReply) => void): grpc.ClientUnaryCall;
    public createBucket(request: api_v1_database_kv_pb.CreateBucketRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.CreateBucketReply) => void): grpc.ClientUnaryCall;
    public deleteBucket(request: api_v1_database_kv_pb.DeleteBucketRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteBucketReply) => void): grpc.ClientUnaryCall;
    public deleteBucket(request: api_v1_database_kv_pb.DeleteBucketRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteBucketReply) => void): grpc.ClientUnaryCall;
    public deleteBucket(request: api_v1_database_kv_pb.DeleteBucketRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteBucketReply) => void): grpc.ClientUnaryCall;
    public getKey(request: api_v1_database_kv_pb.GetKeyRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.GetKeyReply) => void): grpc.ClientUnaryCall;
    public getKey(request: api_v1_database_kv_pb.GetKeyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.GetKeyReply) => void): grpc.ClientUnaryCall;
    public getKey(request: api_v1_database_kv_pb.GetKeyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.GetKeyReply) => void): grpc.ClientUnaryCall;
    public putKey(request: api_v1_database_kv_pb.PutKeyRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.PutKeyReply) => void): grpc.ClientUnaryCall;
    public putKey(request: api_v1_database_kv_pb.PutKeyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.PutKeyReply) => void): grpc.ClientUnaryCall;
    public putKey(request: api_v1_database_kv_pb.PutKeyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.PutKeyReply) => void): grpc.ClientUnaryCall;
    public deleteKey(request: api_v1_database_kv_pb.DeleteKeyRequest, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteKeyReply) => void): grpc.ClientUnaryCall;
    public deleteKey(request: api_v1_database_kv_pb.DeleteKeyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteKeyReply) => void): grpc.ClientUnaryCall;
    public deleteKey(request: api_v1_database_kv_pb.DeleteKeyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_v1_database_kv_pb.DeleteKeyReply) => void): grpc.ClientUnaryCall;
}
