// package: database
// file: api/v1/database/kv.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as api_v1_database_transactions_pb from "../../../api/v1/database/transactions_pb";

export class KVStoreWrapper extends jspb.Message {
  getAccount(): number;
  setAccount(value: number): void;

  getBucket(): string;
  setBucket(value: string): void;

  getTyp(): KVStoreWrapper.RequestTypeMap[keyof KVStoreWrapper.RequestTypeMap];
  setTyp(value: KVStoreWrapper.RequestTypeMap[keyof KVStoreWrapper.RequestTypeMap]): void;

  hasCreateAccountRequest(): boolean;
  clearCreateAccountRequest(): void;
  getCreateAccountRequest(): CreateAccountRequest | undefined;
  setCreateAccountRequest(value?: CreateAccountRequest): void;

  hasCreateAccountReply(): boolean;
  clearCreateAccountReply(): void;
  getCreateAccountReply(): CreateAccountReply | undefined;
  setCreateAccountReply(value?: CreateAccountReply): void;

  hasDeleteAccountRequest(): boolean;
  clearDeleteAccountRequest(): void;
  getDeleteAccountRequest(): DeleteAccountRequest | undefined;
  setDeleteAccountRequest(value?: DeleteAccountRequest): void;

  hasDeleteAccountReply(): boolean;
  clearDeleteAccountReply(): void;
  getDeleteAccountReply(): DeleteAccountReply | undefined;
  setDeleteAccountReply(value?: DeleteAccountReply): void;

  hasGetAccountDescriptorRequest(): boolean;
  clearGetAccountDescriptorRequest(): void;
  getGetAccountDescriptorRequest(): GetAccountDescriptorRequest | undefined;
  setGetAccountDescriptorRequest(value?: GetAccountDescriptorRequest): void;

  hasGetAccountDescriptorReply(): boolean;
  clearGetAccountDescriptorReply(): void;
  getGetAccountDescriptorReply(): GetAccountDescriptorReply | undefined;
  setGetAccountDescriptorReply(value?: GetAccountDescriptorReply): void;

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

  hasGetKeyRequest(): boolean;
  clearGetKeyRequest(): void;
  getGetKeyRequest(): GetKeyRequest | undefined;
  setGetKeyRequest(value?: GetKeyRequest): void;

  hasGetKeyReply(): boolean;
  clearGetKeyReply(): void;
  getGetKeyReply(): GetKeyReply | undefined;
  setGetKeyReply(value?: GetKeyReply): void;

  hasPutKeyRequest(): boolean;
  clearPutKeyRequest(): void;
  getPutKeyRequest(): PutKeyRequest | undefined;
  setPutKeyRequest(value?: PutKeyRequest): void;

  hasPutKeyReply(): boolean;
  clearPutKeyReply(): void;
  getPutKeyReply(): PutKeyReply | undefined;
  setPutKeyReply(value?: PutKeyReply): void;

  hasDeleteKeyRequest(): boolean;
  clearDeleteKeyRequest(): void;
  getDeleteKeyRequest(): DeleteKeyRequest | undefined;
  setDeleteKeyRequest(value?: DeleteKeyRequest): void;

  hasDeleteKeyReply(): boolean;
  clearDeleteKeyReply(): void;
  getDeleteKeyReply(): DeleteKeyReply | undefined;
  setDeleteKeyReply(value?: DeleteKeyReply): void;

  getPayloadCase(): KVStoreWrapper.PayloadCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KVStoreWrapper.AsObject;
  static toObject(includeInstance: boolean, msg: KVStoreWrapper): KVStoreWrapper.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: KVStoreWrapper, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KVStoreWrapper;
  static deserializeBinaryFromReader(message: KVStoreWrapper, reader: jspb.BinaryReader): KVStoreWrapper;
}

export namespace KVStoreWrapper {
  export type AsObject = {
    account: number,
    bucket: string,
    typ: KVStoreWrapper.RequestTypeMap[keyof KVStoreWrapper.RequestTypeMap],
    createAccountRequest?: CreateAccountRequest.AsObject,
    createAccountReply?: CreateAccountReply.AsObject,
    deleteAccountRequest?: DeleteAccountRequest.AsObject,
    deleteAccountReply?: DeleteAccountReply.AsObject,
    getAccountDescriptorRequest?: GetAccountDescriptorRequest.AsObject,
    getAccountDescriptorReply?: GetAccountDescriptorReply.AsObject,
    createBucketRequest?: CreateBucketRequest.AsObject,
    createBucketReply?: CreateBucketReply.AsObject,
    deleteBucketRequest?: DeleteBucketRequest.AsObject,
    deleteBucketReply?: DeleteBucketReply.AsObject,
    getKeyRequest?: GetKeyRequest.AsObject,
    getKeyReply?: GetKeyReply.AsObject,
    putKeyRequest?: PutKeyRequest.AsObject,
    putKeyReply?: PutKeyReply.AsObject,
    deleteKeyRequest?: DeleteKeyRequest.AsObject,
    deleteKeyReply?: DeleteKeyReply.AsObject,
  }

