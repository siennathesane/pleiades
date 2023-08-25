---
title: Contributing
authors:
  - Sienna Lloyd <sienna@linux.com>
tags:
  - contributing
---

# Contributing!

## Trunk-based Development

Trunk-based development is a version control management practice where developers merge small, frequent updates to a core “trunk” or main branch. It’s a common practice among teams and part of the SRE lifecycle since it streamlines merging and integration phases. [Trunk-based development is a required practice of true CI/CD](https://www.atlassian.com/continuous-delivery/continuous-integration/trunk-based-development) [1](https://www.atlassian.com/continuous-delivery/continuous-integration/trunk-based-development).

Rather than relying on feature branches, Trunk Based Development has each developer work locally and independently on their project, and then merge their changes back into the main branch (the trunk) at least once a day. [Merges must occur whether or not feature changes or additions are complete](about:blank#) [2](https://www.gitkraken.com/blog/trunk-based-development).

## Everything Else

* Open an issue with a proposed change if it's larger than a bugfix
* Be very careful about memory allocations
* Use linear commits whenever possible
* Fix-forward, no backports
* Document your code
* Ask questions if you get stuck!
