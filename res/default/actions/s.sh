#!/bin/bash

# Shortcut to grep files / dirs

if [ "$#" -lt 1 ]; then
    echo "Syntax: s <patern> [path]"
    exit 1
fi

path="."

if [ "$#" -gt 1 ]; then
	path="$2"
fi

grep -rni "$1" $path 