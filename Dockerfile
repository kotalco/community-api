FROM golang:1.19-alpine AS builder

RUN apk add --no-cache curl tar
ENV K8S_VERSION=1.21.2
# download kubebuilder tools required by envtest
RUN curl -sSLo envtest-bins.tar.gz "https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-${K8S_VERSION}-$(go env GOOS)-$(go env GOARCH).tar.gz"
RUN tar -vxzf envtest-bins.tar.gz -C /usr/local/

WORKDIR /api

COPY . .

RUN CGO_ENABLED=0 go build -o server

FROM alpine
ENV ETCD_UNSUPPORTED_ARCH=arm64
# Add new non-root user 'kotal'
RUN adduser -D kotal
USER kotal
WORKDIR /home/kotal

# required by api server to determine config/crds path
ENV GOPATH=/go
COPY --from=builder /go/pkg/mod/github.com/kotalco /go/pkg/mod/github.com/kotalco
# tools (etcd, apiserver, and kubectl) required by envtest
COPY --from=builder /usr/local/kubebuilder /usr/local/kubebuilder
COPY --from=builder /api/server /home/kotal/api/server

EXPOSE 5000
ENTRYPOINT [ "./api/server" ]
