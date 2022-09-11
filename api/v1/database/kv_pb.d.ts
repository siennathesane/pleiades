// package: database
// file: api/v1/database/kv.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class PayloadWrapper extends jspb.Message {
  getAccount(): Uint8Array | string;
  getAccount_asU8(): Uint8Array;
  getAccount_asB64(): string;
  setAccount(value: Uint8Array | string): void;

  getBucket(): Uint8Array | string;
  getBucket_asU8(): Uint8Array;
  getBucket_asB64(): string;
  setBucket(value: Uint8Array | string): void;

  getTyp(): PayloadWrapper.RequestTypeMap[keyof PayloadWrapper.RequestTypeMap];
  setTyp(value: PayloadWrapper.RequestTypeMap[keyof PayloadWrapper.RequestTypeMap]): void;

  hasCreateBucketRequest(): boolean;
  clearCreateBucketRequest(): void;
  getCreateBucketRequest(): CreateBucketRequest | undefined;
  setCreateBucketRequest(value?: CreateBucketRequest): void;

  hasCreateBucketReply(): boolean;
  clearCreateBucketReply(): void;
  getCreateBucketReply(): CreateBucketReply | undefined;
  setCreateBucketReply(value?: CreateBucketReply): void;

  hasDeleteBucketRequest(): boolean;
  clearDeleteBucketRequest(): void;
  getDeleteBucketRequest(): DeleteBucketRequest | undefined;
  setDeleteBucketRequest(value?: DeleteBucketRequest): void;

  hasDeleteBucketReply(): boolean;
  clearDeleteBucketReply(): void;
  getDeleteBucketReply(): DeleteBucketReply | undefined;
  setDeleteBucketReply(value?: DeleteBucketReply): void;

  getPayloadCase(): PayloadWrapper.PayloadCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PayloadWrapper.AsObject;
  static toObject(includeInstance: boolean, msg: PayloadWrapper): PayloadWrapper.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PayloadWrapper, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PayloadWrapper;
  static deserializeBinaryFromReader(message: PayloadWrapper, reader: jspb.BinaryReader): PayloadWrapper;
}

export namespace PayloadWrapper {
  export type AsObject = {
    account: Uint8Array | string,
    bucket: Uint8Array | string,
    typ: PayloadWrapper.RequestTypeMap[keyof PayloadWrapper.RequestTypeMap],
    createBucketRequest?: CreateBucketRequest.AsObject,
    createBucketReply?: CreateBucketReply.AsObject,
    deleteBucketRequest?: DeleteBucketRequest.AsObject,
    deleteBucketReply?: DeleteBucketReply.AsObject,
  }

  export interface RequestTypeMap {
    CREATE_BUCKET_REQUEST: 0;
    CREATE_BUCKET_REPLY: 1;
    DELETE_BUCKET_REQUEST: 2;
    DELETE_BUCKET_REPLY: 3;
    GET_KEY_REQUEST: 4;
    GET_KEY_REPLY: 5;
    PUT_KEY_REQUEST: 6;
    PUT_KEY_REPLY: 7;
    DELETE_KEY_REQUEST: 8;
    DELETE_KEY_REPLY: 9;
  }

  export const RequestType: RequestTypeMap;

  export enum PayloadCase {
    PAYLOAD_NOT_SET = 0,
    CREATE_BUCKET_REQUEST = 4,
    CREATE_BUCKET_REPLY = 5,
    DELETE_BUCKET_REQUEST = 6,
    DELETE_BUCKET_REPLY = 7,
  }
}

export class CreateAccountRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getOwner(): string;
  setOwner(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountRequest): CreateAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountRequest;
  static deserializeBinaryFromReader(message: CreateAccountRequest, reader: jspb.BinaryReader): CreateAccountRequest;
}

export namespace CreateAccountRequest {
  export type AsObject = {
    accountId: number,
    owner: string,
  }
}

export class CreateAccountReply extends jspb.Message {
  hasAccountDescriptor(): boolean;
  clearAccountDescriptor(): void;
  getAccountDescriptor(): AccountDescriptor | undefined;
  setAccountDescriptor(value?: AccountDescriptor): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountReply.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountReply): CreateAccountReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountReply;
  static deserializeBinaryFromReader(message: CreateAccountReply, reader: jspb.BinaryReader): CreateAccountReply;
}

export namespace CreateAccountReply {
  export type AsObject = {
    accountDescriptor?: AccountDescriptor.AsObject,
  }
}

export class DeleteAccountRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getOwner(): string;
  setOwner(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAccountRequest): DeleteAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountRequest;
  static deserializeBinaryFromReader(message: DeleteAccountRequest, reader: jspb.BinaryReader): DeleteAccountRequest;
}

