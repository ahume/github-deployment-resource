package resource_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "github.com/google/go-github/github"

  "github.com/ahume/github-deployment-resource"
  "github.com/ahume/github-deployment-resource/fakes"
)

var _ = Describe("Check Command", func() {
  var (
    githubClient *fakes.FakeGitHub
    command      *resource.CheckCommand

    returnedDeployments []*github.Deployment
  )

  BeforeEach(func() {
    githubClient = &fakes.FakeGitHub{}
    command = resource.NewCheckCommand(githubClient)

    returnedDeployments = []*github.Deployment{}
  })

  JustBeforeEach(func() {
    githubClient.ListDeploymentsReturns(returnedDeployments, nil)
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
        command := resource.NewCheckCommand(githubClient)

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
        command := resource.NewCheckCommand(githubClient)

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
})
