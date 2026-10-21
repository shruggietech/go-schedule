// Command release-promotion binds standing-waiver release operations to one authenticated identity.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

func main() { os.Exit(run(os.Args[1:], os.Getenv, os.Stdout, os.Stderr)) }

func run(args []string, getenv func(string) string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "Usage: release-promotion resolve|download|ensure-checksum|publish --state FILE [options]")
		return 2
	}
	command := args[0]
	if command != "resolve" && command != "download" && command != "ensure-checksum" && command != "publish" {
		fmt.Fprintln(stderr, "Choose resolve, download, ensure-checksum or publish.")
		return 2
	}
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	statePath := flags.String("state", "", "Temporary release-selection JSON file")
	repository := flags.String("repository", "", "Repository owner/name (resolve)")
	tag := flags.String("tag", "", "Exact vMAJOR.MINOR.PATCH release tag (resolve)")
	directory := flags.String("dir", "", "New directory for an initial or final download")
	checksum := flags.String("file", "", "Generated SHA256SUMS.txt to compare or upload")
	if err := flags.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 || *statePath == "" || (command == "resolve" && (!repositoryPattern.MatchString(*repository) || !tagPattern.MatchString(*tag))) || (command == "download" && *directory == "") || (command == "ensure-checksum" && *checksum == "") {
		fmt.Fprintln(stderr, "Supply --state and the required command options; use -h for help.")
		return 2
	}
	token := getenv("GH_TOKEN")
	if token == "" {
		fmt.Fprintln(stderr, "GH_TOKEN is required to read authenticated drafts and update the selected release.")
		return 1
	}
	api := getenv("GITHUB_API_URL")
	if api == "" {
		api = "https://api.github.com"
	}
	client, err := newClient(api, token, &http.Client{Timeout: 5 * time.Minute})
	if err == nil && client.base.Scheme != "https" {
		err = fmt.Errorf("GITHUB_API_URL must use HTTPS")
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	if command == "resolve" {
		err = client.resolve(ctx, *repository, *tag, *statePath)
	} else {
		var selected selection
		selected, err = loadSelection(*statePath)
		if err == nil {
			switch command {
			case "download":
				err = client.download(ctx, &selected, *directory, *statePath)
			case "ensure-checksum":
				err = client.ensureChecksum(ctx, &selected, *checksum, *statePath)
			case "publish":
				err = client.publish(ctx, &selected)
			}
		}
	}
	if err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", command, err)
		return 1
	}
	selected, loadErr := loadSelection(*statePath)
	if loadErr != nil {
		fmt.Fprintln(stderr, loadErr)
		return 1
	}
	fmt.Fprintf(stdout, "%s completed for release %d (%s).\n", command, selected.ReleaseID, selected.Tag)
	return 0
}

var repositoryPattern = regexp.MustCompile("^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$")
var tagPattern = regexp.MustCompile("^v[0-9]+[.][0-9]+[.][0-9]+$")
