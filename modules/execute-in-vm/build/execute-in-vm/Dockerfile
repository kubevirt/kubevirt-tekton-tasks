FROM registry.access.redhat.com/ubi9/ubi-minimal AS builder
RUN microdnf install -y  tar gzip && microdnf clean all
ENV TASK_NAME=execute-in-vm \
    GOFLAGS="-mod=vendor" \
    GO111MODULE=on
WORKDIR /src/${TASK_NAME}
RUN curl -L https://go.dev/dl/go1.20.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH=$PATH:/usr/local/go/bin
COPY . .
RUN	CGO_ENABLED=0 GOOS=linux go build -o /${TASK_NAME} cmd/${TASK_NAME}/main.go

FROM registry.access.redhat.com/ubi9/ubi:latest
ENV TASK_NAME=execute-in-vm
ENV ENTRY_CMD=/usr/local/bin/${TASK_NAME} \
    USER_UID=1001 \
    USER_NAME=${TASK_NAME} \
    HOME=/home/${TASK_NAME}

# install ssh client
RUN yum update -y --disableplugin=subscription-manager && \
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
