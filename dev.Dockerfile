FROM golang:1.5.1
MAINTAINER Xavier Bruhiere

ENV TERM xterm
# needed for ghr to work
RUN git config --global user.name hackliff

ADD Makefile /tmp/Makefile
RUN cd /tmp && make install.tools

# godoc
EXPOSE 6060

CMD ["go"]
