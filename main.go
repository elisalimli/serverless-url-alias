package main

import (
	"context"
	"fmt"
	"log"

	"github.com/elisalimli/serverless-url-alias/domain"
	"github.com/elisalimli/serverless-url-alias/internal/sheets"
)

func main() {
	ctx := context.Background()

	client, err := sheets.NewClient(ctx)
	domain := domain.NewDomain(*client)

	if err != nil {
		log.Fatalf("error occured while creating new service : %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := "1Czrh54YBicnvTv6MrIqQ1CsLKlFf9D07EtKQNXTb6w0"
	readRange := "Sheet1!A:B"
	data, err := domain.GetSpreadsheetData(spreadsheetId, readRange)

	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Major:")
		for _, row := range data {
			// Print columns A and E, which correspond to indices 0 and 4.
			fmt.Printf("%#v %#v\n", row[0], row[1])
		}
	}
}
