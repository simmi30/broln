# If you change this value, please change it in the following files as well:
# /.travis.yml
# /Dockerfile
# /make/builder.Dockerfile
# /.github/workflows/main.yml
# /.github/workflows/release.yml
FROM golang:1.17.3-alpine as builder

LABEL maintainer="Olaoluwa Osuntokun <laolu@lightning.engineering>"

# Force Go to use the cgo based DNS resolver. This is required to ensure DNS
# queries required to connect to linked containers succeed.
ENV GODEBUG netdns=cgo

# Install dependencies and install/build broln.
RUN apk add --no-cache --update alpine-sdk \
    git \
    make 

# Copy in the local repository to build from.
COPY . /go/src/github.com/brsuite/broln

RUN cd /go/src/github.com/brsuite/broln \
&&  make \
&&  make install tags="signrpc walletrpc chainrpc invoicesrpc"

# Start a new, final image to reduce size.
FROM alpine as final

# Expose broln ports (server, rpc).
EXPOSE 9782 10019

# Copy the binaries and entrypoint from the builder image.
COPY --from=builder /go/bin/brolncli /bin/
COPY --from=builder /go/bin/broln /bin/

# Add bash.
RUN apk add --no-cache \
    bash

# Copy the entrypoint script.
COPY "docker/broln/start-broln.sh" .
RUN chmod +x start-broln.sh
