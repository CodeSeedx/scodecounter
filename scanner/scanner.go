package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/CodeSeedx/scodecounter/language"
)

// FileStats holds statistics for a single file.
type FileStats struct {
	Path       string
	Language   string
	Lines      int
	CodeLines  int
	BlankLines int
	CommentLines int
}

// DirStats holds aggregated statistics for a directory.
type DirStats struct {
	Path         string
	Files        int
	Lines        int
	CodeLines    int
	BlankLines   int
	CommentLines int
}

// LangStats holds aggregated statistics per language.
type LangStats struct {
	Language     string
	Files        int
	Lines        int
	CodeLines    int
	BlankLines   int
	CommentLines int
}

// ScanResult holds the complete scan result.
type ScanResult struct {
	Root         string
	TotalFiles   int
	TotalLines   int
	TotalCode    int
	TotalBlank   int
	TotalComment int
	ByLanguage   map[string]*LangStats
	ByDirectory  map[string]*DirStats
	Files        []FileStats
}

// Scan walks the directory and collects code statistics.
func Scan(root string) (*ScanResult, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	result := &ScanResult{
		Root:        absRoot,
		ByLanguage:  make(map[string]*LangStats),
		ByDirectory: make(map[string]*DirStats),
	}

	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip hidden directories and common non-source directories
		if info.IsDir() {
			base := info.Name()
			if strings.HasPrefix(base, ".") && base != "." {
				return filepath.SkipDir
			}
			switch base {
			case "node_modules", "vendor", "__pycache__", "venv", ".venv",
				"dist", "build", "target", ".git", ".svn", ".hg":
				return filepath.SkipDir
			}
			return nil
		}

		// Detect language
		lang := language.Detect(path)
		if lang == nil {
			return nil // Skip unknown files
		}

		// Count lines
		fs, err := countLines(path, lang)
		if err != nil {
			return nil // Skip files that can't be read
		}
		fs.Path = path
		fs.Language = lang.Name

		// Update totals
		result.TotalFiles++
		result.TotalLines += fs.Lines
		result.TotalCode += fs.CodeLines
		result.TotalBlank += fs.BlankLines
		result.TotalComment += fs.CommentLines
		result.Files = append(result.Files, *fs)

		// Update per-language stats
		ls, ok := result.ByLanguage[lang.Name]
		if !ok {
			ls = &LangStats{Language: lang.Name}
			result.ByLanguage[lang.Name] = ls
		}
		ls.Files++
		ls.Lines += fs.Lines
		ls.CodeLines += fs.CodeLines
		ls.BlankLines += fs.BlankLines
		ls.CommentLines += fs.CommentLines

		// Update per-directory stats (relative directory from root)
		relDir := filepath.Dir(path)
		relDir, _ = filepath.Rel(absRoot, relDir)
		if relDir == "" || relDir == "." {
			relDir = "."
		}
		ds, ok := result.ByDirectory[relDir]
		if !ok {
			ds = &DirStats{Path: relDir}
			result.ByDirectory[relDir] = ds
		}
		ds.Files++
		ds.Lines += fs.Lines
		ds.CodeLines += fs.CodeLines
		ds.BlankLines += fs.BlankLines
		ds.CommentLines += fs.CommentLines

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func countLines(path string, lang *language.Language) (*FileStats, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fs := &FileStats{}
	scanner := bufio.NewScanner(f)
	inBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		fs.Lines++
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			fs.BlankLines++
			continue
		}

		// Handle block comments
		if lang.BlockStart != "" && lang.BlockEnd != "" {
			if inBlock {
				fs.CommentLines++
				if strings.Contains(trimmed, lang.BlockEnd) {
					inBlock = false
				}
				continue
			}
			if strings.HasPrefix(trimmed, lang.BlockStart) {
				fs.CommentLines++
				if !strings.Contains(trimmed, lang.BlockEnd) {
					inBlock = true
				}
				continue
			}
		}

		// Handle single-line comments
		isComment := false
		for _, prefix := range lang.LineComment {
			if prefix != "" && strings.HasPrefix(trimmed, prefix) {
				isComment = true
				break
			}
		}
		if isComment {
			fs.CommentLines++
			continue
		}

		fs.CodeLines++
	}

	return fs, scanner.Err()
}
