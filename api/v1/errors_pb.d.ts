// package: database
// file: api/v1/errors.proto

import * as jspb from "google-protobuf";

export class DBError extends jspb.Message {
  getType(): DBErrorTypeMap[keyof DBErrorTypeMap];
  setType(value: DBErrorTypeMap[keyof DBErrorTypeMap]): void;

  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DBError.AsObject;
  static toObject(includeInstance: boolean, msg: DBError): DBError.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DBError, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DBError;
  static deserializeBinaryFromReader(message: DBError, reader: jspb.BinaryReader): DBError;
}

export namespace DBError {
  export type AsObject = {
    type: DBErrorTypeMap[keyof DBErrorTypeMap],
    message: string,
  }
}

export interface DBErrorTypeMap {
  SESSION: 0;
  KEY_VALUE: 1;
  RAFT_CONTROL: 2;
  RAFT_CLUSTER: 3;
}

export const DBErrorType: DBErrorTypeMap;

