# If you change this value, please change it in the following files as well:
# /.travis.yml
# /dev.Dockerfile
# /make/builder.Dockerfile
# /.github/workflows/main.yml
# /.github/workflows/release.yml
FROM golang:1.17.3-alpine as builder

# Force Go to use the cgo based DNS resolver. This is required to ensure DNS
# queries required to connect to linked containers succeed.
ENV GODEBUG netdns=cgo

# Pass a tag, branch or a commit using build-arg.  This allows a docker
# image to be built from a specified Git state.  The default image
# will use the Git tip of master by default.
ARG checkout="master"

# Install dependencies and build the binaries.
RUN apk add --no-cache --update alpine-sdk \
    git \
    make \
    gcc \
&&  git clone https://github.com/brsuite/broln /go/src/github.com/brsuite/broln \
&&  cd /go/src/github.com/brsuite/broln \
&&  git checkout $checkout \
&&  make release-install

# Start a new, final image.
FROM alpine as final

# Define a root volume for data persistence.
VOLUME /root/.broln

# Add utilities for quality of life and SSL-related reasons. We also require
# curl and gpg for the signature verification script.
RUN apk --no-cache add \
    bash \
    jq \
    ca-certificates \
    gnupg \
    curl

# Copy the binaries from the builder image.
COPY --from=builder /go/bin/brolncli /bin/
COPY --from=builder /go/bin/broln /bin/
COPY --from=builder /go/src/github.com/brsuite/broln/scripts/verify-install.sh /
COPY --from=builder /go/src/github.com/brsuite/broln/scripts/keys/* /keys/

# Store the SHA256 hash of the binaries that were just produced for later
# verification.
RUN sha256sum /bin/broln /bin/brolncli > /shasums.txt \
  && cat /shasums.txt

# Expose broln ports (p2p, rpc).
EXPOSE 9782 10019

# Specify the start command and entrypoint as the broln daemon.
ENTRYPOINT ["broln"]
