package resource

import (
	"fmt"
	"github.com/bradfitz/slice"
	"io"
	"strconv"
)

type CheckCommand struct {
	github GitHub
	writer io.Writer
}

func NewCheckCommand(github GitHub, writer io.Writer) *CheckCommand {
	return &CheckCommand{
		github: github,
		writer: writer,
	}
}

func (c *CheckCommand) Run(request CheckRequest) ([]Version, error) {
	fmt.Fprintln(c.writer, "getting deployments list")
	deployments, err := c.github.ListDeployments()

	if err != nil {
		return []Version{}, err
	}

	if len(deployments) == 0 {
		return []Version{}, nil
	}

	var latestVersions []Version

	for _, deployment := range deployments {
		id := *deployment.ID
		if strconv.Itoa(id) >= request.Version.ID {
			latestVersions = append(latestVersions, Version{ID: strconv.Itoa(id)})
		}
	}

	slice.Sort(latestVersions[:], func(i, j int) bool {
		return latestVersions[i].ID < latestVersions[j].ID
	})

	latestVersion := latestVersions[len(latestVersions)-1]

	if request.Version.ID == "" {
		return []Version{
			latestVersion,
		}, nil
	}

	return latestVersions, nil
}
