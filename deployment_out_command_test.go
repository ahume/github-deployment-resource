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

var _ = Describe("Deployment Out Command", func() {
	var (
		command      *resource.DeploymentOutCommand
		githubClient *fakes.FakeGitHub

		sourcesDir string

		request resource.OutRequest
	)

	BeforeEach(func() {
		var err error

		githubClient = &fakes.FakeGitHub{}
		command = resource.NewDeploymentOutCommand(githubClient, ioutil.Discard)

		sourcesDir, err = ioutil.TempDir("", "github-deployment")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		Ω(os.RemoveAll(sourcesDir)).Should(Succeed())
	})

	Context("when creating a new deployment", func() {
		Context("with strings in params", func() {
			Context("when all possible params are present", func() {
				BeforeEach(func() {
					githubClient.CreateDeploymentReturns(&github.Deployment{
						ID:          github.Int(1),
						Ref:         github.String("ref"),
						SHA:         github.String("1234"),
						Task:        github.String("task"),
						Description: github.String("desc"),
						Environment: github.String("env"),
						Creator: &github.User{
							Login: github.String("theboss"),
						},
						CreatedAt: &github.Timestamp{time.Date(2016, 01, 20, 15, 15, 15, 0, time.UTC)},
					}, nil)

					request = resource.OutRequest{
						Params: resource.OutParams{
							Ref:         "ref",
							Task:        "task",
							Description: "desc",
							Payload:     []byte("{\"one\": \"two\"}"),
							Environment: "env",
						},
					}
				})
				It("creates a new deployment", func() {
					_, err := command.Run(sourcesDir, request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(githubClient.CreateDeploymentCallCount()).Should(Equal(1))
					deployment := githubClient.CreateDeploymentArgsForCall(0)

					Ω(deployment.Ref).Should(Equal(github.String("ref")))
					Ω(deployment.Task).Should(Equal(github.String("task")))
					Ω(deployment.Description).Should(Equal(github.String("desc")))
					Ω(deployment.Payload).Should(Equal(github.String("{\"one\": \"two\"}")))
					Ω(deployment.Environment).Should(Equal(github.String("env")))
				})

				It("returns some metadata", func() {
					outResponse, err := command.Run(sourcesDir, request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(outResponse.Metadata).Should(ConsistOf(
						resource.MetadataPair{Name: "id", Value: "1"},
						resource.MetadataPair{Name: "ref", Value: "ref"},
						resource.MetadataPair{Name: "sha", Value: "1234"},
						resource.MetadataPair{Name: "task", Value: "task"},
						resource.MetadataPair{Name: "environment", Value: "env"},
						resource.MetadataPair{Name: "description", Value: "desc"},
						resource.MetadataPair{Name: "creator", Value: "theboss"},
						resource.MetadataPair{Name: "created_at", Value: "2016-01-20 15:15:15"},
						resource.MetadataPair{Name: "status_count", Value: "0"},
					))
				})

				It("returns the new version number", func() {
					outResponse, err := command.Run(sourcesDir, request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(outResponse.Version).Should(Equal(
						resource.Version{
							ID: "1",
						},
					))
				})
			})

			Context("when only required params are present", func() {
				BeforeEach(func() {
					githubClient.CreateDeploymentReturns(&github.Deployment{
						ID:  github.Int(1),
						Ref: github.String("ref"),
						SHA: github.String("1234"),
						Creator: &github.User{
							Login: github.String("theboss"),
						},
						CreatedAt: &github.Timestamp{time.Date(2016, 01, 20, 15, 15, 15, 0, time.UTC)},
					}, nil)

					request = resource.OutRequest{
						Params: resource.OutParams{
							Ref: "ref",
						},
					}
				})

				It("returns some metadata", func() {
					outResponse, err := command.Run(sourcesDir, request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(outResponse.Metadata).Should(ConsistOf(
						resource.MetadataPair{Name: "id", Value: "1"},
						resource.MetadataPair{Name: "ref", Value: "ref"},
						resource.MetadataPair{Name: "sha", Value: "1234"},
						resource.MetadataPair{Name: "creator", Value: "theboss"},
						resource.MetadataPair{Name: "created_at", Value: "2016-01-20 15:15:15"},
						resource.MetadataPair{Name: "status_count", Value: "0"},
					))
				})
			})
		})

		Context("with file names in params", func() {
			BeforeEach(func() {
				refPath := filepath.Join(sourcesDir, "ref")
				taskPath := filepath.Join(sourcesDir, "task")
				payloadPath := filepath.Join(sourcesDir, "payload")
				envPath := filepath.Join(sourcesDir, "environment")
				descPath := filepath.Join(sourcesDir, "description")

				file(refPath, "ref-from-file")
				file(taskPath, "task-from-file")
				file(payloadPath, "{\"three\": \"four\"}")
				file(envPath, "environment-from-file")
				file(descPath, "description-from-file")
			})

			Context("when all possible params are present", func() {
				BeforeEach(func() {
					githubClient.CreateDeploymentReturns(&github.Deployment{
						ID:          github.Int(1),
						Ref:         github.String("ref-from-file"),
						SHA:         github.String("1234"),
						Task:        github.String("task-from-file"),
						Description: github.String("description-from-file"),
						Environment: github.String("environment-from-file"),
						Creator: &github.User{
							Login: github.String("theboss"),
						},
						CreatedAt: &github.Timestamp{time.Date(2016, 01, 20, 15, 15, 15, 0, time.UTC)},
					}, nil)

					request = resource.OutRequest{
						Params: resource.OutParams{
							Ref: map[string]interface{}{
								"file": "ref",
							},
							Task: map[string]interface{}{
								"file": "task",
							},
							Payload:     []byte("{\"one\": \"two\"}"),
							PayloadPath: "payload",
							Environment: map[string]interface{}{
								"file": "environment",
							},
							Description: map[string]interface{}{
								"file": "description",
							},
						},
					}
				})

				It("creates a new deployment", func() {
					_, err := command.Run(sourcesDir, request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(githubClient.CreateDeploymentCallCount()).Should(Equal(1))
					deployment := githubClient.CreateDeploymentArgsForCall(0)

					Ω(deployment.Ref).Should(Equal(github.String("ref-from-file")))
					Ω(deployment.Task).Should(Equal(github.String("task-from-file")))
					Ω(deployment.Environment).Should(Equal(github.String("environment-from-file")))
					Ω(deployment.Payload).Should(Equal(github.String("{\"one\":\"two\",\"three\":\"four\"}")))
					Ω(deployment.Description).Should(Equal(github.String("description-from-file")))
				})

				It("returns some metadata", func() {
					outResponse, err := command.Run(sourcesDir, request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(outResponse.Metadata).Should(ConsistOf(
						resource.MetadataPair{Name: "id", Value: "1"},
						resource.MetadataPair{Name: "ref", Value: "ref-from-file"},
						resource.MetadataPair{Name: "sha", Value: "1234"},
						resource.MetadataPair{Name: "task", Value: "task-from-file"},
						resource.MetadataPair{Name: "environment", Value: "environment-from-file"},
						resource.MetadataPair{Name: "description", Value: "description-from-file"},
						resource.MetadataPair{Name: "creator", Value: "theboss"},
						resource.MetadataPair{Name: "created_at", Value: "2016-01-20 15:15:15"},
						resource.MetadataPair{Name: "status_count", Value: "0"},
					))
				})
			})

			Context("when only required params are present", func() {
				BeforeEach(func() {
					githubClient.CreateDeploymentReturns(&github.Deployment{
						ID:  github.Int(1),
						Ref: github.String("ref-from-file"),
						SHA: github.String("1234"),
						Creator: &github.User{
							Login: github.String("theboss"),
						},
						CreatedAt: &github.Timestamp{time.Date(2016, 01, 20, 15, 15, 15, 0, time.UTC)},
					}, nil)
					request = resource.OutRequest{
						Params: resource.OutParams{
							Ref: map[string]interface{}{
								"file": "ref",
							},
						},
					}
				})

				It("returns some metadata", func() {
					outResponse, err := command.Run(sourcesDir, request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(outResponse.Metadata).Should(ConsistOf(
						resource.MetadataPair{Name: "id", Value: "1"},
						resource.MetadataPair{Name: "ref", Value: "ref-from-file"},
						resource.MetadataPair{Name: "sha", Value: "1234"},
						resource.MetadataPair{Name: "creator", Value: "theboss"},
						resource.MetadataPair{Name: "created_at", Value: "2016-01-20 15:15:15"},
						resource.MetadataPair{Name: "status_count", Value: "0"},
					))
				})
			})
		})
	})
})
