package main

import (
	"flag"
	"fmt"
	"log"
	"logtimeexport/app"
	"logtimeexport/pkg/common"
	"os"
	"path/filepath"
)

type cmdLnParams struct {
	issueId    string // Get worklogs by issue-id
	userId     string // Get worklogs by user-id
	avoidExcel bool   // Avoid exporting excel file
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

	// Launch Jira-gathering tasks
	if Params.issueId != "" {
		app.GatherJiraDataByIssueId(cfg, dir, Params.issueId)
	} else if Params.userId != "" {
		app.GatherJiraDataByUserId(cfg, dir, Params.userId)
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

