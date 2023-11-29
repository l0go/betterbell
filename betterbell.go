package main

import (
    "log"
    "net/http"
    "html/template"
    //"codeberg.org/logo/betterbell/audio"
//	"os"
    
//	mp3 "github.com/hajimehoshi/go-mp3"
    "github.com/yobert/alsa"
)

var tpl *template.Template

var cfg config
type config struct {
    Output *alsa.Device
    PlaybackDevices []*alsa.Device
}

func main() {
    log.Println("Running")

    cfg.PlaybackDevices, _ = playbackDevice()

    // HTTP
    var err error
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

type name struct {
    Name string
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
