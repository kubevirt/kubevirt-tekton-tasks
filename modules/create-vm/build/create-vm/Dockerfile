FROM golang:1.14 AS builder
WORKDIR /src/create-vm
COPY . .
ENV GOFLAGS="-mod=readonly"
ENV GO111MODULE=on
RUN	CGO_ENABLED=0 GOOS=linux go build -o /create-vm cmd/create-vm/main.go

FROM registry.access.redhat.com/ubi7/ubi-minimal:latest
ENV CREATE_VM=/usr/local/bin/create-vm \
    USER_UID=1001 \
    USER_NAME=create-vm

# install create-vm binary
COPY --from=builder /create-vm ${CREATE_VM}
COPY build/create-vm/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]
CMD ["--help"]

USER ${USER_UID}
