### GoEd 
Goed is a Terminal based code/text editor.

**CURRENT STATE**

It's currently in **EARLY ALPHA**, meaning it's somewhat usable, and I use it as
my daily editor, but it's probably still buggy and definitely changing a lot.
So far, tested mostly only Linux and a tiny bit on OSX.

### What is it ?
I would say that the main source of inspiration is **Acme**, however it will
be more configurable and not require the mouse quite as much, but still leverage it.

Early screenshot (6/2/2015): 
![Screenshot](https://raw.github.com/tcolar/goed/master/screenshot.png)

### Installation
Prerequities: 
- Have Go(Golang) installed
- Setup a good terminal font and size (ie: Monospace, 10 possibly)

```
export PATH=$PATH:$GOPATH/bin
go get github.com/tcolar/goed
cd $GOPATH/src/github.com/tcolar/goed
./build.sh
```

Note: Might distribute binaries once stable.

To run it: 
```
goed <file(s)>
```

### Supported terminals
In theory it should work on any terminal, however the level of support for things 
like mouse support or extended colors vary a lot, here is a short list of tested 
setups so far:

- Linux / GnomeTerminal : My usual setup, workd great.
- Linux / Lxerminal : Works well but does not support mouse events above 256 columns.
- Linux / Konsole : Works well, seems mouse is offset a bit off maybe.
- Linux / Terminator : Works ok, but can't use right click for "open" action

OSX (barely tested so far) :
- OSX / Iterm2 : Pretty good, can't use right mouse click.
- OSX / Term.app : Not very good, does not seem to evesupport mouse ??n 

### Manual
TODO: How does it work, UI usage, shortcuts etc ....

### FAQ's
TODO

### Contributing
- Reporting issues is welcome.
- PR's are welcome.
- For new functionality a quick discussion first might be best.
    
