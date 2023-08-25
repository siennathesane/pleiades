---
title: Overview
authors:
  - Sienna Lloyd <sienna@linux.com>
tags:
  - architecture
  - aspect
  - component
  - service
---

Pleiades is grouped into a few different classifications:

* Components
* Aspects
* Services

Each of these classifications provides different bits of functionality. Components are self-contained units of functionality that can be reused across different parts of Pleiades. They can be thought of as building blocks that can be combined with other components to create larger, more complex systems. Aspects, on the other hand, are cross-cutting features or functionality that affect multiple components or modules.

Of the three classifications, things are either runtime-centric or library-centric. Runtime-centric pieces are focused on managing the state of a Pleiades node (or larger constellation), but library-centric pieces only provide reusable functionality. For the most part, Pleiades aims to keep the runtime code fairly light with a focal point on event-driven wrappers of library functionality.

# Components

- HLC
- Raft Engine
- ZeroMQ
- RocksDB

# Aspects

* Storage Engine
* Netcode & RPC framework
* Messaging substrate

# Services

- Gossip
- kvstore
- Raft
- Messaging
- System