package internal

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gen2brain/malgo"
	mp3 "github.com/hajimehoshi/go-mp3"
)

func Ring() error {
	log.Println("Ring!")

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

	time.Sleep(4 * time.Second)

	return nil
}
