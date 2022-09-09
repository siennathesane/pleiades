// package: database
// file: api/v1/database/transactions.proto

import * as jspb from "google-protobuf";

export class CloseTransactionRequest extends jspb.Message {
  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): Transaction | undefined;
  setTransaction(value?: Transaction): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseTransactionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CloseTransactionRequest): CloseTransactionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CloseTransactionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseTransactionRequest;
  static deserializeBinaryFromReader(message: CloseTransactionRequest, reader: jspb.BinaryReader): CloseTransactionRequest;
}

export namespace CloseTransactionRequest {
  export type AsObject = {
    transaction?: Transaction.AsObject,
    timeout: number,
  }
}

export class CloseTransactionReply extends jspb.Message {
  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): Transaction | undefined;
  setTransaction(value?: Transaction): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseTransactionReply.AsObject;
  static toObject(includeInstance: boolean, msg: CloseTransactionReply): CloseTransactionReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CloseTransactionReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseTransactionReply;
  static deserializeBinaryFromReader(message: CloseTransactionReply, reader: jspb.BinaryReader): CloseTransactionReply;
}

export namespace CloseTransactionReply {
  export type AsObject = {
    transaction?: Transaction.AsObject,
    timeout: number,
  }
}

export class CommitRequest extends jspb.Message {
  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): Transaction | undefined;
  setTransaction(value?: Transaction): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CommitRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CommitRequest): CommitRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CommitRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CommitRequest;
  static deserializeBinaryFromReader(message: CommitRequest, reader: jspb.BinaryReader): CommitRequest;
}

export namespace CommitRequest {
  export type AsObject = {
    transaction?: Transaction.AsObject,
  }
}

export class CommitReply extends jspb.Message {
  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): Transaction | undefined;
  setTransaction(value?: Transaction): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CommitReply.AsObject;
  static toObject(includeInstance: boolean, msg: CommitReply): CommitReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CommitReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CommitReply;
  static deserializeBinaryFromReader(message: CommitReply, reader: jspb.BinaryReader): CommitReply;
}

export namespace CommitReply {
  export type AsObject = {
    transaction?: Transaction.AsObject,
  }
}

export class Transaction extends jspb.Message {
  getShardId(): number;
  setShardId(value: number): void;

  getClientId(): number;
  setClientId(value: number): void;

  getTransactionId(): number;
  setTransactionId(value: number): void;

  getRespondedTo(): number;
  setRespondedTo(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Transaction.AsObject;
  static toObject(includeInstance: boolean, msg: Transaction): Transaction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Transaction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Transaction;
  static deserializeBinaryFromReader(message: Transaction, reader: jspb.BinaryReader): Transaction;
}

export namespace Transaction {
  export type AsObject = {
    shardId: number,
    clientId: number,
    transactionId: number,
    respondedTo: number,
  }
}

export class NewTransactionRequest extends jspb.Message {
  getShardId(): number;
  setShardId(value: number): void;

  getClientId(): number;
  setClientId(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewTransactionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: NewTransactionRequest): NewTransactionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NewTransactionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewTransactionRequest;
  static deserializeBinaryFromReader(message: NewTransactionRequest, reader: jspb.BinaryReader): NewTransactionRequest;
}

export namespace NewTransactionRequest {
  export type AsObject = {
    shardId: number,
    clientId: number,
  }
}

export class NewTransactionReply extends jspb.Message {
  hasTransaction(): boolean;
  clearTransaction(): void;
  getTransaction(): Transaction | undefined;
  setTransaction(value?: Transaction): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewTransactionReply.AsObject;
  static toObject(includeInstance: boolean, msg: NewTransactionReply): NewTransactionReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NewTransactionReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewTransactionReply;
  static deserializeBinaryFromReader(message: NewTransactionReply, reader: jspb.BinaryReader): NewTransactionReply;
}

export namespace NewTransactionReply {
  export type AsObject = {
    transaction?: Transaction.AsObject,
  }
}

