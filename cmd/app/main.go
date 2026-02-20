package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mohamoundaljadan/screenshot/internal/app"
)

const version = "0.1.0-dev"

func main() {
	svc, err := buildService()
	if err != nil {
		fmt.Fprintln(os.Stderr, "startup failed:", err)
		os.Exit(1)
	}

	state, err := svc.GetAppState()
	if err != nil {
		fmt.Fprintln(os.Stderr, "state failed:", err)
		os.Exit(1)
	}
	b, _ := json.MarshalIndent(state, "", "  ")
	fmt.Println("go-wails-shot backend ready")
	fmt.Println(string(b))
	fmt.Println("Install Wails CLI and bind this service into a desktop shell.")
}

func buildService() (*app.Service, error) {
	return app.NewDefaultService(version)
}
