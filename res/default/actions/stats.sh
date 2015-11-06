#!/bin/bash
echo -- FREE --
free
echo -- DF --
df -x tmpfs
echo -- TOP PROCS --
ps aux | sort -rk 3,3 | head -n 6