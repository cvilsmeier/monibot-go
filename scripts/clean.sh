#!/bin/sh

if [ ! -f go.mod ]
then
    echo "go.mod not found, must be in src directory"
    exit 1
fi

rm -rf _dist
