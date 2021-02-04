package main

import (
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Test string `json:"test"`
	Jira struct {
		Host    	string `json:"host"`
		Login 		string `json:"login"`
		Token 		string `json:"token"`
	} `json:"jira"`
}

type cell struct {
	row 		int
	col 		string
	colIndex	int
}

const LETTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (c * cell) init() {
	c.colIndex = 0
	c.col = string(LETTERS[c.colIndex])
	c.row = 1
}

func (c * cell) initCol() {
	c.colIndex = 0
	c.col = string(LETTERS[c.colIndex])
}

func (c * cell) incRow()  {
	c.row ++
}

func (c * cell) incCol() * cell {
	c.colIndex ++
	if c.colIndex >= len(LETTERS){
		c.colIndex = 0
	}
	c.col = string(LETTERS[c.colIndex])
	return c
}

func (c * cell) getStr() string {
	return fmt.Sprintf("%s%d", c.col, c.row)
}

func LoadConfiguration(file string) Config {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	config := Config{}
	jsonParser.Decode(&config)

	return config
}

func main() {

	// Initialize current path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Load configuration file
	cfg := LoadConfiguration(dir + "/config.json")

	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.Login,
		Password: cfg.Jira.Token,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), cfg.Jira.Host)

	/*
	issue, _, err := jiraClient.Issue.Get("ORANGE-1", nil)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
	*/

	var op * jira.AddWorklogQueryOptions = &jira.AddWorklogQueryOptions{Expand: "properties"}


	issue, _, err := jiraClient.Issue.GetWorklogs("ORANGE-1", jira.WithQueryOptions(op))

	if err != nil {
		panic(err)
	}

	f := excelize.NewFile()

	var cellIndex cell
	cellIndex.init()

	for i := range issue.Worklogs {
		is := issue.Worklogs[i]

		f.SetCellValue("Sheet1", cellIndex.getStr(), time.Time(*is.Updated).String())
		f.SetCellValue("Sheet1", cellIndex.incCol().getStr(), is.Author.DisplayName)
		f.SetCellValue("Sheet1", cellIndex.incCol().getStr(), is.TimeSpent)
		f.SetCellValue("Sheet1", cellIndex.incCol().getStr(), is.Comment)

		cellIndex.incRow()
		cellIndex.initCol()
		//fmt.Printf("Author: %s Time: %s \n", is.Author, is.TimeSpent)
	}

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		println(err.Error())
	}

	// MESOS-3325: Running mesos-slave@0.23 in a container causes slave to be lost after a restart
	// Type: Bug
	// Priority: Critical
}