  export interface RequestTypeMap {
    CREATE_ACCOUNT_REQUEST: 0;
    CREATE_ACCOUNT_REPLY: 1;
    DELETE_ACCOUNT_REQUEST: 2;
    DELETE_ACCOUNT_REPLY: 3;
    GET_ACCOUNT_DESCRIPTOR_REQUEST: 4;
    GET_ACCOUNT_DESCRIPTOR_REPLY: 5;
    CREATE_BUCKET_REQUEST: 6;
    CREATE_BUCKET_REPLY: 7;
    DELETE_BUCKET_REQUEST: 8;
    DELETE_BUCKET_REPLY: 9;
    GET_KEY_REQUEST: 10;
    GET_KEY_REPLY: 11;
    PUT_KEY_REQUEST: 12;
    PUT_KEY_REPLY: 13;
    DELETE_KEY_REQUEST: 14;
    DELETE_KEY_REPLY: 15;
  }

  export const RequestType: RequestTypeMap;

  export enum PayloadCase {
    PAYLOAD_NOT_SET = 0,
    CREATE_ACCOUNT_REQUEST = 4,
    CREATE_ACCOUNT_REPLY = 5,
    DELETE_ACCOUNT_REQUEST = 6,
    DELETE_ACCOUNT_REPLY = 7,
    GET_ACCOUNT_DESCRIPTOR_REQUEST = 8,
    GET_ACCOUNT_DESCRIPTOR_REPLY = 9,
    CREATE_BUCKET_REQUEST = 10,
    CREATE_BUCKET_REPLY = 11,
    DELETE_BUCKET_REQUEST = 12,
    DELETE_BUCKET_REPLY = 13,
    GET_KEY_REQUEST = 14,
    GET_KEY_REPLY = 15,
    PUT_KEY_REQUEST = 16,
    PUT_KEY_REPLY = 17,
    DELETE_KEY_REQUEST = 18,
    DELETE_KEY_REPLY = 19,
  }
}

export class CreateAccountRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getOwner(): string;
  setOwner(value: string): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class CreateAccountReply extends jspb.Message {
  hasAccountDescriptor(): boolean;
  clearAccountDescriptor(): void;
  getAccountDescriptor(): AccountDescriptor | undefined;
  setAccountDescriptor(value?: AccountDescriptor): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class DeleteAccountRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getOwner(): string;
  setOwner(value: string): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class DeleteAccountReply extends jspb.Message {
  getOk(): boolean;
  setOk(value: boolean): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class GetAccountDescriptorRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
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

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class CreateBucketReply extends jspb.Message {
  hasBucketDescriptor(): boolean;
  clearBucketDescriptor(): void;
  getBucketDescriptor(): BucketDescriptor | undefined;
  setBucketDescriptor(value?: BucketDescriptor): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class DeleteBucketRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getName(): string;
  setName(value: string): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class DeleteBucketReply extends jspb.Message {
  getOk(): boolean;
  setOk(value: boolean): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
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

  getBucketName(): string;
  setBucketName(value: string): void;

  getKey(): string;
  setKey(value: string): void;

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
    bucketName: string,
    key: string,
  }
}

export class GetKeyReply extends jspb.Message {
  hasKeyValuePair(): boolean;
  clearKeyValuePair(): void;
  getKeyValuePair(): KeyValue | undefined;
  setKeyValuePair(value?: KeyValue): void;

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
  }
}

export class PutKeyRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getBucketName(): string;
  setBucketName(value: string): void;

  hasKeyValuePair(): boolean;
  clearKeyValuePair(): void;
  getKeyValuePair(): KeyValue | undefined;
  setKeyValuePair(value?: KeyValue): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    bucketName: string,
    keyValuePair?: KeyValue.AsObject,
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class PutKeyReply extends jspb.Message {
  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class DeleteKeyRequest extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): void;

  getBucketName(): string;
  setBucketName(value: string): void;

  getKey(): string;
  setKey(value: string): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    bucketName: string,
    key: string,
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
  }
}

export class DeleteKeyReply extends jspb.Message {
  getOk(): boolean;
  setOk(value: boolean): void;

  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): api_v1_database_transactions_pb.Transaction | undefined;
  setTransaction(value?: api_v1_database_transactions_pb.Transaction): void;

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
    transaction?: api_v1_database_transactions_pb.Transaction.AsObject,
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

