## TODO

### Core / Editor
- Bug: logitech mouse scrollwheel outputtings "A" in the file when moving fast !!
- redo/undo (disallow for large files ?), keep undos on disk ?
- Copy / read theme from ~/.goed (create if not there)
- Actions/Events/Key mapping -> defaults + customs (gofmt, color, run etc...)
- Copy indentation to next line on return
- Deal and/or reject files with CR/LF
- View ID : <goed start timestamp>_id ?

## Extended
- sed/sam like cmd language
- Benchmark of editor scroll/cursor insert etc ... ?
- Normalize 0 vs 1 index also line,col vs col,line & row
- More tests

## Commands
+ execute (gofmt)
- search
- search/replace/next/
- Command bar: up/down, left, right, copy, paste
  -> make commands like that by configuration ?

## WM
- Use extended mouse coordinates(>236) (SGR 1006) http://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
- dbl click to "collapse" / fold (single line)
- fullscreen view option ?
+ reload view / rerun command (TODO: dirty check)

### UI
- use a channel(size 1) for save / view render
- code highlights
- scrollbar indicator
- Open -> if no file, show welcome/help pane
- / foo -> grep -nr foo $curfile
- / foo [path] .... 
- "Local" file search with next/previous ?

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
