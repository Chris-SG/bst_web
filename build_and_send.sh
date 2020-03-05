#!/bin/bash

env GOOS=linux GOARCH=amd64 go build

ssh -t bst@35.196.119.150 "sudo systemctl stop bst"
scp bst_web bst@35.196.119.150:/home/bst/bst_web
scp -r templates bst@35.196.119.150:/home/bst
scp -r dist bst@35.196.119.150:/home/bst
ssh -t bst@35.196.119.150 "sudo systemctl start bst"
