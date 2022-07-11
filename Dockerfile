# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-db-copier"
LABEL REPO="https://github.com/mlibrodo/db-copier"

ENV PROJPATH=/go/src/github.com/mlibrodo/db-copier

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/mlibrodo/db-copier
WORKDIR /go/src/github.com/mlibrodo/db-copier

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/mlibrodo/db-copier"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/db-copier/bin

WORKDIR /opt/db-copier/bin

COPY --from=build-stage /go/src/github.com/mlibrodo/db-copier/bin/db-copier /opt/db-copier/bin/
RUN chmod +x /opt/db-copier/bin/db-copier

# Create appuser
RUN adduser -D -g '' db-copier
USER db-copier

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/db-copier/bin/db-copier"]
