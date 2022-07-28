// package: database
// file: api/v1/database/node_host_config.proto

import * as jspb from "google-protobuf";

export class HasNodeInfoRequest extends jspb.Message {
  getClusterid(): number;
  setClusterid(value: number): void;

  getNodeid(): number;
  setNodeid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HasNodeInfoRequest.AsObject;
  static toObject(includeInstance: boolean, msg: HasNodeInfoRequest): HasNodeInfoRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HasNodeInfoRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HasNodeInfoRequest;
  static deserializeBinaryFromReader(message: HasNodeInfoRequest, reader: jspb.BinaryReader): HasNodeInfoRequest;
}

export namespace HasNodeInfoRequest {
  export type AsObject = {
    clusterid: number,
    nodeid: number,
  }
}

export class HasNodeInfoResponse extends jspb.Message {
  getHasnodeinfo(): boolean;
  setHasnodeinfo(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HasNodeInfoResponse.AsObject;
  static toObject(includeInstance: boolean, msg: HasNodeInfoResponse): HasNodeInfoResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HasNodeInfoResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HasNodeInfoResponse;
  static deserializeBinaryFromReader(message: HasNodeInfoResponse, reader: jspb.BinaryReader): HasNodeInfoResponse;
}

export namespace HasNodeInfoResponse {
  export type AsObject = {
    hasnodeinfo: boolean,
  }
}

export class GetNodeHostInfoRequest extends jspb.Message {
  hasOption(): boolean;
  clearOption(): void;
  getOption(): NodeHostInfoOption | undefined;
  setOption(value?: NodeHostInfoOption): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNodeHostInfoRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetNodeHostInfoRequest): GetNodeHostInfoRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetNodeHostInfoRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNodeHostInfoRequest;
  static deserializeBinaryFromReader(message: GetNodeHostInfoRequest, reader: jspb.BinaryReader): GetNodeHostInfoRequest;
}

export namespace GetNodeHostInfoRequest {
  export type AsObject = {
    option?: NodeHostInfoOption.AsObject,
  }
}

export class GetNodeHostInfoResponse extends jspb.Message {
  hasInfo(): boolean;
  clearInfo(): void;
  getInfo(): NodeHostInfo | undefined;
  setInfo(value?: NodeHostInfo): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNodeHostInfoResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetNodeHostInfoResponse): GetNodeHostInfoResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetNodeHostInfoResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNodeHostInfoResponse;
  static deserializeBinaryFromReader(message: GetNodeHostInfoResponse, reader: jspb.BinaryReader): GetNodeHostInfoResponse;
}

export namespace GetNodeHostInfoResponse {
  export type AsObject = {
    info?: NodeHostInfo.AsObject,
  }
}

export class GetNodeHostConfigRequest extends jspb.Message {
  hasClusterid(): boolean;
  clearClusterid(): void;
  getClusterid(): number;
  setClusterid(value: number): void;

  hasNodeid(): boolean;
  clearNodeid(): void;
  getNodeid(): number;
  setNodeid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNodeHostConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetNodeHostConfigRequest): GetNodeHostConfigRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetNodeHostConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNodeHostConfigRequest;
  static deserializeBinaryFromReader(message: GetNodeHostConfigRequest, reader: jspb.BinaryReader): GetNodeHostConfigRequest;
}

export namespace GetNodeHostConfigRequest {
  export type AsObject = {
    clusterid: number,
    nodeid: number,
  }
}

export class GetNodeHostConfigResponse extends jspb.Message {
  hasNodehostconfig(): boolean;
  clearNodehostconfig(): void;
  getNodehostconfig(): NodeHostConfig | undefined;
  setNodehostconfig(value?: NodeHostConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNodeHostConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetNodeHostConfigResponse): GetNodeHostConfigResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetNodeHostConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNodeHostConfigResponse;
  static deserializeBinaryFromReader(message: GetNodeHostConfigResponse, reader: jspb.BinaryReader): GetNodeHostConfigResponse;
}

