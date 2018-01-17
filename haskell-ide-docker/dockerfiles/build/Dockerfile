FROM ubuntu:16.04

# curl
RUN apt-get update \
 && apt-get install curl -y

# stack
RUN echo "export PATH=~/.local/bin:\$PATH" >> ~/.bashrc
RUN curl -sSL https://get.haskellstack.org/ | sh

# build packages required for the vim configuration
CMD []

