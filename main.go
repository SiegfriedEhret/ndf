package main

import (
	"fmt"

	"flag"
	"os"

	"bytes"
	"io/ioutil"

	"sort"
	"strings"

	"context"
	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
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
	-milestone	The id of the milestone
	-close		Close issues and milestone`
)

var (
	owner       string
	repo        string
	milestone   string
	closeIssues bool
	debug       bool

	help    bool
	version bool
)

func init() {
	flag.StringVar(&owner, "owner", "SiegfriedEhret", "Set the Github username")
	flag.StringVar(&repo, "repo", "ndf", "Set the Github repository")
	flag.StringVar(&milestone, "milestone", "1", "Set milestone to release")
	flag.BoolVar(&closeIssues, "close", false, "Close things")
	flag.BoolVar(&debug, "d", false, "Run in debug mode")

	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&version, "version", false, "Show version")
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
}

func main() {
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
	fmt.Printf("%s %s\n", appName, appVersion)

	client := github.NewClient(nil)

	labels, _, err := client.Issues.ListLabels(context.Background(), owner, repo, nil)

	if err != nil {
		fmt.Println(err)
	}

	sort.Slice(labels, func(i, j int) bool {
		return *labels[i].Name < *labels[j].Name
	})

	var md bytes.Buffer

	for _, label := range labels {
		labelName := *label.Name

		if strings.HasPrefix(labelName, "type/") {
			continue
		}

		opts := &github.IssueListByRepoOptions{
			Milestone: milestone,
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

			md.WriteString("## " + labelToDisplay + "\n")

			for _, issue := range issues {
				body := *issue.Body
				title := *issue.Title

				logrus.Debug(body, title)

				for _, issueLabel := range issue.Labels {
					switch *issueLabel.Name {
					case "type/link":
						md.WriteString("- " + title + ": " + body + "\n")
					case "type/text":
						md.WriteString(title + "\n\n" + body + "\n")
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
