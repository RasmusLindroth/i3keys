# i3keys
A program for the tiling window manager [i3](https://i3wm.org/).
It also supports [Sway](https://swaywm.org/) but then you need to use the flag `-s`. 

This is program lists all the keys that are bound to some action in i3 or sway, and 
all keys that are not bound to any actions. Now you don't have to search 
through your configuration file or going down the track of trial and error 
anymore.

The program will output one separate keyboard for all of your different modifier 
combinations. It will also output each combination for all of your modes.
The keyboard should look like your keymap if everything works as expected :)

You can see the results in three ways. By opening a local web page,
output the keyboards as text in the terminal or generate SVG files.

### Example of the output for Mod4 + any key
![Example image](https://i.imgur.com/4J1fbdQ.png)
* Green = the modifier key(s)
* Red = the binding is occupied
* White = one free key to use


## How to
Currently there is one binary for Linux amd64. You can find it under [releases](https://github.com/RasmusLindroth/i3keys/releases). Download it and make it exacutable with `chmod +x i3keys` and then you can run it with `./i3keys`. If you want to use it "everywhere" you must move it to somewhere in your $PATH.

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

//If starting doesn't work try running it from $HOME/go/bin
```

Example usage
```
//Run web interface on port 8080
i3keys web 8080

//sway usage
i3keys -s web 8080

//or output text to the terminal
i3keys text ISO

//or filter text output
i3keys text ISO Mod4+Ctrl

//or output SVG to the current directory
i3keys svg ISO ./

//or filter SVG
i3keys svg ISO ./ Mod4+Ctrl
```

If you still having problems see the 
[installation guide for Go](https://golang.org/doc/install#install) or open 
up an issue.

#### You have started i3keys

For web output you will need to start your browser and head over to the URL 
printed in  your terminal e.g. http://localhost:8080. This is is the most 
feature-rich option. Because here you can see the command bound to each key.

If you have selected to out SVG files. Open them with a program that supports 
SVG files. For example your program for viewing images, InkScape or a browser.
Note that each mode will be located in a folder named mode\<int\>-name-of-mode.
Inside that folder you will find a file named `mode-full-name` where you can see
the real name of the mode. Because this program removes some characters and 
shortens the name of the folder.

If you have selected text output, you should already have seen the result in 
your terminal.

### Text example

```
$ i3keys text iso

Available keybindings per modifier group

Mod4:
-----
Escape, F1, F4, F6, F7, F8, F9, F10, F11, F12, Scroll_Lock, Pause
section, plus, dead_acute, BackSpace, Insert, Home, Page_Up, Num_Lock, KP_Divide, KP_Multiply, KP_Subtract
i, o, dead_diaeresis, Delete, End, Page_Down, KP_Home, KP_Up, KP_Page_Up, KP_Add
Caps_Lock, g, quoteright, KP_Left, KP_Begin, KP_Right
Shift_L, less, x, Shift_R, KP_End, KP_Down, KP_Page_Down, KP_Enter
Control_L, Super_L, Alt_L, ISO_Level3_Shift, Super_R, Menu, Control_R, KP_Insert, KP_Delete


Mod4+Shift:
-----------
Escape, F1, F2, F4, F5, F6, F7, F8, F9, F10, F11, F12, Scroll_Lock, Pause
section, 9, 0, plus, dead_acute, BackSpace, Insert, Home, Page_Up, Num_Lock, KP_Divide, KP_Multiply, KP_Subtract
Tab, w, t, y, u, i, o, p, dead_diaeresis, Return, Delete, End, Page_Down, KP_Home, KP_Up, KP_Page_Up, KP_Add
Caps_Lock, a, f, g, quoteright, KP_Left, KP_Begin, KP_Right
Shift_L, less, z, x, v, n, comma, period, Shift_R, KP_End, KP_Down, KP_Page_Down, KP_Enter
Control_L, Super_L, Alt_L, ISO_Level3_Shift, Super_R, Menu, Control_R, KP_Insert, KP_Delete
```

### Help message

```
Usage:

	i3keys [-s] <command> [arguments]
	Add the flag -s for sway

The commands are:

	web <port>
		start the web ui and listen on <port>

	text <layout> [mods]
		output available keybindings in the terminal

	svg <layout> [dest] [mods]
		outputs one SVG file for each modifier group

	version
		print i3keys version

Arguments:

	<layout>
		is required. Can be ISO or ANSI

	[mods]
		is optional. Can be a single modifier or a group of modifiers. Group them with a plus sign, e.g. Mod4+Ctrl

	[dest]
		is optional. Where to output files, defaults to the current directory

```

### Disclaimer
* It's only tested with evdev handling input. So maybe you get the wrong 
 mappings. Open an issue in that case and I will look in to it.
* There are no tests right now. So you might run into some issues.
