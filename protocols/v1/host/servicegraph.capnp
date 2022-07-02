using Go = import "/go.capnp";
@0xb1ade260fb3cadf2;
$Go.package("host");
$Go.import("r3t.io/pkg/protocols/v1/host");

using import "config.capnp".ServiceType;
interface Service @0xed78136d1400ca3e extends(Node) {
    setNodeId @0 (nodeId :Int64) -> ();
    getServiceName @1 () -> (name :Text);
    getServiceType @2 () -> (type :ServiceType);
    getDependencies @3 () -> (dependencies :List(Service));
    prepareToRun @4 () -> (error :Text);
    isRunning @5 () -> (running :Bool);
    start @6 (retry :Bool) -> (error :Text);
    stop @7 (retry :Bool, force :Bool) -> (error :Text);
}

interface Node @0xe322c44e33fc17b4 {
    id @0 () -> (id :Int64);
}

interface Edge @0x9d51abb4c18add06 {
    from @0 () -> (node :Node);
    to @1 () -> (node :Node);
    reversedEdge @2 () -> (edge :Edge);
}

interface Graph @0xc20098da5460c109 {
    node @0 (id :Int64) -> (node :Node);
    nodes @1 () -> (nodes :List(Node));
    from @2 (id :Int64) -> (nodes :List(Node));
    hasEdgeBetween @3 (xId :Int64, yId :Int64) -> (connected :Bool);
    edge @4 (uId :Int64, vId :Int64) -> (edge :Edge);
}

interface DirectedGraph @0xca335475fc06b787 extends(Graph) {
    hasEdgeFromTo @0 (uId :Int64, vId :Int64) -> (connected :Bool);
    to @1 (id :Int64) -> (nodes :List(Node));
}

interface ServiceLibrary @0x8930b538011a6d48 {
	addService @0 (svc :Service) -> (error :Text);
	addServices @1 (svcs :List(Service)) -> (error :Text);
	getService @2 (svc :Service) -> (svc: Service, error :Text);
	getServices @3 (svcs :List(Service)) -> (svcs :List(Service), error :Text);
	startService @4 (svc :Service, retry :Bool) -> (error :Text);
	stopService @5 (retry :Bool, force :Bool, svc :Service) -> (error :Text);
	stopServices @6 (retry :Bool, force :Bool, svcs :List(Service)) -> (error :Text);
	getServiceStatus @7 (svc :ServiceStatus) -> (error :Text);
	getServiceStatuses @8 ( svcs :List(ServiceStatus)) -> (error :Text);
}

interface ServiceStatus @0xfe5ad4396565b592 {
    enum SvcState {
        running @0;
        stopped @1;
    }
    status @0 () -> (state :SvcState);
}