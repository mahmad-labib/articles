package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Server struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

func ServerValues() Server {
	initServer := Server{}
	initServer.Port = os.Getenv("PORT")

	initServer.ReadTimeout = ServerTimeout("ReadTimeout", 5)
	initServer.WriteTimeout = ServerTimeout("WriteTimeout", 10)
	initServer.IdleTimeout = ServerTimeout("IdleTimeout", 120)

	return initServer
}

func ServerTimeout(envName string, defaultValue int) int {
	value := os.Getenv(envName)
	if value != "" {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			slog.Error("Error during conversion:", "Msg", err)
			os.Exit(1)
		}
		return valueInt
	}
	return defaultValue
}
