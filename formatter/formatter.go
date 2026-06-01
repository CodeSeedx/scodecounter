package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/CodeSeedx/scodecounter/gitinfo"
	"github.com/CodeSeedx/scodecounter/scanner"
)

// Output holds the complete output data.
type Output struct {
	Root         string              `json:"root"`
	TotalFiles   int                 `json:"total_files"`
	TotalLines   int                 `json:"total_lines"`
	TotalCode    int                 `json:"total_code"`
	TotalBlank   int                 `json:"total_blank"`
	TotalComment int                 `json:"total_comment"`
	ByLanguage   []scanner.LangStats `json:"by_language"`
	ByDirectory  []scanner.DirStats  `json:"by_directory"`
	Git          *gitinfo.CommitInfo `json:"git,omitempty"`
}

// FormatJSON outputs the result as JSON.
func FormatJSON(w io.Writer, result *scanner.ScanResult, gitInfo *gitinfo.CommitInfo) error {
	out := Output{
		Root:         result.Root,
		TotalFiles:   result.TotalFiles,
		TotalLines:   result.TotalLines,
		TotalCode:    result.TotalCode,
		TotalBlank:   result.TotalBlank,
		TotalComment: result.TotalComment,
		Git:          gitInfo,
	}

	// Convert maps to sorted slices
	for _, ls := range result.ByLanguage {
		out.ByLanguage = append(out.ByLanguage, *ls)
	}
	sort.Slice(out.ByLanguage, func(i, j int) bool {
		return out.ByLanguage[i].CodeLines > out.ByLanguage[j].CodeLines
	})

	for _, ds := range result.ByDirectory {
		out.ByDirectory = append(out.ByDirectory, *ds)
	}
	sort.Slice(out.ByDirectory, func(i, j int) bool {
		return out.ByDirectory[i].CodeLines > out.ByDirectory[j].CodeLines
	})

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
}

// FormatTerminal outputs the result as terminal tables.
func FormatTerminal(w io.Writer, result *scanner.ScanResult, gitInfo *gitinfo.CommitInfo) {
	// Summary
	fmt.Fprintf(w, "\n  scodecounter - Code Statistics\n")
	fmt.Fprintf(w, "  Root: %s\n\n", result.Root)
	fmt.Fprintf(w, "  ┌─────────────────────┬──────────┐\n")
	fmt.Fprintf(w, "  │ %-19s │ %8s │\n", "Metric", "Value")
	fmt.Fprintf(w, "  ├─────────────────────┼──────────┤\n")
	fmt.Fprintf(w, "  │ %-19s │ %8d │\n", "Files", result.TotalFiles)
	fmt.Fprintf(w, "  │ %-19s │ %8d │\n", "Total Lines", result.TotalLines)
	fmt.Fprintf(w, "  │ %-19s │ %8d │\n", "Code Lines", result.TotalCode)
	fmt.Fprintf(w, "  │ %-19s │ %8d │\n", "Blank Lines", result.TotalBlank)
	fmt.Fprintf(w, "  │ %-19s │ %8d │\n", "Comment Lines", result.TotalComment)
	fmt.Fprintf(w, "  └─────────────────────┴──────────┘\n")

	// By Language
	if len(result.ByLanguage) > 0 {
		fmt.Fprintf(w, "\n  By Language\n")
		fmt.Fprintf(w, "  %-20s %8s %10s %10s %10s %8s\n", "Language", "Files", "Code", "Blank", "Comment", "Lines")
		fmt.Fprintf(w, "  %-20s %8s %10s %10s %10s %8s\n", "────────", "─────", "────", "─────", "───────", "─────")

		langs := make([]scanner.LangStats, 0, len(result.ByLanguage))
		for _, ls := range result.ByLanguage {
			langs = append(langs, *ls)
		}
		sort.Slice(langs, func(i, j int) bool {
			return langs[i].CodeLines > langs[j].CodeLines
		})

		for _, ls := range langs {
			fmt.Fprintf(w, "  %-20s %8d %10d %10d %10d %8d\n",
				ls.Language, ls.Files, ls.CodeLines, ls.BlankLines, ls.CommentLines, ls.Lines)
		}
	}

	// By Directory
	if len(result.ByDirectory) > 0 {
		fmt.Fprintf(w, "\n  By Directory\n")
		fmt.Fprintf(w, "  %-30s %8s %10s %10s %10s %8s\n", "Directory", "Files", "Code", "Blank", "Comment", "Lines")
		fmt.Fprintf(w, "  %-30s %8s %10s %10s %10s %8s\n", "─────────", "─────", "────", "─────", "───────", "─────")

		dirs := make([]scanner.DirStats, 0, len(result.ByDirectory))
		for _, ds := range result.ByDirectory {
			dirs = append(dirs, *ds)
		}
		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i].CodeLines > dirs[j].CodeLines
		})

		// Limit to top 20 directories for readability
		limit := len(dirs)
		if limit > 20 {
			limit = 20
		}
		for _, ds := range dirs[:limit] {
			fmt.Fprintf(w, "  %-30s %8d %10d %10d %10d %8d\n",
				ds.Path, ds.Files, ds.CodeLines, ds.BlankLines, ds.CommentLines, ds.Lines)
		}
		if len(dirs) > 20 {
			fmt.Fprintf(w, "  ... and %d more directories\n", len(dirs)-20)
		}
	}

	// Git Info
	if gitInfo != nil && gitInfo.TotalCommits > 0 {
		fmt.Fprintf(w, "\n  Git Info\n")
		fmt.Fprintf(w, "  ┌─────────────────────┬──────────┐\n")
		fmt.Fprintf(w, "  │ %-19s │ %8d │\n", "Total Commits", gitInfo.TotalCommits)
		if gitInfo.FirstCommit != "" {
			fmt.Fprintf(w, "  │ %-19s │ %8s │\n", "First Commit", gitInfo.FirstCommit[:10])
		}
		if gitInfo.LastCommit != "" {
			fmt.Fprintf(w, "  │ %-19s │ %8s │\n", "Last Commit", gitInfo.LastCommit[:10])
		}
		fmt.Fprintf(w, "  └─────────────────────┴──────────┘\n")

		if len(gitInfo.ByAuthor) > 0 {
			fmt.Fprintf(w, "\n  Commits by Author\n")
			fmt.Fprintf(w, "  %-30s %8s\n", "Author", "Commits")
			fmt.Fprintf(w, "  %-30s %8s\n", "──────", "───────")

			type authorCount struct {
				Name  string
				Count int
			}
			var authors []authorCount
			for name, count := range gitInfo.ByAuthor {
				authors = append(authors, authorCount{name, count})
			}
			sort.Slice(authors, func(i, j int) bool {
				return authors[i].Count > authors[j].Count
			})
			for _, a := range authors {
				fmt.Fprintf(w, "  %-30s %8d\n", a.Name, a.Count)
			}
		}
	}

	fmt.Fprintln(w)
}
