#!/bin/bash

set -eux

mkdir -p "$HOME"/.ssh

echo "$DEPLOY_PRIVATE_KEY" > "$HOME"/.ssh/deploy
export GIT_SSH_COMMAND='ssh -i ~/.ssh/deploy'

git remote add mirror git@github.com:anthropos-labs/api.git
git push --mirror mirror
