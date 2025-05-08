package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"os"
	"path"
)

type NomadState struct {
	Jobs []api.Job
}

func ReadSnapshot(file string) (NomadState, error) {
	state := NomadState{}
	if file == "" {
		file = ".nomadstate"
	}
	stateFile, err := os.ReadFile(file)
	if err != nil {
		return state, err
	}

	err = json.Unmarshal(stateFile, &state)
	return state, err
}

func GenerateSnapshot(client *api.Client, output string) error {
	if output != "" {
		err := os.MkdirAll(output, 0644)

		if err != nil {
			return fmt.Errorf("failed to create output path: %s", err)
		}
	}

	qo := api.QueryOptions{}
	jobStubs, _, err := client.Jobs().List(&qo)

	if err != nil {
		return fmt.Errorf("failed to get job list: %s", err)
	}

	var state = NomadState{}
	q := &api.QueryOptions{}
	for _, j := range jobStubs {
		var job *api.Job
		q.Namespace = j.Namespace
		job, _, err := client.Jobs().Info(j.ID, q)

		if err != nil {
			return fmt.Errorf("failed to get job definition: %s", err)
		}
		state.Jobs = append(state.Jobs, *job)
	}

	b, err := json.MarshalIndent(state, "", "\t")

	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}

	filepath := path.Join(output, ".nomadstate")
	err = os.WriteFile(filepath, b, 0644)

	if err != nil {
		return fmt.Errorf("failed to create state file: %s", err)
	}

	return nil
}
