package main

import (
	"log/slog"

	"github.com/art6f/goblncr/internal/app"
)

func main() {
	slog.Info("LB Starting...")

	app.Main()
}
