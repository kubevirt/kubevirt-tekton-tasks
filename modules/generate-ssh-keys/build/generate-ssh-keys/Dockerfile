FROM golang:1.15 AS builder
ENV TASK_NAME=generate-ssh-keys \
    GOFLAGS="-mod=vendor" \
    GO111MODULE=on
WORKDIR /src/${TASK_NAME}
COPY . .
RUN	CGO_ENABLED=0 GOOS=linux go build -o /${TASK_NAME} cmd/${TASK_NAME}/main.go

FROM registry.access.redhat.com/ubi8/ubi:latest
ENV TASK_NAME=generate-ssh-keys
ENV ENTRY_CMD=/usr/local/bin/${TASK_NAME} \
    USER_UID=1001 \
    USER_NAME=${TASK_NAME} \
    HOME=/home/${TASK_NAME}

# install ssh client
RUN yum update -y --disableplugin=subscription-manager --exclude=filesystem && \
    yum install openssh-clients -y --disableplugin=subscription-manager && \
    yum clean all --disableplugin=subscription-manager && \
    rm -rf /var/cache/yum /var/cache/dnf /var/lib/rpm

# install task binary
COPY --from=builder /${TASK_NAME} ${ENTRY_CMD}
COPY build/${TASK_NAME}/bin /usr/local/bin

RUN  /usr/local/bin/user_setup
ENTRYPOINT ["/usr/local/bin/entrypoint"]
CMD ["--help"]

USER ${USER_UID}
