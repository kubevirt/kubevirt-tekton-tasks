FROM fedora:41

ENV HOME=/go/src/github.com/kubevirt/kubevirt-tekton-tasks/ \
    USER_UID=1001 \
    USERNAME=generatetasks


RUN curl -L https://go.dev/dl/go1.22.5.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH=$PATH:/usr/local/go/bin

WORKDIR ${HOME}
ADD . .
COPY automation/generatetasks/build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

RUN dnf install -y jq ansible make diffutils which git gcc && rm -rf /var/cache/yum /var/cache/dnf

#set permissions for ansible tmp folder
RUN mkdir -p /.ansible/tmp
RUN chmod -R 777 /.ansible

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
