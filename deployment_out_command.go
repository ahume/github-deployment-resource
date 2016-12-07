package resource

import (
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
  ref := request.Params.Ref;
  if request.Params.RefPath != "" {
    var err error
    ref, err = c.fileContents(filepath.Join(sourceDir, request.Params.RefPath))
    if err != nil {
      return OutResponse{}, err
    }
  }

  task := request.Params.Task
  if request.Params.TaskPath != "" {
    var err error
    task, err = c.fileContents(filepath.Join(sourceDir, request.Params.TaskPath))
    if err != nil {
      return OutResponse{}, err
    }
  }

  payload := ""
  if request.Params.PayloadPath != "" {
    var err error
    payload, err = c.fileContents(filepath.Join(sourceDir, request.Params.PayloadPath))
    if err != nil {
      return OutResponse{}, err
    }
  }

  environment := request.Params.Environment
  if request.Params.EnvironmentPath != "" {
    var err error
    environment, err = c.fileContents(filepath.Join(sourceDir, request.Params.EnvironmentPath))
    if err != nil {
      return OutResponse{}, err
    }
  }

  description := request.Params.Description
  if request.Params.DescriptionPath != "" {
    var err error
    description, err = c.fileContents(filepath.Join(sourceDir, request.Params.DescriptionPath))
    if err != nil {
      return OutResponse{}, err
    }
  }

  newDeployment := &github.DeploymentRequest{
    Ref: github.String(ref),
  }

  if len(task) > 0 {
    newDeployment.Task = github.String(task);
  }
  if len(payload) > 0 {
    newDeployment.Payload = github.String(payload);
  }
  if len(environment) > 0 {
    newDeployment.Environment = github.String(environment);
  }
  if len(description) > 0 {
    newDeployment.Description = github.String(description);
  }

  deployment, err := c.github.CreateDeployment(newDeployment)
  if err != nil {
    return OutResponse{}, err
  }

  return OutResponse{
    Version: Version{ID: strconv.Itoa(*deployment.ID)},
    Metadata: metadataFromDeployment(deployment),
  }, nil
}

func (c *DeploymentOutCommand) fileContents(path string) (string, error) {
  contents, err := ioutil.ReadFile(path)
  if err != nil {
    return "", err
  }

  return strings.TrimSpace(string(contents)), nil
}
