FROM golang:1.12.9
RUN go get -d -v github.com/docker/docker
WORKDIR /scope
ADD . /scope
RUN cd /scope && go build
ENTRYPOINT ./main
