FROM docker.io/golang:1.22
WORKDIR /go/src/app
COPY . .
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
RUN go mod init app
RUN go get -d -v ./...
