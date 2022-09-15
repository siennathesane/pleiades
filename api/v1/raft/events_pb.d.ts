// package: raft
// file: api/v1/raft/events.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class RaftEvent extends jspb.Message {
  getTyp(): EventTypeMap[keyof EventTypeMap];
  setTyp(value: EventTypeMap[keyof EventTypeMap]): void;

  getAction(): EventMap[keyof EventMap];
  setAction(value: EventMap[keyof EventMap]): void;

  hasTimestamp(): boolean;
  clearTimestamp(): void;
  getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasLogEntry(): boolean;
  clearLogEntry(): void;
  getLogEntry(): RaftLogEntryEvent | undefined;
  setLogEntry(value?: RaftLogEntryEvent): void;

  hasSnapshot(): boolean;
  clearSnapshot(): void;
  getSnapshot(): RaftSnapshotEvent | undefined;
  setSnapshot(value?: RaftSnapshotEvent): void;

  hasConnection(): boolean;
  clearConnection(): void;
  getConnection(): RaftConnectionEvent | undefined;
  setConnection(value?: RaftConnectionEvent): void;

  hasNode(): boolean;
  clearNode(): void;
  getNode(): RaftNodeEvent | undefined;
  setNode(value?: RaftNodeEvent): void;

  hasHostShutdown(): boolean;
  clearHostShutdown(): void;
  getHostShutdown(): RaftHostShutdown | undefined;
  setHostShutdown(value?: RaftHostShutdown): void;

  getEventCase(): RaftEvent.EventCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftEvent.AsObject;
  static toObject(includeInstance: boolean, msg: RaftEvent): RaftEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RaftEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftEvent;
  static deserializeBinaryFromReader(message: RaftEvent, reader: jspb.BinaryReader): RaftEvent;
}

export namespace RaftEvent {
  export type AsObject = {
    typ: EventTypeMap[keyof EventTypeMap],
    action: EventMap[keyof EventMap],
    timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    logEntry?: RaftLogEntryEvent.AsObject,
    snapshot?: RaftSnapshotEvent.AsObject,
    connection?: RaftConnectionEvent.AsObject,
    node?: RaftNodeEvent.AsObject,
    hostShutdown?: RaftHostShutdown.AsObject,
  }

  export enum EventCase {
    EVENT_NOT_SET = 0,
    LOG_ENTRY = 4,
    SNAPSHOT = 5,
    CONNECTION = 6,
    NODE = 7,
    HOST_SHUTDOWN = 8,
  }
}

export class RaftLogEntryEvent extends jspb.Message {
  getShardId(): number;
  setShardId(value: number): void;

  getReplicaId(): number;
  setReplicaId(value: number): void;

  getIndex(): number;
  setIndex(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftLogEntryEvent.AsObject;
  static toObject(includeInstance: boolean, msg: RaftLogEntryEvent): RaftLogEntryEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RaftLogEntryEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftLogEntryEvent;
  static deserializeBinaryFromReader(message: RaftLogEntryEvent, reader: jspb.BinaryReader): RaftLogEntryEvent;
}

export namespace RaftLogEntryEvent {
  export type AsObject = {
    shardId: number,
    replicaId: number,
    index: number,
  }
}

export class RaftSnapshotEvent extends jspb.Message {
  getShardId(): number;
  setShardId(value: number): void;

  getReplicaId(): number;
  setReplicaId(value: number): void;

  getFromIndex(): number;
  setFromIndex(value: number): void;

  getToIndex(): number;
  setToIndex(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftSnapshotEvent.AsObject;
  static toObject(includeInstance: boolean, msg: RaftSnapshotEvent): RaftSnapshotEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RaftSnapshotEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftSnapshotEvent;
  static deserializeBinaryFromReader(message: RaftSnapshotEvent, reader: jspb.BinaryReader): RaftSnapshotEvent;
}

export namespace RaftSnapshotEvent {
  export type AsObject = {
    shardId: number,
    replicaId: number,
    fromIndex: number,
    toIndex: number,
  }
}

export class RaftConnectionEvent extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): void;

  getIsSnapshot(): boolean;
  setIsSnapshot(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftConnectionEvent.AsObject;
  static toObject(includeInstance: boolean, msg: RaftConnectionEvent): RaftConnectionEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RaftConnectionEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftConnectionEvent;
  static deserializeBinaryFromReader(message: RaftConnectionEvent, reader: jspb.BinaryReader): RaftConnectionEvent;
}

export namespace RaftConnectionEvent {
  export type AsObject = {
    address: string,
    isSnapshot: boolean,
  }
}

export class RaftNodeEvent extends jspb.Message {
  getShardId(): number;
  setShardId(value: number): void;

  getReplicaId(): number;
  setReplicaId(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftNodeEvent.AsObject;
  static toObject(includeInstance: boolean, msg: RaftNodeEvent): RaftNodeEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RaftNodeEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftNodeEvent;
  static deserializeBinaryFromReader(message: RaftNodeEvent, reader: jspb.BinaryReader): RaftNodeEvent;
}

export namespace RaftNodeEvent {
  export type AsObject = {
    shardId: number,
    replicaId: number,
  }
}

export class RaftHostShutdown extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftHostShutdown.AsObject;
  static toObject(includeInstance: boolean, msg: RaftHostShutdown): RaftHostShutdown.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RaftHostShutdown, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftHostShutdown;
  static deserializeBinaryFromReader(message: RaftHostShutdown, reader: jspb.BinaryReader): RaftHostShutdown;
}

export namespace RaftHostShutdown {
  export type AsObject = {
  }
}

export interface EventTypeMap {
  LOG_ENTRY: 0;
  SNAPSHOT: 1;
  CONNECTION: 2;
  HOST: 3;
  NODE: 4;
}

export const EventType: EventTypeMap;

export interface EventMap {
  CONNECTION_ESTABLISHED: 0;
  CONNECTION_FAILED: 1;
  LOG_COMPACTED: 2;
  LOGDB_COMPACTED: 3;
  MEMBERSHIP_CHANGED: 4;
  NODE_HOST_SHUTTING_DOWN: 5;
  NODE_READY: 6;
  NODE_UNLOADED: 7;
  SEND_SNAPSHOT_ABORTED: 8;
  SEND_SNAPSHOT_COMPLETED: 9;
  SEND_SNAPSHOT_STARTED: 10;
  SNAPSHOT_COMPACTED: 11;
  SNAPSHOT_CREATED: 12;
  SNAPSHOT_RECEIVED: 13;
  SNAPSHOT_RECOVERED: 14;
}

export const Event: EventMap;

