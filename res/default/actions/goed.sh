# Goed bash init

# Search goed custom/builtin tools first 
export PATH=$HOME/.goed/actions/:$HOME/.goed/default/actions/:$PATH

function goed_open() {
# open a file in goed
	goed_api open $GOED_INSTANCE "`pwd`" "$1"
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

