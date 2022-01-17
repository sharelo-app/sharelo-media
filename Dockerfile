FROM golang:1.17
WORKDIR /go/src/github.com/sharelo-app/sharelo-media
RUN apt-get update
RUN apt-get install -y ffmpeg
RUN apt-get install -y s3cmd
COPY ./go.mod ./go.mod
RUN go mod download
COPY . .
CMD ["go", "run", "main/main.go"]