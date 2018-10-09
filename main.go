package main

import (
	"fmt"

	"flag"
	"os"

	"bytes"
	"io/ioutil"

	"strings"

	"context"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"gitlab.com/SiegfriedEhret/ndf/githoub"
)

const (
	appName    string = "ndf"
	appVersion string = "0.1.0"
	helpText   string = `=========
Flags:
	-help		Show this help :-)
	-version	To display the current version

Actually do something:
	-owner		The name of the repository owner
	-repo		The name of the repository
	-token		The Github token, see https://github.com/settings/tokens
	-milestone	The id of the milestone
	-close		Close issues and milestone`
)

var (
	owner       string
	repo        string
	milestone   string
	token       string
	closeIssues bool
	debug       bool

	help    bool
	version bool
)

func init() {
	flag.StringVar(&owner, "owner", "SiegfriedEhret", "Set the Github username")
	flag.StringVar(&repo, "repo", "ndf", "Set the Github repository")
	flag.StringVar(&milestone, "milestone", "test1", "Set milestone to release")
	flag.StringVar(&token, "token", "", "Set Github access token (https://github.com/settings/tokens)")
	flag.BoolVar(&closeIssues, "close", false, "Close milestones and issues")
	flag.BoolVar(&debug, "d", false, "Run in debug mode")

	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&version, "version", false, "Show version")
}

func main() {
	flag.Parse()

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.WithFields(logrus.Fields{
			"owner":       owner,
			"repo":        repo,
			"milestone":   milestone,
			"closeIssues": closeIssues,
			"debug":       debug,
		}).Debug("Cli flags")
	}

	if help {
		fmt.Printf("%s %s\n", appName, appVersion)
		fmt.Println(helpText)
	} else if version {
		fmt.Printf("%s %s\n", appName, appVersion)
		os.Exit(0)
	} else {
		doThings()
	}
}

func doThings() {
	client := githoub.GetGithubClient(token)

	err, milestoneId := githoub.GetMilestone(client, owner, repo, milestone)

	if err != nil {
		logrus.Fatal(err.Error())
	}

	err, labels := githoub.GetLabels(client, owner, repo)

	if err != nil {
		logrus.Debug("Failed to get labels")
		logrus.Fatal(err.Error())
	}

	var md bytes.Buffer

	for _, label := range labels {
		labelName := *label.Name

		if strings.HasPrefix(labelName, "type/") {
			continue
		}

		opts := &github.IssueListByRepoOptions{
			Milestone: milestoneId,
			Labels:    []string{labelName},
		}

		issues, _, err := client.Issues.ListByRepo(context.Background(), owner, repo, opts)

		if err != nil {
			fmt.Println(err)
		} else if len(issues) == 0 {
			fmt.Println("No issues for label", labelName)
		} else {
			logrus.Debug(labelName)

			labelToDisplay := labelName
			index := strings.Index(labelToDisplay, "/")
			if index != -1 {
				labelToDisplay = labelToDisplay[index+1:]
			}

			md.WriteString("## " + labelToDisplay + "\n\n")

			for _, issue := range issues {
				body := *issue.Body
				title := *issue.Title

				logrus.Debug(body, title)

				for _, issueLabel := range issue.Labels {
					switch *issueLabel.Name {
					case "type/link":
						md.WriteString("- " + title + ": " + body + "\n\n")
					case "type/text":
						md.WriteString(title + "\n\n" + body + "\n\n")
					}
				}
			}

			err := ioutil.WriteFile("./ndf-"+milestone+".md", md.Bytes(), 0644)

			if err != nil {
				fmt.Println("Error while creating file", err)
			}
		}
	}

	logrus.Debug(md.String())
}
