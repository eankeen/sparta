package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
)

// Extension contains information about each extension. This may be autogenerated
type Extension struct {
	Name string   `toml:"name"`
	Tags []string `toml:"tags"`
}

// Workspace is the aggregation of extensions which is meant to be used for a specific programming language or other subdomain
type Workspace struct {
	Name string   `toml:"name"`
	Use  []string `toml:"use"`
}

// Config gives information about the whole configuration file
type Config struct {
	Version    string      `toml:"version"`
	Extensions []Extension `toml:"extensions"`
	Workspaces []Workspace `toml:"workspaces"`
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}

// tests if an extension has a version
func extensionHasVersion(str string) bool {
	str = str[len(str)-1:]

	if strings.Contains("1234567890", str) {
		return true
	}

	return false
}

func isFolderEmpty(path string) bool {
	dirs, err := ioutil.ReadDir(path)
	p(err)

	if len(dirs) == 0 {
		return true
	}

	return false
}

func readConfig(opts Options) Config {
	var config Config

	configRaw, err := ioutil.ReadFile(filepath.Join(opts.ConfigFile))
	p(err)

	err = toml.Unmarshal(configRaw, &config)
	p(err)

	return config
}

// returns array of extensions
// example: ["yzhang.markdown-all-in-one@3.3.0"]
func getVscodeExtensions() []string {
	cmd := exec.Command("code", "--list-extensions")

	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	p(err)

	return strings.Split(string(stdout), "\n")
}

func ensureLength(arr []string, minLength int, message string) {
	if len(arr) < minLength {
		log.Fatalln(message)
	}
}

func contains(arr []string, query string) bool {
	for _, el := range arr {
		if el == query {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Println(`sparta

Description:
  Contextual vscode extension management

Commands:
  init
    Initiates an 'extensions.toml' folder that contains all extensions for tagging

  update
    Updates the extensions and resymlinks them

  check
    Prints all extensions mismatches between default globally installed and ones defined in extensions.toml

  launch [workspace]
	 Launches a particular workspace in vscode

  plumbing download-extensions
	 Redownloads your current extensions into a directory

  plumbing remove-extensions
    Removes all current extensions

  plumbing symlink-extensions
	 For each extension in a workspace, symlink it to the extensions downloaded

  plumbing remove-symlinks
    Remove all symlinks`)
}
