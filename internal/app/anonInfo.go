package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"sensQID/internal/pkg/database"
)

type anonInfo struct {
	table       string
	columns     []string
	l           []int
	columnsAndL map[string]int
}

func NewAnonInfo() *anonInfo {
	return &anonInfo{}
}

func (info *anonInfo) getInfo(db *database.DB) error {
	reader := bufio.NewReader(os.Stdin)

	table, err := info.getTableName(reader, db)
	if err != nil {
		return err
	}
	info.table = table

	columns, err := info.getSensQID(reader, db)
	if err != nil {
		return err
	}
	info.columns = columns

	l, err := info.getL(reader)
	if err != nil {
		return err
	}
	info.l = l

	columnsAndL := make(map[string]int)
	for i, column := range columns {
		columnsAndL[column] = l[i]
	}

	info.columnsAndL = columnsAndL

	return nil
}

func (info *anonInfo) getTableName(reader *bufio.Reader, db *database.DB) (string, error) {
	fmt.Println("Enter table name:")
	table, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	table = strings.TrimSuffix(table, "\r\n")
	exist, err := db.IsTableExist(table)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("table " + table + " isn't exist")
	}

	return table, nil
}

func (info *anonInfo) getSensQID(reader *bufio.Reader, db *database.DB) ([]string, error) {
	fmt.Println("Enter sensitive QID column name:")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	input = strings.ReplaceAll(input, " ", "")
	input = strings.TrimSuffix(input, "\r\n")
	QIDs := strings.Split(input, ",")

	for _, QID := range QIDs {
		exist, err := db.IsColumnExist(info.table, QID)
		if err != nil {
			return nil, err
		}

		if !exist {
			return nil, errors.New("column " + QID + " isn't exist")
		}
	}

	return QIDs, nil
}

func (info *anonInfo) getL(reader *bufio.Reader) ([]int, error) {
	fmt.Println("Enter (l_1, ..., l_q):")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	input = strings.ReplaceAll(input, " ", "")
	input = strings.TrimSuffix(input, "\r\n")

	lStr := strings.Split(input, ",")
	if len(lStr) != len(info.columns) {
		return nil, errors.New("count of l's is not equal to number of columns")
	}

	var l []int
	for _, num := range lStr {
		n, err := strconv.Atoi(num)
		if err != nil {
			return nil, err
		}
		l = append(l, n)
	}

	return l, nil
}
