package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/hashicorp/nomad/api"
	"log"
	"os"
)

func main() {
	a := kingpin.New("NomadState", "").UsageWriter(os.Stdout)

	snapCmd := a.Command("snapshot", "Generate current state of jobs from Nomad cluster")
	output := snapCmd.Flag("output", "output path of nomadstate file").Default("").String()

	applyCmd := a.Command("apply", "Schedule jobs from given state")
	file := applyCmd.Flag("file", "State file").String()

	config := api.DefaultConfig()
	config.Namespace = "*"

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("Failed to create nomad client", err)
	}

	switch kingpin.MustParse(a.Parse(os.Args[1:])) {
	case snapCmd.FullCommand():
		err = GenerateSnapshot(client, *output)
	case applyCmd.FullCommand():
		err = ApplySnapshot(client, *file)
	}

	if err != nil {
		log.Fatal(err)
	}
}
