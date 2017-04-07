package main

import (
	"fmt"
	"github.com/sachinkung/xlsx2code/base"

	"github.com/tealeg/xlsx"
)

type TableCheck struct {
}

func (check *TableCheck) CheckIdUnique(sheets []xlsx.Sheet, tableHead *base.TableHead) error {
	//idMap := map[string]int{}
	for i, sheet := range sheets {
		for j, row := range sheet.Rows {
			//忽略空行
			if row == nil {
				continue
			}
			for k, cell := range row.Cells {
				//忽略没有标明字段belong的列
				if k >= tableHead.ColNum {
					break
				}
				fmt.Println(i, j, k, cell.Value)
			}
		}
	}

	return nil
}
