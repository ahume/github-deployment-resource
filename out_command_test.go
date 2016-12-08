package resource_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ahume/go-github/github"

	"github.com/ahume/github-deployment-resource"
	"github.com/ahume/github-deployment-resource/fakes"
)

func file(path, contents string) {
	Ω(ioutil.WriteFile(path, []byte(contents), 0644)).Should(Succeed())
}

var _ = Describe("Status Out Command", func() {
	var (
		command      *resource.OutCommand
		githubClient *fakes.FakeGitHub

		sourcesDir string

		request resource.OutRequest
	)

	BeforeEach(func() {
		var err error

		githubClient = &fakes.FakeGitHub{}
		command = resource.NewOutCommand(githubClient, ioutil.Discard)

		sourcesDir, err = ioutil.TempDir("", "github-deployment")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		Ω(os.RemoveAll(sourcesDir)).Should(Succeed())
	})

	buildDeployment := func(id int, env string, task string) *github.Deployment {
		return &github.Deployment{
			ID:          github.Int(id),
			Environment: github.String(env),
			Task:        github.String(task),
			Ref:         github.String("master"),
			SHA:         github.String("12345"),
			Description: github.String("One more"),
			Creator: &github.User{
				Login: github.String("Something"),
			},
			CreatedAt: &github.Timestamp{time.Date(2016, 01, 20, 15, 15, 15, 0, time.UTC)},
		}
	}

	Context("when creating a new deployment status", func() {
		BeforeEach(func() {
			githubClient.GetDeploymentReturns(buildDeployment(1234, "production", "deploy"), nil)

			githubClient.CreateDeploymentStatusReturns(&github.DeploymentStatus{
				ID:        github.Int(12),
				State:     github.String("success"),
				CreatedAt: &github.Timestamp{time.Date(2016, 01, 20, 20, 20, 20, 0, time.UTC)},
			}, nil)
		})
		Context("with strings in params", func() {
			BeforeEach(func() {
				request = resource.OutRequest{
					Params: resource.OutParams{
						ID:    "1234",
						State: "success",
					},
				}
			})

			It("creates a new status", func() {
				_, err := command.Run(sourcesDir, request)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(githubClient.CreateDeploymentStatusCallCount()).Should(Equal(1))
				id, status := githubClient.CreateDeploymentStatusArgsForCall(0)

				Ω(id).Should(Equal(*github.Int(1234)))
				Ω(status.State).Should(Equal(github.String("success")))
			})

			It("returns some metadata", func() {
				outResponse, err := command.Run(sourcesDir, request)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(outResponse.Metadata).Should(ConsistOf(
					resource.MetadataPair{Name: "id", Value: "1234"},
					resource.MetadataPair{Name: "ref", Value: "master"},
					resource.MetadataPair{Name: "sha", Value: "12345"},
					resource.MetadataPair{Name: "task", Value: "deploy"},
					resource.MetadataPair{Name: "environment", Value: "production"},
					resource.MetadataPair{Name: "description", Value: "One more"},
					resource.MetadataPair{Name: "creator", Value: "Something"},
					resource.MetadataPair{Name: "latest_state", Value: "success"},
					resource.MetadataPair{Name: "status_id", Value: "12"},
					resource.MetadataPair{Name: "created_at", Value: "2016-01-20 15:15:15"},
					resource.MetadataPair{Name: "updated_at", Value: "2016-01-20 20:20:20"},
				))
			})

			It("returns the version number of the deployment, not the status", func() {
				outResponse, err := command.Run(sourcesDir, request)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(outResponse.Version).Should(Equal(
					resource.Version{
						ID: "1234",
					},
				))
			})
		})

		Context("with file names in params", func() {
			BeforeEach(func() {
				idPath := filepath.Join(sourcesDir, "id")
				statePath := filepath.Join(sourcesDir, "state")

				file(idPath, "1234")
				file(statePath, "success")

				request = resource.OutRequest{
					Params: resource.OutParams{
						ID: resource.File{
							File: "id",
						},
						State: resource.File{
							File: "state",
						},
					},
				}
			})

			It("creates a new status", func() {
				_, err := command.Run(sourcesDir, request)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(githubClient.CreateDeploymentStatusCallCount()).Should(Equal(1))
				id, status := githubClient.CreateDeploymentStatusArgsForCall(0)

				Ω(id).Should(Equal(*github.Int(1234)))
				Ω(status.State).Should(Equal(github.String("success")))
			})

			It("returns some metadata", func() {
				outResponse, err := command.Run(sourcesDir, request)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(outResponse.Metadata).Should(ConsistOf(
					resource.MetadataPair{Name: "id", Value: "1234"},
					resource.MetadataPair{Name: "ref", Value: "master"},
					resource.MetadataPair{Name: "sha", Value: "12345"},
					resource.MetadataPair{Name: "task", Value: "deploy"},
					resource.MetadataPair{Name: "environment", Value: "production"},
					resource.MetadataPair{Name: "description", Value: "One more"},
					resource.MetadataPair{Name: "creator", Value: "Something"},
					resource.MetadataPair{Name: "latest_state", Value: "success"},
					resource.MetadataPair{Name: "status_id", Value: "12"},
					resource.MetadataPair{Name: "created_at", Value: "2016-01-20 15:15:15"},
					resource.MetadataPair{Name: "updated_at", Value: "2016-01-20 20:20:20"},
				))
			})
		})
	})
})