export namespace GetNodeHostConfigResponse {
  export type AsObject = {
    nodehostconfig?: NodeHostConfig.AsObject,
  }
}

export class NodeHostConfig extends jspb.Message {
  getDeploymentid(): number;
  setDeploymentid(value: number): void;

  getWaldir(): string;
  setWaldir(value: string): void;

  getNodehostdir(): string;
  setNodehostdir(value: string): void;

  getRoundtriptimemillisecond(): number;
  setRoundtriptimemillisecond(value: number): void;

  getRaftaddress(): string;
  setRaftaddress(value: string): void;

  getAddressbynodehostid(): boolean;
  setAddressbynodehostid(value: boolean): void;

  getListenaddress(): string;
  setListenaddress(value: string): void;

  getMutualtls(): boolean;
  setMutualtls(value: boolean): void;

  getCafile(): string;
  setCafile(value: string): void;

  getCertfile(): string;
  setCertfile(value: string): void;

  getKeyfile(): string;
  setKeyfile(value: string): void;

  getEnablemetrics(): boolean;
  setEnablemetrics(value: boolean): void;

  getMaxsendqueuesize(): number;
  setMaxsendqueuesize(value: number): void;

  getMaxreceivequeuesize(): number;
  setMaxreceivequeuesize(value: number): void;

  getMaxsnapshotsendbytespersecond(): number;
  setMaxsnapshotsendbytespersecond(value: number): void;

  getMaxsnapshotrecvbytespersecond(): number;
  setMaxsnapshotrecvbytespersecond(value: number): void;

  getNotifycommit(): boolean;
  setNotifycommit(value: boolean): void;

  hasGossipconfig(): boolean;
  clearGossipconfig(): void;
  getGossipconfig(): GossipConfig | undefined;
  setGossipconfig(value?: GossipConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NodeHostConfig.AsObject;
  static toObject(includeInstance: boolean, msg: NodeHostConfig): NodeHostConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NodeHostConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NodeHostConfig;
  static deserializeBinaryFromReader(message: NodeHostConfig, reader: jspb.BinaryReader): NodeHostConfig;
}

export namespace NodeHostConfig {
  export type AsObject = {
    deploymentid: number,
    waldir: string,
    nodehostdir: string,
    roundtriptimemillisecond: number,
    raftaddress: string,
    addressbynodehostid: boolean,
    listenaddress: string,
    mutualtls: boolean,
    cafile: string,
    certfile: string,
    keyfile: string,
    enablemetrics: boolean,
    maxsendqueuesize: number,
    maxreceivequeuesize: number,
    maxsnapshotsendbytespersecond: number,
    maxsnapshotrecvbytespersecond: number,
    notifycommit: boolean,
    gossipconfig?: GossipConfig.AsObject,
  }
}

export class GossipConfig extends jspb.Message {
  getBindaddress(): string;
  setBindaddress(value: string): void;

  getAdvertiseaddress(): string;
  setAdvertiseaddress(value: string): void;

  clearSeedList(): void;
  getSeedList(): Array<string>;
  setSeedList(value: Array<string>): void;
  addSeed(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GossipConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GossipConfig): GossipConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GossipConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GossipConfig;
  static deserializeBinaryFromReader(message: GossipConfig, reader: jspb.BinaryReader): GossipConfig;
}

export namespace GossipConfig {
  export type AsObject = {
    bindaddress: string,
    advertiseaddress: string,
    seedList: Array<string>,
  }
}

export class NodeHostInfo extends jspb.Message {
  getNodehostid(): string;
  setNodehostid(value: string): void;

  getRaftaddress(): string;
  setRaftaddress(value: string): void;

  hasGossip(): boolean;
  clearGossip(): void;
  getGossip(): GossipInfo | undefined;
  setGossip(value?: GossipInfo): void;

  clearClusterinfoList(): void;
  getClusterinfoList(): Array<ClusterInfo>;
  setClusterinfoList(value: Array<ClusterInfo>): void;
  addClusterinfo(value?: ClusterInfo, index?: number): ClusterInfo;

