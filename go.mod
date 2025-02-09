module github.com/surfaceyu/edge-tts-go

go 1.20

require golang.org/x/crypto v0.11.0

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/term v0.10.0 // indirect
	golang.org/x/text v0.11.0
)

replace github.com/surfaceyu/edge-tts-go v0.1.0 => github.com/alekssamos/edge-tts-go v0.1.1
