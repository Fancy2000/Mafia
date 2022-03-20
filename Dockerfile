FROM golang:1.16-alpine

WORKDIR go/src/hw5

COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY service ./

EXPOSE 5050
RUN mkdir service
RUN mv service.pb.go ./service
RUN mv service_grpc.pb.go ./service
RUN go mod download
RUN go build main.go
CMD ["./main"]