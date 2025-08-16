#
# 1. Build Container
#
FROM registry.arvan.tech/docker-golang-builder:1.23 AS build

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

ARG GO_PROXY
ENV GOPROXY=${GO_PROXY}

RUN mkdir -p /src

# First add modules list to better utilize caching
COPY go.sum go.mod /src/

WORKDIR /src

COPY . /src

# Build components.
# Put built binaries and runtime resources in /app dir ready to be copied over or used.
RUN make install && \
    mkdir -p /app && \
    cp -r $GOPATH/bin/arvanch /app/

RUN cp -r /src/migrations /app/

#
# 2. Runtime Container
#
FROM registry.arvan.tech/docker-golang-deployer:3.20

COPY --from=build /app /app/

ENTRYPOINT ["./arvanch"]
