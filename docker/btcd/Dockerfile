FROM golang:1.17.3-alpine as builder

LABEL maintainer="Olaoluwa Osuntokun <laolu@lightning.engineering>"

# Install build dependencies such as git and glide.
RUN apk add --no-cache git gcc musl-dev

WORKDIR $GOPATH/src/github.com/brsuite/brond

# Pin down brond to a version that we know works with broln.
ARG brond_VERSION=v0.20.1-beta

# Grab and install the latest version of of brond and all related dependencies.
RUN git clone https://github.com/brsuite/brond.git . \
    && git checkout $brond_VERSION \
    && go install -v . ./cmd/...

# Start a new image
FROM alpine as final

# Expose mainnet ports (server, rpc)
EXPOSE 8333 8334

# Expose testnet ports (server, rpc)
EXPOSE 18333 18334

# Expose simnet ports (server, rpc)
EXPOSE 18555 18556

# Expose segnet ports (server, rpc)
EXPOSE 28901 28902

# Copy the compiled binaries from the builder image.
COPY --from=builder /go/bin/addblock /bin/
COPY --from=builder /go/bin/bronctl /bin/
COPY --from=builder /go/bin/brond /bin/
COPY --from=builder /go/bin/findcheckpoint /bin/
COPY --from=builder /go/bin/gencerts /bin/

COPY "start-bronctl.sh" .
COPY "start-brond.sh" .

RUN apk add --no-cache \
    bash \
    ca-certificates \
&&  mkdir "/rpc" "/root/.brond" "/root/.bronctl" \
&&  touch "/root/.brond/brond.conf" \
&&  chmod +x start-bronctl.sh \
&&  chmod +x start-brond.sh \
# Manually generate certificate and add all domains, it is needed to connect
# "bronctl" and "broln" to "brond" over docker links.
&& "/bin/gencerts" --host="*" --directory="/rpc" --force

# Create a volume to house pregenerated RPC credentials. This will be
# shared with any broln, bronctl containers so they can securely query brond's RPC
# server.
# You should NOT do this before certificate generation!
# Otherwise manually generated certificate will be overridden with shared
# mounted volume! For more info read dockerfile "VOLUME" documentation.
VOLUME ["/rpc"]
