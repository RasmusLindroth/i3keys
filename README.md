# i3keys
A program for the tiling window manager [i3](https://i3wm.org/).

This is program lists all the keys that are bound to some action in i3wm, and 
all keys that are not bound to any actions. Now you don't have to search 
through your configuration file or going down the track of trial and error 
anymore.

The program will ouput one separate keyboard for all of your different modifier 
combinations. And the keyboard will look like your keymap. I hope :)

You can see the results in two ways. Either by opening a local web page or 
output the keyboards as SVG files.

### Example of the output for Mod4 + any key
![Example image](https://i.imgur.com/4J1fbdQ.png)
* Green = the modifier key(s)
* Red = the binding is occupied
* White = one free key to use


## How to
Currently there is no released binary. You'll have to build the program.

### Using Arch?

You can find it in the Arch User Repository (AUR).

https://aur.archlinux.org/packages/i3keys/

#### Go getting and installing

If your version of Go have support for modules
```
//Get this project
git clone https://github.com/RasmusLindroth/i3keys.git

//Install
go install

//Run web interface on port 8080
i3keys web 8080

//or output SVG to the current directory
i3keys svg ISO ./

//If starting doesn't work try running it from $HOME/go/bin
```

Alternative way to install
```
//Get this program
go get -u github.com/RasmusLindroth/i3keys

//Go to the project directory
cd $HOME/go/src/github.com/RasmusLindroth/i3keys

//Install
go install

//Run web interface on port 8080
i3keys web 8080

//or output SVG to the current directory
i3keys svg ISO ./

//If starting doesn't work try running it from $HOME/go/bin
```

If you still having problems see the 
[installation guide for Go](https://golang.org/doc/install#install) or open 
up an issue.

#### You have started i3keys
Now you will need to start your browser and head over to the URL printed in 
your terminal e.g. http://localhost:8080

There you can select your keyboard layout and voilà!

Or if you have selected to out SVG files. Open them with a program that supports 
SVG files. For example your program for viewing images, InkScape or a browser.

### Help message

```
Usage:

	i3keys <command> [arguments]

The commands are:

	web <port>            start the web ui and listen on <port>
	svg <layout> [dest]   outputs one SVG file for each modifier group. <layout> can be ISO or ANSI, [dest] defaults to current directory
	version               print i3keys version
```

### Disclaimer
* It's only tested with evdev handling input. So maybe you get the wrong 
 mappings. Open an issue in that case and I will look in to it.
* There are no tests right now. So you might run into some issues.
