ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="alpine:3.17.0"
FROM ${KINDEST_REPOSITORY}thavlik/t4vd-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY compiler/pkg/api compiler/pkg/api
COPY gateway/pkg/api gateway/pkg/api
COPY sources/pkg/api sources/pkg/api
COPY hound hound
WORKDIR /go/src/github.com/thavlik/t4vd/hound/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o hound
FROM ${RUNTIME_BASE_IMAGE}
COPY --from=builder /go/src/github.com/thavlik/t4vd/hound/cmd/hound /usr/local/bin
CMD ["hound"]
