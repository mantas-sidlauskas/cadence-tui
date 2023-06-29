package config

type Config struct {
	Address string `short:"a" long:"address" description:"Cadence server address" default:"localhost:7833"`
}
