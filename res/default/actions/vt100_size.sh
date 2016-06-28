#!/bin/bash
set -e #fail early

# Set tty size to match goed view size (stty rows, stty cols)

view=$GOED_VIEW
inst=$GOED_INSTANCE

rows=`goed --api view_rows $inst $view`
cols=`goed --api view_cols $inst $view`

goed --api view_set_vt_cols $inst $view $cols
stty rows $rows
stty cols $cols

echo "Set VT size to $rows rows, $cols cols"