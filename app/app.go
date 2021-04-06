package app

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"logtimeexport/pkg/common"
	"os"
	"sort"
	"time"
)


/**
Connect to Jira, extract log-time details from a specific issue and
export it to an Excel file.
*/
func GatherJiraDataByIssueKey(cfg common.Config, dir string, issueKey string) {

	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.Login,
		Password: cfg.Jira.Token,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), cfg.Jira.Host)


	var issueList []jira.Issue
	issue, _, err := jiraClient.Issue.Get(issueKey, nil)
	if err != nil {
		fmt.Printf("Error requesting issue before requesting worklogs. Error: %s \n", err.Error())
		os.Exit(1)
	}
	issueList = append(issueList, *issue)

	// Retrieve worklogs from issue key
	// Please note that Jira API has a known bug and it is not possible to filter worklogs by date.
	// TODO: Add date filtering once the bug is fixed.
	var op *jira.GetWorklogsQueryOptions = &jira.GetWorklogsQueryOptions{Expand: "properties"}
	workLogs, _, err := jiraClient.Issue.GetWorklogs(issueKey, jira.WithQueryOptions(op))

	if err != nil {
		fmt.Printf("Error requesting worklogs from issue. Error: %s \n", err.Error())
		os.Exit(1)
	}

	f := initExcelFile()
	saveIssueWorkLogsToExcelFile(cfg, workLogs.Worklogs, issueList, f)
	saveExcelFile(dir, f)
}

/**
Connect to Jira, extract log-time details from a specific user and export it to an Excel file.
(The request of work-logs from user is made trough JQL query)
*/
func GatherJiraDataByUserID(cfg common.Config, dir string, userID string) {

	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.Login,
		Password: cfg.Jira.Token,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), cfg.Jira.Host)

	jql := fmt.Sprintf("worklogDate >= -365d and worklogAuthor = %s order by updatedDate DESC", userID)

	options := &jira.SearchOptions{Fields: []string{"summary", "status", "assignee", "worklog"}}
	issues, resp, err := jiraClient.Issue.Search(jql, options)

	if err != nil {
		fmt.Printf("Error requesting worklogs from user id. Error: %s \n", err.Error())
		os.Exit(1)
	}

	// TODO Check WarningMessages in response body
	// It is not implemented yet by go-jira
	// Issue https://github.com/andygrunwald/go-jira/issues/368

	if resp.Total == 0 {
		fmt.Printf("No worklogs found for that user id.\n")
		os.Exit(1)
	}

	var op *jira.GetWorklogsQueryOptions = &jira.GetWorklogsQueryOptions{Expand: "properties"}

	fmt.Printf("\n Extracting worklogs from %d issues.\n", resp.Total)

	var workLogs []jira.WorklogRecord

	// Collect all work-logs from retrieved issues
	for _, i := range issues {
		fmt.Printf("%s : %+v\n", i.Key, i.Fields.Summary)

		tempWorkLog, _, err := jiraClient.Issue.GetWorklogs(i.ID, jira.WithQueryOptions(op))
		if err != nil {
			fmt.Printf("Error requesting worklogs from issue. Error: %s \n", err.Error())
			os.Exit(1)
		}

		// Check every single worklog to discard users different from specified one
		// Note that Jira still does not allow to filter on worklog requests
		for _, x := range tempWorkLog.Worklogs {
			if x.Author.AccountID == userID {
				workLogs = append(workLogs, x)
			}
		}
	}

	// Sort worklogs by Started datetime
	sort.Slice(workLogs, func(i, j int) bool {
		return time.Time(*workLogs[i].Started).Before(time.Time(*workLogs[j].Started))
	})

	f := initExcelFile()
	saveIssueWorkLogsToExcelFile(cfg, workLogs, issues, f)
	saveExcelFile(dir, f)
}
