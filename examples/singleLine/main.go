package main

import (
 "log"

 "github.com/alekssamos/edge-tts-go/edgeTTS"
)

func check(e error) {
 if e != nil {
  log.Fatal(e)
 }
}

func main() {
 //edgeTTS.PrintVoices("ru-RU")
 edgeTTS.NewTTS("result.mp3").AddTextWithVoice("Здарово!", "ru-RU-DmitryNeural").Speak()
}