package utils

import (
	"regexp"
	"strconv"

	exc "github.com/360EntSecGroup-Skylar/excelize"
)

// 参考： ttps://xuri.me/excelize/zh-hans/workbook.html

var (
	rgxAll = regexp.MustCompile(`\{\{\s*[\w.]+\s*\}\}`)
	// rgx         = regexp.MustCompile(`\{\{\s*(\w+)\.\w+\s*\}\}`)
	// rangeRgx    = regexp.MustCompile(`\{\{\s*range\s+(\w+)\s*\}\}`)
	// rangeEndRgx = regexp.MustCompile(`\{\{\s*end\s*\}\}`)
)

type XlsxTemplate struct {
	File *exc.File
}

func (m *XlsxTemplate) PrintOut(templatePath string, targetPath string, tmp map[string]interface{}) error {

	// TODO: 如果超出了要复制sheet:
	// 用GetSheetName获得sheet
	// PrintOut 改成打印具体sheet号。 所有遇到数组的地方，都把已经打印的element删掉（剩下的留到下一页打印）。如果有下一页，就继续打印下一页。
	// PrintOnePage 返回一个bool，如果是false停止打印，否则就继续打印
	// 如果要这么做，一开始必须备份一个工作表（因为{{}}信息在填写过程就会被破坏掉）

	f, err := exc.OpenFile(templatePath)

	// 获得sheet的名字
	sheetName := f.GetSheetName(1)

	println("sheetname", sheetName)

	if err != nil {
		return err
	}

	// 循环模板
	for columnName := range tmp {

		// []map[string]interface{}
		//
		cValue := tmp[columnName]

		// ----------------------------------------------------判断记录里面是什么类型：图片，数组，文字
		switch field := cValue.(type) {

		// -------------- 如果是map，说明是引用的外部struct
		case map[string]interface{}:

			for subColumnName := range field {

				kw := "{{" + columnName + "." + subColumnName + "}}"
				resultCells := f.SearchSheet(sheetName, kw)

				if len(resultCells) > 0 {

					// 绝大部分情况是只打印一个。不唯一的情况比如：总金额要显示在不同地方
					for i := 0; i < len(resultCells); i++ {
						println(sheetName, resultCells[i], field[subColumnName])
						f.SetCellValue(sheetName, resultCells[i], field[subColumnName])
					}
				}
			}

		// -------------- 如果是map数组，说明是个列表
		case []map[string]interface{}:

			// 如果是数组, 先判断长度，如果长度大于0才执行
			if len(field) == 0 {
				continue
			}

			// 取出第一条记录，利用它的key来生成搜索条件，比如{{subitem_list.buyerCode}}
			templateRecord := field[0]

			// 每个key搜索一次
			for subColumnName := range templateRecord {

				kw := "{{" + columnName + "." + subColumnName + "}}"
				resultCells := f.SearchSheet(sheetName, kw)

				// 搜索不到就跳过
				if len(resultCells) == 0 {
					continue
				}

				// TODO: 如果空格数量不够，就向下复制，然后再搜索一次
				if len(field) > len(resultCells) {

					_, rowNumStr := ParseFlight(resultCells[0])
					rowNum, _ := strconv.Atoi(rowNumStr)
					f.DuplicateRow("Sheet1", rowNum)
					resultCells = f.SearchSheet(sheetName, kw)
				}

				// 循环体是 field下的数组， 而不是搜索结果 resultCells
				for i := 0; i < len(resultCells); i++ {

					// 如果实际记录不够填. 就填入空字符
					if i > len(field)-1 {
						f.SetCellValue(sheetName, resultCells[i], "")
						continue
					}

					// 往空格里填入(excelize会自动帮我转string)
					subFieldValue := field[i][subColumnName]
					f.SetCellValue(sheetName, resultCells[i], subFieldValue)
				}
			}

		// -------------- 默认是文字
		default:
			// 如果是string或者int
			resultCells := f.SearchSheet(sheetName, "{{"+columnName+"}}")

			if len(resultCells) > 0 {

				// 结果不唯一
				for i := 0; i < len(resultCells); i++ {
					println(sheetName, resultCells[i], field)
					f.SetCellValue(sheetName, resultCells[i], field)
				}
			}
		}

	}
	// 最后清掉search不到的{{}}
	resultCells := f.SearchSheet(sheetName, rgxAll.String(), true)
	for i := 0; i < len(resultCells); i++ {

		f.SetCellValue(sheetName, resultCells[i], "")
	}
	// result := f.SearchSheet("Sheet1", "{{code}}")

	// fmt.Println("search:")
	// fmt.Println(result)

	// for _, row := range rows {
	// 	for _, colCell := range row {

	// 		fmt.Println(strings.TrimRight(strings.TrimLeft(colCell, "{{"), "}}"))
	// 		fmt.Println(colCell)
	// 	}
	// 	// println()
	// }

	// f.SetCellValue("Sheet2", "A2", "Hello world.2")
	// f.SetCellValue("Sheet1", "B2", 100)

	err = f.SaveAs(targetPath)
	if err != nil {
		return err
	}
	return nil
}

// func getListProp(row []string) string {
// 	for _, cell := range row {
// 		if cell == "" {
// 			continue
// 		}
// 		if match := rgx.FindAllStringSubmatch(cell, -1); match != nil {
// 			return match[0][1]
// 		}
// 	}
// 	return ""
// }

// func getRangeProp(row []string) string {
// 	if len(row) != 0 {
// 		match := rangeRgx.FindAllStringSubmatch(row[0], -1)
// 		if match != nil {
// 			return match[0][1]
// 		}
// 	}

// 	return ""
// }

// func getRangeEndIndex(rows [][]string) int {
// 	var nesting int
// 	for idx := 0; idx < len(rows); idx++ {
// 		if len(rows[idx]) == 0 {
// 			continue
// 		}

// 		if rangeEndRgx.MatchString(rows[idx][0]) {
// 			if nesting == 0 {
// 				return idx
// 			}

// 			nesting--
// 			continue
// 		}

// 		if rangeRgx.MatchString(rows[idx][0]) {
// 			nesting++
// 		}
// 	}

// 	return -1
// }
