UNAME := $(shell uname)

SED := sed -i '$$ d'
ifeq ($(UNAME), Darwin)
	SED := sed -i '' -e '$$ d'
endif

VERSION := 0.1.0-rc2
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u)

default: build

test:
	go test -v ./...

run: container-build
	docker run --rm -it \
		-v "$(CURDIR)/bundles:/go/src/github.com/fabiofalci/sconsify/bundles" \
		-v /dev/snd:/dev/snd \
		--privileged \
		sconsify-build \
		bash -c 'make build && ./bundles/sconsify'

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
	$(SED) spotify/key.go && cat spotify/spotify_key_array.key >> spotify/key.go \
		&& go build -ldflags "-X main.version $(VERSION) -X main.commit $(COMMIT) -X main.buildDate '$(BUILD_DATE)'" -o bundles/sconsify \
		; git checkout spotify/key.go

container-build: bundles
	docker build -t sconsify-build .

binary: container-build
	docker run --rm -it -v "$(CURDIR)/bundles:/go/src/github.com/fabiofalci/sconsify/bundles" sconsify-build make build

shell: container-build
	docker run --rm -it \
		-v "$(CURDIR)/bundles:/go/src/github.com/fabiofalci/sconsify/bundles" \
		-v /dev/snd:/dev/snd \
		--privileged \
		sconsify-build \
		bash

bundles:
	mkdir -p bundles
