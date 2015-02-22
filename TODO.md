## TODO

X Open -> if already opened don't reopen
- Open -> if no file, show welcome/help pane
- : 25 -> goto line 25
- / foo -> grep -nr foo $curfile
- / foo [path] .... 
- "Local" file search with next/previous ?
X Move view to close button -> panic !
- Bug: logitech mouse scrollwheel outputtings "A" in the file when moving fast !!
+ reload view / rerun command (TODO: dirty check)
X ctrl+x
X delete/backspace/overwrite selection
- undo/redo
- split event in multiple files, event, ui folder
- copy indentation to next line on return
- Command bar: up/down, left, right, copy, paste
  -> make commands like that by configuration ?
- use a channel(size 1) for save / view render
- Deal and/or reject files with CR/LF
X Broken delete at EOF
X Test backend/view
~ Test editor/Wm
X reimpl memory buffer
X mem buffer cmd (default for dir listings)
- View ID : <goed start timestamp>_id ?
X Replace lines/viewlines etc.. by uing the slice -> performance
- Benchmark of editor scroll/cursor insert etc ... ?
- Normalize 0 vs 1 index also line,col vs col,line & row
- goed package
- tests : view.Save, scrolling, copy/yank, paste
- tests cmdbar : runCmd, paste, yank, open, newview, newcol,  exec
- tests editor : Opendir, setstatus, setstatuserr
- tests event.go ?
- tests io.go : copyfile
- test wm : widgetat, viewMove, viewColumn, ViewIndex, ColIndex, AddCol,
	AddViewSmart, InsertViewSmart, AddView, InsertView, ReplaceView,
	DelCol, DelView, TerminateView, DelColCheck, ActivaeView, ViewById


Most important now:

X Move /resize windo
X Copy/paste
X List / open files
+ GoFmt
- Configured events/shortcuts
X Dirty status should be kept from buffer status (insert, delete etc....)
- Syntax highlighting
X closeView/Col with dirty check allow if twice in row ? 
X closeView icon
- Select section (mouse and/or Kb) -> with scrolling
X Warn on dirty view close
- Undo / Redo
X File buffers

### Core
X Terminal support
X View buffer
X Eventing basics (Mouse / keyboard)
X Navigate buffer (arrows, scroll, pagination etc..)
X Edit buffer -> insert
X Edit buffer -> delete
X save file (C^s)
X open file(edit tag??)
X create buffer copy in folder
X use file copy as buffer (use interface for i/o) Backend interface

->-> Minimum "usable" product

- redo/undo (disallow for large files ?), keep undos on disk ?
X open a folder (listing)
X plumber -> clickable line numbers
X plumber -> click file paths
- Keyboard mapping

## Extended
X theme/colors
- Hexadecimal mode
X large file, in place mode (no copy ?)
- sed/sam like cmd language

## Commands
+ execute (gofmt)
X command creates window
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
X exit -> check dirty
- dbl click to "collapse" / fold (single line)
- fullscreen view option ?

### UI
X 256 colors cant see cursor when on white space
X Color scheme when on white shell
- code highlights
- scrollbar indicator

### Events
- scrolling selection support
X CTRL+ O -> open selection
- ctrl + enter = execute current line
- select + ctrl + enter = execute   ALT+E ?
X -> execute goes to new window, closing that window kills the process ??
+ escape or caps or hh then,
- h help/commands (Ctrl+h)
X copy/paste (ctrl+c, ctrl+v)
X cut (ctrl+x)
- bl ? ^U -> Delete from cursor to start of line.
- bw ? ^W -> Delete word before the cursor.
- bs, delete/backspace on selection -> delete selection
+ home ^A -> Move cursor to start of the line.
+ end ^E -> Move cursor to end of the line.
X o open (Ctrl+o) -> show recent first, tab completion -> if dir then to new window
- g goto (Ctrl+g)
- f find (Ctrl+f)
- next selection ? (Ctrl+n, Ctrl+Shift+n)
X nc newcol 
X nv niewview
X dc, dv delcol delview
X r refresh / reload -> Ctrl+r -> refresh buffer / dir listing ?
- re replace ? -> +CtrlN for replace next ??
- mh, mj, mk, ml move to view left right, up, down  (ctrl h,j,k,l)
- rh, rj, rk, rl relocate the view left right, up, down  (ctrl+shift+ h,j,k,l)
- sh,sj,sk,sl select l,r,u,d (alt h,j,k,l) or shift + arrows
+ e exec -> output to new "shell" window (ctrl +e) (remember prevs ?)
- redo, undo (ctrl+z, ctrl+y)
- sa selectall ctrl+shift+a
+ supports thing like d, d 5, y 3, y, p ?
