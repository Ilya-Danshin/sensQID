package app

import (
	"container/list"
	"errors"
	"math/rand"

	"sensQID/internal/pkg/database"
)

func (info *anonInfo) freql1lqDiv(db *database.DB) error {

	err := db.CreateAnonTable(info.table, info.columns)
	if err != nil {
		return err
	}

	columns, err := db.GetTableColumns(info.table)
	if err != nil {
		return err
	}

	// Get number of all rows in DB
	tableVolume, err := db.GetTableVolume(info.table)
	if err != nil {
		return err
	}

	// For each row
	for i := 0; i < tableVolume; i++ {
		row := list.New()
		// For each column
		for _, column := range columns {
			// Get column type
			columnType, err := db.GetColumnType(info.table, column)
			if err != nil {
				return err
			}

			switch columnType {
			case "integer":
				value, err := db.GetIntValue(info.table, column, i)
				if err != nil {
					return err
				}
				// If current column is anonymized column - get random l-1 values from table and original value
				if isContain(column, info.columns) {
					values, err := db.GetRandomIntValues(info.table, column, info.columnsAndL[column]-1, value)
					if err != nil {
						return err
					}
					values = append(values, value)
					rand.Shuffle(len(values), func(i, j int) {
						values[i], values[j] = values[j], values[i]
					})

					row.PushBack(values)
				} else {
					// Else just copy value in anonymized table
					row.PushBack(value)
				}
			case "text":
				value, err := db.GetTextValue(info.table, column, i)
				if err != nil {
					return err
				}
				// If current column is anonymized column - get random l-1 values from table and original value
				if isContain(column, info.columns) {
					values, err := db.GetRandomStrValues(info.table, column, info.columnsAndL[column]-1, value)
					if err != nil {
						return err
					}
					values = append(values, value)
					rand.Shuffle(len(values), func(i, j int) {
						values[i], values[j] = values[j], values[i]
					})

					row.PushBack(values)
				} else {
					// Else just copy value in anonymized table
					row.PushBack(value)
				}
			default:
				return errors.New("unknown type of column " + column)
			}
		}
		// Insert anonymized aggregated row in table
		err = db.InsertRow(info.table, columns, row)
		if err != nil {
			return err
		}
	}

	return nil
}

func isContain(i string, arr []string) bool {
	for _, a := range arr {
		if a == i {
			return true
		}
	}

	return false
}
