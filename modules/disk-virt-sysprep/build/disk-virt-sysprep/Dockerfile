FROM registry.access.redhat.com/ubi8/ubi-minimal AS taskBuilder
RUN microdnf install -y golang-1.15.* && microdnf clean all
ENV TASK_NAME=disk-virt-sysprep \
    GOFLAGS="-mod=vendor" \
    GO111MODULE=on
WORKDIR /src/${TASK_NAME}
COPY . .
RUN	CGO_ENABLED=0 GOOS=linux go build -o /${TASK_NAME} cmd/${TASK_NAME}/main.go

FROM registry.access.redhat.com/ubi8/ubi:latest AS rhsrvanyBuilder
ENV TASK_NAME=disk-virt-sysprep
COPY build/${TASK_NAME}/repos/CentOS-Stream-rhsrvany.repo /etc/yum.repos.d/CentOS-Stream.repo
COPY build/${TASK_NAME}/repos/RPM-GPG-KEY-centosofficial /etc/pki/rpm-gpg/RPM-GPG-KEY-centosofficial
RUN yum install git make autoconf automake mingw32-gcc -y --disableplugin=subscription-manager
RUN git clone https://github.com/rwmjones/rhsrvany.git
WORKDIR /rhsrvany
RUN autoreconf --install && autoconf && mingw32-configure --disable-dependency-tracking && make

FROM registry.access.redhat.com/ubi8/ubi:latest
ENV TASK_NAME=disk-virt-sysprep
ENV ENTRY_CMD=/usr/local/bin/${TASK_NAME} \
    USER_UID=1001 \
    USER_NAME=${TASK_NAME} \
    HOME=/home/${TASK_NAME} \
    LIBGUESTFS_APPLIANCE_VERSION=1.40.1


# HACK: CI injects centos7 images instead of ubi, so setup repos accordingly
COPY build/${TASK_NAME}/repos/CentOS-Stream.repo /etc/yum.repos.d/CentOS-Stream.repo
COPY build/${TASK_NAME}/repos/RPM-GPG-KEY-centosofficial /etc/pki/rpm-gpg/RPM-GPG-KEY-centosofficial
RUN if grep -iq centos /etc/*release; then \
        rm -rf /etc/yum.repos.d/CentOS-Stream.repo /etc/pki/rpm-gpg/RPM-GPG-KEY-centosofficial; \
    else \
        find /etc/yum.repos.d/ -maxdepth 1 ! -name CentOS-Stream.repo -type f -exec rm -f {} +; \
    fi

# xz will be used to unpack libguestfs appliance later when running the container
RUN yum install xz libguestfs-tools-c -y --disableplugin=subscription-manager && \
    yum clean all --disableplugin=subscription-manager && \
    rm -rf /etc/yum.repos.d/* /etc/pki/rpm-gpg/RPM-GPG-KEY-centosofficial /var/cache/yum /var/cache/dnf /var/lib/rpm
RUN mkdir /data && curl https://download.libguestfs.org/binaries/appliance/appliance-${LIBGUESTFS_APPLIANCE_VERSION}.tar.xz --output /data/appliance.tar.xz

# install libguestfs rhsrvany.exe win dependency for virt-sysprep
COPY --from=rhsrvanyBuilder /rhsrvany/RHSrvAny/rhsrvany.exe /usr/share/virt-tools/rhsrvany.exe

# install task binary
COPY --from=taskBuilder /${TASK_NAME} ${ENTRY_CMD}
COPY build/${TASK_NAME}/bin /usr/local/bin

RUN  /usr/local/bin/user_setup
ENTRYPOINT ["/usr/local/bin/entrypoint"]
CMD ["--help"]

USER ${USER_UID}
