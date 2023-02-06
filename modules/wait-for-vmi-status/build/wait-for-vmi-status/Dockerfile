FROM registry.access.redhat.com/ubi9/go-toolset:1.18.9 AS builder
ENV TASK_NAME=wait-for-vmi-status \
    GOFLAGS="-mod=vendor" \
    GO111MODULE=on
WORKDIR /src/${TASK_NAME}
COPY . .
RUN	CGO_ENABLED=0 GOOS=linux go build -o /${TASK_NAME} cmd/${TASK_NAME}/main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
ENV TASK_NAME=wait-for-vmi-status
ENV ENTRY_CMD=/usr/local/bin/${TASK_NAME} \
    USER_UID=1001 \
    USER_NAME=${TASK_NAME}

# install task binary
COPY --from=builder /${TASK_NAME} ${ENTRY_CMD}
COPY build/${TASK_NAME}/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
