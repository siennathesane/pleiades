// package: raft
// file: api/v1/raft/raft_shard.proto

import * as jspb from "google-protobuf";

export class AddReplicaRequest extends jspb.Message {
  getReplicaid(): number;
  setReplicaid(value: number): void;

  getShardid(): number;
  setShardid(value: number): void;

  getHostname(): string;
  setHostname(value: string): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddReplicaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddReplicaRequest): AddReplicaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddReplicaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddReplicaRequest;
  static deserializeBinaryFromReader(message: AddReplicaRequest, reader: jspb.BinaryReader): AddReplicaRequest;
}

export namespace AddReplicaRequest {
  export type AsObject = {
    replicaid: number,
    shardid: number,
    hostname: string,
    timeout: number,
  }
}

export class AddReplicaReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddReplicaReply.AsObject;
  static toObject(includeInstance: boolean, msg: AddReplicaReply): AddReplicaReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddReplicaReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddReplicaReply;
  static deserializeBinaryFromReader(message: AddReplicaReply, reader: jspb.BinaryReader): AddReplicaReply;
}

export namespace AddReplicaReply {
  export type AsObject = {
  }
}

export class AddReplicaObserverRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  getReplicaid(): number;
  setReplicaid(value: number): void;

  getHostname(): string;
  setHostname(value: string): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddReplicaObserverRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddReplicaObserverRequest): AddReplicaObserverRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddReplicaObserverRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddReplicaObserverRequest;
  static deserializeBinaryFromReader(message: AddReplicaObserverRequest, reader: jspb.BinaryReader): AddReplicaObserverRequest;
}

export namespace AddReplicaObserverRequest {
  export type AsObject = {
    shardid: number,
    replicaid: number,
    hostname: string,
    timeout: number,
  }
}

export class AddReplicaObserverReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddReplicaObserverReply.AsObject;
  static toObject(includeInstance: boolean, msg: AddReplicaObserverReply): AddReplicaObserverReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddReplicaObserverReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddReplicaObserverReply;
  static deserializeBinaryFromReader(message: AddReplicaObserverReply, reader: jspb.BinaryReader): AddReplicaObserverReply;
}

export namespace AddReplicaObserverReply {
  export type AsObject = {
  }
}

export class AddReplicaWitnessRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  getReplicaid(): number;
  setReplicaid(value: number): void;

  getHostname(): string;
  setHostname(value: string): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddReplicaWitnessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddReplicaWitnessRequest): AddReplicaWitnessRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddReplicaWitnessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddReplicaWitnessRequest;
  static deserializeBinaryFromReader(message: AddReplicaWitnessRequest, reader: jspb.BinaryReader): AddReplicaWitnessRequest;
}

export namespace AddReplicaWitnessRequest {
  export type AsObject = {
    shardid: number,
    replicaid: number,
    hostname: string,
    timeout: number,
  }
}

export class AddReplicaWitnessReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddReplicaWitnessReply.AsObject;
  static toObject(includeInstance: boolean, msg: AddReplicaWitnessReply): AddReplicaWitnessReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddReplicaWitnessReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddReplicaWitnessReply;
  static deserializeBinaryFromReader(message: AddReplicaWitnessReply, reader: jspb.BinaryReader): AddReplicaWitnessReply;
}

export namespace AddReplicaWitnessReply {
  export type AsObject = {
  }
}

export class DeleteReplicaRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  getReplicaid(): number;
  setReplicaid(value: number): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteReplicaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteReplicaRequest): DeleteReplicaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteReplicaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteReplicaRequest;
  static deserializeBinaryFromReader(message: DeleteReplicaRequest, reader: jspb.BinaryReader): DeleteReplicaRequest;
}

export namespace DeleteReplicaRequest {
  export type AsObject = {
    shardid: number,
    replicaid: number,
    timeout: number,
  }
}

export class DeleteReplicaReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteReplicaReply.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteReplicaReply): DeleteReplicaReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteReplicaReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteReplicaReply;
  static deserializeBinaryFromReader(message: DeleteReplicaReply, reader: jspb.BinaryReader): DeleteReplicaReply;
}

export namespace DeleteReplicaReply {
  export type AsObject = {
  }
}

export class GetLeaderIdRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  getReplicaid(): number;
  setReplicaid(value: number): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLeaderIdRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLeaderIdRequest): GetLeaderIdRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetLeaderIdRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLeaderIdRequest;
  static deserializeBinaryFromReader(message: GetLeaderIdRequest, reader: jspb.BinaryReader): GetLeaderIdRequest;
}

export namespace GetLeaderIdRequest {
  export type AsObject = {
    shardid: number,
    replicaid: number,
    timeout: number,
  }
}

export class GetLeaderIdReply extends jspb.Message {
  getLeader(): number;
  setLeader(value: number): void;

  getAvailable(): boolean;
  setAvailable(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLeaderIdReply.AsObject;
  static toObject(includeInstance: boolean, msg: GetLeaderIdReply): GetLeaderIdReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetLeaderIdReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLeaderIdReply;
  static deserializeBinaryFromReader(message: GetLeaderIdReply, reader: jspb.BinaryReader): GetLeaderIdReply;
}

export namespace GetLeaderIdReply {
  export type AsObject = {
    leader: number,
    available: boolean,
  }
}

export class GetShardMembersRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetShardMembersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetShardMembersRequest): GetShardMembersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetShardMembersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetShardMembersRequest;
  static deserializeBinaryFromReader(message: GetShardMembersRequest, reader: jspb.BinaryReader): GetShardMembersRequest;
}

