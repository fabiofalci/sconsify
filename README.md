A spotify console app
---------------------

A very early stage of a spotify console application.

Requirements: [Libspotify SDK](https://developer.spotify.com/technologies/libspotify/) & [PortAudio](http://www.portaudio.com/).

### ArchLinux

    pacman -S portaudio
    yaourt -S libspotify


Dev: How to run
---------------

* Get a Spotify application key and copy to `/sconsify/spotify_appkey.key`

* `make run`

![alt tag](https://raw.githubusercontent.com/wiki/fabiofalci/sconsify/sconsify.png)


Binary build
------------

* Copy Spotify application key as a byte array to `/sconsify/spotify/spotify_key_array.key`.

* `make build`

* Run `$GOPATH/bin/sconsify` 
