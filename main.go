package main

import (
	"flag"
	"github.com/iamharvey/regi/internal/command"
	"github.com/spf13/pflag"
	"log"
)

func main() {
	// Remove timestamp from log
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Create Regi.
	cmd := command.NewRegiCommand()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Run Regi.
	if err := cmd.Execute(); err != nil {
		log.Fatal(err, "unable to run Regi")
	}
}
