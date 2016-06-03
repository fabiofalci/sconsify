UNAME := $(shell uname)

SED := sed -i '$$ d'
ifeq ($(UNAME), Darwin)
	SED := sed -i '' -e '$$ d'
endif

VERSION := 0.3.0-local
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date +"%s")
SPOTIFY_CLIENT_ID := 4e1fa8c468ce42c2a45c7c9e40e6d9d0
AUTH_REDIRECT_URL := https://fabiofalci.github.io/sconsify/auth/

default: build

test:
	go test -v ./...

test-ui:
	go run test/test_sconsify.go -run-test |&pp

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
	$(GOPATH)/bin/glide -q install
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
	docker run --rm -v "$(CURDIR)/bundles:/go/src/github.com/fabiofalci/sconsify/bundles" sconsify-build make build

shell: container-build
	docker run --rm -it \
		-v "$(CURDIR)/bundles:/go/src/github.com/fabiofalci/sconsify/bundles" \
		-v /dev/snd:/dev/snd \
		--privileged \
		sconsify-build \
		bash

bundles:
	mkdir -p bundles
