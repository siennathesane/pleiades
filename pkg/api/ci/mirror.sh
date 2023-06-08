#!/bin/bash

#
# Copyright (c) 2023 Sienna Lloyd
#
# Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License here:
#  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
#

set -eux

mkdir -p "$HOME"/.ssh

echo "$DEPLOY_PRIVATE_KEY" > "$HOME"/.ssh/deploy
chmod 0600 "$HOME"/.ssh/deploy

ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

export GIT_SSH_COMMAND='ssh -i ~/.ssh/deploy'
git remote add mirror git@github.com:mxplusb/api.git
git push --mirror mirror
