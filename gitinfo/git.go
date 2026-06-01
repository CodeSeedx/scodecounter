package gitinfo

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// CommitInfo holds git commit statistics.
type CommitInfo struct {
	TotalCommits int            `json:"total_commits"`
	ByAuthor     map[string]int `json:"by_author"`
	FirstCommit  string         `json:"first_commit,omitempty"`
	LastCommit   string         `json:"last_commit,omitempty"`
}

// GetCommitInfo returns git commit statistics for a repository.
func GetCommitInfo(repoPath string) (*CommitInfo, error) {
	info := &CommitInfo{
		ByAuthor: make(map[string]int),
	}

	// Get total commit count
	out, err := gitExec(repoPath, "rev-list", "--count", "HEAD")
	if err != nil {
		return info, nil // Not a git repo or no commits
	}
	info.TotalCommits, _ = strconv.Atoi(strings.TrimSpace(string(out)))

	// Get commit count by author
	out, err = gitExec(repoPath, "shortlog", "-sn", "--no-merges", "HEAD")
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Format: "   123  Author Name"
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) == 2 {
				count, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
				author := strings.TrimSpace(parts[1])
				info.ByAuthor[author] = count
			}
		}
	}

	// Get first commit date
	out, err = gitExec(repoPath, "log", "--reverse", "--format=%ai", "--max-count=1")
	if err == nil {
		info.FirstCommit = strings.TrimSpace(string(out))
	}

	// Get last commit date
	out, err = gitExec(repoPath, "log", "--format=%ai", "--max-count=1")
	if err == nil {
		info.LastCommit = strings.TrimSpace(string(out))
	}

	return info, nil
}

func gitExec(repoPath string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("git %s: %s", args[0], string(exitErr.Stderr))
		}
		return nil, err
	}
	return out, nil
}
