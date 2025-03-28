package edgeTTS

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	terminal "golang.org/x/term"
)

type EdgeTTS struct {
	communicator *Communicate
	texts        []*CommunicateTextTask
	outCome      io.WriteCloser
	LastError    error
}

type Args struct {
	Text           string
	Voice          string
	Proxy          string
	Rate           string
	Volume         string
	Pitch          string
	WordsInCue     float64
	WriteMedia     string
	WriteSubtitles string
}

func isTerminal(file *os.File) bool {
	return terminal.IsTerminal(int(file.Fd()))
}

func PrintVoices(locale string) error {
	// Print all available voices.
	voices, err := listVoices()
	if err != nil {
		return err
	}
	sort.Slice(voices, func(i, j int) bool {
		return voices[i].ShortName < voices[j].ShortName
	})

	filterFieldName := map[string]bool{
		"SuggestedCodec": true,
		"FriendlyName":   true,
		"Status":         true,
		"VoiceTag":       true,
		"Language":       true,
	}

	printed := false
	for _, voice := range voices {
		if locale != "" && strings.ToLower(voice.Locale) != strings.ToLower(locale) {
			continue
		}
		fmt.Printf("\n")
		t := reflect.TypeOf(voice)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldName := field.Name
			if filterFieldName[fieldName] {
				continue
			}
			fieldValue := reflect.ValueOf(voice).Field(i).Interface()
			fmt.Printf("%s: %v\n", fieldName, fieldValue)
			printed = true
		}
	}
	if !printed {
		return fmt.Errorf("No matching voice with locale %q", locale)
	}
	return nil
}

func NewTTS(writeMedia string) *EdgeTTS {
	args := &Args{
		WriteMedia: writeMedia,
		Voice: "zh-CN-XiaoxiaoNeural",
		Volume: "+0%",
		Pitch: "+0Hz",
		Rate: "+0%",
	}
	if isTerminal(os.Stdin) && isTerminal(os.Stdout) && args.WriteMedia == "" {
		fmt.Fprintln(os.Stderr, "Warning: TTS output will be written to the terminal. Use --write-media to write to a file.")
		fmt.Fprintln(os.Stderr, "Press Ctrl+C to cancel the operation. Press Enter to continue.")
		_, err := fmt.Scanln()
		if err != nil {
			os.Exit(1)
		}
	}
	if _, err := os.Stat(args.WriteMedia); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(args.WriteMedia), 0755)
		if err != nil {
			return &EdgeTTS{LastError: fmt.Errorf("Failed to create dir: %w", err)}
		}
	}
	tts := NewCommunicate().WithVoice(args.Voice).WithRate(args.Rate).WithVolume(args.Volume).WithPitch(args.Pitch)
	file, err := os.OpenFile(args.WriteMedia, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return &EdgeTTS{LastError: fmt.Errorf("Failed to open file: %w", err)}
	}
	tts.openWs()
	return &EdgeTTS{
		communicator: tts,
		outCome:      file,
		texts:        []*CommunicateTextTask{},
	}
}

func (eTTS *EdgeTTS) task(text string, voice string, rate string, volume string, pitch string) *CommunicateTextTask {
	return &CommunicateTextTask{
		text: text,
		option: CommunicateTextOption{
			voice:  voice,
			rate:   rate,
			volume: volume,
			pitch:  pitch,
		},
	}
}

func (eTTS *EdgeTTS) AddTextDefault(text string) *EdgeTTS {
	eTTS.texts = append(eTTS.texts, eTTS.task(text, "", "", "", ""))
	return eTTS
}

func (eTTS *EdgeTTS) AddTextWithVoice(text string, voice string) *EdgeTTS {
	eTTS.texts = append(eTTS.texts, eTTS.task(text, voice, "", "", ""))
	return eTTS
}

func (eTTS *EdgeTTS) AddText(text string, voice string, rate string, volume string, pitch string) *EdgeTTS {
	eTTS.texts = append(eTTS.texts, eTTS.task(text, voice, rate, volume, pitch))
	return eTTS
}

func (eTTS *EdgeTTS) Speak() error {
	defer eTTS.communicator.close()
	defer eTTS.outCome.Close()

	go eTTS.communicator.allocateTask(eTTS.texts)
	eTTS.communicator.createPool()
	for _, text := range eTTS.texts {
		eTTS.outCome.Write(text.speechData)
	}
	if eTTS.LastError == nil {
		eTTS.LastError = eTTS.communicator.lastError
	}
	return eTTS.communicator.lastError
}
