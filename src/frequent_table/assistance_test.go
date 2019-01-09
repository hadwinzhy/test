package frequent_table

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestWeekDate(tests *testing.T) {
	tt := []struct {
		data string
	}{
		{
			data: strconv.Itoa(int(time.Now().Unix())),
		},
		{
			data: "1547056800",
		},
	}
	for _, t := range tt {
		fmt.Println(weekDate(t.data))

	}
}
