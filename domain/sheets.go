package domain

import (
	"fmt"
)

func (d Domain) GetSpreadsheetData(spreadsheetID, readRange string) ([][]interface{}, error) {
	resp, err := d.client.Service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		return nil, nil
	}

	return resp.Values, nil
}
