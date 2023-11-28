#!/usr/bin/env bash
#

for x in {1..100}; 
do 

    packer-sdc -v
    echo "Run ${x}"
    go generate ./...
    git diff --exit-code > /dev/null
    if [[ $? -ne 0  ]]
    then
        echo "bad generate on ${x}"
        git status 
        exit 1
    fi
    echo "=========="
done
