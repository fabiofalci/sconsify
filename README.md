A spotify console app
---------------------

A very early stage of a spotify console application.

Requirements: [Libspotify SDK](https://developer.spotify.com/technologies/libspotify/) & [PortAudio](http://www.portaudio.com/).


Installation
------------

* Download current version [0.1.0-rc1](https://github.com/fabiofalci/sconsify/releases) 

* Install dependencies:

`Archlinux`

	$ pacman -S portaudio
	$ yaourt -S libspotify

`Ubuntu` & `debian` - please ubuntu/debian users we need your confirmation on this

	$ curl http://apt.mopidy.com/mopidy.gpg | apt-key add - && curl -o /etc/apt/sources.list.d/mopidy.list http://apt.mopidy.com/mopidy.list
	$ apt-get update && apt-get install -y libportaudio2 libspotify12 --no-install-recommends 

`OSX`

	$ brew tap homebrew/binary
	$ brew install portaudio
	$ brew install libspotify
	$ cd /usr/local/opt/libspotify/lib/
	$ ln -s libspotify.dylib libspotify

* Run `./sconsify`

![alt tag](https://raw.githubusercontent.com/wiki/fabiofalci/sconsify/sconsify.png)

Modes
-----

There are 2 modes: 

* `Console user interface` mode: it presents a text user interface with playlists and tracks.

* `No user interface` mode: it doesn't present user interface and just random tracks between playlists.


Parameters
----------

* `-username=""`: Spotify username. If not present username will be asked.

* Password will be asked. To not be asked you can set an environment variable with your password `export SCONSIFY_PASSWORD=password`. Be aware your password will be exposed as plain text.

* `-ui=true/false`: Run Sconsify with Console User Interface. If false then no User Interface will be presented and it'll only random between Playlists.

* `-playlists=""`: Select just some playlists to play. Comma separated list.


No UI Parameters
----------------

* `-noui-repeat-on=true/false`: Play your playlist and repeat it after the last track.

* `-noui-silent=true/false`: Silent mode when no UI is used.

* `-noui-random=true/false`: Random between tracks or follow playlist order.


UI mode keyboard 
----------------

* &larr; &darr; &uarr; &rarr; for navigation.

* `space` or `enter`: play selected track.

* `>`: play next track.

* `p`: pause.

* `/`: open a search field.

* `r`: random tracks in the current playlist. Press again to go back to normal mode.

* `R`: random tracks in all playlists. Press again to go back to normal mode.

* `u`: queue selected track to play next.

* `d`: delete selected track from the queue or delete selected search.

* `D`: delete all tracks from the queue if the focus is on the queue.

* `PageUp` `PageDown` `Home` `End`. 

* `Control C` or `q`: exit.

Vi navigation style:

* `h` `j` `k` `l` for navigation.

* `Nj` and `Nk` where N is a number: repeat the command N times.

* `gg`: go to first element. 

* `G`: go to last element.

* `Ngg` and `NG` where N is a number: go to element at position N. 


No UI mode keyboard 
-------------------

* `>`: play next track.

* `Control C`: exit.


sconsifyrc
----------

Similar to [.ackrc](http://beyondgrep.com/documentation/) you can define default parameters in `~/.sconsify/sconsifyrc`:

	-username=your-username
	-noui-silent=true 
	-noui-repeat-on=false


How to build and run using docker
---------------------------------

Get a Spotify application key and copy as a byte array to `/sconsify/spotify/spotify_key_array.key`.

	var key = []byte{
	    0x02, 0xA2, ...
	    ...
	    0xA1}

* `make run`

When building for OSX you may face an issue where it doesn't get you application key. Just retry the build that eventually it will get the key.