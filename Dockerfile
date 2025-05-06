ARG VERSION="dev"
ARG GO_BUILDER=golang:1.24.2
ARG RUNTIME=gcr.io/distroless/static-debian12:nonroot


FROM $GO_BUILDER AS builder
WORKDIR /go/src/github.com/mcp-tekton
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=${VERSION} -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -mod=vendor -o /tmp/tekton-mcp-server ./cmd/tekton-mcp-server

FROM $RUNTIME
COPY --from=builder /tmp/tekton-mcp-server /usr/local/bin/tekton-mcp-server 

CMD ["/usr/local/bin/tekton-mcp-server", "stdio"]

