package resource

import (
	"github.com/ahume/go-github/github"
	"strconv"
)

func metadataFromDeployment(deployment *github.Deployment) []MetadataPair {
	metadata := []MetadataPair{}

	if deployment.ID != nil {
		id := *deployment.ID
		nameMeta := MetadataPair{
			Name:  "id",
			Value: strconv.Itoa(id),
		}
		metadata = append(metadata, nameMeta)
	}

	if deployment.Environment != nil {
		envMeta := MetadataPair{
			Name:  "environment",
			Value: *deployment.Environment,
		}
		metadata = append(metadata, envMeta)
	}

	if deployment.Creator != nil {
		creatorMeta := MetadataPair{
			Name:  "creator",
			Value: *deployment.Creator.Login,
		}
		metadata = append(metadata, creatorMeta)
	}

	return metadata
}
