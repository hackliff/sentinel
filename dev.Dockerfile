FROM golang:1.5.1
MAINTAINER Xavier Bruhiere

ENV TERM xterm
# fix https://github.com/Masterminds/glide/issues/135
ENV GLIDE_HOME /root

# -- Dev tools
# project vendoring
ENV GO15VENDOREXPERIMENT 1
ENV GLIDE_VERSION 0.7.0
ENV TARGET glide-${GLIDE_VERSION}-linux-amd64
RUN curl -LkOs "https://github.com/Masterminds/glide/releases/download/${GLIDE_VERSION}/${TARGET}.tar.gz" && \
  tar xvzf ${TARGET}.tar.gz && \
  mv linux-amd64/glide $GOPATH/bin/ && \
  rm -r *linux-amd64*

# code linting
RUN go get github.com/alecthomas/gometalinter && \
  gometalinter --install --update

# A dead simple, no frills Go cross compile tool
RUN go get github.com/mitchellh/gox

RUN git config --global user.name hackliff

# godoc
EXPOSE 6060

CMD ["go"]
