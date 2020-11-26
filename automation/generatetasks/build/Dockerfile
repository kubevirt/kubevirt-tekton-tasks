FROM fedora:31

ENV HOME=/go/src/github.com/kubevirt/kubevirt-tekton-tasks/ \
    USER_UID=1001 \
    USER_NAME=generatetasks

WORKDIR ${HOME}
ADD . .
COPY automation/generatetasks/build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

RUN dnf install -y jq ansible make diffutils && rm -rf /var/cache/yum /var/cache/dnf

#set permissions for ansible tmp folder
RUN mkdir -p /.ansible/tmp
RUN chmod -R 777 /.ansible

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
