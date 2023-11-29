package main

import (
    "io/ioutil"
    "os"
	"html/template"
	"log"
	"net/http"
	//"codeberg.org/logo/betterbell/audio"
	//	"os"

	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/yobert/alsa"
)

var tpl *template.Template

var cfg config

type config struct {
	Output          *alsa.Device
	PlaybackDevices []*alsa.Device
}

func main() {
	log.Println("Running")

	cfg.PlaybackDevices, _ = playbackDevice()
    err := beepDevice(cfg.PlaybackDevices[0])
    if err != nil {
        log.Fatal(err)
    }

	// HTTP
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/login", getLogin)
	http.HandleFunc("/select-device", postSelectDevice)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", cfg)
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "templates/login.html")
	}
}

func postSelectDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		for _, device := range cfg.PlaybackDevices {
			if device.Title == r.FormValue("devices") {
				cfg.Output = device
				break
			}
		}
	}
	beepDevice(cfg.Output)
}

func playbackDevice() ([]*alsa.Device, error) {
	cards, err := alsa.OpenCards()
	if err != nil {
		return nil, err
	}
	defer alsa.CloseCards(cards)

	var playback_device []*alsa.Device
	for _, card := range cards {
		devices, err := card.Devices()
		if err != nil {
			return nil, err
		}

		for _, device := range devices {
			if device.Type != alsa.PCM {
				continue
			}
			if device.Play {
				playback_device = append(playback_device, device)
			}
		}
	}

	return playback_device, nil
}

func beepDevice(device *alsa.Device) error {
	var err error
	if err = device.Open(); err != nil {
		return err
	}

    f, err := os.Open("/src/static/bell.mp3")
	if err != nil {
        return err
	}
    defer f.Close()

    dec, err := mp3.NewDecoder(f)
    if err != nil {
        return err
    }

	_, err = device.NegotiateFormat(alsa.S16_LE)
	if err != nil {
		return err
	}

	_, err = device.NegotiateChannels(2)
	if err != nil {
		return err
	}

	_, err = device.NegotiateRate(dec.SampleRate())
	if err != nil {
		return err
	}

    _, err = device.NegotiatePeriodSize(1024*2)
	if err != nil {
		return err
	}

    _, err = device.NegotiateBufferSize(2048*2)
    if err != nil {
		return err
	}

	if err = device.Prepare(); err != nil {
		return err
	}
    
	// Wait for playback to complete.
    buf, err := ioutil.ReadAll(dec)
	if err != nil {
		panic(err.Error())
	}

    device.Write(buf, int(dec.Length()) / device.BytesPerFrame())

	return nil
}
