#!/bin/zsh

# fib
fibonacci() {
	local a c
	local -F1 b
	a=0 ; b=1
	print $a
	repeat 100
	do
		print "${b%.*}"
		c=$a
		a=$b
		((b = c + b))
	done
}


if [[ $# -ne 0 ]]; then
	echo "Too many args"
	exit 1
fi

fibonacci

echo "done"

