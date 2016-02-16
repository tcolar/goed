#!/bin/bash

# Search text pattern in files / dirs (grep)

if [ "$#" -lt 1 ]; then
    echo "Syntax: s <patern> [path]"
    exit 1
fi

path="."

if [ "$#" -gt 1 ]; then
	path="$2"
fi

grep -rni --color "$1" $path 