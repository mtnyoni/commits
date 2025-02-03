package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
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

func (c *Codecommit) GetRepos() (codecommit.ListRepositoriesOutput, error) {
	repos, err := c.client.ListRepositories(context.TODO(), &codecommit.ListRepositoriesInput{})
	if err != nil {
		log.Fatalf("Unable to list repositories, %v", err)
		return codecommit.ListRepositoriesOutput{}, err
	}
	return *repos, nil
}

func main() {
	cc := CodecommitClient()

	repos, err := cc.GetRepos()
	if err != nil {
		log.Fatalf("Unable to list repositories, %v", err)
	}

	for _, repo := range repos.Repositories {
		log.Println(*repo.RepositoryName)
	}
}
