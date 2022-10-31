# Pleiades Internal Message Bus (IMB)

## Design

Currently, the Internal Message Bus (IMB) uses a in-memory implementation of [NATS](https://nats.io). NATS is completely overpowered for Pleiades, but it's also exactly what Pleiades needs in the short-term. It currently runs in-memory with socket-only connections to the Pleiades process, so various
processes in Pleiades can open and close clients as needed without having to worry about security concerns.

## Subjects

Pleiades splits up the various messaging needs into separate Subjects. For the most part, the subjects are pretty straightforward: a hierarchal ordering of granularity.

There is a top-level queue called `SYSTEM` which coalesces all system Subjects into a unified event stream so it can be subscribed to later. Internal consumers can also subscribe to specific Subjects instead of the system Subject.

| Subject                    | Purpose                                                         | Owner    | Type    |
|----------------------------|-----------------------------------------------------------------|----------|---------|
| `SYSTEM`                   | The root system Subject                                         | _SYSTEM_ | Queue   |
| `system.raftv1`            | Top-level subject for Raft messages                             | Raft     | Subject |
| `system.raftv1.connection` | Connection alerts for Raft                                      | Raft     | Subject |
| `system.raftv1.host`       | Host alerts for Raft                                            | Raft     | Subject |
| `system.raftv1.log`        | Log notifications for Raft                                      | Raft     | Subject |
| `system.raftv1.node`       | Node events for Raft                                            | Raft     | Subject |
| `system.raftv1.raft`       | Raft events                                                     | Raft     | Subject |
| `system.raftv1.snapshot`   | Snapshot events for Raft                                        | Raft     | Subject |
| `system.raftv1.<123>`      | Shard-specific events, where `<123>` is a specific shard number | Raft     | Subject |

## Ownership

The ownership of the stream and subject namespaces are handled by the [Embedded Messaging](../pkg/messaging/embedded_messaging.go) component. It handles the creation of the necessary Subjects and Queues.

## Logging

The IMB is completely asynchronous, which can make it hard to develop, hard to debug, and just generally difficult to work with. To save yourself, your teammates, and pretty much everyone a difficult time, _LOG ALL ERRORS_. This will help troubleshoot everything from tests to live errors.