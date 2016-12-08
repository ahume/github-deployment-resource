package resource

import (
	"github.com/ahume/go-github/github"
	"strconv"
)

func metadataFromDeployment(deployment *github.Deployment, statuses []*github.DeploymentStatus) []MetadataPair {
	metadata := []MetadataPair{}

	if deployment.ID != nil {
		id := *deployment.ID
		nameMeta := MetadataPair{
			Name:  "id",
			Value: strconv.Itoa(id),
		}
		metadata = append(metadata, nameMeta)
	}

	if deployment.Ref != nil {
		refMeta := MetadataPair{
			Name:  "ref",
			Value: *deployment.Ref,
		}
		metadata = append(metadata, refMeta)
	}

	if deployment.SHA != nil {
		shaMeta := MetadataPair{
			Name:  "sha",
			Value: *deployment.SHA,
		}
		metadata = append(metadata, shaMeta)
	}

	if deployment.Task != nil {
		taskMeta := MetadataPair{
			Name:  "task",
			Value: *deployment.Task,
		}
		metadata = append(metadata, taskMeta)
	}

	if deployment.Environment != nil {
		envMeta := MetadataPair{
			Name:  "environment",
			Value: *deployment.Environment,
		}
		metadata = append(metadata, envMeta)
	}

	if deployment.Description != nil {
		descMeta := MetadataPair{
			Name:  "description",
			Value: *deployment.Description,
		}
		metadata = append(metadata, descMeta)
	}

	if deployment.Creator != nil {
		creatorMeta := MetadataPair{
			Name:  "creator",
			Value: *deployment.Creator.Login,
		}
		metadata = append(metadata, creatorMeta)
	}

	if deployment.CreatedAt != nil {
		createdtAtMeta := MetadataPair{
			Name:  "created_at",
			Value: deployment.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		metadata = append(metadata, createdtAtMeta)
	}

	if len(statuses) > 0 {
		if statuses[0].ID != nil {
			id := *statuses[0].ID
			nameMeta := MetadataPair{
				Name:  "status_id",
				Value: strconv.Itoa(id),
			}
			metadata = append(metadata, nameMeta)
		}

		if statuses[0].State != nil {
			envMeta := MetadataPair{
				Name:  "status",
				Value: *statuses[0].State,
			}
			metadata = append(metadata, envMeta)
		}

		if statuses[0].CreatedAt != nil {
			createdtAtMeta := MetadataPair{
				Name:  "status_created_at",
				Value: statuses[0].CreatedAt.Format("2006-01-02 15:04:05"),
			}
			metadata = append(metadata, createdtAtMeta)
		}
	}

	metadata = append(metadata, MetadataPair{
		Name:  "status_count",
		Value: strconv.Itoa(len(statuses)),
	})

	return metadata
}
