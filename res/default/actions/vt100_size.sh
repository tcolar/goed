#!/bin/bash
set -e #fail early

# Set tty size to match goed view size (stty rows, stty cols)

view=$GOED_VIEW
inst=$GOED_INSTANCE

rows=`goed_api view_rows $inst $view`
cols=`goed_api view_cols $inst $view`

set -x
stty rows $rows
stty cols $cols