export namespace DeleteAccountRequest {
  export type AsObject = {
    accountId: number,
    owner: string,
  }
}

export class DeleteAccountReply extends jspb.Message {
  getOk(): boolean;
  setOk(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountReply.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAccountReply): DeleteAccountReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteAccountReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountReply;
  static deserializeBinaryFromReader(message: DeleteAccountReply, reader: jspb.BinaryReader): DeleteAccountReply;
}

export namespace DeleteAccountReply {
  export type AsObject = {
    ok: boolean,
  }
}

export class GetAccountDescriptorRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountDescriptorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountDescriptorRequest): GetAccountDescriptorRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAccountDescriptorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountDescriptorRequest;
  static deserializeBinaryFromReader(message: GetAccountDescriptorRequest, reader: jspb.BinaryReader): GetAccountDescriptorRequest;
}

export namespace GetAccountDescriptorRequest {
  export type AsObject = {
    accountId: number,
  }
}

export class GetAccountDescriptorReply extends jspb.Message {
  hasAccountDescriptor(): boolean;
  clearAccountDescriptor(): void;
  getAccountDescriptor(): AccountDescriptor | undefined;
  setAccountDescriptor(value?: AccountDescriptor): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountDescriptorReply.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountDescriptorReply): GetAccountDescriptorReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetAccountDescriptorReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountDescriptorReply;
  static deserializeBinaryFromReader(message: GetAccountDescriptorReply, reader: jspb.BinaryReader): GetAccountDescriptorReply;
}

export namespace GetAccountDescriptorReply {
  export type AsObject = {
    accountDescriptor?: AccountDescriptor.AsObject,
  }
}

export class AccountDescriptor extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getOwner(): string;
  setOwner(value: string): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasLastUpdated(): boolean;
  clearLastUpdated(): void;
  getLastUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getBucketCount(): number;
  setBucketCount(value: number): void;

  clearBucketsList(): void;
  getBucketsList(): Array<string>;
  setBucketsList(value: Array<string>): void;
  addBuckets(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountDescriptor.AsObject;
  static toObject(includeInstance: boolean, msg: AccountDescriptor): AccountDescriptor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AccountDescriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccountDescriptor;
  static deserializeBinaryFromReader(message: AccountDescriptor, reader: jspb.BinaryReader): AccountDescriptor;
}

export namespace AccountDescriptor {
  export type AsObject = {
    accountId: number,
    owner: string,
    created?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lastUpdated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    bucketCount: number,
    bucketsList: Array<string>,
  }
}

export class CreateBucketRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getName(): string;
  setName(value: string): void;

  getOwner(): string;
  setOwner(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateBucketRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateBucketRequest): CreateBucketRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateBucketRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateBucketRequest;
  static deserializeBinaryFromReader(message: CreateBucketRequest, reader: jspb.BinaryReader): CreateBucketRequest;
}

export namespace CreateBucketRequest {
  export type AsObject = {
    accountId: number,
    name: string,
    owner: string,
  }
}

export class CreateBucketReply extends jspb.Message {
  hasBucketDescriptor(): boolean;
  clearBucketDescriptor(): void;
  getBucketDescriptor(): BucketDescriptor | undefined;
  setBucketDescriptor(value?: BucketDescriptor): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateBucketReply.AsObject;
  static toObject(includeInstance: boolean, msg: CreateBucketReply): CreateBucketReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateBucketReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateBucketReply;
  static deserializeBinaryFromReader(message: CreateBucketReply, reader: jspb.BinaryReader): CreateBucketReply;
}

export namespace CreateBucketReply {
  export type AsObject = {
    bucketDescriptor?: BucketDescriptor.AsObject,
  }
}

export class DeleteBucketRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteBucketRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteBucketRequest): DeleteBucketRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteBucketRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteBucketRequest;
  static deserializeBinaryFromReader(message: DeleteBucketRequest, reader: jspb.BinaryReader): DeleteBucketRequest;
}

export namespace DeleteBucketRequest {
  export type AsObject = {
    accountId: number,
    name: string,
  }
}

export class DeleteBucketReply extends jspb.Message {
  getOk(): boolean;
  setOk(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteBucketReply.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteBucketReply): DeleteBucketReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteBucketReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteBucketReply;
  static deserializeBinaryFromReader(message: DeleteBucketReply, reader: jspb.BinaryReader): DeleteBucketReply;
}

export namespace DeleteBucketReply {
  export type AsObject = {
    ok: boolean,
  }
}

export class BucketDescriptor extends jspb.Message {
  getOwner(): string;
  setOwner(value: string): void;

  getSize(): number;
  setSize(value: number): void;

