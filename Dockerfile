#----
# Build stage
#----
FROM golang:1.12.4-alpine3.9 as builder

WORKDIR /$GOPATH/src/github.com/gcammarata/gc-simple-home

# Install git
RUN apk --update add git

# Install build required packages
RUN apk --update add make gcc tar rsync

ENV GO111MODULE=on

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
RUN go env

COPY . /MultiStage
WORKDIR /MultiStage

RUN CGO_ENABLED=0 GOOS=linux go build -o gcsh main.go

RUN ls
#----
# GCSH stage
#----

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /MultiStage/gcsh .
COPY --from=builder /MultiStage/conf ./conf
RUN ls

CMD ["./gcsh"]
