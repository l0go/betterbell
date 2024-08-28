package internal

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

func Ring(id int, peers *PeerState) error {
	log.Printf("%d: Ring!", id)

	f, err := os.Open("./static/bell.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	
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

	go func() error {
		for _, peer := range peers.Get() {
			url, err := url.Parse(peer.Endpoint)
			if err != nil {
				return err
			}

			_, err = http.Head(fmt.Sprintf("%s://%s/ring/?id=%d&secret=%s", url.Scheme, url.Host, id, peer.Secret))
			if err != nil {
				return err
			}
		}
		return nil
	}()

	<-done

	return nil
}
