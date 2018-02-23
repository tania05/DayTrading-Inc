FROM golang:latest

ARG buildpath
ENV GOPATH=$GOPATH:/app
ENV PATH=$PATH:/usr/local/go/bin

# Run our go-gets here so they are cached
RUN go get github.com/gorilla/mux
RUN go get github.com/BurntSushi/toml
RUN go get github.com/go-redis/redis
RUN go get github.com/go-redis/cache
RUN go get github.com/vmihailenco/msgpack


ADD . /app/src

RUN cd /app/src/$buildpath && go build && go install

# Yes, this is in the root directory
ADD config.toml config.toml

EXPOSE 8080
CMD ["/bin/bash"]