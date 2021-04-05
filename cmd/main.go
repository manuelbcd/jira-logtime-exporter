package main

import (
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/andygrunwald/go-jira"
	"log"
	"logtimeexport/pkg/cellmanager"
	"logtimeexport/pkg/common"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type cmdLnParams struct {
	issueId    string // Get worklogs by issue-id
	userId     string // Get worklogs by user-id
	avoidExcel bool   // Avoid exporting excel file
}

const _excelFormulaTime string = "=IF(ISNUMBER(FIND(\"d\",[COL][ROW])),LEFT([COL][ROW],FIND(\"d\",[COL][ROW])-1)*24)+IF(ISNUMBER(FIND(\"h\",[COL][ROW])),MID(0&[COL][ROW],MAX(1,FIND(\"h\",0&[COL][ROW])-2),2))+IFERROR(MID(0&[COL][ROW],MAX(1,FIND(\"m\",0&[COL][ROW])-2),2)/60,0)"
const _excelFormulaCount string = "=SUMIF(TimeLog!$B:$B,Totals![COL][ROW],TimeLog!$D:$D)"

func main() {

	// Initialize current application path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Capture command line options
	Params := captureCommandLine()

	// Load configuration file
	cfg := common.LoadConfiguration(dir + "/config.json")

	// Launch Jira-gathering tasks
	if Params.issueId != "" {
		gatherJiraDataByIssueId(cfg, dir, Params.issueId)
	} else if Params.userId != "" {
		gatherJiraDataByUserId(cfg, dir, Params.userId)
	} else {
		fmt.Println("No parameters detected (Method IssueID or UserID?")
		os.Exit(1)
	}

}

/**
Capture command line arguments and return within an structure
*/
func captureCommandLine() cmdLnParams {
	issueStrPtr := flag.String("issueid", "", "Issue ID to gather log-time from. (i.e. -issueid=PROJ-21)")
	userIdPtr := flag.String("userid", "", "User ID from whom you want to extract log-time (i.e. -userid=3a8273c90fa-3b9a483720)")
	avoidExcelStrPtr := flag.Bool("avoidexcel", false, "Avoid excel file creation")
	helpPtr := flag.Bool("help", false, "Help")
	flag.Parse()

	if (*issueStrPtr == "" && *userIdPtr == "") || *helpPtr == true {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Show captured configuration values
	if *issueStrPtr != "" {
		fmt.Printf("Issue: %s", *issueStrPtr)
	}
	if *userIdPtr != "" {
		if *issueStrPtr != "" {
			fmt.Printf("Error: Can't set both issueID and userID. You must choose one of them.")
			os.Exit(1)
		}
		fmt.Printf("User ID: %s", *userIdPtr)
	}
	fmt.Printf(", AvoidExcel: %t", *avoidExcelStrPtr)
	fmt.Printf("\n")

	return cmdLnParams{
		*issueStrPtr,
		*userIdPtr,
		*avoidExcelStrPtr,
	}
}

/**
Save jira work-log to excel file and adds formulas
*/
func saveIssueWorkLogsToExcelFile(issue *jira.Worklog, dir string) {
	f := excelize.NewFile()
	f.NewSheet("TimeLog")
	f.NewSheet("Totals")
	f.DeleteSheet("Sheet1")
	var cellIndex cellmanager.Cell
	cellIndex.Init()

	var nList = common.Names{
		NameList: make([]string, 0),
	}

	// Iterate over received work-logs to insert them into excel rows
	for i := range issue.Worklogs {
		is := issue.Worklogs[i]
		nList.AddName(is.Author.DisplayName)

		f.SetCellValue("TimeLog", cellIndex.GetStr(), time.Time(*is.Started).String())
		f.SetCellValue("TimeLog", cellIndex.IncCol().GetStr(), is.Author.DisplayName)
		timeCol := cellIndex.IncCol().GetStr()
		f.SetCellValue("TimeLog", timeCol, is.TimeSpent)
		// Formula to calculate numeric time from "d h m s" format
		f.SetCellFormula(
			"TimeLog",
			cellIndex.IncCol().GetStr(),
			strings.ReplaceAll(_excelFormulaTime, "[COL][ROW]", timeCol))
		f.SetCellValue("TimeLog", cellIndex.IncCol().GetStr(), is.Comment)

		// Increment row and initialize column
		cellIndex.IncRow()
		cellIndex.InitCol()
	}

	cellIndex.Init()
	auxCellIndex := cellIndex
	auxCellIndex.IncRow()
	// Iterate over names to build "Totals" sheet
	for i := range nList.NameList {
		f.SetCellValue("Totals", cellIndex.IncCol().GetStr(), nList.NameList[i])
		f.SetCellFormula(
			"Totals",
			auxCellIndex.IncCol().GetStr(),
			strings.ReplaceAll(_excelFormulaCount, "[COL][ROW]", cellIndex.GetStr()))
	}

	// Save excel file
	if err := f.SaveAs(dir + "/Book1.xlsx"); err != nil {
		println(err.Error())
	} else {
		fmt.Printf("\nFile created with success\n")
	}
}



/**
Connect to Jira, extract log-time details from a specific issue and
export it to an Excel file.
*/
func gatherJiraDataByIssueId(cfg common.Config, dir string, issueid string) {

	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.Login,
		Password: cfg.Jira.Token,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), cfg.Jira.Host)

	// Removed dateStart property since Jira API has a known bug. TODO: reactivate once the bug is fixed.
	// dateStart := int64(time.Now().Unix())
	// var op * jira.GetWorklogsQueryOptions = &jira.GetWorklogsQueryOptions{Expand: "properties", StartedAfter: dateStart}

	var op *jira.GetWorklogsQueryOptions = &jira.GetWorklogsQueryOptions{Expand: "properties"}
	workLogs, _, err := jiraClient.Issue.GetWorklogs(issueid, jira.WithQueryOptions(op))

	if err != nil {
		fmt.Printf("Error requesting worklogs from issue. Error: %s \n", err.Error())
		os.Exit(1)
	}

	saveIssueWorkLogsToExcelFile(workLogs, dir)
}

/**
Connect to Jira, extract log-time details from a specific user and export it to an Excel file.
(The request of work-logs from user is made trough JQL query)
*/
func gatherJiraDataByUserId(cfg common.Config, dir string, userId string) {

	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.Login,
		Password: cfg.Jira.Token,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), cfg.Jira.Host)

	jql := fmt.Sprintf("worklogDate >= -365d and worklogAuthor = %s order by updatedDate DESC", userId)

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

	fmt.Printf("\n TOTAL ELEMENTS: %d \n", resp.Total)

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

