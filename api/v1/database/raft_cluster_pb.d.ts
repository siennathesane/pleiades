// package: database
// file: api/v1/database/raft_cluster.proto

import * as jspb from "google-protobuf";

export class RaftClusterPayload extends jspb.Message {
  hasGetclustermembershiprequest(): boolean;
  clearGetclustermembershiprequest(): void;
  getGetclustermembershiprequest(): GetClusterMembershipRequest | undefined;
  setGetclustermembershiprequest(value?: GetClusterMembershipRequest): void;

  hasGetclustermembershipresponse(): boolean;
  clearGetclustermembershipresponse(): void;
  getGetclustermembershipresponse(): GetClusterMembershipResponse | undefined;
  setGetclustermembershipresponse(value?: GetClusterMembershipResponse): void;

  hasStopclusterrequest(): boolean;
  clearStopclusterrequest(): void;
  getStopclusterrequest(): StopClusterRequest | undefined;
  setStopclusterrequest(value?: StopClusterRequest): void;

  hasStopclusterresponse(): boolean;
  clearStopclusterresponse(): void;
  getStopclusterresponse(): StopClusterResponse | undefined;
  setStopclusterresponse(value?: StopClusterResponse): void;

  hasStartclusterrequest(): boolean;
  clearStartclusterrequest(): void;
  getStartclusterrequest(): StartClusterRequest | undefined;
  setStartclusterrequest(value?: StartClusterRequest): void;

  hasStartclusterresponse(): boolean;
  clearStartclusterresponse(): void;
  getStartclusterresponse(): StartClusterResponse | undefined;
  setStartclusterresponse(value?: StartClusterResponse): void;

  getMethod(): RaftClusterPayload.MethodNameMap[keyof RaftClusterPayload.MethodNameMap];
  setMethod(value: RaftClusterPayload.MethodNameMap[keyof RaftClusterPayload.MethodNameMap]): void;

  getTypeCase(): RaftClusterPayload.TypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftClusterPayload.AsObject;
  static toObject(includeInstance: boolean, msg: RaftClusterPayload): RaftClusterPayload.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RaftClusterPayload, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftClusterPayload;
  static deserializeBinaryFromReader(message: RaftClusterPayload, reader: jspb.BinaryReader): RaftClusterPayload;
}

export namespace RaftClusterPayload {
  export type AsObject = {
    getclustermembershiprequest?: GetClusterMembershipRequest.AsObject,
    getclustermembershipresponse?: GetClusterMembershipResponse.AsObject,
    stopclusterrequest?: StopClusterRequest.AsObject,
    stopclusterresponse?: StopClusterResponse.AsObject,
    startclusterrequest?: StartClusterRequest.AsObject,
    startclusterresponse?: StartClusterResponse.AsObject,
    method: RaftClusterPayload.MethodNameMap[keyof RaftClusterPayload.MethodNameMap],
  }

  export interface MethodNameMap {
    START_CLUSTER: 0;
    STOP_CLUSTER: 1;
    GET_CLUSTER_MEMBERSHIP: 2;
  }

  export const MethodName: MethodNameMap;

  export enum TypeCase {
    TYPE_NOT_SET = 0,
    GETCLUSTERMEMBERSHIPREQUEST = 1,
    GETCLUSTERMEMBERSHIPRESPONSE = 2,
    STOPCLUSTERREQUEST = 3,
    STOPCLUSTERRESPONSE = 4,
    STARTCLUSTERREQUEST = 5,
    STARTCLUSTERRESPONSE = 6,
  }
}

export class GetClusterMembershipRequest extends jspb.Message {
  getClusterid(): number;
  setClusterid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetClusterMembershipRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetClusterMembershipRequest): GetClusterMembershipRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetClusterMembershipRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetClusterMembershipRequest;
  static deserializeBinaryFromReader(message: GetClusterMembershipRequest, reader: jspb.BinaryReader): GetClusterMembershipRequest;
}

export namespace GetClusterMembershipRequest {
  export type AsObject = {
    clusterid: number,
  }
}

export class GetClusterMembershipResponse extends jspb.Message {
  hasMembership(): boolean;
  clearMembership(): void;
  getMembership(): RaftMembership | undefined;
  setMembership(value?: RaftMembership): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetClusterMembershipResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetClusterMembershipResponse): GetClusterMembershipResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetClusterMembershipResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetClusterMembershipResponse;
  static deserializeBinaryFromReader(message: GetClusterMembershipResponse, reader: jspb.BinaryReader): GetClusterMembershipResponse;
}

export namespace GetClusterMembershipResponse {
  export type AsObject = {
    membership?: RaftMembership.AsObject,
  }
}

export class RaftMembership extends jspb.Message {
  getConfigchangeid(): number;
  setConfigchangeid(value: number): void;

  getNodesMap(): jspb.Map<number, string>;
  clearNodesMap(): void;
  getObserversMap(): jspb.Map<number, string>;
  clearObserversMap(): void;
  getWitnessesMap(): jspb.Map<number, string>;
  clearWitnessesMap(): void;
  getRemovedMap(): jspb.Map<number, NilVal>;
  clearRemovedMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftMembership.AsObject;
  static toObject(includeInstance: boolean, msg: RaftMembership): RaftMembership.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RaftMembership, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftMembership;
  static deserializeBinaryFromReader(message: RaftMembership, reader: jspb.BinaryReader): RaftMembership;
}

