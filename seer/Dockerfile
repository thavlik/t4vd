ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="ghcr.io/puppeteer/puppeteer:latest"
FROM ${KINDEST_REPOSITORY}thavlik/bjjv-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY sources/pkg/api sources/pkg/api
COPY seer seer
WORKDIR /go/src/github.com/thavlik/bjjvb/seer/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o seer
FROM ${RUNTIME_BASE_IMAGE}
USER root
RUN apt-get update \
    && apt-get install -y \
        python3 \
        python3-distutils \
        ffmpeg \
    && rm -rf /var/lib/apt/lists/* \
    && wget -q -O - https://bootstrap.pypa.io/get-pip.py | python3 - \
    && python3 -m pip install youtube_dl
COPY --from=builder /go/src/github.com/thavlik/bjjvb/seer/cmd/seer /usr/local/bin
COPY seer/scripts /scripts
CMD ["seer"]

