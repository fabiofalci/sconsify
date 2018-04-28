FROM ubuntu:17.10

RUN apt-get update && apt-get install -y curl build-essential git mercurial portaudio19-dev ca-certificates --no-install-recommends 

# Libspotify
RUN curl http://apt.mopidy.com/mopidy.gpg | apt-key add - && curl -o /etc/apt/sources.list.d/mopidy.list http://apt.mopidy.com/mopidy.list
RUN apt-get update && apt-get install -y libspotify-dev

# Install Go
RUN curl -sSL https://storage.googleapis.com/golang/go1.10.1.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /go

# Install glide for Go	
RUN go get github.com/Masterminds/glide

WORKDIR /go/src/github.com/fabiofalci/sconsify

# Upload sconsify source
COPY . /go/src/github.com/fabiofalci/sconsify

ENV PATH=$PATH:$GOPATH/bin