export namespace RaftMembership {
  export type AsObject = {
    configchangeid: number,
    nodesMap: Array<[number, string]>,
    observersMap: Array<[number, string]>,
    witnessesMap: Array<[number, string]>,
    removedMap: Array<[number, NilVal.AsObject]>,
  }
}

export class NilVal extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NilVal.AsObject;
  static toObject(includeInstance: boolean, msg: NilVal): NilVal.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NilVal, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NilVal;
  static deserializeBinaryFromReader(message: NilVal, reader: jspb.BinaryReader): NilVal;
}

export namespace NilVal {
  export type AsObject = {
  }
}

export class StopClusterRequest extends jspb.Message {
  getClusterid(): number;
  setClusterid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopClusterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StopClusterRequest): StopClusterRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StopClusterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopClusterRequest;
  static deserializeBinaryFromReader(message: StopClusterRequest, reader: jspb.BinaryReader): StopClusterRequest;
}

export namespace StopClusterRequest {
  export type AsObject = {
    clusterid: number,
  }
}

export class StopClusterResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopClusterResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StopClusterResponse): StopClusterResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StopClusterResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopClusterResponse;
  static deserializeBinaryFromReader(message: StopClusterResponse, reader: jspb.BinaryReader): StopClusterResponse;
}

export namespace StopClusterResponse {
  export type AsObject = {
  }
}

export class StartClusterRequest extends jspb.Message {
  getInitialmembersMap(): jspb.Map<number, string>;
  clearInitialmembersMap(): void;
  getJoin(): boolean;
  setJoin(value: boolean): void;

  hasRaftconfig(): boolean;
  clearRaftconfig(): void;
  getRaftconfig(): RaftConfig | undefined;
  setRaftconfig(value?: RaftConfig): void;

  hasConcurrent(): boolean;
  clearConcurrent(): void;
  getConcurrent(): boolean;
  setConcurrent(value: boolean): void;

  hasOndisk(): boolean;
  clearOndisk(): void;
  getOndisk(): boolean;
  setOndisk(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartClusterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartClusterRequest): StartClusterRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StartClusterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartClusterRequest;
  static deserializeBinaryFromReader(message: StartClusterRequest, reader: jspb.BinaryReader): StartClusterRequest;
}

export namespace StartClusterRequest {
  export type AsObject = {
    initialmembersMap: Array<[number, string]>,
    join: boolean,
    raftconfig?: RaftConfig.AsObject,
    concurrent: boolean,
    ondisk: boolean,
  }
}

export class StartClusterResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartClusterResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StartClusterResponse): StartClusterResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StartClusterResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartClusterResponse;
  static deserializeBinaryFromReader(message: StartClusterResponse, reader: jspb.BinaryReader): StartClusterResponse;
}

export namespace StartClusterResponse {
  export type AsObject = {
  }
}

export class RaftConfig extends jspb.Message {
  getNodeid(): number;
  setNodeid(value: number): void;

  getClusterid(): number;
  setClusterid(value: number): void;

  getCheckquorum(): boolean;
  setCheckquorum(value: boolean): void;

  getElectionroundtriptime(): number;
  setElectionroundtriptime(value: number): void;

  getHeartbeatroundtriptime(): number;
  setHeartbeatroundtriptime(value: number): void;

  getSnapshotentries(): number;
  setSnapshotentries(value: number): void;

  getCompactionoverhead(): number;
  setCompactionoverhead(value: number): void;

  getOrderedconfigchange(): boolean;
  setOrderedconfigchange(value: boolean): void;

  getMaxinmemlogsize(): number;
  setMaxinmemlogsize(value: number): void;

  getSnapshotcompressiontype(): number;
  setSnapshotcompressiontype(value: number): void;

  getEntrycompressiontype(): number;
  setEntrycompressiontype(value: number): void;

  getDisableautocompactions(): boolean;
  setDisableautocompactions(value: boolean): void;

  getIsobserver(): boolean;
  setIsobserver(value: boolean): void;

  getIswitness(): boolean;
  setIswitness(value: boolean): void;

  getQuiesce(): boolean;
  setQuiesce(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftConfig.AsObject;
  static toObject(includeInstance: boolean, msg: RaftConfig): RaftConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RaftConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftConfig;
  static deserializeBinaryFromReader(message: RaftConfig, reader: jspb.BinaryReader): RaftConfig;
}

export namespace RaftConfig {
  export type AsObject = {
    nodeid: number,
    clusterid: number,
    checkquorum: boolean,
    electionroundtriptime: number,
    heartbeatroundtriptime: number,
    snapshotentries: number,
    compactionoverhead: number,
    orderedconfigchange: boolean,
    maxinmemlogsize: number,
    snapshotcompressiontype: number,
    entrycompressiontype: number,
    disableautocompactions: boolean,
    isobserver: boolean,
    iswitness: boolean,
    quiesce: boolean,
  }
}

