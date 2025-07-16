# Args
ARG BASE_IMG="quay.io/samba.org/samba-server:latest"
ARG GIT_VERSION="(unset)"
ARG COMMIT_ID="(unset)"
ARG ARCH=""

# Build smbmetrics
FROM docker.io/golang:1.24 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd cmd
COPY internal internal

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} GO111MODULE=on \
    go build -a \
    -ldflags "-X main.Version=${GIT_VERSION} -X main.CommitID=${COMMIT_ID}" \
    -o smbmetrics cmd/main.go

# Use samba-server (with its smb.conf and samba utils) as base image
FROM $BASE_IMG
COPY --from=builder /workspace/smbmetrics /bin/smbmetrics

ENTRYPOINT ["/bin/smbmetrics"]
EXPOSE 9922
