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

	// Cleanup device when done or force cleanup after 3 seconds.
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	childCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	go func(ctx context.Context) {
		defer device.Close()
		<-ctx.Done()
		log.Println("Closing device.")
		wg.Done()
	}(childCtx)

	channels, err := device.NegotiateChannels(1, 2)
	if err != nil {
		return err
	}

	rate, err := device.NegotiateRate(44100)
	if err != nil {
		return err
	}

	format, err := device.NegotiateFormat(alsa.S16_LE, alsa.S32_LE)
	if err != nil {
		return err
	}

	// A 50ms period is a sensible value to test low-ish latency.
	// We adjust the buffer so it's of minimal size (period * 2) since it appear ALSA won't
	// start playback until the buffer has been filled to a certain degree and the automatic
	// buffer size can be quite large.
	// Some devices only accept even periods while others want powers of 2.
	wantPeriodSize := 2048 // 46ms @ 44100Hz

	periodSize, err := device.NegotiatePeriodSize(wantPeriodSize)
	if err != nil {
		return err
	}

	bufferSize, err := device.NegotiateBufferSize(wantPeriodSize * 2)
	if err != nil {
		return err
	}

	if err = device.Prepare(); err != nil {
		return err
	}

	log.Printf("Negotiated parameters: %d channels, %d hz, %v, %d period size, %d buffer size\n",
		channels, rate, format, periodSize, bufferSize)

	dec, err := mp3.NewDecoder(f)
	if err != nil {
		return fmt.Errorf("could not decode %s: %s", filename, err.Error())
	}

    
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open %s: %s", filename, err.Error())
	}

	defer f.Close()
	var buf bytes.Buffer
	dec, err := mp3.NewDecoder(f)
	if err != nil {
		return fmt.Errorf("could not decode %s: %s", filename, err.Error())
	}

	// Play 2 seconds of beep.
	duration := 2 * time.Second
	t := time.NewTimer(duration)
	for t := 0.; t < duration.Seconds(); {
		var buf bytes.Buffer

		for i := 0; i < periodSize; i++ {
			v := math.Sin(t * 2 * math.Pi * 440) // A4
			v *= 0.1                             // make a little quieter

			switch format {
			case alsa.S16_LE:
				sample := int16(v * math.MaxInt16)

				for c := 0; c < channels; c++ {
					binary.Write(&buf, binary.LittleEndian, sample)
				}

			case alsa.S32_LE:
				sample := int32(v * math.MaxInt32)

				for c := 0; c < channels; c++ {
					binary.Write(&buf, binary.LittleEndian, sample)
				}

			default:
                log.Fatalf("Unhandled sample format: %v", format)
			}

			t += 1 / float64(rate)
		}

		if err := device.Write(buf.Bytes(), periodSize); err != nil {
			return err
		}
	}
	// Wait for playback to complete.
	<-t.C
	log.Printf("Playback should be complete now.\n")

	return nil
}
