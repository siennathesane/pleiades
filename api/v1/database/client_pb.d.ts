// package: database
// file: api/v1/database/client.proto

import * as jspb from "google-protobuf";
import * as api_v1_database_kv_pb from "../../../api/v1/database/kv_pb";

export class Session extends jspb.Message {
  getClusterId(): number;
  setClusterId(value: number): void;

  getClientId(): number;
  setClientId(value: number): void;

  getSessionId(): number;
  setSessionId(value: number): void;

  getRespondedTo(): number;
  setRespondedTo(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Session.AsObject;
  static toObject(includeInstance: boolean, msg: Session): Session.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Session, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Session;
  static deserializeBinaryFromReader(message: Session, reader: jspb.BinaryReader): Session;
}

export namespace Session {
  export type AsObject = {
    clusterId: number,
    clientId: number,
    sessionId: number,
    respondedTo: number,
  }
}

export class NewSessionRequest extends jspb.Message {
  getClusterId(): number;
  setClusterId(value: number): void;

  getClientId(): number;
  setClientId(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewSessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: NewSessionRequest): NewSessionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NewSessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewSessionRequest;
  static deserializeBinaryFromReader(message: NewSessionRequest, reader: jspb.BinaryReader): NewSessionRequest;
}

export namespace NewSessionRequest {
  export type AsObject = {
    clusterId: number,
    clientId: number,
  }
}

export class NewSessionResponse extends jspb.Message {
  getSessionId(): number;
  setSessionId(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewSessionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: NewSessionResponse): NewSessionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NewSessionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewSessionResponse;
  static deserializeBinaryFromReader(message: NewSessionResponse, reader: jspb.BinaryReader): NewSessionResponse;
}

export namespace NewSessionResponse {
  export type AsObject = {
    sessionId: number,
  }
}

export class ProposeRequest extends jspb.Message {
  hasSession(): boolean;
  clearSession(): void;
  getSession(): Session | undefined;
  setSession(value?: Session): void;

  hasCommand(): boolean;
  clearCommand(): void;
  getCommand(): api_v1_database_kv_pb.KeyValue | undefined;
  setCommand(value?: api_v1_database_kv_pb.KeyValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProposeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProposeRequest): ProposeRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProposeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProposeRequest;
  static deserializeBinaryFromReader(message: ProposeRequest, reader: jspb.BinaryReader): ProposeRequest;
}

export namespace ProposeRequest {
  export type AsObject = {
    session?: Session.AsObject,
    command?: api_v1_database_kv_pb.KeyValue.AsObject,
  }
}

export class ProposeResponse extends jspb.Message {
  getCommandId(): number;
  setCommandId(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProposeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ProposeResponse): ProposeResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProposeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProposeResponse;
  static deserializeBinaryFromReader(message: ProposeResponse, reader: jspb.BinaryReader): ProposeResponse;
}

export namespace ProposeResponse {
  export type AsObject = {
    commandId: number,
  }
}

