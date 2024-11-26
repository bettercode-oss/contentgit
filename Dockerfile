FROM golang:1.23.3-bullseye AS builder
WORKDIR /go/src/contentgit
COPY . .
RUN go mod download
RUN go install -ldflags '-w -extldflags "-static"'

# make application docker image use alpine
FROM alpine:3.10
# using timezone
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Seoul
RUN apk add -U tzdata

WORKDIR /go/bin/
# copy config files to image
COPY --from=builder /go/src/contentgit/config/*.yaml ./config/
# copy execute file to image
COPY --from=builder /go/bin/contentgit .
EXPOSE 2022
CMD ["./contentgit"]
