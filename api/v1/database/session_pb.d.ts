// package: database
// file: api/v1/database/session.proto

import * as jspb from "google-protobuf";

export class SessionPayload extends jspb.Message {
  hasNewsessionrequest(): boolean;
  clearNewsessionrequest(): void;
  getNewsessionrequest(): NewSessionRequest | undefined;
  setNewsessionrequest(value?: NewSessionRequest): void;

  hasNewsessionresponse(): boolean;
  clearNewsessionresponse(): void;
  getNewsessionresponse(): NewSessionResponse | undefined;
  setNewsessionresponse(value?: NewSessionResponse): void;

  hasClosesessionrequest(): boolean;
  clearClosesessionrequest(): void;
  getClosesessionrequest(): CloseSessionRequest | undefined;
  setClosesessionrequest(value?: CloseSessionRequest): void;

  hasClosesessionresponse(): boolean;
  clearClosesessionresponse(): void;
  getClosesessionresponse(): CloseSessionResponse | undefined;
  setClosesessionresponse(value?: CloseSessionResponse): void;

  getMethod(): SessionPayload.MethodNameMap[keyof SessionPayload.MethodNameMap];
  setMethod(value: SessionPayload.MethodNameMap[keyof SessionPayload.MethodNameMap]): void;

  getTypeCase(): SessionPayload.TypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SessionPayload.AsObject;
  static toObject(includeInstance: boolean, msg: SessionPayload): SessionPayload.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SessionPayload, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SessionPayload;
  static deserializeBinaryFromReader(message: SessionPayload, reader: jspb.BinaryReader): SessionPayload;
}

export namespace SessionPayload {
  export type AsObject = {
    newsessionrequest?: NewSessionRequest.AsObject,
    newsessionresponse?: NewSessionResponse.AsObject,
    closesessionrequest?: CloseSessionRequest.AsObject,
    closesessionresponse?: CloseSessionResponse.AsObject,
    method: SessionPayload.MethodNameMap[keyof SessionPayload.MethodNameMap],
  }

  export interface MethodNameMap {
    NEW_SESSION: 0;
    CLOSE_SESSION: 1;
  }

  export const MethodName: MethodNameMap;

  export enum TypeCase {
    TYPE_NOT_SET = 0,
    NEWSESSIONREQUEST = 1,
    NEWSESSIONRESPONSE = 2,
    CLOSESESSIONREQUEST = 3,
    CLOSESESSIONRESPONSE = 4,
  }
}

export class CloseSessionRequest extends jspb.Message {
  hasSession(): boolean;
  clearSession(): void;
  getSession(): Session | undefined;
  setSession(value?: Session): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseSessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CloseSessionRequest): CloseSessionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CloseSessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseSessionRequest;
  static deserializeBinaryFromReader(message: CloseSessionRequest, reader: jspb.BinaryReader): CloseSessionRequest;
}

export namespace CloseSessionRequest {
  export type AsObject = {
    session?: Session.AsObject,
    timeout: number,
  }
}

export class CloseSessionResponse extends jspb.Message {
  hasSession(): boolean;
  clearSession(): void;
  getSession(): Session | undefined;
  setSession(value?: Session): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseSessionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CloseSessionResponse): CloseSessionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CloseSessionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseSessionResponse;
  static deserializeBinaryFromReader(message: CloseSessionResponse, reader: jspb.BinaryReader): CloseSessionResponse;
}

export namespace CloseSessionResponse {
  export type AsObject = {
    session?: Session.AsObject,
    timeout: number,
  }
}

export class ProposeSessionRequest extends jspb.Message {
  hasSession(): boolean;
  clearSession(): void;
  getSession(): Session | undefined;
  setSession(value?: Session): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProposeSessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProposeSessionRequest): ProposeSessionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProposeSessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProposeSessionRequest;
  static deserializeBinaryFromReader(message: ProposeSessionRequest, reader: jspb.BinaryReader): ProposeSessionRequest;
}

export namespace ProposeSessionRequest {
  export type AsObject = {
    session?: Session.AsObject,
    timeout: number,
  }
}

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

