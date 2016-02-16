# Goed bash init

export GOED_INSTANCE=$1
export GOED_VIEW=$2

# Search goed custom/builtin tools first 
export PATH=$HOME/.goed/actions/:$HOME/.goed/default/actions/:$PATH

function goed_cd() {
# cd into a directory, and notify goed of the new dir
	builtin cd $@
	goed_api view_cwd $GOED_INSTANCE $GOED_VIEW "`pwd`"
} 

#gapi=`which goed_api`
#echo $gapi
export EDITOR="goed_api edit $GOED_INSTANCE `pwd`" # open a file/dir

alias cd="goed_cd"

alias o="goed_api open $GOED_INSTANCE `pwd`" # open a file/dir
alias s="search_text.sh"
alias f="find_files.sh"
alias sz="vt100_size.sh"

