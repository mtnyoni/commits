package main

import (
	"context"
	"log"
	"main/logger"
	"os"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
)

func main() {
	logger := logger.SetupLogger()
	cc := CodecommitClient(logger)

	repos, err := cc.GetRepos(logger)
	if err != nil {
		log.Fatalf("Unable to list repositories, %v", err)
	}

	for _, repo := range repos {
		branches, err := cc.GetBranches(*repo.RepositoryName, logger)
		if err != nil {
			log.Fatalf("Unable to list branches, %v", err)
		}

		for _, branch := range branches {
			log.Println(branch)
		}
	}
}

type Codecommit struct {
	client *codecommit.Client
}

func CodecommitClient(logger *logger.Logger) Codecommit {
	logger.Info("Get AWS config")
	cfg, err := aws_config.LoadDefaultConfig(context.TODO())

	if err != nil {
		logger.Error("Unable to load SDK config, %v", err)
		os.Exit(1)
	}

	csc := codecommit.NewFromConfig(cfg)
	return Codecommit{client: csc}
}

func (c *Codecommit) GetRepos(logger *logger.Logger) ([]types.RepositoryNameIdPair, error) {

	logger.Info("Get list of repos")
	repos, err := c.client.ListRepositories(context.TODO(), &codecommit.ListRepositoriesInput{})

	if err != nil {
		logger.Error("Unable to list repositories, %v", err)
		return nil, err
	}

	logger.Info("Got list of repos")
	return repos.Repositories, nil
}

func (c *Codecommit) GetBranches(repo string, logger *logger.Logger) ([]string, error) {
	logger.Info("Get list of branches for repo: %v", repo)
	branches, err := c.client.ListBranches(context.TODO(), &codecommit.ListBranchesInput{
		RepositoryName: &repo,
	})

	if err != nil {
		logger.Error("Unable to list branches, %v", err)
		return nil, err
	}

	logger.Info("Got list of branches for repo: %v", repo)
	return branches.Branches, nil
}

func (c *Codecommit) GetBranchCommitsNumber(repo string, branch string, logger *logger.Logger) (*string, error) {
	logger.Info("Get branch commits number for repo: %v and branch: %v", repo, branch)
	commit, err := c.client.GetBranch(context.TODO(), &codecommit.GetBranchInput{
		RepositoryName: &repo,
		BranchName:     &branch,
	})

	if err != nil {
		logger.Error("Unable to get branch information, %v", err)
		return nil, err
	}

	logger.Info("Got branch commits number for repo: %v and branch: %v", repo, branch)
	return commit.Branch.CommitId, nil
}

// func GetBranchNumberOfCommits(branchName string, repoName string, c *codecommit.Client) int {
// 	branch, err := c.GetBranch(context.TO(), &codecommit.GetBranchInput{
// 		RepositoryName: &repoName,
// 		BranchName:     &branchName,
// 	})
// 	if err != nil {
// 		log.Fatalf("Unable to get branch information, %v", err)
// 	}

// 	commitCount := 0
// 	commitID := branch.Branch.CommitId

// 	for commitID != nil {
// 		commitCount++
// 		commit, err := getCommitWithRetry(c, repoName, *commitID)
// 		if err != nil {
// 			log.Fatalf("Unable to get commit information, %v", err)
// 		}

// 		if len(commit.Commit.Parents) > 0 {
// 			commitID = &commit.Commit.Parents[0]
// 		} else {
// 			commitID = nil
// 		}
// 	}

// 	return commitCount
// }

// func getCommitWithRetry(c *codecommit.Client, repoName string, commitID string) (*codecommit.GetCommitOutput, error) {
// 	var commit *codecommit.GetCommitOutput
// 	var err error
// 	maxRetries := 5
// 	for i := 0; i < maxRetries; i++ {
// 		commit, err = c.GetCommit(context.TO(), &codecommit.GetCommitInput{
// 			RepositoryName: &repoName,
// 			CommitId:       &commitID,
// 		})
// 		if err == nil {
// 			return commit, nil
// 		}

// 		if awsErr, ok := err(aws); ok && awsErr.ErrorCode() == "ThrottlingException" {
// 			time.Sleep(time.Duration(2^i) * time.Second)
// 		} else {
// 			return nil, err
// 		}
// 	}
// 	return nil, err
// }

// func something() {
// 	// Initialize a session that the SDK uses to load credentials from the shared credentials file ~/.aws/credentials
// 	cfg, err := config.LoadDefaultConfig(context.TO())
// 	if err != nil {
// 		log.Fatalf("Unable to load AWS SDK config, %v", err)
// 	}

// 	// Create CodeCommit service client
// 	csc := codecommit.NewFromConfig(cfg)

// 	repos, err := csc.ListRepositories(context.TO(), &codecommit.ListRepositoriesInput{})
// 	if err != nil {
// 		log.Fatalf("Unable to list repositories, %v", err)
// 	}

// 	for _, repo := range repos.Repositories {
// 		branches := GetRepoBranches(*repo.RepositoryName, csc)
// 		for _, branch := range branches {
// 			fmt.Println(GetBranchNumberOfCommits(branch, *repo.RepositoryName, csc))
// 		}
// 	}

// 	// // List commits
// 	// input := &codecommit.ListR{
// 	// 	RepositoryName: aws.String(repoName),
// 	// }

// 	// result, err := cvc.ListCommits(input)
// 	// if err != nil {
// 	// 	log.Fatalf("Unable to list commits, %v", err)
// 	// }

// 	// // Aggregate commit counts by author
// 	// commitCounts := make(map[string]int)
// 	// for _, commit := range result.Commits {
// 	// 	author := *commit.Author.Name
// 	// 	commitCounts[author]++
// 	// }

// 	// // Print commit counts
// 	// for author, count := range commitCounts {
// 	// 	fmt.Printf("Author: %s, Commits: %d\n", author, count)
// 	// }
// }
