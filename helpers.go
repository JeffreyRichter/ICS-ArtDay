package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func readCSV(filename string, fieldsPerRecord int) [][]string {
	f, err := os.Open(filename)
	PanicOnErr(err)
	defer f.Close()

	r := csv.NewReader(f)               // https://pkg.go.dev/encoding/csv#Reader
	r.FieldsPerRecord = fieldsPerRecord // Each row can have a different # of fields
	r.TrimLeadingSpace = true
	records, err := r.ReadAll()
	PanicOnErr(err)
	return records
}

func atoi(s string) int {
	if s == "" { // If no value, make it 0
		return 0
	}
	n, err := strconv.Atoi(s)
	PanicOnErr(err)
	return n
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
