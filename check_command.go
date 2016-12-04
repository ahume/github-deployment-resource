package resource

import (
	"github.com/bradfitz/slice"
	"strconv"
)

type CheckCommand struct {
	github GitHub
}

func NewCheckCommand(github GitHub) *CheckCommand {
	return &CheckCommand{
		github: github,
	}
}

func (c *CheckCommand) Run(request CheckRequest) ([]Version, error) {
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

	if (request.Version == Version{}) {
		return []Version{
			latestVersion,
		}, nil
	}

	return latestVersions, nil
}
