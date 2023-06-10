# Pleiades

Pleiades is a globally distributed, [constellation mesh database](#what-is-a-constellation-mesh) designed to enable complex and systemically critical workloads.

## Current State

Pleiades has been under active development since May 2022, and has undergone many design iterations and test implementations. These iterations are critical to finding and designing the right internal architectural patterns and implementations to meet systemically critical workload requirements.

Pleiades is currently pre-alpha and not yet ready for operational deployment. Check the project tab for the latest details.

# Vision

_Pleiades' globally distributed runtime, low-latency operations, exabyte-scale data storage, and straightforward multi-modal interface allows complex, systemically critical data workloads to be run safely in highly-regulated environments._

As the software industry moves forward, there is a massive hole that large systemic enterprises, financial institutions, government agencies, large scientific or educational institutions, and other massive, complex organizations have: data management at scale. Web2 and its related technologies are powerful and reliable, but they're oftentimes not up to par for the needs of systemically critical systems. Extensions exist and abound for many solutions, but ultimately, even pyramids get replaced by mausoleums. Web3 is capable, but many of the technologies and directions Web3 is taking will never serve the real world in meaningful ways.

Pleiades is not better, worse, or ultimately comparable to a lot of systems which exist, and this is by design. While the goal is to have Pleiades be a semi-easy drop-in migration from some existing solutions, it is nothing like most existing solutions. Pleiades makes trade-offs between a simple user experience and a complex internal architecture. While Pleiades must be easy to use for end-users (operators are also end users), it can't sacrifice end-user simplicity for functional capability. See the [constellation mesh](#what-is-a-constellation-mesh) section for more details.

At its core, Pleiades is informed by solutions like TiKV, CockroachDB, MongoDB, RocksDB, PostgreSQL, Trino, Redis, Azure CosmosDB, Google Spanner, Neo4J, Ethereum, LibP2P, and others. Each of these different solutions provides different bits and pieces of research and design insight that ultimately informs the types of things that makes Pleiades, well, Pleiades. For example, Pleiades should have change streaming, similar to MongoDB, but it also needs to scale like Spanner while having the same performance characteristics of TiKV and enabling self-scaling like CockroachDB.

Part of this vision also means keeping it accessible to all organizations. You can use Pleiades for anything you want, so long as you're not selling it. Contributions from consumers aren't mandatory, but a rising tide raises all ships.

## What is a Constellation Mesh?

Because Pleiades is like nothing else which exists, by design, it's hard to describe how it is different. Initially, Pleiades was compared to a globally distributed data fabric, but that was missing a core aspect of Pleiades: autonomy. We tried to compare it to a data mesh, but there's no architectural alignment with the specific business domains expected with a data mesh. After many different iterations, and comparing Pleiades to many different system models, @mxplusb landed on _constellation mesh_.

So what is a constellation mesh? A mesh, [as defined by Oracle](https://www.oracle.com/integration/what-is-data-mesh/), is a distributed architecture for data management. That fits into the distributed architecture model for Pleiades without defining domain alignment but doesn't define the type of mesh. In systems engineering, the term "constellation" is commonly used to refer to autonomous satellites which work in coordination to provide different aspects of a distributed, unified, and autonomous data set. Pleiades is a distributed, autonomous system working in independent coordination focusing on data management as it's primary feature (re: fancy distributed database).

So what makes Pleiades a distributed, autonomous system working in coordination? Pleiades is designed from the ground up to be able to handle an exabyte's worth of data while only having a single operator. This means internal architectures require a mixture of distribution and autonomy whenever possible, giving visibility to the operator, but also not requiring hands-on operational management. Configurations, workload scheduling, and many other internal operations in the constellation must happen independently, and the scale requires a leader-less, decentralized design. This also increases
the complexity, but only if modeled incorrectly.

One of the most useful models for understanding the systemic impacts of automated decision-making in a decentralized
network is a force-directed graph. In force-directed graphs, attraction is generally modeled with $F_s = kx$ and
repulsion is generally modeled with $|F| = \frac{1}{4\pi\epsilon_0}\frac{|q_1q_2|}{r^2}$, and iterative simulations
demonstrate mechanical equilibrium can be achieved across the entire graph. To simplify, if you have a very large
spider web, any force applied to one part of the web will proportionally affect the others. Automated
decision-making in a distributed system can have very similar sets of characteristics, but instead of a _mechanical_
equilibrium being simulated, _virtual change equilibrium_ would be achieved through network propagation. Virtual
Change Equilibrium (VCE) is the result of a constellation runtime event (CRE, re: a change) being successfully propagated
throughout the entire constellation.

Pleiades' internal clustering model is modeled with graph connectivity instead of a centralized membership. Each node in the constellation is only aware of its neighbours, some top-level metadata, and how to handle CREs. A CRE is _a neighbour-only broadcast_. An example would be when a node is shutting down, it will broadcast the CRE of `leave` to its local neighbours who will repeat the same message, so on and so forth, until the constellation has reached virtual change equilibrium. Other CREs are things like `join`, `query`, and `update`, all of which follow the same propagation model to achieve VCE.

The constellation model is enabled through [SWIM](https://www.cs.cornell.edu/projects/Quicksilver/public_pdfs/SWIM.pdf), and Pleiades currently uses Hashicorp's implementation with their [Lifeguard extensions](https://arxiv.org/abs/1707.00788). While SWIM allows for the constellation's members to be modeled concretely, Pleiades also leverages network tomography to handle VCE. Pleiades currently uses the Hashicorp network tomography library, [`coordinate`](https://github.com/hashicorp/serf/tree/master/coordinate), to provide real-time computed network coordinates for each node in the constellation. SWIM enables constellation membership, `coordinate` provides locality, and the internal lifecycle state machines together make Pleiades an autonomous distributed system.

Using a workload adjustment event, it can be easier to understand the implementation nuances. When a node containing the leader of a range replica receives an internal scaling event (re: scale up or scale down), it will send a CRE with some change metadata, and the constellation will quickly achieve VCE asychronously. From the triggering node's perspective, VCE happens once the broadcast call returns. After triggering the CRE, the node will broadcast a `query` CRE asking for the nearest neighbours with available capacity to create a new range replica. Once the neighbours have been identified, the node will communicate directly with the closest identified neighbour to instantiate a range replica. The identified neighbour will broadcast the new range replica CRE as part of the initial VCE, and the original node will start the relevant scaling workflow. One the workflow has finished, the original node which triggered the CRE will broadcast a final CRE with the final change metadata.

While there are many things not covered in that example, the internal autonomy of the constellation allows for things like scaling events to be handled by any node, while keeping the entire constellation aware of the change. This is what makes Pleiades a _constellation mesh database_. Hopefully that helps! Please feel free to open a discussion if you have questions.

# Goals

Here are a set of guiding goals that Pleiades strives to meet. They may expand or shrink over the years, but overall, these are core aspects of what makes Pleiades unique.

* 10ms write latency regardless of scale
* Can handle an exabyte's worth of data the same as a megabyte
* There's no functional difference between one node or a thousand
* Everything is observable
* A single, unified global cluster can be operated by a single systems engineer

# Use Cases

Every systemically critical workload is different, but they all generally consist of requirements around data governance, high throughput reads and writes, low latency operations, multi-regional scalability, highly regulated environments, and complex BCDR expectations. If you have this kind of use
case, please let leave (as much of) the details (as you can) in an issue, it's important to understand what the needs of the community are.

# Documentation

There will be some documentation at some point! This will always be a work in progress. If you need info, create a discussion and ping @mxplusb; she'll make sure you get what you need.

# Contributing

Check out [CONTRIBUTING.md](./CONTRIBUTING.md) for more information.