package main

import (
	"fmt"
	"main/codecommit"
	"main/logger"
)

func main() {
	logger := logger.SetupLogger()
	cc := codecommit.NewClient(logger)

	repos, err := cc.GetRepos(logger)
	if err != nil {
		logger.Error("Failed to retrieve repositories: %v", err)
		return
	}

	for _, repo := range repos {
		branches, err := cc.GetBranches(*repo.RepositoryName, logger)
		if err != nil {
			logger.Error("Failed to retrieve branches for repository %s: %v", *repo.RepositoryName, err)
			continue
		}

		for _, branch := range branches {
			commits, err := cc.GetCommitsOnBranch(*repo.RepositoryName, &branch)
			if err != nil {
				fmt.Printf("Error getting commits: %v\n", err)
				return
			}

			fmt.Println("Length of commits", len(commits))
		}
	}

}
