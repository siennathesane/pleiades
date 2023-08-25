## Core Technical Goal
The core technical goals inform the internal bits of logic

- 10ms write latency regardless of scale
- Can handle an exabyte's worth of data the same as a megabyte
- There's no functional difference between one node or a thousand
- Everything is observable
- A single, unified global cluster can be operated by a single systems engineer

## Core Use Case Goals
These core use cases are designed to inform what types of workloads Pleiades must support.

- Realtime financial reconciliation
- ECS workloads for large-scale simulations
- DistSQL-like interface, PostgreSQL-compatible
- Key-value interface
