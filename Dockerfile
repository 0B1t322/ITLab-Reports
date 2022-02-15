FROM golang:1.13 AS builder

WORKDIR /go/src/ITLab-Reports/
ENV CGO_ENABLED=0
COPY ./src/ITLabReports/api/go.* ./
RUN go mod download
COPY ./src/ITLabReports/api .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

WORKDIR /app
COPY --from=builder /go/src/ITLab-Reports/main .
CMD ["./main"]