  getKeyCount(): number;
  setKeyCount(value: number): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasLastUpdated(): boolean;
  clearLastUpdated(): void;
  getLastUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BucketDescriptor.AsObject;
  static toObject(includeInstance: boolean, msg: BucketDescriptor): BucketDescriptor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BucketDescriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BucketDescriptor;
  static deserializeBinaryFromReader(message: BucketDescriptor, reader: jspb.BinaryReader): BucketDescriptor;
}

export namespace BucketDescriptor {
  export type AsObject = {
    owner: string,
    size: number,
    keyCount: number,
    created?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lastUpdated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class GetKeyRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getBucketId(): string;
  setBucketId(value: string): void;

  getKeyId(): string;
  setKeyId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetKeyRequest): GetKeyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetKeyRequest;
  static deserializeBinaryFromReader(message: GetKeyRequest, reader: jspb.BinaryReader): GetKeyRequest;
}

export namespace GetKeyRequest {
  export type AsObject = {
    accountId: number,
    bucketId: string,
    keyId: string,
  }
}

export class GetKeyReply extends jspb.Message {
  hasKeyValuePair(): boolean;
  clearKeyValuePair(): void;
  getKeyValuePair(): KeyValue | undefined;
  setKeyValuePair(value?: KeyValue): void;

  getSize(): number;
  setSize(value: number): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasLastUpdated(): boolean;
  clearLastUpdated(): void;
  getLastUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetKeyReply.AsObject;
  static toObject(includeInstance: boolean, msg: GetKeyReply): GetKeyReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetKeyReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetKeyReply;
  static deserializeBinaryFromReader(message: GetKeyReply, reader: jspb.BinaryReader): GetKeyReply;
}

export namespace GetKeyReply {
  export type AsObject = {
    keyValuePair?: KeyValue.AsObject,
    size: number,
    created?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lastUpdated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class PutKeyRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getBucketId(): string;
  setBucketId(value: string): void;

  hasKeyValuePair(): boolean;
  clearKeyValuePair(): void;
  getKeyValuePair(): KeyValue | undefined;
  setKeyValuePair(value?: KeyValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PutKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PutKeyRequest): PutKeyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PutKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PutKeyRequest;
  static deserializeBinaryFromReader(message: PutKeyRequest, reader: jspb.BinaryReader): PutKeyRequest;
}

export namespace PutKeyRequest {
  export type AsObject = {
    accountId: number,
    bucketId: string,
    keyValuePair?: KeyValue.AsObject,
  }
}

export class PutKeyReply extends jspb.Message {
  hasKeyValuePair(): boolean;
  clearKeyValuePair(): void;
  getKeyValuePair(): KeyValue | undefined;
  setKeyValuePair(value?: KeyValue): void;

  getSize(): number;
  setSize(value: number): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PutKeyReply.AsObject;
  static toObject(includeInstance: boolean, msg: PutKeyReply): PutKeyReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PutKeyReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PutKeyReply;
  static deserializeBinaryFromReader(message: PutKeyReply, reader: jspb.BinaryReader): PutKeyReply;
}

export namespace PutKeyReply {
  export type AsObject = {
    keyValuePair?: KeyValue.AsObject,
    size: number,
    created?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class DeleteKeyRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getBucketId(): string;
  setBucketId(value: string): void;

  getKeyId(): string;
  setKeyId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteKeyRequest): DeleteKeyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteKeyRequest;
  static deserializeBinaryFromReader(message: DeleteKeyRequest, reader: jspb.BinaryReader): DeleteKeyRequest;
}

export namespace DeleteKeyRequest {
  export type AsObject = {
    accountId: number,
    bucketId: string,
    keyId: string,
  }
}

export class DeleteKeyReply extends jspb.Message {
  getOk(): boolean;
  setOk(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteKeyReply.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteKeyReply): DeleteKeyReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteKeyReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteKeyReply;
  static deserializeBinaryFromReader(message: DeleteKeyReply, reader: jspb.BinaryReader): DeleteKeyReply;
}

export namespace DeleteKeyReply {
  export type AsObject = {
    ok: boolean,
  }
}

export class KeyValue extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

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
    key: string,
    createRevision: number,
    modRevision: number,
    version: number,
    value: Uint8Array | string,
    lease: number,
  }
}

export class Event extends jspb.Message {
  getType(): KeyOperationTypeMap[keyof KeyOperationTypeMap];
  setType(value: KeyOperationTypeMap[keyof KeyOperationTypeMap]): void;

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
    type: KeyOperationTypeMap[keyof KeyOperationTypeMap],
    kv?: KeyValue.AsObject,
    prevKv?: KeyValue.AsObject,
  }
}

export interface KeyOperationTypeMap {
  GET: 0;
  PUT: 1;
  DELETE: 2;
}

export const KeyOperationType: KeyOperationTypeMap;

