package controllers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ExcelExport(fileName string, headers []string, values [][]string, c *gin.Context) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	// utf-8
	buf.WriteString("\xEF\xBB\xBF")

	// headers
	w.Write(headers)

	// writer content
	for _, value := range values {
		w.Write(value)
		w.Flush()
	}

	fileFullName := fmt.Sprintf("%s_%s.xls", strconv.Itoa(int(time.Now().Unix())), fileName)

	// header
	c.Header("Content-Type", "application/vnd.ms-excel")
	c.Header("Content-Disposition", "attachment; filename="+fileFullName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	response, err := ioutil.ReadAll(buf)

	if err != nil {
		return
	}
	c.Data(http.StatusOK, c.GetHeader("Content-Type"), response)

}
