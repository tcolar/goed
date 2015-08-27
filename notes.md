
o [foo.go] -> open / create action
oe [foo.pdf] -> open external : open / xdg-open
b "os.exec" -> web search action
/ [text] [path] -> search file/folder contents action
f [pattern] -> search file name action (find)

ctrl+f -> "/ [selection]" in cmdbar ?
ctrl+g -> next ?? (right click in acme)

ctrl+Enter -> execute line
ctrl+Tab -> path completion
ctrl+b -> b selection

-------------
command

name="s"
#pattern="^s/.*"
enabled=true
action="foo.sh"
# action = go run foo.go
# onEvent="save"
# shortcuts="ctrl+f"
