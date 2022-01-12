package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sharelo-app/sharelo-media/services/pipeline"
	"github.com/sharelo-app/sharelo-media/services/transcode"
	"github.com/sharelo-app/sharelo-media/services/uploader"
)

type FileInfo struct {
	UserId        string
	VideoUploadId string
	Url           string
	FileName      string
}

type TranscodeInfo struct {
	UserId        string `json:"user_id"`
	VideoUploadId string `json:"video_upload_id"`
	TranscodedUrl string `json:"transcoded_url"`
	DownloadSize  string `json:"download_size"`
	StreamUrl     string `json:"stream_url"`
	PreviewUrl    string `json:"preview_url"`
	Duration      string `json:"duration"`
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func GetFileSize(fileUrl string) (string, error) {
	fi, err := os.Stat(fileUrl)
	if err != nil {
		return "0", err
	}
	inStr := strconv.FormatInt(fi.Size(), 10)
	return inStr, nil
}

func main() {
	conn, err := amqp.Dial("amqp://shareloadmin:shareloadmin@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"transcodings_queue", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var fileInfo FileInfo
			body := d.Body
			json.Unmarshal(body, &fileInfo)
			fmt.Printf("Url: %s, FileName: %s\n", fileInfo.Url, fileInfo.FileName)
			savedFileUrl := transcode.ConvertToMp4(fileInfo.Url, fileInfo.FileName)
			duration := transcode.GetVideoDuration(savedFileUrl)
			size, err := GetFileSize(savedFileUrl)
			if err != nil {
				FailOnError(err, "Failed to get file size")
			}
			fmt.Printf("converted: %s\n", savedFileUrl)
			transcode.GenShortClip(savedFileUrl, fileInfo.FileName)
			transcode.VideoToMultiBitrates(savedFileUrl, fileInfo.FileName)
			transcode.GenMasterPlaylist(fileInfo.FileName)
			uploader.UploadDir(fileInfo.UserId, fileInfo.FileName)
			fmt.Println("Finished")
			streamUrl := uploader.GetStreamUrl(fileInfo.UserId, fileInfo.FileName)
			transcodedUrl := uploader.GetTranscodedUrl(fileInfo.UserId, fileInfo.FileName)
			previewUrl := uploader.GetPreviewUrl(fileInfo.UserId, fileInfo.FileName)
			info := &TranscodeInfo{
				UserId:        fileInfo.UserId,
				VideoUploadId: fileInfo.VideoUploadId,
				TranscodedUrl: transcodedUrl,
				DownloadSize:  size,
				StreamUrl:     streamUrl,
				PreviewUrl:    previewUrl,
				Duration:      duration,
			}
			infoJson, err := json.Marshal(info)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			pipeline.PublishAfterTranscoded(infoJson)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
