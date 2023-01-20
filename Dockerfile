FROM  golang:1.18-stretch AS builder
ADD . /src
WORKDIR /src
RUN cd /src && echo $(ls -1 /src)
RUN go mod download
RUN go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" .
CMD ["./uploader"]