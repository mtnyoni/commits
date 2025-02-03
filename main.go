package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
)

type Codecommit struct {
	client *codecommit.Client
}

func CodecommitClient() Codecommit {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config, %v", err)
	}

	csc := codecommit.NewFromConfig(cfg)
	return Codecommit{client: csc}
}

func (c *Codecommit) GetRepos() ([]types.RepositoryNameIdPair, error) {
	repos, err := c.client.ListRepositories(context.TODO(), &codecommit.ListRepositoriesInput{})
	if err != nil {
		return nil, err
	}

	return repos.Repositories, nil
}

func (c *Codecommit) GetBranches(repo string) ([]string, error) {
	branches, err := c.client.ListBranches(context.TODO(), &codecommit.ListBranchesInput{
		RepositoryName: &repo,
	})

	if err != nil {
		return nil, err
	}

	return branches.Branches, nil
}

func main() {
	cc := CodecommitClient()

	repos, err := cc.GetRepos()
	if err != nil {
		log.Fatalf("Unable to list repositories, %v", err)
	}

	for _, repo := range repos {
		branches, err := cc.GetBranches(*repo.RepositoryName)
		if err != nil {
			log.Fatalf("Unable to list branches, %v", err)
		}

		for _, branch := range branches {
			log.Println(branch)
		}
	}
}
