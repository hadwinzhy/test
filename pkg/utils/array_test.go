package utils

import (
	"fmt"
	"testing"
	"time"
)

type exampleType struct {
	Date  time.Time
	Value int
}

func TestInsertMissing(t *testing.T) {
	fmt.Println("start")
	testSlice := make([]exampleType, 10)
	InsertMissing(&testSlice, "Date", "day", time.Now(), time.Now())
	fmt.Println(len(testSlice), cap(testSlice))
}
