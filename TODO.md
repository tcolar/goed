## todo

Most important now:

X Move /resize window
- Select section (mouse and/or Kb)
- Copy/paste
- Warn on dirty view close
- Undo / Redo
- GoFmt
- File buffers ?

### Core
X Terminal support
X View buffer
X Eventing basics (Mouse / keyboard)
X Navigate buffer (arrows, scroll, pagination etc..)
X Edit buffer -> insert
X Edit buffer -> delete
X save file (C^s)
X open file(edit tag??)
- create buffer copy in folder
- use file copy as buffer (use interface for i/o) Backend interface

->-> Minimum "usable" product

- redo/undo (disallow for large files ?), keep undos on disk ?
- open a folder (listing)
- plumber -> clickable line numbers
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
X newcol
X  delcol -> check dirty
X newview
X delview > check dirty
- move views
- resize cols
~ exit -> check dirty

### UI
- 256 colors cant see cursor when on white space
- Color scheme when on white shell
- code highlights
- scrollbar indicator

### Events
- select + enter = execute
- select + middle click = execute
+ escape or caps or hh then,
- h help/commands (Ctrl+h)
- copy/paste (ctrl+c, ctrl+v)
- bl ? ^U -> Delete from cursor to start of line.
- bw ? ^W -> Delete word before the cursor.
+ home ^A -> Move cursor to start of the line.
- end ^E -> Move cursor to end of the line.
+ o open (Ctrl+o) -> show recent first, tab completion -> if dir then to new window
- g goto (Ctrl+g)
- f find (Ctrl+f)
X nc newcol 
X nv niewview
X dc, dv delcol delview
- r refresh / reload
- mh, mj, mk, ml move to view left right, up, down  (ctrl h,j,k,l)
- rh, rj, rk, rl relocate the view left right, up, down  (ctrl+shift+ h,j,k,l)
- sh,sj,sk,sl select l,r,u,d (alt h,j,k,l) or shift + arrows
- e exec -> output to new "shell" window (ctrl +e) (remember prevs ?)
- redo, undo (ctrl+z, ctrl+y)
- sa selectall ctrl+shift+a

