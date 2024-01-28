package main

import (
	"Distributed/loadConfig"
	"log/slog"
)

func main() {
	loadConfig.LogInit()

	slog.Debug(
		"executing database query",
		slog.String("query", "SELECT * FROM users %d"),
	)

	select {}

}

type MyHandler struct {
	slog.Handler
}
