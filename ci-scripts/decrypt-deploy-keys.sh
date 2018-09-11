#!/bin/bash
set -x # display output of commands

# decrypt our deployment ssh key
openssl aes-256-cbc -K $encrypted_c76c477e8187_key -iv $encrypted_c76c477e8187_iv -in poem.computer-deploy-key.enc -out poem.computer-deploy-key -d

# clean up (we no longer need the encrypted key)
rm poem.computer-deploy-key.enc

# set the ssh key as the default on this ci server instance
chmod 600 poem.computer-deploy-key
mv poem.computer-deploy-key ~/.ssh/id_rsa
