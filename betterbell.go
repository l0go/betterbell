package main

import (
    "context"
    "log"
)

func main() {
    log.Println("Running")
    ctx, _ := context.WithCancel(context.Background())
    playMp3(ctx, "/etc/bell.mp3")
}
