FROM registry.access.redhat.com/ubi8/ubi-minimal AS taskBuilder
RUN microdnf install -y tar gzip && microdnf clean all
ENV TASK_NAME=disk-virt-customize \
    GOFLAGS="-mod=vendor" \
    GO111MODULE=on
WORKDIR /src/${TASK_NAME}
RUN curl -L https://go.dev/dl/go1.18.1.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH=$PATH:/usr/local/go/bin
COPY . .
RUN	CGO_ENABLED=0 GOOS=linux go build -o /${TASK_NAME} cmd/${TASK_NAME}/main.go

FROM registry.access.redhat.com/ubi8/ubi:latest AS rhsrvanyBuilder
ENV TASK_NAME=disk-virt-customize
COPY build/${TASK_NAME}/repos/CentOS-Stream-rhsrvany.repo /etc/yum.repos.d/CentOS-Stream.repo
COPY build/${TASK_NAME}/repos/RPM-GPG-KEY-centosofficial /etc/pki/rpm-gpg/RPM-GPG-KEY-centosofficial
RUN yum install git make autoconf automake mingw32-gcc -y --disableplugin=subscription-manager
RUN git clone https://github.com/rwmjones/rhsrvany.git
WORKDIR /rhsrvany
RUN autoreconf --install && autoconf && mingw32-configure --disable-dependency-tracking && make

FROM quay.io/kubevirt/libguestfs-tools:v0.52.0
ENV TASK_NAME=disk-virt-customize
ENV ENTRY_CMD=/usr/local/bin/${TASK_NAME} \
    USER_UID=1001 \
    USER_NAME=${TASK_NAME} \
    HOME=/home/${TASK_NAME}

# install libguestfs rhsrvany.exe win dependency for virt-customize
COPY --from=rhsrvanyBuilder /rhsrvany/RHSrvAny/rhsrvany.exe /usr/share/virt-tools/rhsrvany.exe

# install task binary
COPY --from=taskBuilder /${TASK_NAME} ${ENTRY_CMD}
COPY build/${TASK_NAME}/bin /usr/local/bin

#RUN  /usr/local/bin/user_setup
ENTRYPOINT ["/usr/local/bin/entrypoint"]
CMD ["--help"]

USER ${USER_UID}
