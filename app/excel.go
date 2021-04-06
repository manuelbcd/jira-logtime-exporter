package app

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/andygrunwald/go-jira"
	"logtimeexport/pkg/cellmanager"
	"logtimeexport/pkg/common"
	"strings"
	"time"
)

const _excelFormulaTime string = "=IF(ISNUMBER(FIND(\"d\",[COL][ROW])),LEFT([COL][ROW],FIND(\"d\",[COL][ROW])-1)*24)+IF(ISNUMBER(FIND(\"h\",[COL][ROW])),MID(0&[COL][ROW],MAX(1,FIND(\"h\",0&[COL][ROW])-2),2))+IFERROR(MID(0&[COL][ROW],MAX(1,FIND(\"m\",0&[COL][ROW])-2),2)/60,0)"
const _excelFormulaCount string = "=SUMIF(TimeLog!$C:$C,Totals![COL][ROW],TimeLog!$E:$E)"

func initExcelFile() *excelize.File {
	f := excelize.NewFile()
	f.NewSheet("TimeLog")
	f.NewSheet("Totals")
	f.NewSheet("Issues")
	f.DeleteSheet("Sheet1")

	return f
}

// saveExcelFile stores excelince.File to disk
func saveExcelFile(dir string, f *excelize.File) {
	// Save excel file
	if err := f.SaveAs(dir + "/Book1.xlsx"); err != nil {
		println(err.Error())
	} else {
		fmt.Printf("\nFile created with success\n")
	}
}

// Save jira work-log to excel file and adds formulas
func saveIssueWorkLogsToExcelFile(workLogs []jira.WorklogRecord, issueList []jira.Issue, f *excelize.File) {
	var cellIndex cellmanager.Cell
	cellIndex.Init()

	var nList = common.Names{
		NameList: make([]string, 0),
	}

	// TimeLog sheet.
	// Iterate over received work-logs to insert them into excel rows
	for i := range workLogs {
		is := workLogs[i]
		nList.AddName(is.Author.DisplayName)

		f.SetCellValue("TimeLog", cellIndex.GetStr(), time.Time(*is.Started).String())
		f.SetCellValue("TimeLog", cellIndex.IncCol().GetStr(), findIssueKeyByID(issueList, is.IssueID))
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

	// Totals sheet.
	// Iterate over names to build "Totals" sheet
	for i := range nList.NameList {
		f.SetCellValue("Totals", cellIndex.IncCol().GetStr(), nList.NameList[i])
		f.SetCellFormula(
			"Totals",
			auxCellIndex.IncCol().GetStr(),
			strings.ReplaceAll(_excelFormulaCount, "[COL][ROW]", cellIndex.GetStr()))
	}

	cellIndex.Init()

	// Issues sheet
	// Iterate over names to build "Totals" sheet
	for _, i := range issueList {
		f.SetCellValue("Issues", cellIndex.GetStr(), i.Key)
		f.SetCellValue("Issues", cellIndex.IncCol().GetStr(), i.Fields.Summary)
		// Increment row and initialize column
		cellIndex.IncRow()
		cellIndex.InitCol()
	}

}

// findIssueKeyByID seeks an issue Key by issue ID and returns the Key if found
func findIssueKeyByID(issueList []jira.Issue, issueID string) string {
	result := ""
	for _, i := range issueList {
		if i.ID == issueID {
			result = i.Key
			break
		}
	}
	return result
}