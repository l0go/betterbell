package audio

import (
	"context"
	"fmt"
    "io"
	"os"

	"github.com/anisse/alsa"
	mp3 "github.com/hajimehoshi/go-mp3"
)

func PlayMP3(ctx context.Context, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open %s: %s", filename, err.Error())
	}
	defer f.Close()

    defer fmt.Println(f)

	dec, err := mp3.NewDecoder(f)
	if err != nil {
		return fmt.Errorf("could not decode %s: %s", filename, err.Error())
	}

	sampleRate := dec.SampleRate()
	p, err := alsa.NewPlayer(sampleRate, 2, 2, 4096)
	if err != nil {
		return fmt.Errorf("could not init alsa: %s", err.Error())
	}
	defer p.Close()

	_, err = copyCtx(ctx, p, dec)

	return err
}

// Code below from http://ixday.github.io/post/golang-cancel-copy/
type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

func copyCtx(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	n, err := io.Copy(dst, readerFunc(func(p []byte) (int, error) {

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return src.Read(p)
		}
	}))
	return n, err
}
