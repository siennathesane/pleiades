using Go = import "/go.capnp";
@0x9f0a9459eb121bf0;
$Go.package("database");
$Go.import("r3t.io/pkg/protocols/v1/database");

struct KeyValue @0x9fc77743d79c134f $Go.doc("KeyValue is a key-value pair used for the database") {
    key @0 :Data;
    value @1 :Data;
    createRevision @2 :Int64;
    modifyRevision @3 :Int64;
    version @4 :Int64;
    lease @5 :Int64;
}

struct Event @0xb7f920001018ddbb {
    enum EventType {
        put @0;
        delete @1;
    }

    type @0 :EventType;
    keyValue @1 :KeyValue;
    previousKeyValue @2 :KeyValue;
}
