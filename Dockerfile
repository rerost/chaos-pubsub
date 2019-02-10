# FROM golang:1.11
FROM rerost/chaos-pubsub:latest

ENV GOPATH /go

ENV APP_ROOT $GOPATH/src/github.com/rerost/chaos-pubsub
RUN ln -s $APP_ROOT/ /app
WORKDIR /app

# Install Dependency
COPY Makefile ${APP_ROOT}
COPY Gopkg.toml ${APP_ROOT}/
COPY Gopkg.lock ${APP_ROOT}

RUN make vendor

# Build Binary
COPY . ${APP_ROOT}/
RUN gex grapi build

CMD ["/app/bin/server"]
