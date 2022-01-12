package uploader

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/sharelo-app/sharelo-media/services/config"
)

func GetBaseRemotePath(userId string) string {
	config, err := config.ReadConfigFile("")
	if err != nil {
		log.Fatal("Unable to read the config file: ", err)
	}
	return "http://127.0.0.1:9000/" + config.WasabiConfig.Bucket + "/users/" + userId
}

func GetStreamUrl(userId string, fileName string) string {
	basePath := GetBaseRemotePath(userId)
	return basePath + "/output/" + fileName + "/master.m3u8"
}

func GetTranscodedUrl(userId string, fileName string) string {
	basePath := GetBaseRemotePath(userId)
	return basePath + "/output/" + fileName + "/" + fileName + ".mp4"
}

func GetPreviewUrl(userId string, fileName string) string {
	basePath := GetBaseRemotePath(userId)
	return basePath + "/output/" + fileName + "/" + fileName + "_preview.mp4"
}

func UploadDir(userId string, fileName string) {
	config, err := config.ReadConfigFile("")
	if err != nil {
		log.Fatal("Unable to read the config file: ", err)
		return
	}
	fromDir := "/tmp/" + fileName
	fmt.Printf("fromDir: %s\n", fromDir)
	toDir := "s3://" + config.WasabiConfig.Bucket + "/users/" + userId + "/output/"
	fmt.Printf("toDir: %s\n", toDir)
	cmd := exec.Command("s3cmd", "sync", fromDir, toDir, "-c", "./.s3cfg") // Relative path from main.go
	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	log.Fatalf("S3 stdout error: %v\n", err)
	// }
	// stderr, err := cmd.StderrPipe()
	// if err != nil {
	// 	log.Fatalf("S3 stderr error: %v\n", err)
	// }
	err = cmd.Run()
	if err != nil {
		log.Fatalf("S3 start error: %v\n", err)
	}
	os.RemoveAll(fromDir)
	// in := bufio.NewScanner(stdout)
	// for in.Scan() {
	// 	log.Printf(in.Text()) // write each line to your log, or anything you need
	// }
	// if err := in.Err(); err != nil {
	// 	log.Printf("error: %s", err)
	// }
	// in_err := bufio.NewScanner(stderr)
	// for in_err.Scan() {
	// 	log.Printf(in_err.Text()) // write each line to your log, or anything you need
	// }
}
