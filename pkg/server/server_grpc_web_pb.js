/**
 * @fileoverview gRPC-Web generated client stub for server
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var api_v1_database_kv_pb = require('../../api/v1/database/kv_pb.js')

var api_v1_database_transactions_pb = require('../../api/v1/database/transactions_pb.js')

var api_v1_raft_raft_shard_pb = require('../../api/v1/raft/raft_shard_pb.js')

var api_v1_raft_raft_host_pb = require('../../api/v1/raft/raft_host_pb.js')
const proto = {};
proto.server = require('./server_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.ShardManagerClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.ShardManagerPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.AddReplicaRequest,
 *   !proto.raft.AddReplicaReply>}
 */
const methodDescriptor_ShardManager_AddReplica = new grpc.web.MethodDescriptor(
  '/server.ShardManager/AddReplica',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.AddReplicaRequest,
  api_v1_raft_raft_shard_pb.AddReplicaReply,
  /**
   * @param {!proto.raft.AddReplicaRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.AddReplicaReply.deserializeBinary
);


/**
 * @param {!proto.raft.AddReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.AddReplicaReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.AddReplicaReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.addReplica =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/AddReplica',
      request,
      metadata || {},
      methodDescriptor_ShardManager_AddReplica,
      callback);
};


/**
 * @param {!proto.raft.AddReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.AddReplicaReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.addReplica =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/AddReplica',
      request,
      metadata || {},
      methodDescriptor_ShardManager_AddReplica);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.AddReplicaObserverRequest,
 *   !proto.raft.AddReplicaObserverReply>}
 */
const methodDescriptor_ShardManager_AddReplicaObserver = new grpc.web.MethodDescriptor(
  '/server.ShardManager/AddReplicaObserver',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.AddReplicaObserverRequest,
  api_v1_raft_raft_shard_pb.AddReplicaObserverReply,
  /**
   * @param {!proto.raft.AddReplicaObserverRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.AddReplicaObserverReply.deserializeBinary
);


/**
 * @param {!proto.raft.AddReplicaObserverRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.AddReplicaObserverReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.AddReplicaObserverReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.addReplicaObserver =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/AddReplicaObserver',
      request,
      metadata || {},
      methodDescriptor_ShardManager_AddReplicaObserver,
      callback);
};


/**
 * @param {!proto.raft.AddReplicaObserverRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.AddReplicaObserverReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.addReplicaObserver =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/AddReplicaObserver',
      request,
      metadata || {},
      methodDescriptor_ShardManager_AddReplicaObserver);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.AddReplicaWitnessRequest,
 *   !proto.raft.AddReplicaWitnessReply>}
 */
const methodDescriptor_ShardManager_AddReplicaWitness = new grpc.web.MethodDescriptor(
  '/server.ShardManager/AddReplicaWitness',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.AddReplicaWitnessRequest,
  api_v1_raft_raft_shard_pb.AddReplicaWitnessReply,
  /**
   * @param {!proto.raft.AddReplicaWitnessRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.AddReplicaWitnessReply.deserializeBinary
);


/**
 * @param {!proto.raft.AddReplicaWitnessRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.AddReplicaWitnessReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.AddReplicaWitnessReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.addReplicaWitness =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/AddReplicaWitness',
      request,
      metadata || {},
      methodDescriptor_ShardManager_AddReplicaWitness,
      callback);
};


/**
 * @param {!proto.raft.AddReplicaWitnessRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.AddReplicaWitnessReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.addReplicaWitness =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/AddReplicaWitness',
      request,
      metadata || {},
      methodDescriptor_ShardManager_AddReplicaWitness);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.GetLeaderIdRequest,
 *   !proto.raft.GetLeaderIdReply>}
 */
const methodDescriptor_ShardManager_GetLeaderId = new grpc.web.MethodDescriptor(
  '/server.ShardManager/GetLeaderId',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.GetLeaderIdRequest,
  api_v1_raft_raft_shard_pb.GetLeaderIdReply,
  /**
   * @param {!proto.raft.GetLeaderIdRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.GetLeaderIdReply.deserializeBinary
);


/**
 * @param {!proto.raft.GetLeaderIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.GetLeaderIdReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.GetLeaderIdReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.getLeaderId =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/GetLeaderId',
      request,
      metadata || {},
      methodDescriptor_ShardManager_GetLeaderId,
      callback);
};


/**
 * @param {!proto.raft.GetLeaderIdRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.GetLeaderIdReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.getLeaderId =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/GetLeaderId',
      request,
      metadata || {},
      methodDescriptor_ShardManager_GetLeaderId);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.GetShardMembersRequest,
 *   !proto.raft.GetShardMembersReply>}
 */
const methodDescriptor_ShardManager_GetShardMembers = new grpc.web.MethodDescriptor(
  '/server.ShardManager/GetShardMembers',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.GetShardMembersRequest,
  api_v1_raft_raft_shard_pb.GetShardMembersReply,
  /**
   * @param {!proto.raft.GetShardMembersRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.GetShardMembersReply.deserializeBinary
);


/**
 * @param {!proto.raft.GetShardMembersRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.GetShardMembersReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.GetShardMembersReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.getShardMembers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/GetShardMembers',
      request,
      metadata || {},
      methodDescriptor_ShardManager_GetShardMembers,
      callback);
};


/**
 * @param {!proto.raft.GetShardMembersRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.GetShardMembersReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.getShardMembers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/GetShardMembers',
      request,
      metadata || {},
      methodDescriptor_ShardManager_GetShardMembers);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.NewShardRequest,
 *   !proto.raft.NewShardReply>}
 */
const methodDescriptor_ShardManager_NewShard = new grpc.web.MethodDescriptor(
  '/server.ShardManager/NewShard',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.NewShardRequest,
  api_v1_raft_raft_shard_pb.NewShardReply,
  /**
   * @param {!proto.raft.NewShardRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.NewShardReply.deserializeBinary
);


/**
 * @param {!proto.raft.NewShardRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.NewShardReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.NewShardReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.newShard =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/NewShard',
      request,
      metadata || {},
      methodDescriptor_ShardManager_NewShard,
      callback);
};


/**
 * @param {!proto.raft.NewShardRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.NewShardReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.newShard =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/NewShard',
      request,
      metadata || {},
      methodDescriptor_ShardManager_NewShard);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.RemoveDataRequest,
 *   !proto.raft.RemoveDataReply>}
 */
const methodDescriptor_ShardManager_RemoveData = new grpc.web.MethodDescriptor(
  '/server.ShardManager/RemoveData',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.RemoveDataRequest,
  api_v1_raft_raft_shard_pb.RemoveDataReply,
  /**
   * @param {!proto.raft.RemoveDataRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.RemoveDataReply.deserializeBinary
);


/**
 * @param {!proto.raft.RemoveDataRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.RemoveDataReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.RemoveDataReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.removeData =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/RemoveData',
      request,
      metadata || {},
      methodDescriptor_ShardManager_RemoveData,
      callback);
};


/**
 * @param {!proto.raft.RemoveDataRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.RemoveDataReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.removeData =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/RemoveData',
      request,
      metadata || {},
      methodDescriptor_ShardManager_RemoveData);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.DeleteReplicaRequest,
 *   !proto.raft.DeleteReplicaReply>}
 */
const methodDescriptor_ShardManager_RemoveReplica = new grpc.web.MethodDescriptor(
  '/server.ShardManager/RemoveReplica',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.DeleteReplicaRequest,
  api_v1_raft_raft_shard_pb.DeleteReplicaReply,
  /**
   * @param {!proto.raft.DeleteReplicaRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.DeleteReplicaReply.deserializeBinary
);


/**
 * @param {!proto.raft.DeleteReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.DeleteReplicaReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.DeleteReplicaReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.removeReplica =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/RemoveReplica',
      request,
      metadata || {},
      methodDescriptor_ShardManager_RemoveReplica,
      callback);
};


/**
 * @param {!proto.raft.DeleteReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.DeleteReplicaReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.removeReplica =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/RemoveReplica',
      request,
      metadata || {},
      methodDescriptor_ShardManager_RemoveReplica);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.StartReplicaRequest,
 *   !proto.raft.StartReplicaReply>}
 */
const methodDescriptor_ShardManager_StartReplica = new grpc.web.MethodDescriptor(
  '/server.ShardManager/StartReplica',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.StartReplicaRequest,
  api_v1_raft_raft_shard_pb.StartReplicaReply,
  /**
   * @param {!proto.raft.StartReplicaRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.StartReplicaReply.deserializeBinary
);


/**
 * @param {!proto.raft.StartReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.StartReplicaReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.StartReplicaReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.startReplica =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/StartReplica',
      request,
      metadata || {},
      methodDescriptor_ShardManager_StartReplica,
      callback);
};


/**
 * @param {!proto.raft.StartReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.StartReplicaReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.startReplica =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/StartReplica',
      request,
      metadata || {},
      methodDescriptor_ShardManager_StartReplica);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.StartReplicaRequest,
 *   !proto.raft.StartReplicaReply>}
 */
const methodDescriptor_ShardManager_StartReplicaObserver = new grpc.web.MethodDescriptor(
  '/server.ShardManager/StartReplicaObserver',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.StartReplicaRequest,
  api_v1_raft_raft_shard_pb.StartReplicaReply,
  /**
   * @param {!proto.raft.StartReplicaRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.StartReplicaReply.deserializeBinary
);


/**
 * @param {!proto.raft.StartReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.StartReplicaReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.StartReplicaReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.startReplicaObserver =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/StartReplicaObserver',
      request,
      metadata || {},
      methodDescriptor_ShardManager_StartReplicaObserver,
      callback);
};


/**
 * @param {!proto.raft.StartReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.StartReplicaReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.startReplicaObserver =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/StartReplicaObserver',
      request,
      metadata || {},
      methodDescriptor_ShardManager_StartReplicaObserver);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.StopReplicaRequest,
 *   !proto.raft.StopReplicaReply>}
 */
const methodDescriptor_ShardManager_StopReplica = new grpc.web.MethodDescriptor(
  '/server.ShardManager/StopReplica',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_shard_pb.StopReplicaRequest,
  api_v1_raft_raft_shard_pb.StopReplicaReply,
  /**
   * @param {!proto.raft.StopReplicaRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_shard_pb.StopReplicaReply.deserializeBinary
);


/**
 * @param {!proto.raft.StopReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.StopReplicaReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.StopReplicaReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.ShardManagerClient.prototype.stopReplica =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.ShardManager/StopReplica',
      request,
      metadata || {},
      methodDescriptor_ShardManager_StopReplica,
      callback);
};


/**
 * @param {!proto.raft.StopReplicaRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.StopReplicaReply>}
 *     Promise that resolves to the response
 */
proto.server.ShardManagerPromiseClient.prototype.stopReplica =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.ShardManager/StopReplica',
      request,
      metadata || {},
      methodDescriptor_ShardManager_StopReplica);
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.RaftHostClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.RaftHostPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.CompactRequest,
 *   !proto.raft.CompactReply>}
 */
const methodDescriptor_RaftHost_Compact = new grpc.web.MethodDescriptor(
  '/server.RaftHost/Compact',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_host_pb.CompactRequest,
  api_v1_raft_raft_host_pb.CompactReply,
  /**
   * @param {!proto.raft.CompactRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_host_pb.CompactReply.deserializeBinary
);


/**
 * @param {!proto.raft.CompactRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.CompactReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.CompactReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.RaftHostClient.prototype.compact =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.RaftHost/Compact',
      request,
      metadata || {},
      methodDescriptor_RaftHost_Compact,
      callback);
};


/**
 * @param {!proto.raft.CompactRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.CompactReply>}
 *     Promise that resolves to the response
 */
proto.server.RaftHostPromiseClient.prototype.compact =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.RaftHost/Compact',
      request,
      metadata || {},
      methodDescriptor_RaftHost_Compact);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.GetHostConfigRequest,
 *   !proto.raft.GetHostConfigReply>}
 */
const methodDescriptor_RaftHost_GetHostConfig = new grpc.web.MethodDescriptor(
  '/server.RaftHost/GetHostConfig',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_host_pb.GetHostConfigRequest,
  api_v1_raft_raft_host_pb.GetHostConfigReply,
  /**
   * @param {!proto.raft.GetHostConfigRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_host_pb.GetHostConfigReply.deserializeBinary
);


/**
 * @param {!proto.raft.GetHostConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.GetHostConfigReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.GetHostConfigReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.RaftHostClient.prototype.getHostConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.RaftHost/GetHostConfig',
      request,
      metadata || {},
      methodDescriptor_RaftHost_GetHostConfig,
      callback);
};


/**
 * @param {!proto.raft.GetHostConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.GetHostConfigReply>}
 *     Promise that resolves to the response
 */
proto.server.RaftHostPromiseClient.prototype.getHostConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.RaftHost/GetHostConfig',
      request,
      metadata || {},
      methodDescriptor_RaftHost_GetHostConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.SnapshotRequest,
 *   !proto.raft.SnapshotReply>}
 */
const methodDescriptor_RaftHost_Snapshot = new grpc.web.MethodDescriptor(
  '/server.RaftHost/Snapshot',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_host_pb.SnapshotRequest,
  api_v1_raft_raft_host_pb.SnapshotReply,
  /**
   * @param {!proto.raft.SnapshotRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_host_pb.SnapshotReply.deserializeBinary
);


/**
 * @param {!proto.raft.SnapshotRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.SnapshotReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.SnapshotReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.RaftHostClient.prototype.snapshot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.RaftHost/Snapshot',
      request,
      metadata || {},
      methodDescriptor_RaftHost_Snapshot,
      callback);
};


/**
 * @param {!proto.raft.SnapshotRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.SnapshotReply>}
 *     Promise that resolves to the response
 */
proto.server.RaftHostPromiseClient.prototype.snapshot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.RaftHost/Snapshot',
      request,
      metadata || {},
      methodDescriptor_RaftHost_Snapshot);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.raft.StopRequest,
 *   !proto.raft.StopReply>}
 */
const methodDescriptor_RaftHost_Stop = new grpc.web.MethodDescriptor(
  '/server.RaftHost/Stop',
  grpc.web.MethodType.UNARY,
  api_v1_raft_raft_host_pb.StopRequest,
  api_v1_raft_raft_host_pb.StopReply,
  /**
   * @param {!proto.raft.StopRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_raft_raft_host_pb.StopReply.deserializeBinary
);


/**
 * @param {!proto.raft.StopRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.raft.StopReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.raft.StopReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.RaftHostClient.prototype.stop =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.RaftHost/Stop',
      request,
      metadata || {},
      methodDescriptor_RaftHost_Stop,
      callback);
};


/**
 * @param {!proto.raft.StopRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.raft.StopReply>}
 *     Promise that resolves to the response
 */
proto.server.RaftHostPromiseClient.prototype.stop =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.RaftHost/Stop',
      request,
      metadata || {},
      methodDescriptor_RaftHost_Stop);
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.TransactionsClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.TransactionsPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.NewTransactionRequest,
 *   !proto.database.NewTransactionReply>}
 */
const methodDescriptor_Transactions_NewTransaction = new grpc.web.MethodDescriptor(
  '/server.Transactions/NewTransaction',
  grpc.web.MethodType.UNARY,
  api_v1_database_transactions_pb.NewTransactionRequest,
  api_v1_database_transactions_pb.NewTransactionReply,
  /**
   * @param {!proto.database.NewTransactionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_transactions_pb.NewTransactionReply.deserializeBinary
);


/**
 * @param {!proto.database.NewTransactionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.NewTransactionReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.NewTransactionReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.TransactionsClient.prototype.newTransaction =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.Transactions/NewTransaction',
      request,
      metadata || {},
      methodDescriptor_Transactions_NewTransaction,
      callback);
};


/**
 * @param {!proto.database.NewTransactionRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.NewTransactionReply>}
 *     Promise that resolves to the response
 */
proto.server.TransactionsPromiseClient.prototype.newTransaction =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.Transactions/NewTransaction',
      request,
      metadata || {},
      methodDescriptor_Transactions_NewTransaction);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.CloseTransactionRequest,
 *   !proto.database.CloseTransactionReply>}
 */
const methodDescriptor_Transactions_CloseTransaction = new grpc.web.MethodDescriptor(
  '/server.Transactions/CloseTransaction',
  grpc.web.MethodType.UNARY,
  api_v1_database_transactions_pb.CloseTransactionRequest,
  api_v1_database_transactions_pb.CloseTransactionReply,
  /**
   * @param {!proto.database.CloseTransactionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_transactions_pb.CloseTransactionReply.deserializeBinary
);


/**
 * @param {!proto.database.CloseTransactionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.CloseTransactionReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.CloseTransactionReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.TransactionsClient.prototype.closeTransaction =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.Transactions/CloseTransaction',
      request,
      metadata || {},
      methodDescriptor_Transactions_CloseTransaction,
      callback);
};


/**
 * @param {!proto.database.CloseTransactionRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.CloseTransactionReply>}
 *     Promise that resolves to the response
 */
proto.server.TransactionsPromiseClient.prototype.closeTransaction =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.Transactions/CloseTransaction',
      request,
      metadata || {},
      methodDescriptor_Transactions_CloseTransaction);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.CommitRequest,
 *   !proto.database.CommitReply>}
 */
const methodDescriptor_Transactions_Commit = new grpc.web.MethodDescriptor(
  '/server.Transactions/Commit',
  grpc.web.MethodType.UNARY,
  api_v1_database_transactions_pb.CommitRequest,
  api_v1_database_transactions_pb.CommitReply,
  /**
   * @param {!proto.database.CommitRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_transactions_pb.CommitReply.deserializeBinary
);


/**
 * @param {!proto.database.CommitRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.CommitReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.CommitReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.TransactionsClient.prototype.commit =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.Transactions/Commit',
      request,
      metadata || {},
      methodDescriptor_Transactions_Commit,
      callback);
};


/**
 * @param {!proto.database.CommitRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.CommitReply>}
 *     Promise that resolves to the response
 */
proto.server.TransactionsPromiseClient.prototype.commit =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.Transactions/Commit',
      request,
      metadata || {},
      methodDescriptor_Transactions_Commit);
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.KVStoreServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.KVStoreServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.CreateAccountRequest,
 *   !proto.database.CreateAccountReply>}
 */
const methodDescriptor_KVStoreService_CreateAccount = new grpc.web.MethodDescriptor(
  '/server.KVStoreService/CreateAccount',
  grpc.web.MethodType.UNARY,
  api_v1_database_kv_pb.CreateAccountRequest,
  api_v1_database_kv_pb.CreateAccountReply,
  /**
   * @param {!proto.database.CreateAccountRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_kv_pb.CreateAccountReply.deserializeBinary
);


/**
 * @param {!proto.database.CreateAccountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.CreateAccountReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.CreateAccountReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.KVStoreServiceClient.prototype.createAccount =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.KVStoreService/CreateAccount',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_CreateAccount,
      callback);
};


/**
 * @param {!proto.database.CreateAccountRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.CreateAccountReply>}
 *     Promise that resolves to the response
 */
proto.server.KVStoreServicePromiseClient.prototype.createAccount =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.KVStoreService/CreateAccount',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_CreateAccount);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.DeleteAccountRequest,
 *   !proto.database.DeleteAccountReply>}
 */
const methodDescriptor_KVStoreService_DeleteAccount = new grpc.web.MethodDescriptor(
  '/server.KVStoreService/DeleteAccount',
  grpc.web.MethodType.UNARY,
  api_v1_database_kv_pb.DeleteAccountRequest,
  api_v1_database_kv_pb.DeleteAccountReply,
  /**
   * @param {!proto.database.DeleteAccountRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_kv_pb.DeleteAccountReply.deserializeBinary
);


/**
 * @param {!proto.database.DeleteAccountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.DeleteAccountReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.DeleteAccountReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.KVStoreServiceClient.prototype.deleteAccount =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.KVStoreService/DeleteAccount',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_DeleteAccount,
      callback);
};


/**
 * @param {!proto.database.DeleteAccountRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.DeleteAccountReply>}
 *     Promise that resolves to the response
 */
proto.server.KVStoreServicePromiseClient.prototype.deleteAccount =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.KVStoreService/DeleteAccount',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_DeleteAccount);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.CreateBucketRequest,
 *   !proto.database.CreateBucketReply>}
 */
const methodDescriptor_KVStoreService_CreateBucket = new grpc.web.MethodDescriptor(
  '/server.KVStoreService/CreateBucket',
  grpc.web.MethodType.UNARY,
  api_v1_database_kv_pb.CreateBucketRequest,
  api_v1_database_kv_pb.CreateBucketReply,
  /**
   * @param {!proto.database.CreateBucketRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_kv_pb.CreateBucketReply.deserializeBinary
);


/**
 * @param {!proto.database.CreateBucketRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.CreateBucketReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.CreateBucketReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.KVStoreServiceClient.prototype.createBucket =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.KVStoreService/CreateBucket',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_CreateBucket,
      callback);
};


/**
 * @param {!proto.database.CreateBucketRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.CreateBucketReply>}
 *     Promise that resolves to the response
 */
proto.server.KVStoreServicePromiseClient.prototype.createBucket =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.KVStoreService/CreateBucket',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_CreateBucket);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.DeleteBucketRequest,
 *   !proto.database.DeleteBucketReply>}
 */
const methodDescriptor_KVStoreService_DeleteBucket = new grpc.web.MethodDescriptor(
  '/server.KVStoreService/DeleteBucket',
  grpc.web.MethodType.UNARY,
  api_v1_database_kv_pb.DeleteBucketRequest,
  api_v1_database_kv_pb.DeleteBucketReply,
  /**
   * @param {!proto.database.DeleteBucketRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_kv_pb.DeleteBucketReply.deserializeBinary
);


/**
 * @param {!proto.database.DeleteBucketRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.DeleteBucketReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.DeleteBucketReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.KVStoreServiceClient.prototype.deleteBucket =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.KVStoreService/DeleteBucket',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_DeleteBucket,
      callback);
};


/**
 * @param {!proto.database.DeleteBucketRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.DeleteBucketReply>}
 *     Promise that resolves to the response
 */
proto.server.KVStoreServicePromiseClient.prototype.deleteBucket =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.KVStoreService/DeleteBucket',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_DeleteBucket);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.GetKeyRequest,
 *   !proto.database.GetKeyReply>}
 */
const methodDescriptor_KVStoreService_GetKey = new grpc.web.MethodDescriptor(
  '/server.KVStoreService/GetKey',
  grpc.web.MethodType.UNARY,
  api_v1_database_kv_pb.GetKeyRequest,
  api_v1_database_kv_pb.GetKeyReply,
  /**
   * @param {!proto.database.GetKeyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_kv_pb.GetKeyReply.deserializeBinary
);


/**
 * @param {!proto.database.GetKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.GetKeyReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.GetKeyReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.KVStoreServiceClient.prototype.getKey =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.KVStoreService/GetKey',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_GetKey,
      callback);
};


/**
 * @param {!proto.database.GetKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.GetKeyReply>}
 *     Promise that resolves to the response
 */
proto.server.KVStoreServicePromiseClient.prototype.getKey =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.KVStoreService/GetKey',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_GetKey);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.PutKeyRequest,
 *   !proto.database.PutKeyReply>}
 */
const methodDescriptor_KVStoreService_PutKey = new grpc.web.MethodDescriptor(
  '/server.KVStoreService/PutKey',
  grpc.web.MethodType.UNARY,
  api_v1_database_kv_pb.PutKeyRequest,
  api_v1_database_kv_pb.PutKeyReply,
  /**
   * @param {!proto.database.PutKeyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_kv_pb.PutKeyReply.deserializeBinary
);


/**
 * @param {!proto.database.PutKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.PutKeyReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.PutKeyReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.KVStoreServiceClient.prototype.putKey =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.KVStoreService/PutKey',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_PutKey,
      callback);
};


/**
 * @param {!proto.database.PutKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.PutKeyReply>}
 *     Promise that resolves to the response
 */
proto.server.KVStoreServicePromiseClient.prototype.putKey =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.KVStoreService/PutKey',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_PutKey);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.database.DeleteKeyRequest,
 *   !proto.database.DeleteKeyReply>}
 */
const methodDescriptor_KVStoreService_DeleteKey = new grpc.web.MethodDescriptor(
  '/server.KVStoreService/DeleteKey',
  grpc.web.MethodType.UNARY,
  api_v1_database_kv_pb.DeleteKeyRequest,
  api_v1_database_kv_pb.DeleteKeyReply,
  /**
   * @param {!proto.database.DeleteKeyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_v1_database_kv_pb.DeleteKeyReply.deserializeBinary
);


/**
 * @param {!proto.database.DeleteKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.database.DeleteKeyReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.database.DeleteKeyReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.KVStoreServiceClient.prototype.deleteKey =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.KVStoreService/DeleteKey',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_DeleteKey,
      callback);
};


/**
 * @param {!proto.database.DeleteKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.database.DeleteKeyReply>}
 *     Promise that resolves to the response
 */
proto.server.KVStoreServicePromiseClient.prototype.deleteKey =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.KVStoreService/DeleteKey',
      request,
      metadata || {},
      methodDescriptor_KVStoreService_DeleteKey);
};


module.exports = proto.server;

