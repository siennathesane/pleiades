---
title: Contributing
authors:
  - Sienna Lloyd <sienna@linux.com>
tags:
  - contributing
  - virtual-machine
  - multipass
---

# Contributing!

## Server-side Development

Pleiades is only formally supported on Ubuntu LTS, starting with 22.04. As part of that, all of our development tools are based around that and [Visual Studio Code](https://code.visualstudio.com/). Several of the bits of code require Linux with the `PREEMPT_RT` patch for the runtime server due to the hard real-time requirements. To make it easier to develop, we use [`multipass`](https://multipass.run/) from Canonical. You'll also need a [free Ubuntu Pro](https://ubuntu.com/pro) subscription to get the `PREEMPT_RT` kernel.

Once you've installed `multipass`, it's recommended to create an alias of `mp` to make it easier to work with.

## Setup

```sh
# make sure you're in the root of the repo
# create the vm. this uses 2 cpu, 8gib of ram, and 20gb of disk space.
# adjust these values as needed
mp launch -c 2 -m 8G -d 20G --cloud-init build/cloud-config.yaml -n primary --mount $PWD:/home/ubuntu/pleiades

# copy your git stuff to make it easier
mp transfer -r ~/.ssh primary:/home/ubuntu/
mp transfer -r ~/.git-credentials primary:/home/ubuntu/
mp transfer -r ~/.gitconfig primary:/home/ubuntu/
```

Once the VM has been provisioned, it will say something along the lines of `Launched: primary` and a few mounting notes. While you wait, [log into your Ubuntu Pro account](https://ubuntu.com/pro/dashboard) and grab your token. Once the machine has been provisioned, you can access it and finish the installation:

```sh
# get into the vm
mp shell primary

# (optional) install the recommended vscode extensions
code --install-extension "GitHub.copilot"
code --install-extension "minherz.copyright-inserter"
code --install-extension "ms-vscode-remote.vscode-remote-extensionpack"
code --install-extension "ms-vsliveshare.vsliveshare"
code --install-extension "redhat.vscode-commons"
code --install-extension "redhat.vscode-xml"
code --install-extension "redhat.vscode-yaml"
code --install-extension "remcohaszing.schemastore"
code --install-extension "rust-lang.rust-analyzer"
code --install-extension "timonwong.shellcheck"

# attach your pro subscription
sudo pro attach <token>

# install the rt kernel
sudo pro enable realtime-kernel

# reboot
sudo reboot now
```

> [!warning]
> Even if you have an Intel processor, do not install the `intel-iotg` variant of the RT kernel as Pleiades is currently supporting ARM platforms.

At this point, the core VM is set up. Hooking it up to vscode is pretty simple at this point:

1. Get the IP of the VM
	1. You can use `mp info primary --format json | jq -r '.info.primary.ipv4[0]'` to make it simpler
2. Add a new SSH host to vscode
3. Connect to the SSH host
4. Select the `pleiades` folder from the `ubuntu` user's home directory.

At this point you are all set up!

> [!info]
> When you created the VM, it mounts the Pleiades code base at `/home/ubuntu/pleiades`, and the instructions help you transfer your SSH key to make git work. If you follow the instructions, then you'll be able to work in the VM exclusively while still maintaining the source code locally.
> 
> This might seem strange to remote into a VM to work on code that's local to your laptop, but it's to overcome the `PREEMPT_RT` needs for the server-side code, otherwise a devcontainer would work.
> 
> This is intended to be a quality of life thing, but do what feels right for you.

### OpenPGP Keys

If you're like me and you use OpenPGP keys for git security, you'll also want to add those. I use Keybase to manage my keys just for ease of use. You'll want to add those to your devbox as well.

```sh
# export the keys
keybase pgp export > key.asc
keybase pgp export -s > priv-key.asc

# transfer the keys to the devbox
mp transfer key.asc primary:/home/ubuntu
mp transfer priv-key.asc primary:/home/ubuntu

# import the keys in the devbox
gpg --import key.asc
gpg --import priv-key.asc

```
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
