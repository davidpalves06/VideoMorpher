FROM golang:1.23.4-bullseye AS build

WORKDIR /app

COPY ./internal/ ./internal
COPY ./cmd/ ./cmd
COPY ./go.mod .

RUN go build -o videomorpher cmd/videomorpher/videomorpher.go

FROM ubuntu:25.04

WORKDIR /videomorpher

RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get install -y ffmpeg

COPY config.json .
COPY ./static/ ./static

COPY --from=build /app/videomorpher videomorpher

ENTRYPOINT [ "./videomorpher"]