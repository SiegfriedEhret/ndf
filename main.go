package main

import (
	"fmt"

	"flag"
	"os"

	"bytes"
	"io/ioutil"

	"github.com/google/go-github/github"
	"strings"
)

const (
	appName    string = "ndf"
	appVersion string = "0.1.0"
	helpText   string = `=========
Flags:
	--help		Show this help :-)
	--version	To display the current version
Actually do something:
	--owner		The name of the repository owner
	--repo		The name of the repository
	--milestone	The id of the milestone
	--close		Close issues and milestone`
)

var (
	owner       string
	repo        string
	milestone   string
	closeIssues bool

	help    bool
	version bool
)

func init() {
	flag.StringVar(&owner, "owner", "SiegfriedEhret", "Set the Github username")
	flag.StringVar(&repo, "repo", "ndf", "Set the Github repository")
	flag.StringVar(&milestone, "milestone", "1", "Set milestone to release")
	flag.BoolVar(&closeIssues, "close", false, "Close things")

	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&version, "version", false, "Show version")
	flag.Parse()
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

	labels, _, err := client.Issues.ListLabels(owner, repo, nil)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(labels)

	var md bytes.Buffer

	for _, label := range labels {
		labelName := *label.Name

		if strings.HasPrefix(labelName, "type/") {
			continue
		}

		opts := &github.IssueListByRepoOptions{
			Milestone: "1",
			Labels:    []string{labelName},
		}

		issues, _, err := client.Issues.ListByRepo(owner, repo, opts)

		if err != nil {
			fmt.Println(err)
		} else if len(issues) == 0 {
			fmt.Println("No issues for label", labelName)
		} else {
			fmt.Println(labelName)

			md.WriteString("## " + labelName + "\n")

			for _, issue := range issues {
				body := *issue.Body
				title := *issue.Title

				for _, issueLabel := range issue.Labels {
					switch *issueLabel.Name {
					case "type/link":
						md.WriteString("- " + title + ": " + body + "\n")
					case "type/text":
						md.WriteString(title + "\n\n" + body + "\n")
					}
				}
			}

			err := ioutil.WriteFile("./ndf-" + milestone + ".md", md.Bytes(), 0644)

			if err != nil {
				fmt.Println("Error while creating file", err)
			}
		}
	}

	fmt.Println(md.String())

	//issues, _, err := client.Issues.ListByRepo(owner, repo, nil)
	//
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(issues)
	//}
}

//func getMilestones(client *github.Client) ([]*github.Milestone, *github.Response, error) {
//	return client.Issues.ListMilestones(owner, repo, nil)
//}
