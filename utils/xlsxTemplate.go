package utils

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	exc "github.com/360EntSecGroup-Skylar/excelize"
)

// 参考： ttps://xuri.me/excelize/zh-hans/workbook.html

var (
	// 清除用
	rgxAll  = regexp.MustCompile(`\{\{\s*[\w.\/]+\s*\}\}`)
	rgxFile = regexp.MustCompile(`\{\{\&\s*[\w.\/]+\s*\}\}`)

	// rgx         = regexp.MustCompile(`\{\{\s*(\w+)\.\w+\s*\}\}`)
	// rangeRgx    = regexp.MustCompile(`\{\{\s*range\s+(\w+)\s*\}\}`)
	// rangeEndRgx = regexp.MustCompile(`\{\{\s*end\s*\}\}`)
)

type XlsxTemplate struct {
	File *exc.File
}

type Empty struct{}

// 把json树平摊开，用aa/bb/cc作为key返回.
// !! 注意，这里没有考虑nested里面的图片。如果nested里面有图片，直接抛弃不用
// 以下代码第一个swtich先把要检查的数据源摊开，再用递归检查里面的节点
func (m *XlsxTemplate) matchMapKey(targetMap map[string]Empty, subKey string) string {

	for wholeKey := range targetMap {
		if strings.Contains(wholeKey, subKey) {
			return wholeKey
		}
	}

	return ""
}

// 这里不处理list[]，只处理nested的对象。
func (m *XlsxTemplate) generateMapping(
	objPath string,
	node interface{},
	mapping map[string]interface{},
	targetCells map[string]Empty) (tableFields map[string]interface{}) {

	switch field := node.(type) {
	case map[string]interface{}:
		for subColumnName := range field {

			subObjPath := ""

			// 如果是根节点，就不需要"/"
			if objPath == "" {
				subObjPath = subColumnName
			} else {
				subObjPath = objPath + "/" + subColumnName
			}

			// 不考虑子列表里的图片
			mapping = m.generateMapping(subObjPath, field[subColumnName], mapping, targetCells)
		}
	case []interface{}:
	case nil:
		// 如果是数组，直接忽略
	default:

		// 如果到底了，检查excel里有没有目标格，没有的话就不管： 检查到 {{&}}的格式就是图片，否则{{}}是普通值
		if wholeKey := m.matchMapKey(targetCells, "{{&"+objPath+"}}"); wholeKey != "" {
			mapping[wholeKey] = field
		} else if wholeKey := m.matchMapKey(targetCells, "{{"+objPath+"}}"); wholeKey != "" {

			mapping[wholeKey] = field
		}
	}

	return mapping
}

func (m *XlsxTemplate) DeleteLastSheet(f *exc.File) {
	count := 1
	sheetName := f.GetSheetName(2)
	for sheetName != "" {
		count++
		sheetName = f.GetSheetName(count)
	}

	if count > 2 {
		sheetReadme := f.GetSheetName(count - 1)

		if sheetReadme == "说明" {
			f.DeleteSheet(sheetReadme)
		}
	}
}

