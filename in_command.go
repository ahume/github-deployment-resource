package resource

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type InCommand struct {
	github GitHub
	writer io.Writer
}

func NewInCommand(github GitHub, writer io.Writer) *InCommand {
	return &InCommand{
		github: github,
		writer: writer,
	}
}

func (c *InCommand) Run(destDir string, request InRequest) (InResponse, error) {
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return InResponse{}, err
	}

	id, _ := strconv.Atoi(request.Version.ID)
	deployment, err := c.github.GetDeployment(id)
	if err != nil {
		return InResponse{}, err
	}

	if deployment == nil {
		return InResponse{}, errors.New("no deployment")
	}

	idPath := filepath.Join(destDir, "id")
	err = ioutil.WriteFile(idPath, []byte(request.Version.ID), 0644)
	if err != nil {
		return InResponse{}, err
	}

	refPath := filepath.Join(destDir, "ref")
	err = ioutil.WriteFile(refPath, []byte(*deployment.Ref), 0644)
	if err != nil {
		return InResponse{}, err
	}

	shaPath := filepath.Join(destDir, "sha")
	err = ioutil.WriteFile(shaPath, []byte(*deployment.SHA), 0644)
	if err != nil {
		return InResponse{}, err
	}

	if deployment.Task != nil {
		taskPath := filepath.Join(destDir, "task")
		err = ioutil.WriteFile(taskPath, []byte(*deployment.Task), 0644)
		if err != nil {
			return InResponse{}, err
		}
	}

	if deployment.Environment != nil {
		envPath := filepath.Join(destDir, "environment")
		err = ioutil.WriteFile(envPath, []byte(*deployment.Environment), 0644)
		if err != nil {
			return InResponse{}, err
		}
	}

	if deployment.Description != nil {
		descPath := filepath.Join(destDir, "description")
		err = ioutil.WriteFile(descPath, []byte(*deployment.Description), 0644)
		if err != nil {
			return InResponse{}, err
		}
	}

	// Save the whole deployment too I guess.
	deploymentPath := filepath.Join(destDir, "deploymentJSON")
	deploymentJSON, _ := json.Marshal(deployment)
	err = ioutil.WriteFile(deploymentPath, deploymentJSON, 0644)
	if err != nil {
		return InResponse{}, err
	}

	return InResponse{
		Version:  Version{ID: strconv.Itoa(*deployment.ID)},
		Metadata: metadataFromDeployment(deployment),
	}, nil
}
