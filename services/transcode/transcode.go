package transcode

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type ScaleConfig struct {
	dimention         string
	width             string
	height            string
	v_target_bit_rate string
	maxrate           string
	bufsize           string
	a_target_bit_rate string
}

func GenMasterPlaylist(fileName string) {
	fmt.Println("Create master playlist...")
	f, err := os.Create("/tmp/" + fileName + "/" + "master.m3u8")
	if err != nil {
		log.Fatalf("Some error: %v\n", err)
	}
	f.WriteString(`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360
hls360p/360p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=1400000,RESOLUTION=842x480
hls480p/480p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2800000,RESOLUTION=1280x720
hls720p/720p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080
hls1080p/1080p.m3u8`)
	f.Close()
	fmt.Println("Created master playlist")
}

func VideoToMultiBitrates(fileUrl string, fileName string) string {
	fmt.Println("Start transcoding")
	// cmd := exec.Command("mkdir", "/tmp/"+fileName)
	// _, err := cmd.Output()
	// if err != nil {
	// 	log.Fatalf("Some error: %v\n", err)
	// }
	scale_configs := [...]ScaleConfig{
		{
			dimention:         "360p",
			width:             "640",
			height:            "360",
			v_target_bit_rate: "800k",
			maxrate:           "800k",
			bufsize:           "1200k",
			a_target_bit_rate: "96k",
		},
		{
			dimention:         "480p",
			width:             "842",
			height:            "480",
			v_target_bit_rate: "1400k",
			maxrate:           "1498k",
			bufsize:           "2100k",
			a_target_bit_rate: "128k",
		},
		{
			dimention:         "720p",
			width:             "1280",
			height:            "720",
			v_target_bit_rate: "2800k",
			maxrate:           "2996k",
			bufsize:           "4200k",
			a_target_bit_rate: "128k",
		},
		{
			dimention:         "1080p",
			width:             "1920",
			height:            "1080",
			v_target_bit_rate: "5000k",
			maxrate:           "5350k",
			bufsize:           "7500k",
			a_target_bit_rate: "192k",
		},
	}
	base_args := []string{
		"-hide_banner",
		"-y",
		"-i",
		fileUrl,
	}

	var config_args []string
	for _, conf := range scale_configs {
		cmd := exec.Command("mkdir", "/tmp/"+fileName+"/hls"+conf.dimention)
		_, err := cmd.Output()
		if err != nil {
			log.Fatalf("Create dir error: %v\n", err)
		}
		append_args := []string{
			"-vf",
			"scale=w=" + conf.width + ":h=" + conf.height + ":force_original_aspect_ratio=decrease",
			"-c:a",
			"aac",
			"-ar",
			"48000",
			"-c:v",
			"h264",
			"-profile:v",
			"main",
			"-crf",
			"20",
			"-sc_threshold",
			"0",
			"-g",
			"48",
			"-keyint_min",
			"48",
			"-hls_time",
			"2",
			"-hls_playlist_type",
			"vod",
			"-b:v",
			conf.v_target_bit_rate,
			"-maxrate",
			conf.maxrate,
			"-bufsize",
			conf.bufsize,
			"-b:a",
			conf.a_target_bit_rate,
			"-hls_segment_filename",
			"/tmp/" + fileName + "/hls" + conf.dimention + "/" + conf.dimention + "_%03d.ts",
			"/tmp/" + fileName + "/hls" + conf.dimention + "/" + conf.dimention + ".m3u8",
		}
		config_args = append(config_args, append_args...)
	}
	final_args := append(base_args, config_args...)
	fmt.Println("Run ffmpeg")
	cmd := exec.Command("ffmpeg", final_args...)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Some error: %v\n", err)
	}
	fmt.Println("Transcoded")
	return "success"
}

func ConvertToMp4(fileUrl string, fileName string) string {
	fmt.Println("Start converting")
	cmd := exec.Command("mkdir", "/tmp/"+fileName)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Create dir error: %v\n", err)
	}
	// Should convert to 720p or 1080p
	savedFileUrl := "/tmp/" + fileName + "/" + fileName + ".mp4"
	fmt.Printf("Url: %s, FileName: %s\n, SavingUrl: %s\n", fileUrl, fileName, savedFileUrl)
	cmd = exec.Command("ffmpeg",
		"-i",
		fileUrl,
		"-fflags",
		"+genpts",
		"-r",
		"25",
		savedFileUrl)
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Convert error: %v\n", err)
	}
	return savedFileUrl
}

// Generate 10 secs preview clip
func GenShortClip(fileUrl string, fileName string) string {
	fmt.Println("Start generating short clip")
	name := fileName + "_preview"
	outpath := "/tmp/" + fileName + "/" + name + ".mp4"
	cmd := exec.Command("ffmpeg",
		"-i",
		fileUrl,
		"-ss",
		"00:00:0.0",
		"-t",
		"10",
		"-fflags",
		"+genpts",
		"-r",
		"25",
		outpath)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Convert error: %v\n", err)
	}
	return name
}

func GetVideoDuration(fileUrl string) string {
	fmt.Printf("Get the duration of the file: %s", fileUrl)
	cmd := exec.Command("ffprobe",
		"-v",
		"error",
		"-show_entries",
		"format=duration",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		fileUrl,
	)
	dur, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to get duration: %v\n", err)
	}
	return string(dur)
}
