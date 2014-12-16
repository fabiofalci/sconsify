UNAME := $(shell uname)

SED := sed -i '$$ d'
ifeq ($(UNAME), Darwin)
	SED := sed -i '' -e '$$ d'
endif

default: build

test:
	go test -v ./...

run: build
	bundles/sconsify

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
	go get ./...
	$(SED) spotify/key.go && cat spotify/spotify_key_array.key >> spotify/key.go && go build -o bundles/sconsify ; git checkout spotify/key.go

container-build: bundles
	docker build -t sconsify-build .

binary: container-build
	docker run --rm -it -v "$(CURDIR)/bundles:/go/src/github.com/fabiofalci/sconsify/bundles" sconsify-build make build

shell: container-build
	docker run --rm -it -v "$(CURDIR)/bundles:/go/src/github.com/fabiofalci/sconsify/bundles" sconsify-build bash

bundles:
	mkdir -p bundles
