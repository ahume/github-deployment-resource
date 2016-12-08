package resource

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ahume/go-github/github"
)

type DeploymentOutCommand struct {
	github GitHub
	writer io.Writer
}

func NewDeploymentOutCommand(github GitHub, writer io.Writer) *DeploymentOutCommand {
	return &DeploymentOutCommand{
		github: github,
		writer: writer,
	}
}

func (c *DeploymentOutCommand) Run(sourceDir string, request OutRequest) (OutResponse, error) {
	// TODO: Ref is required, so what happens if it's missing?
	ref, ok := request.Params.Ref.(string)
	if ok != true {
		var err error
		v := request.Params.Ref.(map[string]interface{})
		ref, err = c.fileContents(filepath.Join(sourceDir, v["file"].(string)))
		if err != nil {
			return OutResponse{}, err
		}
	}

	task, ok := request.Params.Task.(string)
	if ok != true {
		v, ok := request.Params.Task.(map[string]interface{})
		if ok == true {
			var err error
			task, err = c.fileContents(filepath.Join(sourceDir, v["file"].(string)))
			if err != nil {
				return OutResponse{}, err
			}
		}
	}

	payload, ok := request.Params.Payload.(string)
	if ok != true {
		v, ok := request.Params.Payload.(map[string]interface{})
		if ok == true {
			var err error
			payload, err = c.fileContents(filepath.Join(sourceDir, v["file"].(string)))
			if err != nil {
				return OutResponse{}, err
			}
		}
	}

	environment, ok := request.Params.Environment.(string)
	if ok != true {
		v, ok := request.Params.Environment.(map[string]interface{})
		if ok == true {
			var err error
			environment, err = c.fileContents(filepath.Join(sourceDir, v["file"].(string)))
			if err != nil {
				return OutResponse{}, err
			}
		}
	}

	description, ok := request.Params.Description.(string)
	if ok != true {
		v, ok := request.Params.Description.(map[string]interface{})
		if ok == true {
			var err error
			description, err = c.fileContents(filepath.Join(sourceDir, v["file"].(string)))
			if err != nil {
				return OutResponse{}, err
			}
		}
	}

	newDeployment := &github.DeploymentRequest{
		Ref:              github.String(ref),
		RequiredContexts: &[]string{},
	}

	if len(task) > 0 {
		newDeployment.Task = github.String(task)
	}
	if len(payload) > 0 {
		newDeployment.Payload = github.String(payload)
	}
	if len(environment) > 0 {
		newDeployment.Environment = github.String(environment)
	}
	if len(description) > 0 {
		newDeployment.Description = github.String(description)
	}

	fmt.Fprintln(c.writer, "creating deployment")
	deployment, err := c.github.CreateDeployment(newDeployment)
	if err != nil {
		return OutResponse{}, err
	}

	return OutResponse{
		Version:  Version{ID: strconv.Itoa(*deployment.ID)},
		Metadata: metadataFromDeployment(deployment, []*github.DeploymentStatus{}),
	}, nil
}

func (c *DeploymentOutCommand) fileContents(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contents)), nil
}