export namespace GetShardMembersRequest {
  export type AsObject = {
    shardid: number,
  }
}

export class GetShardMembersReply extends jspb.Message {
  getConfigchangeid(): number;
  setConfigchangeid(value: number): void;

  getReplicasMap(): jspb.Map<number, string>;
  clearReplicasMap(): void;
  getObserversMap(): jspb.Map<number, string>;
  clearObserversMap(): void;
  getWitnessesMap(): jspb.Map<number, string>;
  clearWitnessesMap(): void;
  getRemovedMap(): jspb.Map<number, string>;
  clearRemovedMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetShardMembersReply.AsObject;
  static toObject(includeInstance: boolean, msg: GetShardMembersReply): GetShardMembersReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetShardMembersReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetShardMembersReply;
  static deserializeBinaryFromReader(message: GetShardMembersReply, reader: jspb.BinaryReader): GetShardMembersReply;
}

export namespace GetShardMembersReply {
  export type AsObject = {
    configchangeid: number,
    replicasMap: Array<[number, string]>,
    observersMap: Array<[number, string]>,
    witnessesMap: Array<[number, string]>,
    removedMap: Array<[number, string]>,
  }
}

export class NewShardRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  getReplicaid(): number;
  setReplicaid(value: number): void;

  getType(): StateMachineTypeMap[keyof StateMachineTypeMap];
  setType(value: StateMachineTypeMap[keyof StateMachineTypeMap]): void;

  getHostname(): string;
  setHostname(value: string): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewShardRequest.AsObject;
  static toObject(includeInstance: boolean, msg: NewShardRequest): NewShardRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NewShardRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewShardRequest;
  static deserializeBinaryFromReader(message: NewShardRequest, reader: jspb.BinaryReader): NewShardRequest;
}

export namespace NewShardRequest {
  export type AsObject = {
    shardid: number,
    replicaid: number,
    type: StateMachineTypeMap[keyof StateMachineTypeMap],
    hostname: string,
    timeout: number,
  }
}

export class NewShardReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewShardReply.AsObject;
  static toObject(includeInstance: boolean, msg: NewShardReply): NewShardReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NewShardReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewShardReply;
  static deserializeBinaryFromReader(message: NewShardReply, reader: jspb.BinaryReader): NewShardReply;
}

export namespace NewShardReply {
  export type AsObject = {
  }
}

export class RemoveDataRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  getReplicaid(): number;
  setReplicaid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveDataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveDataRequest): RemoveDataRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RemoveDataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveDataRequest;
  static deserializeBinaryFromReader(message: RemoveDataRequest, reader: jspb.BinaryReader): RemoveDataRequest;
}

export namespace RemoveDataRequest {
  export type AsObject = {
    shardid: number,
    replicaid: number,
  }
}

export class RemoveDataReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveDataReply.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveDataReply): RemoveDataReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RemoveDataReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveDataReply;
  static deserializeBinaryFromReader(message: RemoveDataReply, reader: jspb.BinaryReader): RemoveDataReply;
}

export namespace RemoveDataReply {
  export type AsObject = {
  }
}

export class StartReplicaRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  getReplicaid(): number;
  setReplicaid(value: number): void;

  getType(): StateMachineTypeMap[keyof StateMachineTypeMap];
  setType(value: StateMachineTypeMap[keyof StateMachineTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartReplicaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartReplicaRequest): StartReplicaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StartReplicaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartReplicaRequest;
  static deserializeBinaryFromReader(message: StartReplicaRequest, reader: jspb.BinaryReader): StartReplicaRequest;
}

export namespace StartReplicaRequest {
  export type AsObject = {
    shardid: number,
    replicaid: number,
    type: StateMachineTypeMap[keyof StateMachineTypeMap],
  }
}

export class StartReplicaReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartReplicaReply.AsObject;
  static toObject(includeInstance: boolean, msg: StartReplicaReply): StartReplicaReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StartReplicaReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartReplicaReply;
  static deserializeBinaryFromReader(message: StartReplicaReply, reader: jspb.BinaryReader): StartReplicaReply;
}

export namespace StartReplicaReply {
  export type AsObject = {
  }
}

export class StopReplicaRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopReplicaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StopReplicaRequest): StopReplicaRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StopReplicaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopReplicaRequest;
  static deserializeBinaryFromReader(message: StopReplicaRequest, reader: jspb.BinaryReader): StopReplicaRequest;
}

export namespace StopReplicaRequest {
  export type AsObject = {
    shardid: number,
  }
}

export class StopReplicaReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopReplicaReply.AsObject;
  static toObject(includeInstance: boolean, msg: StopReplicaReply): StopReplicaReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StopReplicaReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopReplicaReply;
  static deserializeBinaryFromReader(message: StopReplicaReply, reader: jspb.BinaryReader): StopReplicaReply;
}

export namespace StopReplicaReply {
  export type AsObject = {
  }
}

export interface StateMachineTypeMap {
  TEST: 0;
  KV: 1;
}

export const StateMachineType: StateMachineTypeMap;

