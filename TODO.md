## todo

### Core
X Terminal support
X View buffer
X Eventing basics (Mouse / keyboard)
X Navigate buffer (arrows, scroll, pagination etc..)
X Edit buffer -> insert
X Edit buffer -> delete
~ save file (C^s)
- open file(edit tag??) + create buffer copy in folder
- use file copy as buffer (use interface for i/o) Backend interface

->-> Minimum "usable" product

- redo/undo (disallow for large files ?), keep undos on disk ?
- open a folder (listing)
- plumber -> clickabl line numbers
- plumber -> click file paths
- Keyboard mapping

## Extended
X theme/colors
- Hexadecimal mode
- large file, in place mode (no copy ?)

## Commands
- execute (gofmt)
- command creates window
- search
- search/replace/next/

## WM
X support & display multiple columns / views
- newcol
- delcol
- newview
- delview
- move views
- resize cols
~ exit

### UI
- 256 colors can't see cursor when on white space
- Color scheme when on white shell
- code highlights
- scrollbar indicator

### Events
- cut(C^x)
- paste(C^v)
- No select + R click = paste
- copy(C^c)
- select + L click = copy
- putall
- scrolling (mouse wheel, scrollbar clicks)
- ctrl +c = copy
- shift + arrows = select
- select + enter = execute
- select + middle click = execute

- ctrl+a : select all or line beginning ?
- ctrl+g : goto ?

acme shortcuts:
alt+H,J,K,L -> vim nav ??
^U -> Delete from cursor to start of line.
^W -> Delete word before the cursor.
^H -> Delete character before the cursor.
^A -> Move cursor to start of the line.
^E -> Move cursor to end of the line.