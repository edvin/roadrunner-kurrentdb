FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/roadrunner-server/velox:latest as velox

# app version and build date must be passed during image building (version without any prefix).
# e.g.: `docker build --build-arg "APP_VERSION=1.2.3" --build-arg "BUILD_TIME=$(date +%FT%T%z)" .`
ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

# copy your configuration into the docker
COPY velox.toml .

# we don't need CGO
ENV CGO_ENABLED=0

# RUN build
RUN vx build -c velox.toml -o /usr/bin/

FROM --platform=${TARGETPLATFORM:-linux/amd64} php:8.3-cli

# copy required files from builder image
COPY --from=velox /usr/bin/rr /usr/bin/rr

# use roadrunner binary as image entrypoint
CMD ["/usr/bin/rr"]