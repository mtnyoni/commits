package codecommit

import (
	"context"
	"errors"
	"fmt"
	"main/logger"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/aws/smithy-go"
)

type Codecommit struct {
	client *codecommit.Client
}

func CodecommitClient(logger *logger.Logger) Codecommit {
	logger.Info("Initializing AWS CodeCommit client")
	cfg, err := aws_config.LoadDefaultConfig(context.TODO())

	if err != nil {
		logger.Error("Failed to load AWS SDK config: %v", err)
		os.Exit(1)
	}

	csc := codecommit.NewFromConfig(cfg)
	logger.Info("AWS CodeCommit client initialized successfully")
	return Codecommit{client: csc}
}

func (c *Codecommit) GetRepos(logger *logger.Logger) ([]types.RepositoryNameIdPair, error) {
	logger.Info("Fetching list of repositories")
	repos, err := c.client.ListRepositories(context.TODO(), &codecommit.ListRepositoriesInput{})

	if err != nil {
		logger.Error("Failed to list repositories: %v", err)
		return nil, err
	}

	logger.Info("Successfully retrieved %d repositories", len(repos.Repositories))
	return repos.Repositories, nil
}

func (c *Codecommit) GetBranches(repo string, logger *logger.Logger) ([]string, error) {
	logger.Info("Fetching branches for repository: %s", repo)
	branches, err := c.client.ListBranches(context.TODO(), &codecommit.ListBranchesInput{
		RepositoryName: &repo,
	})

	if err != nil {
		logger.Error("Failed to list branches for repository %s: %v", repo, err)
		return nil, err
	}

	logger.Info("Successfully retrieved %d branches for repository: %s", len(branches.Branches), repo)
	return branches.Branches, nil
}

type CommitInfo struct {
	CommitID string
	Author   string
	Date     string
	Message  string
	Tags     []string
}

func (c *Codecommit) GetCommitsOnBranch(repoName, branchName string) ([]CommitInfo, error) {
	branchOutput, err := c.client.GetBranch(context.TODO(), &codecommit.GetBranchInput{
		RepositoryName: aws.String(repoName),
		BranchName:     aws.String(branchName),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get branch info: %w", err)
	}

	headCommitID := branchOutput.Branch.CommitId
	if headCommitID == nil {
		return nil, fmt.Errorf("branch %s has no commit", branchName)
	}

	// Traverse commit history starting from the HEAD.
	var allCommits []types.Commit
	currentCommitID := headCommitID

	for currentCommitID != nil {
		var commitOutput *codecommit.GetCommitOutput
		maxRetries := 3
		var lastErr error

		// Retry loop for handling throttling errors.
		for attempt := 0; attempt < maxRetries; attempt++ {
			commitOutput, err = c.client.GetCommit(context.TODO(), &codecommit.GetCommitInput{
				CommitId:       currentCommitID,
				RepositoryName: aws.String(repoName),
			})
			if err != nil {
				var apiErr smithy.APIError
				// Check if the error is a throttling exception.
				if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ThrottlingException" {
					// Exponential backoff.
					time.Sleep(time.Duration((attempt+1)*200) * time.Millisecond)
					lastErr = err
					continue
				}
				return nil, fmt.Errorf("failed to get commit %s: %w", aws.ToString(currentCommitID), err)
			}
			lastErr = nil
			break
		}
		if lastErr != nil {
			return nil, fmt.Errorf("failed to get commit %s after retries: %w", aws.ToString(currentCommitID), lastErr)
		}

		commit := commitOutput.Commit
		allCommits = append(allCommits, *commit)

		// Follow the first parent for linear history.
		if len(commit.Parents) > 0 {
			currentCommitID = aws.String(commit.Parents[0])
		} else {
			currentCommitID = nil
		}
	}

	// Map commits to CommitInfo, including associated tags.
	var commitInfos []CommitInfo
	for _, commit := range allCommits {
		info := CommitInfo{
			CommitID: aws.ToString(commit.CommitId),
			Message:  aws.ToString(commit.Message),
			Tags:     []string{},
		}

		if commit.Author != nil {
			info.Author = aws.ToString(commit.Author.Name)
			info.Date = aws.ToString(commit.Author.Date)
		}

		commitInfos = append(commitInfos, info)
	}

	return commitInfos, nil
}
