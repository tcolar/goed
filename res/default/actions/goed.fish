# Goed (fi)sh(y) init

set -x GOED_INSTANCE $args[0]
set -x GOED_VIEW $args[1]

# Search goed custom/builtin tools first 
set -x PATH $HOME/.goed/actions/ $HOME/.goed/default/actions/ $PATH

function goed_cd
	# cd into a directory, and notify goed of the new dir
	builtin cd $argv
	goed --api view_set_work_dir $GOED_INSTANCE $GOED_VIEW "(pwd)"
end 

function o
	goed --api open $GOED_INSTANCE "(pwd)" $argv[0]
end 

set -x EDITOR "goed --api edit $GOED_INSTANCE (pwd)" # edit a file

alias cd="goed_cd"

alias s="search_text.sh"
alias f="find_files.sh"
alias sz="vt100_size.sh"

./vt100_size.sh
