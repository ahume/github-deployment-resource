package resource_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ahume/go-github/github"

	"github.com/ahume/github-deployment-resource"
	"github.com/ahume/github-deployment-resource/fakes"
)

var _ = Describe("Check Command", func() {
	var (
		command      *resource.CheckCommand
		githubClient *fakes.FakeGitHub

		returnedDeployments        []*github.Deployment
		returnedDeploymentStatuses []*github.DeploymentStatus

		requestedEnvironment string
		unwantedEnvironment  string
	)

	BeforeEach(func() {
		githubClient = &fakes.FakeGitHub{}
		command = resource.NewCheckCommand(githubClient, ioutil.Discard)

		returnedDeployments = []*github.Deployment{}
		returnedDeploymentStatuses = []*github.DeploymentStatus{}
	})

	JustBeforeEach(func() {
		githubClient.ListDeploymentsReturns(returnedDeployments, nil)
		githubClient.ListDeploymentStatusesReturns(returnedDeploymentStatuses, nil)
	})

	Context("when this is the first time that the resource has been run", func() {
		Context("when there are no deployments", func() {
			BeforeEach(func() {
				returnedDeployments = []*github.Deployment{}
			})

			It("returns no versions", func() {
				versions, err := command.Run(resource.CheckRequest{})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).Should(BeEmpty())
			})
		})

		Context("when there are deployments", func() {
			BeforeEach(func() {
				returnedDeployments = []*github.Deployment{
					newDeployment(3),
					newDeployment(2),
					newDeployment(1),
				}
			})

			It("outputs the most recent version if there is no current version", func() {
				versions, err := command.Run(resource.CheckRequest{})
				Ω(err).ShouldNot(HaveOccurred())

				Ω(versions).Should(HaveLen(1))
				Ω(versions[0]).Should(Equal(resource.Version{
					ID: "3",
				}))
			})
		})
	})

	Context("when there is a current version", func() {
		Context("when there are no deployments", func() {
			BeforeEach(func() {
				returnedDeployments = []*github.Deployment{}
			})

			It("returns no versions", func() {
				versions, err := command.Run(resource.CheckRequest{
					Version: resource.Version{
						ID: "3",
					},
				})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).Should(BeEmpty())
			})
		})

		Context("when there are deployments", func() {
			BeforeEach(func() {
				returnedDeployments = []*github.Deployment{
					newDeployment(3),
					newDeployment(2),
					newDeployment(1),
				}
			})

			It("outputs the most recent version if it matches the current version", func() {
				command := resource.NewCheckCommand(githubClient, ioutil.Discard)

				versions, err := command.Run(resource.CheckRequest{
					Version: resource.Version{
						ID: "3",
					},
				})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).Should(HaveLen(1))
				Ω(versions[0]).Should(Equal(resource.Version{
					ID: "3",
				}))
			})

			It("outputs versions later than and including the current", func() {
				command := resource.NewCheckCommand(githubClient, ioutil.Discard)

				versions, err := command.Run(resource.CheckRequest{
					Version: resource.Version{
						ID: "2",
					},
				})

				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).Should(HaveLen(2))

				Ω(versions[0]).Should(Equal(resource.Version{
					ID: "2",
				}))
				Ω(versions[1]).Should(Equal(resource.Version{
					ID: "3",
				}))
			})
		})

	})

	Context("when environment provided to filter on", func() {
		Context("when there are no deployments", func() {
			BeforeEach(func() {
				returnedDeployments = []*github.Deployment{}
			})

			It("returns no versions", func() {
				versions, err := command.Run(resource.CheckRequest{
					Version: resource.Version{
						ID: "3",
					},
				})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).Should(BeEmpty())
			})
		})

		Context("when there are deployments but not related to filtered environment", func() {
			BeforeEach(func() {
				requestedEnvironment = "production"
				unwantedEnvironment = "dev"
				returnedDeployments = []*github.Deployment{
					newDeploymentWithEnvironment(3, &unwantedEnvironment),
					newDeploymentWithEnvironment(2, &unwantedEnvironment),
					newDeploymentWithEnvironment(1, &unwantedEnvironment),
				}
			})

			It("returns no versions", func() {
				command := resource.NewCheckCommand(githubClient, ioutil.Discard)

				versions, err := command.Run(resource.CheckRequest{
					Source: resource.Source{
						Environment: requestedEnvironment,
					},
				})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).Should(BeEmpty())
			})
		})

		Context("when there are deployments related to filtered environment", func() {
			BeforeEach(func() {
				requestedEnvironment = "production"
				unwantedEnvironment = "dev"
				returnedDeployments = []*github.Deployment{
					newDeploymentWithEnvironment(5, &unwantedEnvironment),
					newDeploymentWithEnvironment(4, &requestedEnvironment),
					newDeploymentWithEnvironment(3, &requestedEnvironment),
					newDeploymentWithEnvironment(2, &requestedEnvironment),
					newDeploymentWithEnvironment(1, &unwantedEnvironment),
				}
			})

			Context("when this is the first time that the resource has been run", func() {

				It("outputs the most recent version related to filtered environment if there is no current version", func() {
					command := resource.NewCheckCommand(githubClient, ioutil.Discard)

					versions, err := command.Run(resource.CheckRequest{
						Source: resource.Source{
							Environment: requestedEnvironment,
						},
					})

					Ω(err).ShouldNot(HaveOccurred())
					Ω(versions).Should(HaveLen(1))

					Ω(versions[0]).Should(Equal(resource.Version{
						ID: "4",
					}))
				})

			})

			Context("when there is a current version", func() {
				It("outputs the most recent version related to the filtered environment if it matches the current version", func() {
					command := resource.NewCheckCommand(githubClient, ioutil.Discard)

					versions, err := command.Run(resource.CheckRequest{
						Source: resource.Source{
							Environment: requestedEnvironment,
						},
						Version: resource.Version{
							ID: "4",
						},
					})

					Ω(err).ShouldNot(HaveOccurred())
					Ω(versions).Should(HaveLen(1))

					Ω(versions[0]).Should(Equal(resource.Version{
						ID: "4",
					}))
				})

				It("outputs versions later than and including the current", func() {
					command := resource.NewCheckCommand(githubClient, ioutil.Discard)

					versions, err := command.Run(resource.CheckRequest{
						Source: resource.Source{
							Environment: requestedEnvironment,
						},
						Version: resource.Version{
							ID: "3",
						},
					})

					Ω(err).ShouldNot(HaveOccurred())
					Ω(versions).Should(HaveLen(2))

					Ω(versions[0]).Should(Equal(resource.Version{
						ID: "3",
					}))
					Ω(versions[1]).Should(Equal(resource.Version{
						ID: "4",
					}))
				})

			})
		})

	})
})
