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
	issueId		string
	avoidExcel	bool
}

const EXCELFORMULA_TIME string = "=IF(ISNUMBER(FIND(\"d\",[COL][ROW])),LEFT([COL][ROW],FIND(\"d\",[COL][ROW])-1)*24)+IF(ISNUMBER(FIND(\"h\",[COL][ROW])),MID(0&[COL][ROW],MAX(1,FIND(\"h\",0&[COL][ROW])-2),2))+IFERROR(MID(0&[COL][ROW],MAX(1,FIND(\"m\",0&[COL][ROW])-2),2)/60,0)"
const EXCELFORMULA_CONT string = "=SUMIF(TimeLog!$B:$B,Totals![COL][ROW],TimeLog!$D:$D)"

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
	Save jira work-log to excel file and adds formulas
 */
func saveToExcelFile(issue * jira.Worklog, dir string) {
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
			strings.ReplaceAll(EXCELFORMULA_TIME, "[COL][ROW]", timeCol))
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
			strings.ReplaceAll(EXCELFORMULA_CONT, "[COL][ROW]", cellIndex.GetStr()))
	}

	// Save excel file
	if err := f.SaveAs(dir + "/Book1.xlsx"); err != nil {
		println(err.Error())
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

	saveToExcelFile(issue, dir)
}

