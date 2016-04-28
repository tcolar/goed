# Goed bash init

export GOED_INSTANCE=$1
export GOED_VIEW=$2

# Search goed custom/builtin tools first 
export PATH=$HOME/.goed/actions/:$HOME/.goed/default/actions/:$PATH

function goed_cd() {
	# cd into a directory, and notify goed of the new dir
	builtin cd $@
	goed_api view_set_work_dir $GOED_INSTANCE $GOED_VIEW "`pwd`"
} 

function o() {
	goed_api open $GOED_INSTANCE "`pwd`" $@
} 

export EDITOR="goed_api edit $GOED_INSTANCE `pwd`" # edit a file

alias cd="goed_cd"

alias s="search_text.sh"
alias f="find_files.sh"
alias sz="vt100_size.sh"

