package resource

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	if request.Params.Ref == "" {
		return OutResponse{}, errors.New("ref is a required parameter")
	}

	newDeployment := &github.DeploymentRequest{
		Ref:              github.String(request.Params.Ref),
		RequiredContexts: &[]string{},
	}

	if len(request.Params.Task) > 0 {
		newDeployment.Task = github.String(request.Params.Task)
	}
	if len(request.Params.Payload) > 0 {
		newDeployment.Payload = github.String(request.Params.Payload)
	}
	if len(request.Params.Environment) > 0 {
		newDeployment.Environment = github.String(request.Params.Environment)
	}
	if len(request.Params.Description) > 0 {
		newDeployment.Description = github.String(request.Params.Description)
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
