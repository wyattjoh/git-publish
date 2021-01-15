package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Action will perform the publish operation.
func Action(c *cli.Context) error {
	noFetch := c.Bool("no-fetch")

	remoteName := "origin"
	if flag.NArg() == 1 {
		remoteName = flag.Arg(0)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "can't get the current working directory")
	}

	// Validate that the remote exists.
	out, err := exec.Command("git", "remote").Output()
	if err != nil {
		return errors.Wrap(err, "can't validate that the remote exists")
	}

	var foundRemote bool

	for _, remote := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if remote == remoteName {
			foundRemote = true
		}
	}

	if !foundRemote {
		return errors.Errorf("can't find the remote %s in %s", remoteName, cwd)
	}

	if !noFetch {
		// Update remotes.
		fmt.Println("$ git fetch")

		if err := exec.Command("git", "fetch", remoteName).Run(); err != nil {
			return errors.Wrap(err, "can't fetch remote")
		}
	}

	// Get the current branch name.
	fmt.Println("$ git rev-parse --abbrev-ref HEAD")

	out, err = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return errors.Wrap(err, "can't get the current branch")
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
		return errors.Wrap(err, "can't push upstream")
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "git-publish"
	app.Version = fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "no-fetch",
			Usage: "disable featching the remote",
		},
	}
	app.Action = Action

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "can't publish to origin: %v\n", err)
		os.Exit(1)
	}
}
