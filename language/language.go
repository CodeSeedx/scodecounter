package language

import "path/filepath"

// Language represents a programming language with its properties.
type Language struct {
	Name       string
	Extensions []string
	LineComment []string // Single-line comment prefixes
	BlockStart  string   // Block comment start
	BlockEnd    string   // Block comment end
}

var languages = []Language{
	{"Go", []string{".go"}, []string{"//"}, "/*", "*/"},
	{"Python", []string{".py", ".pyw"}, []string{"#"}, "", ""},
	{"JavaScript", []string{".js", ".mjs", ".cjs"}, []string{"//"}, "/*", "*/"},
	{"TypeScript", []string{".ts", ".tsx"}, []string{"//"}, "/*", "*/"},
	{"Java", []string{".java"}, []string{"//"}, "/*", "*/"},
	{"C", []string{".c", ".h"}, []string{"//"}, "/*", "*/"},
	{"C++", []string{".cpp", ".cxx", ".cc", ".hpp", ".hxx", ".hh"}, []string{"//"}, "/*", "*/"},
	{"C#", []string{".cs"}, []string{"//"}, "/*", "*/"},
	{"Rust", []string{".rs"}, []string{"//"}, "/*", "*/"},
	{"Ruby", []string{".rb"}, []string{"#"}, "=begin", "=end"},
	{"PHP", []string{".php", ".phtml"}, []string{"//", "#"}, "/*", "*/"},
	{"Swift", []string{".swift"}, []string{"//"}, "/*", "*/"},
	{"Kotlin", []string{".kt", ".kts"}, []string{"//"}, "/*", "*/"},
	{"Scala", []string{".scala", ".sc"}, []string{"//"}, "/*", "*/"},
	{"Shell", []string{".sh", ".bash", ".zsh"}, []string{"#"}, "", ""},
	{"Lua", []string{".lua"}, []string{"--"}, "--[[", "]]"},
	{"Perl", []string{".pl", ".pm", ".t"}, []string{"#"}, "", ""},
	{"R", []string{".r", ".R"}, []string{"#"}, "", ""},
	{"Julia", []string{".jl"}, []string{"#"}, "#=", "=#"},
	{"Haskell", []string{".hs", ".lhs"}, []string{"--"}, "{-", "-}"},
	{"Erlang", []string{".erl", ".hrl"}, []string{"%"}, "", ""},
	{"Elixir", []string{".ex", ".exs"}, []string{"#"}, "", ""},
	{"Clojure", []string{".clj", ".cljs", ".cljc", ".edn"}, []string{";", "#_"}, "", ""},
	{"Lisp", []string{".lisp", ".lsp", ".l", ".fasl"}, []string{";"}, "#|", "|#"},
	{"Scheme", []string{".scm", ".ss"}, []string{";"}, "#|", "|#"},
	{"Dart", []string{".dart"}, []string{"//"}, "/*", "*/"},
	{"Objective-C", []string{".m", ".mm"}, []string{"//"}, "/*", "*/"},
	{"Assembly", []string{".asm", ".s", ".S"}, []string{"#", ";", "//"}, "", ""},
	{"SQL", []string{".sql"}, []string{"--"}, "/*", "*/"},
	{"HTML", []string{".html", ".htm"}, []string{}, "<!--", "-->"},
	{"CSS", []string{".css"}, []string{}, "/*", "*/"},
	{"SCSS", []string{".scss"}, []string{"//"}, "/*", "*/"},
	{"SASS", []string{".sass"}, []string{"//"}, "", ""},
	{"Less", []string{".less"}, []string{"//"}, "/*", "*/"},
	{"XML", []string{".xml", ".xsl", ".xsd", ".svg"}, []string{}, "<!--", "-->"},
	{"JSON", []string{".json"}, []string{}, "", ""},
	{"YAML", []string{".yaml", ".yml"}, []string{"#"}, "", ""},
	{"TOML", []string{".toml"}, []string{"#"}, "", ""},
	{"Markdown", []string{".md", ".markdown"}, []string{}, "", ""},
	{"LaTeX", []string{".tex", ".sty", ".cls"}, []string{"%"}, "", ""},
	{"Vim", []string{".vim"}, []string{`"`}, "", ""},
	{"Emacs Lisp", []string{".el"}, []string{";"}, "", ""},
	{"Dockerfile", []string{"Dockerfile"}, []string{"#"}, "", ""},
	{"Makefile", []string{"Makefile", "makefile", "GNUmakefile"}, []string{"#"}, "", ""},
	{"CMake", []string{"CMakeLists.txt", ".cmake"}, []string{"#"}, "", ""},
	{"Terraform", []string{".tf", ".tfvars"}, []string{"#"}, "", ""},
	{"Protocol Buffers", []string{".proto"}, []string{"//"}, "/*", "*/"},
	{"GraphQL", []string{".graphql", ".gql"}, []string{"#"}, "", ""},
	{"Zig", []string{".zig"}, []string{"//"}, "", ""},
	{"Nim", []string{".nim"}, []string{"#"}, "#[", "]#"},
	{"V", []string{".v"}, []string{"//"}, "/*", "*/"},
	{"OCaml", []string{".ml", ".mli"}, []string{}, "(*", "*)"},
	{"F#", []string{".fs", ".fsx", ".fsi"}, []string{"//"}, "(*", "*)"},
	{"Pascal", []string{".pas", ".pp"}, []string{"//"}, "{", "}"},
	{"Ada", []string{".adb", ".ads"}, []string{"--"}, "", ""},
	{"Fortran", []string{".f", ".for", ".f90", ".f95", ".f03"}, []string{"!", "c", "C"}, "", ""},
	{"COBOL", []string{".cob", ".cbl"}, []string{"*"}, "", ""},
	{"Groovy", []string{".groovy", ".gvy"}, []string{"//"}, "/*", "*/"},
	{"PowerShell", []string{".ps1", ".psm1"}, []string{"#"}, "<#", "#>"},
	{"Batch", []string{".bat", ".cmd"}, []string{"REM", "rem", "@"}, "", ""},
	{"VBScript", []string{".vbs"}, []string{"'"}, "", ""},
	{"Awk", []string{".awk"}, []string{"#"}, "", ""},
	{"Sed", []string{".sed"}, []string{"#"}, "", ""},
	{"Diff", []string{".diff", ".patch"}, []string{}, "", ""},
	{"INI", []string{".ini", ".cfg", ".conf"}, []string{"#", ";"}, "", ""},
	{"Properties", []string{".properties"}, []string{"#", "!"}, "", ""},
}

// extensionMap maps file extensions to language names.
var extensionMap map[string]*Language

// filenameMap maps exact filenames to language names.
var filenameMap map[string]*Language

func init() {
	extensionMap = make(map[string]*Language)
	filenameMap = make(map[string]*Language)
	for i := range languages {
		lang := &languages[i]
		for _, ext := range lang.Extensions {
			if ext[0] == '.' {
				extensionMap[ext] = lang
			} else {
				filenameMap[ext] = lang
			}
		}
	}
}

// Detect detects the programming language of a file based on its path.
func Detect(filePath string) *Language {
	// First check exact filename match
	base := filepath.Base(filePath)
	if lang, ok := filenameMap[base]; ok {
		return lang
	}
	// Then check extension
	ext := filepath.Ext(base)
	if lang, ok := extensionMap[ext]; ok {
		return lang
	}
	return nil
}

// All returns all registered languages.
func All() []Language {
	return languages
}
