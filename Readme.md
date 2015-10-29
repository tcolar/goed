### GoEd 
Goed is a Terminal based code/text editor.

**CURRENT STATE**:

It's currently in **EARLY ALPHA**, I use it as my day to day editor and it "works on my machine" (TM)
It has not been tested much beyond that yet, there are many open bugs and documentation is lacking.

The goal is to realease more stable and usable alphas & betas soon.

### What is it ?
I would say that the main source of inspiration is **Acme**, however it will
be more configurable and not require the mouse quite as much, but still leverage it.

Early screenshot (6/2/2015): 
![Screenshot](https://raw.github.com/tcolar/goed/master/screenshot.png)

### Installation
Prerequities: 
- Have Go(Golang) installed
- Setup a good monospace terminal font and size (ie: Deja Vu Sans Mono or Monaco, size 10 or so.)

```
export PATH=$PATH:$GOPATH/bin
go get github.com/tcolar/goed
cd $GOPATH/src/github.com/tcolar/goed
./build.sh
```

Note: Might distribute binaries once stable.

To run it: 
```
goed <path(s)>
```

### Supported terminals
In theory it should work with any terminal, however the level of support for things 
like mouse support or extended colors vary a lot, here is a short list of tested 
setups so far:

Linux:
- Linux / GnomeTerminal : My usual setup, works great.
- Linux / Lxerminal : Works well but does not support mouse events above 256 columns.
- Linux / Konsole : Works well, seems mouse events is slightly offset ?
- Linux / Terminator : Works ok, but can't use right click for "open" action.

OSX (barely tested so far) :
- OSX / Iterm2 : Pretty good, can't use right mouse click, scolling artifacts.
- OSX / Term.app : Don't bother ...

Windows:
- Windows / what : Don't tell me you write code in a terminal on windows !

### Manual
TODO: How does it work, UI usage, shortcuts etc ...

### FAQ's
TODO

### Contributing
- Reporting issues is welcome.
- PR's are welcome.
- For new functionality a quick discussion first might be best.
    
