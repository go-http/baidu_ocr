package ocr

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

//识别出的表单单元格
type FormCell struct {
	Row    []int
	Column []int
	Word   string
}

//表单识别结果
type FormOcrResult struct {
	FileUrl string `json:"file_url"` //Excel格式返回时提供
	FormNum int    `json:"form_num"` //表单数量
	Forms   []struct {
		Header []FormCell
		Footer []FormCell
		Body   []FormCell
	}
}

//输出OCR识别的表格结果
func ParseFormResult(str string) (FormOcrResult, error) {
	var result FormOcrResult
	err := json.Unmarshal([]byte(str), &result)

	return result, err
}

//格式化输出OCR识别的表单
func (data FormOcrResult) FormatDump(wordLength int) {
	form := data.Forms[0]

	fmt.Println("跨行／列数据如下：")
	for _, cell := range form.Body {
		if len(cell.Row) != 1 || len(cell.Column) != 1 {
			fmt.Println("\t行：", cell.Row, "列：", cell.Column, "内容：", cell.Word)
		}
	}

	fmt.Println("表头:\n\t", form.Header)
	fmt.Println("表尾:\n\t", form.Footer)
	fmt.Println("数据")

	//获取行列清单
	rowMap := map[int]struct{}{}
	columnMap := map[int]struct{}{}
	dataMap := map[int]map[int]string{}
	for _, cell := range form.Body {
		if len(cell.Row) != 1 || len(cell.Column) != 1 {
			continue
		}
		rowMap[cell.Row[0]] = struct{}{}
		columnMap[cell.Column[0]] = struct{}{}

		if dataMap[cell.Row[0]] == nil {
			dataMap[cell.Row[0]] = map[int]string{}
		}
		dataMap[cell.Row[0]][cell.Column[0]] = cell.Word
	}

	//排序行列清单
	rows := []int{}
	for row, _ := range rowMap {
		rows = append(rows, row)
	}
	sort.Ints(rows)
	columns := []int{}
	for column, _ := range columnMap {
		columns = append(columns, column)
	}
	sort.Ints(columns)

	for _, row := range rows {
		fmt.Printf("\t第%2d行", row)
		if dataMap[row] == nil {
			continue
		}

		for _, column := range columns {
			word := dataMap[row][column]
			if word == "" {
				fmt.Printf(padString("-", wordLength))
			} else {
				fmt.Printf(padString(word, wordLength))
			}
		}
		fmt.Printf("\n")
	}
}

//UTF8字符（按两个ASCII计算宽度）按ASCII宽度对齐
func padString(str string, length int) string {
	var padLeft bool

	if length < 0 {
		padLeft = true
		length *= -1
	}

	runeCount := utf8.RuneCountInString(str)
	if runeCount == 0 || runeCount >= length {
		return str
	}

	length -= runeCount
	length -= (len(str) - runeCount) / 2

	//左
	if length <= 0 {
		return str
	} else if padLeft {
		return str + strings.Repeat(" ", length)
	} else {
		return strings.Repeat(" ", length) + str
	}
}
