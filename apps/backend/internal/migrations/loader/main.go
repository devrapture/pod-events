//go:build atlas

package main

import (
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/devrapture/pod-events/internal/models"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		// list all your model structs here
		&models.User{},
		&models.SpotifyToken{},
		&models.NotificationChannel{},
	// add all model...
	)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
