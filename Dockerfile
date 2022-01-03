FROM golang:1.17
WORKDIR /go/src/github.com/sharelo-app/sharelo-media
COPY . .
RUN go mod download
RUN sudo apt-get update
RUN sudo add-apt-repository ppa:jonathonf/ffmpeg-4
RUN sudo apt install ffmpeg
RUN sudo apt-get install s3cmd