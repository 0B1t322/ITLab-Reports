FROM golang:1.18 AS builder

WORKDIR /go/src/ITLabReports/

ENV CGO_ENABLED=0
COPY go.* ./
# Copy submodule
COPY ./pkg/ITLab/go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main/main.go  

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /go/src/ITLabReports/main .
CMD ["./main"]