/*
	如果是图片, 替换符的格式 "{{&ds_image}}, 100, 200"。后面两个数字是尺寸上限
*/
func (m *XlsxTemplate) PrintOut(templatePath string, targetPath string, tmp map[string]interface{}) error {

	/*
	 DONETODO(客户不要这个):
	 如果超出了要复制sheet:
	 用GetSheetName获得sheet
	 PrintOut 改成打印具体sheet号。 所有遇到数组的地方，都把已经打印的element删掉（剩下的留到下一页打印）。如果有下一页，就继续打印下一页。
	 PrintOnePage 返回一个bool，如果是false停止打印，否则就继续打印
	 如果要这么做，一开始必须备份一个工作表（因为{{}}信息在填写过程就会被破坏掉）*/

	f, err := exc.OpenFile(templatePath)

	// 获得sheet的名字
	sheetName := f.GetSheetName(1)

	// 删除最后一个sheet(作为说明)
	// DeleteSheet
	m.DeleteLastSheet(f)

	if err != nil {
		return err
	}
	// targetCells := f.SearchSheet(sheetName, rgxAll.String(), true)

	targetCellsSet := map[string]Empty{}
	rows := f.GetRows(sheetName)

	for _, row := range rows {
		for _, colCell := range row {
			// reg. 如果带有{{}} 或者{{&}}，则该单元格成为容器，等会用来接受dataSource的值
			if rgxAll.MatchString(colCell) || rgxFile.MatchString(colCell) {
				targetCellsSet[colCell] = Empty{}
			}
		}
	}

	//  利用map做一个set，传进去
	columnMapping := make(map[string]interface{})

	// 区分普通值，列表，图片
	columnMapping = m.generateMapping("", tmp, columnMapping, targetCellsSet)

	// 先看keyword，有的话才根据目标地址读图片。图片都是png
	// 插入图片是一个特殊动作，是不是可以特殊处理?
	// pictures map[string]byte[] //这样子多一个变量. 识别byte有点傻逼。应该识别key

	// 循环第一次，处理图片
	// for imageName := range imageFields {

	// 	imagePath := imageFields[imageName]

	// 	resultCells := f.SearchSheet(sheetName, imageName)

	// 	if len(resultCells) > 0 {
	// 		for i := 0; i < len(resultCells); i++ {

	// 			f.SetCellValue(sheetName, resultCells[i], imagePath)
	// 		}
	// 	}
	// }

	// 循环第二次，处理普通字段
	for columnName := range columnMapping {

		// 变量用逗号表示，所以这里要切出变量

		columnValue := columnMapping[columnName]

		resultCells := f.SearchSheet(sheetName, columnName)

		if len(resultCells) > 0 {

			// 结果有时候不唯一（比如总价格需要显示在列表里，也要显示在合同内容中）
			for i := 0; i < len(resultCells); i++ {

				// 判断如果是图片
				if rgxFile.MatchString(columnName) {

					// 数据源除了path以外，用逗号表示params：宽和高
					dataParams := strings.Split(columnValue.(string), ",")
					dataImage := dataParams[0]

					// 从 columnName 切出参数
					params := strings.Split(columnName, ",")

					// 控制图片大小用：如果逗号切出来长度是3，说明有两个参数且一一对应。其他情况不予考虑
					options := `{"positioning": "Location"}`

					if len(params) == 3 && len(dataParams) == 3 {

						limitW, err := strconv.Atoi(params[1])
						if err != nil {
							fmt.Println("参数出错", err)
							return err
						}

						limitH, err := strconv.Atoi(params[2])
						if err != nil {
							fmt.Println("参数出错", err)
							return err
						}

						dataW, err := strconv.Atoi(dataParams[1])
						if err != nil {
							fmt.Println("参数出错", err)
							return err
						}

						dataH, err := strconv.Atoi(dataParams[2])
						if err != nil {
							fmt.Println("参数出错", err)
							return err
						}

						// 如果高超过，缩减多少，如果宽超过，缩减多少。两个缩减比例对比取最小
						scaleW := float64(limitW) / float64(dataW)
						scaleH := float64(limitH) / float64(dataH)

						scale := math.Min(math.Min(scaleW, scaleH), 1.0)

						options = fmt.Sprintf(`{
							"x_scale": %f,
							"y_scale": %f,
							"positioning": "Location"
						}`, scale, scale)
					}

					file, err := ioutil.ReadFile("uploads/" + dataImage)
					if err != nil {
						fmt.Println("读图片错误", err)
					}
					if err := f.AddPictureFromBytes(sheetName, resultCells[i], options, resultCells[i], ".png", file); err != nil {
						fmt.Println("插入图片错误", err)
					}

					// fmt.Println("找到这格", columnName)

					f.SetCellValue(sheetName, resultCells[i], "")

					// f.SetCellValue(sheetName, resultCells[i], columnValue)
				} else {
					f.SetCellValue(sheetName, resultCells[i], columnValue)
				}

			}
		}
	}

	// 循环第三次，搜索并处理列表.(一般只有一个)
	for columnName := range tmp {

		list := tmp[columnName]

		// 假如字段值是个list，且长度大于0
		if list, ok := list.([]interface{}); ok && len(list) > 0 {

			// 取出第一条，单独拆包，为了取所有需要搜索的字段
			subColumnNameMapping := make(map[string]interface{})
			subColumnNameMapping = m.generateMapping(columnName, list[0], subColumnNameMapping, targetCellsSet)

			// 每个字段都搜索一次。为了在空间不够的时候提前复制（等填入值以后就搜不到了）
			for subColumnName := range subColumnNameMapping {

				resultCells := f.SearchSheet(sheetName, subColumnName)

				// 搜索不到就跳过
				if len(resultCells) == 0 {
					continue
				}

				// 如果空格数量不够，就把最后一行向下复制，然后再搜索一次。
				if len(list) > len(resultCells) {
					rowLetterStr, rowNumStr := ParseFlight(resultCells[len(resultCells)-1])
					rowNum, _ := strconv.Atoi(rowNumStr)

					numDiff := len(list) - len(resultCells)

					// 多复制几次
					for j := 0; j < numDiff; j++ {
						f.DuplicateRow(sheetName, rowNum)
						resultCells = append(resultCells, rowLetterStr+strconv.Itoa(rowNum+j+1))
					}
				}

				// 最后循环列表项, 插数据。经过上面的处理，必然是：len(resultCells) >= len(list)。所以循环小的那个
				for k := 0; k < len(list); k++ {

					// 如果实际记录不够填. 就填入空字符
					if k > len(list)-1 {
						f.SetCellValue(sheetName, resultCells[k], "")
						continue
					}

					// 拆包一次取真实值
					subValueMapping := make(map[string]interface{})
					subValueMapping = m.generateMapping(columnName, list[k], subValueMapping, targetCellsSet)

					// 根据目前搜索的子字段（按约定，每一列用的都是相同的字段名）
					value := subValueMapping[subColumnName]

					// 往空格里填入(excelize会自动帮我转string)
					f.SetCellValue(sheetName, resultCells[k], value)
				}
			}
		}
	}

	// 最后清掉search不到的{{}}
	resultCells := f.SearchSheet(sheetName, rgxAll.String(), true)
	for i := 0; i < len(resultCells); i++ {
		f.SetCellValue(sheetName, resultCells[i], "")
	}

	// for _, row := range rows {
	// 	for _, colCell := range row {

	// 		fmt.Println(strings.TrimRight(strings.TrimLeft(colCell, "{{"), "}}"))
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
