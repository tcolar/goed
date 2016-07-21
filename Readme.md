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

### Download binaries

You can [download prebuilt standalone binaries from bintray here](https://bintray.com/tcolar/Goed/Goed#files), built using [release.sh](release.sh).

### Build from source
If you rather build yourself :

Prerequities: 
- Have Go(Golang) installed
- Setup a good monospace terminal font and size (ie: Deja Vu Sans Mono or Monaco, size 10 or so.)

```
# Have your GOPATH set properly
go get -u github.com/tcolar/goed/...
```

### Running goed
Note : goed **MUST** be in your $PATH !!

```
which goed        # must be found in your path
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

For this reason there an **experimental** GUI (**currently very slow**),
based on go.wde, in the works, not quite ready yet. 

### Keyboard / Mouse shortcuts
Here are the standard key shortcuts, you can modify those to your liking, note
however that terminals support a limited set, in particular on OSX, basically
only CTRL combos work properly. Alt and Command combos are not reported by the
termbox library used by Goed.

You may use `goed --term-events` to find out what events work in your given terminal.

You can customize the mouse/keyboard shortcuts in ~/.goed/bindings.toml
Here are the [standard mouse/keyboard bindings](res/default/bindings.toml)
  
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

You may create your own actions in ~/.goed/ations/ (Work In progress)

See [res/Readme.md](res/Readme.md).

### Command bar
The command bar is at the top of the screen. you can toggle it by clicking it or
using the <ESC> key, think of it as a minimal one line terminal.

Currently it supports a few things:
  - `o <path>` : Opens a file or directory.
  - `: <linenumber>` : Goes to the secified line.
  - `/ <pattern>` : Search pattern (grep)
  
Anything else will just be executed (via shell) into a new view.

Eventually this will allow for custom defined actions based on patterns.

### Configuration
The config file can be edited at ~/.goed/config.toml (The original is under ~/.goed/default/) 

Key/Mouse bindings can be customized at ~/.goed/bindings.toml (original under ~/.goed/default/bindings.toml)

You may create custom themes under ~/.goed/themes/ (originals under ~/.goed/default/themes/)

You may create/override actions under ~/.goed/actions/

### Reporting issues
Report in github, try not to create duplicates.

If possible try to provide the most recent log found in ~/.goed/log/
  
### Contributing
- Reporting issues is welcome.
- PR's are even better.
- For new functionality a quick discussion first might be best.
    
