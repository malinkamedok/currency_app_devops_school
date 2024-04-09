package main

import (
	"p.solovev/internal/app"
	"p.solovev/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		return
	}

	app.Run(cfg)
}
