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
	"time"
)

type cmdLnParams struct {
	issueId		string
	avoidExcel	bool
}

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

	// Launch Jira gathering tasks
	gatherJiraData(cfg, dir, Params.issueId)
}

/**
	Capture command line arguments and return within an structure
 */
func captureCommandLine() cmdLnParams {
	issuePtr := flag.String("issue", "", "Issue ID to gather log-time")
	avoidExcelPtr := flag.Bool("avoidexcel", false, "Avoid excel file creation" )
	helpPtr := flag.Bool("help", false, "Help")
	flag.Parse()

	if *issuePtr == "" || *helpPtr == true {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Print captured values
	fmt.Printf("Issue: %s", *issuePtr)
	fmt.Printf(", AvoidExcel: %t", *avoidExcelPtr)
	fmt.Printf("\n")

	return cmdLnParams{
		*issuePtr,
		*avoidExcelPtr,
	}
}

/**
	Connect to Jira, extract log-time details from a specific issue and
	export it to an Excel file.
 */
func gatherJiraData(cfg common.Config, dir string, issueid string) {

	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.Login,
		Password: cfg.Jira.Token,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), cfg.Jira.Host)

	// Removed dateStart property since Jira API has a known bug. TODO: reactivate once the bug is fixed.
	// dateStart := int64(time.Now().Unix())
	// var op * jira.GetWorklogsQueryOptions = &jira.GetWorklogsQueryOptions{Expand: "properties", StartedAfter: dateStart}

	var op * jira.GetWorklogsQueryOptions = &jira.GetWorklogsQueryOptions{Expand: "properties"}
	issue, _, err := jiraClient.Issue.GetWorklogs(issueid, jira.WithQueryOptions(op))

	if err != nil {
		fmt.Printf("issuePtr: %s \n",err.Error())
		os.Exit(1)
	}

	f := excelize.NewFile()
	var cellIndex cellmanager.Cell
	cellIndex.Init()

	for i := range issue.Worklogs {
		is := issue.Worklogs[i]

		f.SetCellValue("Sheet1", cellIndex.GetStr(), time.Time(*is.Started).String())
		f.SetCellValue("Sheet1", cellIndex.IncCol().GetStr(), is.Author.DisplayName)
		f.SetCellValue("Sheet1", cellIndex.IncCol().GetStr(), is.TimeSpent)
		f.SetCellValue("Sheet1", cellIndex.IncCol().GetStr(), is.Comment)

		cellIndex.IncRow()
		cellIndex.InitCol()
	}

	// Save excel file
	if err := f.SaveAs(dir + "/Book1.xlsx"); err != nil {
		println(err.Error())
	}
}

