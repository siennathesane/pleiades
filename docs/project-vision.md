---
title: Project Vision
---

*Pleiades' globally distributed runtime, low-latency operations, exabyte-scale data storage, and straightforward multi-modal interface allows complex, systemically critical data workloads to be safely run in real-time within highly-regulated environments.*

As the software industry moves forward, there is a massive hole that large systemic enterprises, financial institutions, government agencies, large scientific or educational institutions, and other massive, complex organizations have: data management at scale. Web2 and its related technologies are powerful and reliable, but they're oftentimes not up to par for the needs of systemically critical systems. Extensions exist and abound for many solutions, but ultimately, even pyramids get replaced by mausoleums. Web3 is capable, but many of the technologies and directions Web3 is taking will never serve the real world in meaningful ways.

Pleiades is not better, worse, or ultimately comparable to a lot of systems which exist, and this is by design. While the goal is to have Pleiades be a semi-easy drop-in migration from some existing solutions, it is nothing like most existing solutions. Pleiades makes trade-offs between a simple user experience and a complex internal architecture. While Pleiades must be easy to use for end-users (operators are also end users), it can't sacrifice end-user simplicity for functional capability. See the [constellation mesh]() section for more details.

At its core, Pleiades is informed by solutions like TiKV, CockroachDB, MongoDB, RocksDB, PostgreSQL, Trino, Redis, Azure CosmosDB, Google Spanner, Neo4J, Ethereum, LibP2P, and others. Each of these different solutions provides different bits and pieces of research and design insight that ultimately informs the types of things that makes Pleiades, well, Pleiades. For example, Pleiades should have change streaming, similar to MongoDB, but it also needs to scale like Spanner while having the same performance characteristics of TiKV and enabling self-scaling like CockroachDB.

Part of this vision also means keeping it accessible to all organizations. You can use Pleiades for anything you want, so long as you're not selling it. Contributions from consumers aren't mandatory, but a rising tide raises all ships.