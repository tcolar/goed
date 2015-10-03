# Goed bash init

export GOED_INSTANCE=$1
export GOED_VIEW=$2

# Search goed custom/builtin tools first 
export PATH=~/.goed/actions/:~/.goed/default/actions/:$PATH

function goed_open() {
# open a file in goed (make relative first)
	path="$1" 
	if [ -d "$(dirname "$1")" ]; then
		path="$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
	fi
	goed_api open $GOED_INSTANCE $path
}

function goed_cd() {
# cd into a directory, and notify goed of the new dir
	builtin cd $@
	goed_api view_cwd $GOED_INSTANCE $GOED_VIEW "`pwd`"
} 

export EDITOR="goed_open" # When within goed, $EDITOR is goed

alias cd="goed_cd"

alias o="goed_open" # open a file/dir
alias s="s.sh" # search text (=~ grep) 
alias f="f.sh" # search files (=~ find)
