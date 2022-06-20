FROM registry.access.redhat.com/ubi8/ubi-minimal AS builder
RUN microdnf install -y golang-1.16.* && microdnf clean all

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
