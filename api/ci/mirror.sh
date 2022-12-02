#!/bin/bash

mkdir -p "$HOME"/.ssh

echo $DEPLOY_PRIVATE_KEY > $HOME/.ssh/deploy
echo $DEPLOY_PUBLIC_KEY > $HOME/.ssh/deploy.pub

cat $HOME/.ssh/config <<EOF
Host github.com
  HostName github.com
  User git
  IdentityFile $HOME/.ssh/deploy
  IdentitiesOnly yes
EOF

git remote add mirror git@github.com:anthropos-labs/pleiades.git
git push --mirror mirror
