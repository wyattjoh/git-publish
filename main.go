package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	noFetch := flag.Bool("no-fetch", false, "disable fetching the remote")

	flag.Parse()

	remoteName := "origin"
	if flag.NArg() == 1 {
		remoteName = flag.Arg(0)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get the current working directory: %v\n", err)
		os.Exit(1)
	}

	// Validate that the remote exists.
	out, err := exec.Command("git", "remote").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't validate the remote: %v\n", err)
		os.Exit(1)
	}

	var foundRemote bool
	for _, remote := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if remote == remoteName {
			foundRemote = true
		}
	}
	if !foundRemote {
		fmt.Fprintf(os.Stderr, "can't find the remote %s in %s\n", remoteName, cwd)
		os.Exit(1)
	}

	if !*noFetch {
		// Update remotes.
		fmt.Println("$ git fetch")
		if err := exec.Command("git", "fetch", remoteName).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "can't fetch remote %s: %v\n", remoteName, err)
			os.Exit(1)
		}
	}

	// Get the current branch name.
	fmt.Println("$ git rev-parse --abbrev-ref HEAD")
	out, err = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get current branch: %v\n", err)
		os.Exit(1)
	}

	// Strip off the extra space and newlines.
	activeBranch := strings.TrimSpace(string(out))

	fmt.Printf("$ git push --set-upstream %s %s\n", remoteName, activeBranch)
	cmd := exec.Command("git", "push", "--set-upstream", remoteName, activeBranch)

	// Redirect all the stdout/stderr to the output of this process.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run it.
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "can't push upstream: %v\n", err)
		os.Exit(1)
	}
}
