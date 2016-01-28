### GoEd 
Goed is a Terminal based code/text editor, somewhat inspired by Acme.

**CURRENT STATE**:
It's currently in **ALPHA** and has not been spread around much yet.

I use it as my day to day editor and it "works on my machine"(TM)
It has not been tested much beyond that yet, there are [many open bugs and TODO's](https://github.com/tcolar/goed/issues).

Early screenshot (6/2/2015): 
![Screenshot](https://raw.github.com/tcolar/goed/master/screenshot.png)

Of course typically I have a much large window and resolution such as [this screenshot](https://raw.githubusercontent.com/tcolar/goed/master/screenshot_hd.png).

### Installation
Note: Might distribute binaries once stable.

Prerequities: 
- Have Go(Golang) installed
- Setup a good monospace terminal font and size (ie: Deja Vu Sans Mono or Monaco, size 10 or so.)

```
# Have your GOPATH set properly
go get github.com/tcolar/goed/...
```


To run it: 
```
export PATH=$PATH:$GOPATH/bin #(goed & goed_api must be in your $PATH)
goed <path(s)>

```

Quick start:
- Use right click or Ctrl+O to open a dir/file.
- Use CTRL+T to start a terminal view ($SHELL)

### Supported terminals
In theory it should work with any terminal, however the level of support for things 
like mouse support or extended colors vary a lot, here is a short list of tested 
setups so far:

Long story short is **use gnome-terminal on Linux, iTerm2 on OsX (with right click menu disabled)**.

Linux:
- Linux / GnomeTerminal : My usual setup, works great.
- Linux / Lxerminal : Works well but does not support mouse events above 256 columns.
- Linux / Konsole : Works well, seems mouse events are slightly offset ?
- Linux / Terminator : Works ok, but can't use right click for "open" action.

OSX (barely tested so far) :
- OSX / Iterm2 : Pretty good, **better if disabling right click context menu under "preferences ->  pointer"**. 
- OSX / Term.app : **Don't bother ...**

### Keyboard / Mouse shortcuts

Eventually the plan is to make all the Keyboard and mouse events configurable, 
but for now here they are:

| Action                                  | Linux             | Mac (iTerm2)  |
| --------------------------------------- | ----------------- | ------------- |
| Close view / window                     | LeftClick on 'x'  | same          |
|                                         | or CTRL+W         | same          |
| Copy line/selection text           (1)  | CTRL+C            | same          |
| Cut+Copy line/selection text            | CTRL+X            | same          |
| Execute command bar                     | Enter             | same          |
| Navigate views                          | CTRL+Arrows       | SHIFT+Arrrows |
| Open Terminal view                      | CTRL+T            | same          |
| Open path (in New view)                 | RightClick        | same          |
|                                         | or CTRL+O         | same          |
| Open path (in current view)             | ALT+O             | **TBD**       |
| Paste text                              | CTRL+V            | same          |
| Quit Goed                               | CTRL+Q            | same          |
| Redo                                    | CTRL+Y            | same          | 
| Reload view (not saving changes)        | CTRL+R            | same          |
| Resize / Move view & columns       (2)  | LeftClick top left| same          |
| Save view (to file)                     | CTRL+S            | same          |
| Select all                              | CTRL+A            | same          |
| Select word                             | Dbl LeftClick     | same          |
| Swap views                         (4)  | Dbl LC top left   | same          |
| Undo                                    | CTRL+Z            | same          |
| Text selection                     (3)  | Mouse drag        | same          |
|                                         | or SHIFT+<move>   | same          |
| Toogle command bar                      | ESC               | same          |


  - (1) : In terminal mode, copies if there is a selection, otherwise pass CTRL+C
  - (2) : You click the top left corner of a view (checkmark icon) to initiate the move, then click again where you want to drop it.
  - (3) : Can use Shift + [Arrows, PgUp, PgDown, Home, End].
  - (4) : Double click the top left corner of the view (checkmark icon) to switch that view with the previous current view.
  
### Terminal

Start a new Terminal with CTRL+T, it will be started in the same path as the current view.

The terminal implements basic vt100 support, enough for things such as top and 
interactive git to work.

Note that while in a terminal a limited number of global shortcuts are enabled.

### Terminal actions

The Terminal provides a few builtin shortcuts, such as:
  - "o <path>" : To open a given path/location in goed (or just right click it)
  - "s [-i] <pattern> [path]" : Search text (grep -rn[i] <pattern> [path])
  - "f <pattern> [path]" : Find files (find <path> -name *pattern*) 
  
See [res/default/actions](res/default/actions) for more info.

You may create your own actions in ~/.goed/ations/
See [res/Readme.md](res/Readme.md).

### Command bar
The command bar is at the top of the screen. you can toggle it by clicking it or
using the <ESC> key.

Note: The command bar will likely go away in favor of terminal actions or be replaced 
by something more useful such as sam like commands.

Currently it supports a few things:
  - "o <path>" : Opens a file or directory.
  - ": <linenumber>" : Goes to the secified line.
  - "/ <pattern>" : Search pattern (grep)
  Anything else will just be executed (via shell) in a new view.
  
### Contributing
- Reporting issues is welcome.
- PR's are even better.
- For new functionality a quick discussion first might be best.
    
