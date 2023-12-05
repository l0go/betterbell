# Betterbell

Rings a bell at specified intervals with a pretty UI. The goal is to eventually manage an instance with [Gokrazy](https://gokrazy.org/).

## Building

You need the following dependencies:

- Go 1.21+
- gcc

Run: `CGO_ENABLED=1 go build codeberg.org/logo/betterbell`. You can set the port with the glorious environment variable `PORT` and the location of the sqlite file with `DB_LOCATION`.
