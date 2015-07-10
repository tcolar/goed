#!/bin/bash
set -ex #fail early

view=$GOED_VIEW
inst=$GOED_INSTANCE

# get file location
src=`goed_api view_src_loc $inst $view` 
# goimports on file
goimports -w $src
# reload view buffer from file
goed_api view_reload $inst $view 
