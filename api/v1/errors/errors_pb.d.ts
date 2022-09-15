// package: apiv1
// file: api/v1/errors/errors.proto

import * as jspb from "google-protobuf";
import * as api_v1_errors_error_codes_pb from "../../../api/v1/errors/error_codes_pb";

export class Error extends jspb.Message {
  getCode(): api_v1_errors_error_codes_pb.CodeMap[keyof api_v1_errors_error_codes_pb.CodeMap];
  setCode(value: api_v1_errors_error_codes_pb.CodeMap[keyof api_v1_errors_error_codes_pb.CodeMap]): void;

  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Error.AsObject;
  static toObject(includeInstance: boolean, msg: Error): Error.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Error, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Error;
  static deserializeBinaryFromReader(message: Error, reader: jspb.BinaryReader): Error;
}

export namespace Error {
  export type AsObject = {
    code: api_v1_errors_error_codes_pb.CodeMap[keyof api_v1_errors_error_codes_pb.CodeMap],
    message: string,
  }
}

