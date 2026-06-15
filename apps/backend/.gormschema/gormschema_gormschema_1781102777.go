package main

import (
	"ariga.io/atlas-provider-gorm/gormschema"
	"fmt"
	"io"
	"os"
)

func main() {
	stmts, err := gormschema.New("postgres", gormschema.WithModelPosition(map[any]string{})).Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
