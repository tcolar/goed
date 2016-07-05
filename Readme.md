### GoEd 
Goed is a code/text editor, somewhat inspired by Acme.
It can run within a terminal or as a standalone lightweight gui.

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
go get -u github.com/tcolar/goed/...
```


To run it: 
```
export PATH=$PATH:$GOPATH/bin #(goed **MUST** be in your $PATH)
goed <path(s)>

```

Quick start:
- Use right click or Ctrl+N to open a dir/file.
- Use CTRL+T to start a terminal view ($SHELL)

### Terminal use
In theory it should work with any terminal, however the level of support for things 
like mouse support or extended colors vary a lot.

#### Terminal - Linux
I recommend GnomeTerminal as it has the best support, but Konsole or Lxterminal
should work as well.

#### Terminal - OSX
I highly recommend a real mouse(2+ buttons) and using ITerm2, **do not bother wth Term.app** 
as it has very poor eventing support.

For ITerm2 those settingw work best:
  - Under preferences / pointer, disable right click context menu.
  - Under prefs / pfofiles / default / terminal set ter type to "xterm-256color"
 
#### GUI (Experimental)
The eventing support in terminals varies immensely, some don't support mouse
events, some only support some CTRL, ALT sequences and almost none support any
type of advanced chording.

For this reason there is now an **experimental** GUI (**currently very slow**),
It's based on go.wde and has much much better and reliable eventing support. 

### Keyboard / Mouse shortcuts
Here are the standard key shortcuts, you can modify those to your liking, note
however that terminals support a limited set, in particular on OSX, basically
only CTRL combos work properly. Alt and Command combos are not reported.

You may use `goed --term-events` to find out what events work in your given terminal.

| Action                                  | Linux             | Mac (iTerm2)  |
| --------------------------------------- | ----------------- | ------------- |
| Close view / window                     | LeftClick on 'x'  | same          |
|                                         | or CTRL+W         | same          |
| Copy line/selection text           (1)  | CTRL+C            | same          |
| Cut+Copy line/selection text            | CTRL+X            | same          |
| Delete from start of line to cursor     | CTRL+U            | same          |
| Execute command bar                     | Enter             | same          |
| Move to syart of line                   | CTRL+A or home    | same          |
| Move to end of line                     | CTRL+E or end     | same          |
| Move down                               | CTRL+L or arrow   | same          |
| Move left                               | CTRL+H or arrow   | same          |
| Move right                              | CTRL+J or arrow   | same          |
| Move up                                 | CTRL+K or arrow   | same          |
| Navigate views                          | ALT+Arrows        | same          |
| Open Terminal view                      | CTRL+T            | same          |
| Open path (in New view)                 | RightClick        | same          |
|                                         | or CTRL+N         | same          |
| Open path (in current view)             | CTRL+O            | same          |
| Page up / Page down                     | PgUp, PgDown      | Fn + up / down|
| Paste text                              | CTRL+V            | same          |
| Quit Goed                               | CTRL+Q            | same          |
| Redo                                    | CTRL+Y            | same          |
| Reload view (not saving changes)        | CTRL+R            | same          |
| Resize / Move view & columns       (2)  | LeftClick top left| same          |
| Save view (to file)                     | CTRL+S            | same          |
| Select all                              | CTRL+B            | same          |
| Select word                             | Dbl LeftClick     | same          |
| Swap views                         (4)  | Dbl LC top left   | same          |
| Undo                                    | CTRL+Z            | same          |
| Text selection                     (3)  | Mouse drag        | same          |
|                                         | or SHIFT+move     | same          |
| Toogle command bar                      | ESC               | same          |


  - (1) : In terminal mode, copies if there is a selection, otherwise pass CTRL+C
  - (2) : You click the top left corner of a view (checkmark icon) to initiate the move, then click again where you want to drop it.
  - (3) : Can use Shift + [Arrows, PgUp, PgDown, Home, End].
  - (4) : Double click the top left corner of the view (checkmark icon) to switch that view with the previous current view.
  
### Terminal usage

Start a new Terminal with CTRL+T, it will be started in the same path as the current view.

The terminal implements basic vt100 support, enough for things such as top and 
interactive git to work.

Note that while in a terminal a limited number of global shortcuts are enabled.

### Terminal actions

The Terminal provides a few builtin shortcuts, such as:
  - `sz` : Set the shell tty rows/cols to match the current goed view size.
  - `o <path>` : To open a given path/location in goed (or just right click it)
  - `s [-i] <pattern> [path]` : Search text (grep -rn[i] <pattern> [path])
  - `f <pattern> [path]` : Find files (find <path> -name *pattern*) 
  
See [res/default/actions](res/default/actions) for more info.

You may create your own actions in ~/.goed/ations/
See [res/Readme.md](res/Readme.md).

### Command bar
The command bar is at the top of the screen. you can toggle it by clicking it or
using the <ESC> key.

Note: The command bar will likely go away in favor of terminal actions or be replaced 
by something more useful such as sam like commands.

Currently it supports a few things:
  - `o <path>` : Opens a file or directory.
  - `: <linenumber>` : Goes to the secified line.
  - `/ <pattern>` : Search pattern (grep)
  Anything else will just be executed (via shell) in a new view.

### Configuration
The config file can be edited at ~/.goed/config.toml (The original is under ~/.goed/default/) 

You may create custom themes under ~/.goed/themes/ (originals under ~/.goed/default/themes/)

You may create/override actions under ~/.goed/actions/ .

TODO : key/mouse actions

### Reporting issues
Report in github, try not to create duplicates.

If possible try to provide the most recent log found in ~/.goed/log/
  
### Contributing
- Reporting issues is welcome.
- PR's are even better.
- For new functionality a quick discussion first might be best.
    
