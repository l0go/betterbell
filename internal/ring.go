package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gen2brain/malgo"
	mp3 "github.com/hajimehoshi/go-mp3"
)

func Ring(id int, peers *PeerState) error {
	log.Printf("%d: Ring!", id)

	f, err := os.Open("./static/bell.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	dec, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	// Initialize malgo
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = uint32(dec.SampleRate())
	deviceConfig.Alsa.NoMMap = 1

	onSamples := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		io.ReadFull(dec, pOutputSample)
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onSamples,
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		return err
	}
	defer device.Uninit()

	err = device.Start()
	if err != nil {
		return err
	}

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

	// TODO: Make this less hardcoded
	time.Sleep(4 * time.Second)

	return nil
}
