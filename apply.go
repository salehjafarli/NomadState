package main

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"log"
)

func ApplySnapshot(client *api.Client, file string) error {
	state, err := ReadSnapshot(file)

	if err != nil {
		return fmt.Errorf("failed to get state file: %s", err)
	}

	var qw api.WriteOptions
	for _, job := range state.Jobs {
		log.Printf("Applying job %s", *job.ID)
		qw.Namespace = *job.Namespace
		_, _, err := client.Jobs().Register(&job, &qw)

		if err != nil {
			return fmt.Errorf("failed to register job: %s", err)
		}
	}

	return nil
}
