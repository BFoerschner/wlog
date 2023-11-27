package main

import (
	"context"
	jira "github.com/andygrunwald/go-jira/v2/onpremise"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	jiraURL := os.Args[1]
	token := os.Args[2]
	filePath := os.Args[3]

	tp := jira.BearerAuthTransport{Token: token}
	client, err := jira.NewClient(jiraURL, tp.Client())
	if err != nil {
		panic(err)
	}

	var fileContent, _ = os.ReadFile(filePath)
	addWorklogs(client, string(fileContent))
}

// yyyy-mm-dd TICKETID-1234 DURATION-IN-H COMMENT
func getWorklogRecord(recordString string) jira.WorklogRecord {
	var (
		splitRecordString   = strings.Fields(recordString)
		splitDate           = strings.Split(splitRecordString[0], "-")
		year, _             = strconv.Atoi(splitDate[0])
		month, _            = strconv.Atoi(splitDate[1])
		day, _              = strconv.Atoi(splitDate[2])
		ticketID            = splitRecordString[1]
		timeSpentInHours, _ = strconv.Atoi(splitRecordString[2])
		timeSpentInSeconds  = timeSpentInHours * 60 * 60
		comment             = strings.Join(splitRecordString[3:], " ")
		timeStarted         = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	)

	return jira.WorklogRecord{
		Comment:          comment,
		Started:          (*jira.Time)(&timeStarted),
		TimeSpentSeconds: timeSpentInSeconds,
		IssueID:          ticketID,
	}
}

func getWorklogRecords(recordStrings string) []jira.WorklogRecord {
	var records []jira.WorklogRecord
	for _, recordString := range strings.Split(recordStrings, "\n") {
		records = append(records, getWorklogRecord(recordString))
	}

	return records
}

func addWorklogs(client *jira.Client, recordString string) {
	for _, record := range getWorklogRecords(recordString) {
		client.Issue.AddWorklogRecord(
			context.Background(),
			record.IssueID,
			&record,
		)
	}
}
