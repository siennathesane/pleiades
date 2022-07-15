using Go = import "/go.capnp";
@0x882cc5e81e24c654;
$Go.package("database");
$Go.import("r3t.io/pkg/protocols/v1/database");

interface Client @0xc212c5427766d750 {
    newSession @0 (clusterId :UInt64) -> (session :Session);
    closeSession @1 (session :Session) -> ();

    getClusterId @2 () -> (clusterId :UInt64);
    getClientId @3 () -> (clientId :UInt64);

    propose @4 (session :Session, cmd :Data, timeout :Int64, synchronous :Bool) -> (response :RequestResult);
    read @5 (query :Data) -> (payload :Data);
}

struct Session @0xbb0748b8b81c5da8 $Go.doc("Session is a client-facing database session") {
    clusterId @0 :UInt64;
    clientId @1 :UInt64;
    seriesId @2 :UInt64;
    respondedTo @3 :UInt64;
}

interface RequestState @0xb3cd040205e8ec8c {
    completed @0 () -> (iterator :RequestResultIterator);
}

interface RequestResult @0x83bb5f78b82ac5dc {
    completed @0 () -> (done :Bool);
    getResult @1 () -> (result :UInt64);
    rejected @2 () -> (rejected :Bool);
    terminated @3 () -> (terminated :Bool);
    timeout @4 () -> (timedout :Bool);
}

interface RequestResultIterator @0xba9a01628ce6fb28 {
    get @0 () -> (value :RequestResult);
    next @1 () -> (more :Bool);
}
