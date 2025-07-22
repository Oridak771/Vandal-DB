# Build the masking job binary
FROM golang:1.18 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/masking-job/main.go cmd/masking-job/main.go
COPY masking/ masking/
COPY apis/ apis/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o masking-job cmd/masking-job/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/masking-job .
USER 65532:65532

ENTRYPOINT ["/masking-job"]
