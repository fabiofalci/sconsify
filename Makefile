UNAME := $(shell uname)

SED := sed -i '$$ d'
ifeq ($(UNAME), Darwin)
	SED := sed -i '' -e '$$ d'
endif

VERSION := 0.6.0-next
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date +"%s")
SPOTIFY_CLIENT_ID := 4e1fa8c468ce42c2a45c7c9e40e6d9d0
AUTH_REDIRECT_URL := https://fabiofalci.github.io/sconsify-auth/index.html

default: build

.PHONY: test
test:
	go test -v ./...

test-ui:
	go run test/test_sconsify.go -run-test

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
	glide -q install
	$(SED) spotify/key.go && cat spotify/spotify_key_array.key >> spotify/key.go \
		&& GO15VENDOREXPERIMENT=1 go build -ldflags "\
		 -X main.version=$(VERSION) \
		 -X main.commit=$(COMMIT) \
		 -X main.buildDate=$(BUILD_DATE) \
		 -X main.spotifyClientId=$(SPOTIFY_CLIENT_ID) \
		 -X main.authRedirectUrl=$(AUTH_REDIRECT_URL)" \
		 -o bundles/sconsify \
		; git checkout spotify/key.go

#
# pp: Crash your app in style (https://github.com/maruel/panicparse)
#
build-run: build
	./bundles/sconsify -debug |&pp

container-build: bundles
	docker build -t sconsify-build .

binary: container-build
	docker run --rm -v "$(CURDIR)/bundles/container:/go/src/github.com/fabiofalci/sconsify/bundles" sconsify-build make build

shell: container-build
	docker run --rm -it \
		-v "$(CURDIR)/bundles/container:/go/src/github.com/fabiofalci/sconsify/bundles" \
		-v /dev/snd:/dev/snd \
		--privileged \
		sconsify-build \
		bash

clean:
	rm -rf bundles/

bundles:
	mkdir -p bundles/container

#
# Only works on osx as we can generate both osx and linux binaries in one go.
#
release: clean binary build
	mkdir -p bundles/release/{linux,osx}
	cp bundles/sconsify bundles/release/osx/
	cp bundles/container/sconsify bundles/release/linux/
	zip -j bundles/release/osx-x86_64-sconsify-$(VERSION).zip bundles/release/osx/sconsify
	zip -j bundles/release/linux-x86_64-sconsify-$(VERSION).zip bundles/release/linux/sconsify

