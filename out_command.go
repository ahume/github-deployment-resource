package resource

import (
  "io"
  "io/ioutil"
  "fmt"
  "path/filepath"
  "strconv"
  "strings"

  "github.com/ahume/go-github/github"
)

type OutCommand struct {
  github GitHub
  writer io.Writer
}

func NewOutCommand(github GitHub, writer io.Writer) *OutCommand {
  return &OutCommand{
    github: github,
    writer: writer,
  }
}

func (c *OutCommand) RunForDeployment(sourceDir string, request OutRequest) (OutResponse, error) {
  ref, err := c.fileContents(filepath.Join(sourceDir, request.Params.RefPath))
  if err != nil {
    return OutResponse{}, err
  }

  task, err := c.fileContents(filepath.Join(sourceDir, request.Params.TaskPath))
  if err != nil {
    return OutResponse{}, err
  }

  payload, err := c.fileContents(filepath.Join(sourceDir, request.Params.PayloadPath))
  if err != nil {
    return OutResponse{}, err
  }

  env, err := c.fileContents(filepath.Join(sourceDir, request.Params.EnvironmentPath))
  if err != nil {
    return OutResponse{}, err
  }

  description, err := c.fileContents(filepath.Join(sourceDir, request.Params.DescriptionPath))
  if err != nil {
    return OutResponse{}, err
  }

  newDeployment := &github.DeploymentRequest{
    Ref:          github.String(ref),
    Task:         github.String(task),
    Payload:      github.String(payload),
    Environment:  github.String(env),
    Description:  github.String(description),
  }

  fmt.Fprintf(c.writer, "creating Deployment")
  deployment, err := c.github.CreateDeployment(newDeployment)
  if err != nil {
    return OutResponse{}, err
  }

  return OutResponse{
    Version:  Version{ID: strconv.Itoa(*deployment.ID)},
    Metadata: metadataFromDeployment(deployment),
  }, nil
}

func (c *OutCommand) RunForStatus(sourceDir string, request OutRequest) (OutResponse, error) {
  id, err := c.fileContents(filepath.Join(sourceDir, request.Params.IDPath))
  if err != nil {
    return OutResponse{}, err
  }

  state, err := c.fileContents(filepath.Join(sourceDir, request.Params.StatePath))
  if err != nil {
    return OutResponse{}, err
  }

  newStatus := &github.DeploymentStatusRequest{
    State:          github.String(state),
  }

  fmt.Fprintf(c.writer, "creating DeploymentStatus")
  idString, err := strconv.Atoi(id)
  if err != nil {
    return OutResponse{}, err
  }
  status, err := c.github.CreateDeploymentStatus(idString, newStatus)
  if err != nil {
    return OutResponse{}, err
  }

  return OutResponse{
    Version:  Version{ID: id},
    Metadata: metadataFromStatus(status),
  }, nil
}

func (c *OutCommand) fileContents(path string) (string, error) {
  contents, err := ioutil.ReadFile(path)
  if err != nil {
    return "", err
  }

  return strings.TrimSpace(string(contents)), nil
}
