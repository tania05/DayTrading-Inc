FROM golang:1.9

ARG buildpath
ENV GOPATH=$GOPATH:/app
ENV PATH=$PATH:/usr/local/go/bin

# Run our go-gets here so they are cached
RUN go get github.com/gorilla/mux
RUN go get github.com/BurntSushi/toml
RUN go get github.com/go-redis/redis
RUN go get github.com/go-redis/cache
RUN go get github.com/vmihailenco/msgpack
RUN go get github.com/lib/pq
RUN go get github.com/valyala/gorpc

# ADD . /app/src
ADD webserver /app/src/webserver
ADD trigger /app/src/trigger
ADD load-balancer /app/src/load-balancer
ADD common /app/src/common

RUN cd /app/src/$buildpath && go build && go install


FROM ubuntu:16.04
COPY --from=0 /app/bin/$buildpath /app/bin/$buildpath
# Yes, this is in the root directory
ADD config.toml config.toml

EXPOSE 8080
CMD ["/bin/sh"]
