FROM registry.access.redhat.com/ubi9/go-toolset:1.18.9 AS builder
ENV TASK_NAME=modify-data-object \
    GOFLAGS="-mod=vendor"
WORKDIR /src/${TASK_NAME}
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /${TASK_NAME} cmd/${TASK_NAME}/main.go

FROM quay.io/centos/centos:stream9
ENV TASK_NAME=modify-data-object
ENV ENTRY_CMD=/usr/local/bin/${TASK_NAME}

RUN dnf install -y xorriso
# install task binary
COPY --from=builder /${TASK_NAME} ${ENTRY_CMD}

ENTRYPOINT ["/usr/local/bin/${TASK_NAME}"]
CMD ["--help"]
