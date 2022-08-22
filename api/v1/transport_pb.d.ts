// package: v1
// file: api/v1/transport.proto

import * as jspb from "google-protobuf";

export class State extends jspb.Message {
  getState(): number;
  setState(value: number): void;

  getHeadertofollow(): number;
  setHeadertofollow(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): State.AsObject;
  static toObject(includeInstance: boolean, msg: State): State.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: State, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): State;
  static deserializeBinaryFromReader(message: State, reader: jspb.BinaryReader): State;
}

export namespace State {
  export type AsObject = {
    state: number,
    headertofollow: number,
  }
}

export class Header extends jspb.Message {
  getSize(): number;
  setSize(value: number): void;

  getChecksum(): number;
  setChecksum(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Header.AsObject;
  static toObject(includeInstance: boolean, msg: Header): Header.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Header, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Header;
  static deserializeBinaryFromReader(message: Header, reader: jspb.BinaryReader): Header;
}

export namespace Header {
  export type AsObject = {
    size: number,
    checksum: number,
  }
}

export class HeaderTest extends jspb.Message {
  getSize(): number;
  setSize(value: number): void;

  getChecksum(): number;
  setChecksum(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HeaderTest.AsObject;
  static toObject(includeInstance: boolean, msg: HeaderTest): HeaderTest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HeaderTest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HeaderTest;
  static deserializeBinaryFromReader(message: HeaderTest, reader: jspb.BinaryReader): HeaderTest;
}

export namespace HeaderTest {
  export type AsObject = {
    size: number,
    checksum: number,
  }
}

