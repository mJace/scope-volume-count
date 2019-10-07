FROM golang:1.12.9
WORKDIR /scope
ADD . /scope
RUN cd /scope && go build
ENTRYPOINT ./main
