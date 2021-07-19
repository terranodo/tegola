# Use Amazon Linux so tegola_linux CGO requirements are dynmically linked correctly
FROM amazonlinux:latest

RUN yum install -y go 

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
