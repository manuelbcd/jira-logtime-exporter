package app

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"logtimeexport/pkg/common"
	"os"
)

/**
Connect to Jira, extract log-time details from a specific issue and
export it to an Excel file.
*/
func GatherJiraDataByIssueID(cfg common.Config, dir string, issueID string) {

	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.Login,
		Password: cfg.Jira.Token,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), cfg.Jira.Host)

	// Removed dateStart property since Jira API has a known bug. TODO: reactivate once the bug is fixed.
	// dateStart := int64(time.Now().Unix())
	// var op * jira.GetWorklogsQueryOptions = &jira.GetWorklogsQueryOptions{Expand: "properties", StartedAfter: dateStart}

	var op *jira.GetWorklogsQueryOptions = &jira.GetWorklogsQueryOptions{Expand: "properties"}
	workLogs, _, err := jiraClient.Issue.GetWorklogs(issueID, jira.WithQueryOptions(op))

	if err != nil {
		fmt.Printf("Error requesting worklogs from issue. Error: %s \n", err.Error())
		os.Exit(1)
	}

	f := initExcelFile()
	saveIssueWorkLogsToExcelFile(issueID, workLogs, f)
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

		for _, x := range tempWorkLog.Worklogs {
			workLogs = append(workLogs, x)
		}
	}

	for _, z := range workLogs {
		fmt.Printf("%s : %+v\n", z.Comment, z.TimeSpentSeconds)
	}

}
