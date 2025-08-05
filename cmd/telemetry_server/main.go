package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.eu.org/envloader"
	"telemetry.ex.gosuda.org/telemetry/internal/core/server"
)

var pid = os.Getpid()

func init() {
	go func() {
		for {
			ppid := os.Getppid()
			if ppid == 1 {
				// send sigterm to self
				cmd := exec.Command("kill", "-TERM", strconv.Itoa(pid))
				cmd.Run()
				time.Sleep(time.Second * 60)
				os.Exit(1)
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

func main() {
	envloader.LoadEnvFile(".env")

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("{\"port\":", ln.Addr().(*net.TCPAddr).Port, "}")

	log.Info().Msgf("Server starting on port %d", ln.Addr().(*net.TCPAddr).Port)

	server := server.NewServer()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := server.Serve(ln)
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Server failed")
		}
	}()

	<-sigCh
	log.Info().Msg("Received shutdown signal, shutting down server...")
	server.Shutdown()
	log.Info().Msg("Server shutdown complete")
}
