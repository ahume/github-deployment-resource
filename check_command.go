package resource

import (
	"fmt"
	"io"
	"strconv"

	"github.com/bradfitz/slice"
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

	var latestVersions []Version

	for _, deployment := range deployments {
		if len(request.Source.Environments) > 0 {
			found := false
			for _, env := range request.Source.Environments {
				if env == *deployment.Environment {
					found = true
				}
			}

			if !found {
				continue
			}
		}

		id := *deployment.ID
		lastID, err := strconv.Atoi(request.Version.ID)
		if err != nil {
			latestVersions = append(latestVersions, Version{ID: strconv.Itoa(id)})
		} else {
			if id >= lastID {
				latestVersions = append(latestVersions, Version{ID: strconv.Itoa(id)})
			}
		}
	}

	if len(latestVersions) == 0 {
		return []Version{}, nil
	}

	slice.Sort(latestVersions[:], func(i, j int) bool {
		iID, _ := strconv.Atoi(latestVersions[i].ID)
		jID, _ := strconv.Atoi(latestVersions[j].ID)
		return iID < jID
	})

	latestVersion := latestVersions[len(latestVersions)-1]

	if request.Version.ID == "" {
		return []Version{
			latestVersion,
		}, nil
	}

	return latestVersions, nil
}
