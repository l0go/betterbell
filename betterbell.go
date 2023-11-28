package main

import (
    "log"
    "os"
    "time"
    "fmt"
    "github.com/gopxl/beep"
    "github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func main() {
    fmt.Println("Running")
    f, err := os.Open("audio/bell.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
