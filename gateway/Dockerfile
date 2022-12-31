ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="alpine:3.17.0"
FROM ${KINDEST_REPOSITORY}thavlik/t4vd-base:${KINDEST_TAG} AS builder
COPY base/cmd/iam base/cmd/iam
COPY base/pkg base/pkg
COPY sources/pkg/api sources/pkg/api
COPY seer/pkg/api seer/pkg/api
COPY seer/pkg/infocache/infocache.go seer/pkg/infocache/infocache.go
COPY sources/pkg/store/store.go sources/pkg/store/store.go
COPY compiler/pkg/api compiler/pkg/api
COPY compiler/pkg/datastore/datastore.go compiler/pkg/datastore/datastore.go
COPY filter/pkg/api filter/pkg/api
COPY slideshow/pkg/api slideshow/pkg/api
COPY gateway gateway
WORKDIR /go/src/github.com/thavlik/t4vd/gateway/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gateway
FROM ${RUNTIME_BASE_IMAGE}
COPY --from=builder /go/src/github.com/thavlik/t4vd/gateway/cmd/gateway /usr/local/bin
CMD ["gateway"]
