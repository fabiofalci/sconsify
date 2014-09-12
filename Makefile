test:
	go test -v ./...

run:
	go build && ./sconsify

#
# To build create spotify/spotify_key_array.key containing your application key
# as a go array:
#
#    var key = []byte{
#        0x02, 0xA2, ...
#        ...
#        0xA1}
#
build:
	sed -i '$$ d' spotify/key.go && cat spotify/spotify_key_array.key >> spotify/key.go && go install && git checkout spotify/key.go