using Go = import "/go.capnp";
@0x90858bf9ae63c319;
$Go.package("host");
$Go.import("r3t.io/pkg/protocols/v1/host");

using Config = import "config.capnp";
interface Negotiator @0xe35a52b4e5c60a15 {
    configService @0 () -> (svc :Config.ConfigService);
}
