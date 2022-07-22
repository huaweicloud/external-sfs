# Based on centos
FROM centos:7.6.1810
LABEL maintainers="Kubernetes Authors"
LABEL description="SFS Provisioner"

# Copy from build directory
COPY sfs-provisioner /sfs-provisioner

# Grant execute permission
RUN chmod +x /sfs-provisioner

# Update
RUN yum -y update

# Define default command
ENTRYPOINT ["/sfs-provisioner"]