  clearLoginfoList(): void;
  getLoginfoList(): Array<NodeInfo>;
  setLoginfoList(value: Array<NodeInfo>): void;
  addLoginfo(value?: NodeInfo, index?: number): NodeInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NodeHostInfo.AsObject;
  static toObject(includeInstance: boolean, msg: NodeHostInfo): NodeHostInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NodeHostInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NodeHostInfo;
  static deserializeBinaryFromReader(message: NodeHostInfo, reader: jspb.BinaryReader): NodeHostInfo;
}

export namespace NodeHostInfo {
  export type AsObject = {
    nodehostid: string,
    raftaddress: string,
    gossip?: GossipInfo.AsObject,
    clusterinfoList: Array<ClusterInfo.AsObject>,
    loginfoList: Array<NodeInfo.AsObject>,
  }
}

export class ClusterInfo extends jspb.Message {
  getClusterid(): number;
  setClusterid(value: number): void;

  getNodeid(): number;
  setNodeid(value: number): void;

  getNodesMap(): jspb.Map<number, string>;
  clearNodesMap(): void;
  getConfigchangeindex(): number;
  setConfigchangeindex(value: number): void;

  getStatemachinetype(): number;
  setStatemachinetype(value: number): void;

  getIsleader(): boolean;
  setIsleader(value: boolean): void;

  getIsobserver(): boolean;
  setIsobserver(value: boolean): void;

  getIswitness(): boolean;
  setIswitness(value: boolean): void;

  getPending(): boolean;
  setPending(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClusterInfo.AsObject;
  static toObject(includeInstance: boolean, msg: ClusterInfo): ClusterInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ClusterInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClusterInfo;
  static deserializeBinaryFromReader(message: ClusterInfo, reader: jspb.BinaryReader): ClusterInfo;
}

export namespace ClusterInfo {
  export type AsObject = {
    clusterid: number,
    nodeid: number,
    nodesMap: Array<[number, string]>,
    configchangeindex: number,
    statemachinetype: number,
    isleader: boolean,
    isobserver: boolean,
    iswitness: boolean,
    pending: boolean,
  }
}

export class GossipInfo extends jspb.Message {
  getEnabled(): boolean;
  setEnabled(value: boolean): void;

  getAdvertiseaddress(): string;
  setAdvertiseaddress(value: string): void;

  getNumofknownnodehosts(): number;
  setNumofknownnodehosts(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GossipInfo.AsObject;
  static toObject(includeInstance: boolean, msg: GossipInfo): GossipInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GossipInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GossipInfo;
  static deserializeBinaryFromReader(message: GossipInfo, reader: jspb.BinaryReader): GossipInfo;
}

export namespace GossipInfo {
  export type AsObject = {
    enabled: boolean,
    advertiseaddress: string,
    numofknownnodehosts: number,
  }
}

export class NodeInfo extends jspb.Message {
  getClusterid(): number;
  setClusterid(value: number): void;

  getNodeid(): number;
  setNodeid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NodeInfo.AsObject;
  static toObject(includeInstance: boolean, msg: NodeInfo): NodeInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NodeInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NodeInfo;
  static deserializeBinaryFromReader(message: NodeInfo, reader: jspb.BinaryReader): NodeInfo;
}

export namespace NodeInfo {
  export type AsObject = {
    clusterid: number,
    nodeid: number,
  }
}

export class NodeHostInfoOption extends jspb.Message {
  getSkiploginfo(): boolean;
  setSkiploginfo(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NodeHostInfoOption.AsObject;
  static toObject(includeInstance: boolean, msg: NodeHostInfoOption): NodeHostInfoOption.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NodeHostInfoOption, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NodeHostInfoOption;
  static deserializeBinaryFromReader(message: NodeHostInfoOption, reader: jspb.BinaryReader): NodeHostInfoOption;
}

export namespace NodeHostInfoOption {
  export type AsObject = {
    skiploginfo: boolean,
  }
}

