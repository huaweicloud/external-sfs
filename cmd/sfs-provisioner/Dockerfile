# Based on centos
FROM centos:7.4.1708
LABEL maintainers="Kubernetes Authors"
LABEL description="SFS Provisioner"

# Copy from build directory
COPY sfs-provisioner /sfs-provisioner

# Update
RUN yum -y update

# Define default command
ENTRYPOINT ["/sfs-provisioner"]
