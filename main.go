package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alekssamos/edge-tts-go/edgeTTS"
	"github.com/spf13/pflag"
)

func usage() {
	fmt.Println("Microsoft Edge TTS")
	pflag.PrintDefaults()
}

func main() {
	executeAction := false
	listVoices := pflag.BoolP("list-voices", "l", false, "lists available voices and exits")
	locale := pflag.StringP("locale", "L", "", "locale for voice lists ex: zh-CN, en-US")
	text := pflag.StringP("text", "t", "", "what TTS will say")
	file := pflag.StringP("file", "f", "", "same as --text but read from file")
	voice := pflag.StringP("voice", "v", "zh-CN-XiaoxiaoNeural", "voice for TTS")
	volume := pflag.String("volume", "+0%", "set TTS volume")
	pitch := pflag.String("pitch", "+0Hz", "set TTS pitch")
	rate := pflag.String("rate", "+0%", "set TTS rate")
	writeMedia := pflag.StringP("write-media", "w", "", "send media output to file instead of stdout")
	// proxy := pflag.String("proxy", "", "use a proxy for TTS and voice list")
	pflag.Usage = usage
	pflag.Parse()

	if *listVoices {
		err := edgeTTS.PrintVoices(*locale)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *file != "" {
		executeAction = true
		if *file == "/dev/stdin" {
			reader := bufio.NewReader(os.Stdin)
			*text, _ = reader.ReadString('\n')
		} else {
			data, _ := os.ReadFile(*file)
			*text = string(data)
		}
	}
	if *text != "" {
		executeAction = true
		edgeTTS.NewTTS(*writeMedia).AddText(*text, *voice, *rate, *volume, *pitch).Speak()
	}
	if !executeAction {
		usage()
	}
}
