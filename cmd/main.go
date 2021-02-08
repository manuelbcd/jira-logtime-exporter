package main

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/andygrunwald/go-jira"
	"log"
	"logtimeexport/pkg/cellmanager"
	"logtimeexport/pkg/common"
	"os"
	"path/filepath"
	"time"
)

func main() {

	// Initialize current path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Load configuration file
	cfg := common.LoadConfiguration(dir + "/config.json")

	// Launch Jira gathering tasks
	gatherJiraData(cfg)
}

func gatherJiraData(cfg common.Config) {

	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.Login,
		Password: cfg.Jira.Token,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), cfg.Jira.Host)

	var op * jira.AddWorklogQueryOptions = &jira.AddWorklogQueryOptions{Expand: "properties"}

	issue, _, err := jiraClient.Issue.GetWorklogs("ORANGE-1", jira.WithQueryOptions(op))

	if err != nil {
		panic(err)
	}

	f := excelize.NewFile()
	var cellIndex cellmanager.Cell
	cellIndex.Init()

	for i := range issue.Worklogs {
		is := issue.Worklogs[i]

		f.SetCellValue("Sheet1", cellIndex.GetStr(), time.Time(*is.Updated).String())
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

