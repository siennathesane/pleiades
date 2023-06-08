# Pleiades

Pleiades is a globally distributed database designed to enable complex and systemically critical workloads.

## Current State

Pleiades has been under active development since May 2022, and has undergone many design iterations and test implementations. These iterations are critical to finding and designing the right internal architectural patterns and implementations to meet systemically critical workload requirements.

While Pleiades is mostly functional, it's very buggy and considered to be pre-alpha. Some things work, some things don't, it's a work in progress towards the types stability required for systemically critical workloads.

# Vision

_Pleiades' globally distributed runtime, low-latency operations, exabyte-scale data storage, and straightforward multi-modal interface allows complex, systemically critical data workloads to be run safely in highly-regulated environments._

As the software industry moves forward, there is a massive hole that large systemic enterprises, financial institutions, government agencies, large scientific or educational institutions, and other massive, complex organizations have: data management at scale. Web2 and it's related technologies are powerful and reliable, but they are oftentimes not up to par for the needs of systemically critical systems. Extensions exist and abound for many solutions, but ultimately, even pyramids get replaced by mausoleums. Web3 is capable, but many of the technologies and directions of Web3 will never serve the real world in meaningful ways.

Pleiades is not better, worse, or ultimately comparable to a lot of systems which exist, and this is by design. While the goal is to have Pleiades be a semi-easy drop-in migration from some existing solutions, it is nothing like most existing solutions. Pleiades makes trade-offs between a simple user experience and a complex internal architecture. While Pleiades must be easy to use for end-users (operators are also end users), it can't sacrifice end-user simplicity for functional capability.

At it's core, Pleiades is informed by solutions like TiKV, CockroachDB, MongoDB, RocksDB, PostgreSQL, Trino, Redis, Azure CosmosDB, Google Spanner, Neo4J, Ethereum, LibP2P, and others. Each of these different solutions provides different bits and pieces of research and design insight that ultimately informs the types of things that makes Pleiades, well, Pleiades. For example, Pleiades should have change streaming, similar to MongoDB, but it also needs to scale like Spanner while having the same performance characteristics of TiKV and enabling self-scaling like CockroachDB.

Part of this vision also means keeping it accessible to all organizations. You can use Pleiades for anything you want, so long as you're not selling it. Contributions from consumers aren't mandatory, but a rising tide raises all ships.

# Goals

Here are a set of guiding goals that Pleiades strives to meet. They may expand or shrink over the years, but overall, these are core aspects of what makes Pleiades unique.

* 10ms write latency regardless of scale
* Can handle an exabyte's worth of data the same as a megabyte
* There's no functional difference between one node or a thousand
* Everything is observable

# Use Cases

Every systemically critical workload is different, but they all generally consist of requirements around data governance, high throughput reads and writes, low latency operations, multi-regional scalability, highly regulated environments, and complex BCDR expectations. If you have this kind of use case, please let leave (as much of) the details (as you can) in an issue, it's important to understand what the needs of the community are.

# Documentation

There will be some documentation at some point!

# Contributing

Check out [CONTRIBUTING.md](./CONTRIBUTING.md) for more information.