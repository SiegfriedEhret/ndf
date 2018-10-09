package githoub

import (
	"context"
	"errors"
	"sort"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func GetGithubClient(token string) *github.Client {
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "... your access token ..."},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		return github.NewClient(tc)
	} else {
		return github.NewClient(nil)
	}
}

func GetMilestone(client *github.Client, owner, repo, milestone string) (error, string) {
	milestones, _, err := client.Issues.ListMilestones(context.Background(), owner, repo, nil)

	if err != nil {
		logrus.Fatal("Couldn't get milestones", err)
	} else {
		logrus.Debug("Milestones", milestones)
	}

	for _, mst := range milestones {
		logrus.Debug(mst.Title, milestone)
		if *mst.Title == milestone {
			milestoneId := strconv.Itoa(*mst.Number)
			logrus.Debugf("Found milestone ! %s", milestoneId)
			return nil, milestoneId
		}
	}

	return errors.New("Milestone not found"), ""

}

func GetLabels(client *github.Client, owner, repo string) (error, []*github.Label) {
	labels, _, err := client.Issues.ListLabels(context.Background(), owner, repo, nil)

	if err != nil {
		return err, []*github.Label{}
	}

	sort.Slice(labels, func(i, j int) bool {
		return *labels[i].Name < *labels[j].Name
	})

	return nil, labels
}
