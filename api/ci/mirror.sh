#!/bin/bash

set -eux

mkdir -p "$HOME"/.ssh

echo "$DEPLOY_PRIVATE_KEY" > "$HOME"/.ssh/deploy
chmod 0600 "$HOME"/.ssh/deploy

ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

export GIT_SSH_COMMAND='ssh -i ~/.ssh/deploy'
git remote add mirror git@github.com:anthropos-labs/api.git
git push --mirror mirror
