//go:build ignore

// Bootstrap configures a new project created from templatte-go, then deletes itself.
//
// Usage:
//
//	go run bootstrap.go [-module path] [-description text]
//
// Without flags, it derives the module path from the origin remote and
// prompts for anything it cannot derive.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	templateModule  = "github.com/andrewwkimm/templatte-go"
	templatePackage = "templatte"
	templateDoc     = "// Package templatte is an opinionated Go library template.\n"
)

var (
	modulePattern = regexp.MustCompile(`^[a-z0-9.-]+(/[A-Za-z0-9._~-]+)+$`)
	majorPattern  = regexp.MustCompile(`^v[0-9]+$`)
	remotePattern = regexp.MustCompile(`^(?:https://|git@)([^/:]+)[/:](.+?)(?:\.git)?$`)
)

// project identifies the module being configured.
type project struct {
	module      string
	packageName string
	description string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "bootstrap:", err)
		os.Exit(1)
	}
}

func run() error {
	module := flag.String("module", "", "module path, such as github.com/you/widget")
	description := flag.String("description", "", `package doc sentence that completes "Package <name> ..."`)
	flag.Parse()

	in := bufio.NewReader(os.Stdin)
	if *module == "" {
		*module = prompt(in, "Module path", originModule())
	}
	p := project{module: *module, packageName: packageName(*module)}
	if *description == "" {
		*description = prompt(in, fmt.Sprintf("Complete the sentence %q", "Package "+p.packageName+" ..."), "")
	}
	p.description = sentence(*description)

	if err := validate(p); err != nil {
		return err
	}
	if err := configure(p); err != nil {
		return err
	}
	if err := os.Remove("bootstrap.go"); err != nil {
		return fmt.Errorf("remove bootstrap.go: %w", err)
	}
	fmt.Printf("Configured %s (package %s).\n", p.module, p.packageName)
	return nil
}

// prompt reads one line from in, returning fallback when the line is empty.
func prompt(in *bufio.Reader, label, fallback string) string {
	if fallback != "" {
		fmt.Printf("%s [%s]: ", label, fallback)
	} else {
		fmt.Printf("%s: ", label)
	}
	line, _ := in.ReadString('\n') // A read error leaves an empty answer, which validate rejects.
	if answer := strings.TrimSpace(line); answer != "" {
		return answer
	}
	return fallback
}

// originModule returns the module path implied by the origin remote, or "" when there is none.
func originModule() string {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return ""
	}
	m := remotePattern.FindStringSubmatch(strings.TrimSpace(string(out)))
	if m == nil {
		return ""
	}
	return strings.ToLower(m[1] + "/" + m[2])
}

// packageName derives a Go package name from the last meaningful element of module.
func packageName(module string) string {
	elems := strings.Split(module, "/")
	name := elems[len(elems)-1]
	if majorPattern.MatchString(name) && len(elems) > 1 {
		name = elems[len(elems)-2]
	}
	name = strings.ToLower(name)
	name = strings.TrimPrefix(name, "go-")
	name = strings.TrimSuffix(name, "-go")
	name = strings.TrimSuffix(name, ".go")
	return strings.NewReplacer("-", "", ".", "", "_", "").Replace(name)
}

// sentence trims s and ends it with a period.
func sentence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasSuffix(s, ".") {
		return s
	}
	return s + "."
}

func validate(p project) error {
	if !modulePattern.MatchString(p.module) {
		return fmt.Errorf("invalid module path %q", p.module)
	}
	if !token.IsIdentifier(p.packageName) || token.IsKeyword(p.packageName) {
		return fmt.Errorf("module path %q yields invalid package name %q", p.module, p.packageName)
	}
	if p.description == "" {
		return errors.New("description cannot be empty")
	}
	required := []string{
		"go.mod", "tools/go.mod", "tools/gopls/go.mod", ".golangci.yml", ".vscode/settings.json",
		"doc.go", templatePackage + ".go", templatePackage + "_test.go",
	}
	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("template file missing: %w", err)
		}
	}
	return nil
}

func configure(p project) error {
	if err := replaceOnce("go.mod", "module "+templateModule+"\n", "module "+p.module+"\n"); err != nil {
		return err
	}
	for _, path := range []string{"tools/go.mod", "tools/gopls/go.mod"} {
		dir := strings.TrimSuffix(path, "/go.mod")
		if err := replaceOnce(path, "module "+templateModule+"/"+dir+"\n", "module "+p.module+"/"+dir+"\n"); err != nil {
			return err
		}
	}
	if err := replaceOnce("doc.go", templateDoc, fmt.Sprintf("// Package %s %s\n", p.packageName, p.description)); err != nil {
		return err
	}
	for _, path := range []string{".golangci.yml", ".vscode/settings.json"} {
		if err := replaceOnce(path, templateModule, p.module); err != nil {
			return err
		}
	}
	if err := configureGoFiles(p); err != nil {
		return err
	}
	readme := fmt.Sprintf("# %s\n\nPackage %s %s\n", filepath.Base(p.module), p.packageName, p.description)
	if err := os.WriteFile("README.md", []byte(readme), 0o644); err != nil {
		return fmt.Errorf("write README.md: %w", err)
	}
	return nil
}

// configureGoFiles rewrites the module path and package name in every Go
// file, then renames the seed files after the package.
func configureGoFiles(p project) error {
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != "." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "bin" || d.Name() == "tools") {
			return filepath.SkipDir
		}
		if d.IsDir() || filepath.Ext(path) != ".go" || path == "bootstrap.go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := strings.ReplaceAll(string(data), templateModule, p.module)
		text = strings.ReplaceAll(text, templatePackage, p.packageName)
		return os.WriteFile(path, []byte(text), 0o644)
	})
	if err != nil {
		return fmt.Errorf("rewrite Go files: %w", err)
	}
	for _, suffix := range []string{".go", "_test.go"} {
		if err := os.Rename(templatePackage+suffix, p.packageName+suffix); err != nil {
			return fmt.Errorf("rename seed file: %w", err)
		}
	}
	return nil
}

// replaceOnce replaces old with replacement in the file at path, failing unless old occurs exactly once.
func replaceOnce(path, old, replacement string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	text := string(data)
	if n := strings.Count(text, old); n != 1 {
		return fmt.Errorf("%s: expected one %q, found %d", path, old, n)
	}
	text = strings.Replace(text, old, replacement, 1)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
