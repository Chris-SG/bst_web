#!/bin/bash

sha=$(git ls-remote git://github.com/chris-sg/bst_server_models.git HEAD | awk '{ print $1}')
go get github.com/chris-sg/bst_server_models@$sha