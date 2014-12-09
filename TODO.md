## TODO

Command bar: up/down, left, right, copy, paste
Expand selection + line number etc...
Open -> if already opened don't reopen
/ foo -> grep -nr foo $curfile
: 25 -> goto line 25
-> make commands like that by configuration ?

Most important now:

X Move /resize window
X Copy/paste
X List / open files
- GoFmt
- Configured events/shortcuts
- Dirty status should be kept from buffer status (insert, delete etc....)
- Syntax highlighting
- closeView/Col with dirty check allow if twice in row ? 
- closeView icon
+ Select section (mouse and/or Kb) -> with scrolling
- Warn on dirty view close
- Undo / Redo
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
- sed/sam like cmd language

## Commands
- execute (gofmt)
- command creates window
- search
- search/replace/next/

## WM
X support & display multiple columns / views
X move view 2+ of column to "top bar" -> create new column and put view in it
X Proportional move view / col (ratio)
X replaceView(drop old create new in place) with dirtyCheck, if dirty -> new view
X newcol
X  delcol -> check dirty
X newview
X delview > check dirty
X resize cols
~ exit -> check dirty
- dbl click to "collapse" / fold (single line)
- fullscreen view option ?

### UI
X 256 colors cant see cursor when on white space
- Color scheme when on white shell
- code highlights
- scrollbar indicator

### Events
- scrolling selection support
- CTRL+ O -> open selection
- ctrl + enter = execute current line
- select + ctrl + enter = execute   ALT+E ?
- -> execute goes to new window, closing that window kills the process ??
+ escape or caps or hh then,
- h help/commands (Ctrl+h)
X copy/paste (ctrl+c, ctrl+v)
- cut (ctrl+x)
- bl ? ^U -> Delete from cursor to start of line.
- bw ? ^W -> Delete word before the cursor.
- bs, delete/backspace on selection -> delete selection
+ home ^A -> Move cursor to start of the line.
+ end ^E -> Move cursor to end of the line.
+ o open (Ctrl+o) -> show recent first, tab completion -> if dir then to new window
- g goto (Ctrl+g)
- f find (Ctrl+f)
- next selection ? (Ctrl+n, Ctrl+Shift+n)
X nc newcol 
X nv niewview
X dc, dv delcol delview
- r refresh / reload -> Ctrl+r -> refresh buffer / dir listing ?
- re replace ? -> +CtrlN for replace next ??
- mh, mj, mk, ml move to view left right, up, down  (ctrl h,j,k,l)
- rh, rj, rk, rl relocate the view left right, up, down  (ctrl+shift+ h,j,k,l)
- sh,sj,sk,sl select l,r,u,d (alt h,j,k,l) or shift + arrows
- e exec -> output to new "shell" window (ctrl +e) (remember prevs ?)
- redo, undo (ctrl+z, ctrl+y)
- sa selectall ctrl+shift+a
+ supports thing like d, d 5, y 3, y, p ?
