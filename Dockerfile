FROM golang:1.13 AS builder

WORKDIR /go/src/ITLab-Reports/
ENV CGO_ENABLED=0
COPY ./src/ITLabReports/api/go.* ./
RUN go mod download
COPY ./src/ITLabReports/api .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=builder /go/src/ITLab-Reports/main .
CMD ["./main"]
