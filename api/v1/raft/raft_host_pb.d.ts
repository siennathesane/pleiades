// package: raft
// file: api/v1/raft/raft_host.proto

import * as jspb from "google-protobuf";

export class CompactRequest extends jspb.Message {
  getReplicaid(): number;
  setReplicaid(value: number): void;

  getShardid(): number;
  setShardid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompactRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompactRequest): CompactRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CompactRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompactRequest;
  static deserializeBinaryFromReader(message: CompactRequest, reader: jspb.BinaryReader): CompactRequest;
}

export namespace CompactRequest {
  export type AsObject = {
    replicaid: number,
    shardid: number,
  }
}

export class CompactReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompactReply.AsObject;
  static toObject(includeInstance: boolean, msg: CompactReply): CompactReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CompactReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompactReply;
  static deserializeBinaryFromReader(message: CompactReply, reader: jspb.BinaryReader): CompactReply;
}

export namespace CompactReply {
  export type AsObject = {
  }
}

export class LeaderTransferRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  getTargetnodeid(): string;
  setTargetnodeid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LeaderTransferRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LeaderTransferRequest): LeaderTransferRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LeaderTransferRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LeaderTransferRequest;
  static deserializeBinaryFromReader(message: LeaderTransferRequest, reader: jspb.BinaryReader): LeaderTransferRequest;
}

export namespace LeaderTransferRequest {
  export type AsObject = {
    shardid: number,
    targetnodeid: string,
  }
}

export class LeaderTransferReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LeaderTransferReply.AsObject;
  static toObject(includeInstance: boolean, msg: LeaderTransferReply): LeaderTransferReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LeaderTransferReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LeaderTransferReply;
  static deserializeBinaryFromReader(message: LeaderTransferReply, reader: jspb.BinaryReader): LeaderTransferReply;
}

export namespace LeaderTransferReply {
  export type AsObject = {
  }
}

export class SnapshotRequest extends jspb.Message {
  getShardid(): number;
  setShardid(value: number): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SnapshotRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SnapshotRequest): SnapshotRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SnapshotRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SnapshotRequest;
  static deserializeBinaryFromReader(message: SnapshotRequest, reader: jspb.BinaryReader): SnapshotRequest;
}

export namespace SnapshotRequest {
  export type AsObject = {
    shardid: number,
    timeout: number,
  }
}

export class SnapshotReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SnapshotReply.AsObject;
  static toObject(includeInstance: boolean, msg: SnapshotReply): SnapshotReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SnapshotReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SnapshotReply;
  static deserializeBinaryFromReader(message: SnapshotReply, reader: jspb.BinaryReader): SnapshotReply;
}

export namespace SnapshotReply {
  export type AsObject = {
  }
}

export class StopRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StopRequest): StopRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StopRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopRequest;
  static deserializeBinaryFromReader(message: StopRequest, reader: jspb.BinaryReader): StopRequest;
}

export namespace StopRequest {
  export type AsObject = {
  }
}

export class StopReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopReply.AsObject;
  static toObject(includeInstance: boolean, msg: StopReply): StopReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StopReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopReply;
  static deserializeBinaryFromReader(message: StopReply, reader: jspb.BinaryReader): StopReply;
}

export namespace StopReply {
  export type AsObject = {
  }
}

export class GetHostConfigRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHostConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetHostConfigRequest): GetHostConfigRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetHostConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHostConfigRequest;
  static deserializeBinaryFromReader(message: GetHostConfigRequest, reader: jspb.BinaryReader): GetHostConfigRequest;
}

export namespace GetHostConfigRequest {
  export type AsObject = {
  }
}

export class GetHostConfigReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHostConfigReply.AsObject;
  static toObject(includeInstance: boolean, msg: GetHostConfigReply): GetHostConfigReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetHostConfigReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHostConfigReply;
  static deserializeBinaryFromReader(message: GetHostConfigReply, reader: jspb.BinaryReader): GetHostConfigReply;
}

export namespace GetHostConfigReply {
  export type AsObject = {
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

