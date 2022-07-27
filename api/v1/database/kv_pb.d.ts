// package: database
// file: api/v1/database/kv.proto

import * as jspb from "google-protobuf";

export class KeyValue extends jspb.Message {
  getKey(): Uint8Array | string;
  getKey_asU8(): Uint8Array;
  getKey_asB64(): string;
  setKey(value: Uint8Array | string): void;

  getCreateRevision(): number;
  setCreateRevision(value: number): void;

  getModRevision(): number;
  setModRevision(value: number): void;

  getVersion(): number;
  setVersion(value: number): void;

  getValue(): Uint8Array | string;
  getValue_asU8(): Uint8Array;
  getValue_asB64(): string;
  setValue(value: Uint8Array | string): void;

  getLease(): number;
  setLease(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KeyValue.AsObject;
  static toObject(includeInstance: boolean, msg: KeyValue): KeyValue.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: KeyValue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KeyValue;
  static deserializeBinaryFromReader(message: KeyValue, reader: jspb.BinaryReader): KeyValue;
}

export namespace KeyValue {
  export type AsObject = {
    key: Uint8Array | string,
    createRevision: number,
    modRevision: number,
    version: number,
    value: Uint8Array | string,
    lease: number,
  }
}

export class Event extends jspb.Message {
  getType(): Event.EventTypeMap[keyof Event.EventTypeMap];
  setType(value: Event.EventTypeMap[keyof Event.EventTypeMap]): void;

  hasKv(): boolean;
  clearKv(): void;
  getKv(): KeyValue | undefined;
  setKv(value?: KeyValue): void;

  hasPrevKv(): boolean;
  clearPrevKv(): void;
  getPrevKv(): KeyValue | undefined;
  setPrevKv(value?: KeyValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Event.AsObject;
  static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Event;
  static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
  export type AsObject = {
    type: Event.EventTypeMap[keyof Event.EventTypeMap],
    kv?: KeyValue.AsObject,
    prevKv?: KeyValue.AsObject,
  }

  export interface EventTypeMap {
    PUT: 0;
    DELETE: 1;
  }

  export const EventType: EventTypeMap;
}

