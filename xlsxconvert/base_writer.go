package xlsxconvert

import (
	"errors"
	"strings"

	"fmt"
	"github.com/sachinkung/xlsx2code/base"

	"github.com/tealeg/xlsx"
)

type BaseWriter struct {
	InFilePath    string
	OutputFmt     string
	OutFilePath   string
	IDataWriter   IDataStreamWriter
	TableHeader   *base.TableHead
	StripString   bool
	fmtWriterFunc func(sVal string, cellType int, stripString bool) error
}

func (baseW *BaseWriter) BaseInit() {
	baseW.TableHeader = base.NewTableHead()

	if baseW.OutputFmt == "txt" {
		baseW.fmtWriterFunc = baseW.txtWriteField
		baseW.IDataWriter = &TxtDataStreamWriter{}
	} else {
		baseW.fmtWriterFunc = baseW.xcWriteField
		baseW.IDataWriter = &XcDataStreamWriter{}
	}
}

func (baseW *BaseWriter) checkHasInvalidCharSymbol(str string) error {
	if strings.Contains(str, "\t") || strings.Contains(str, "\r") || strings.Contains(str, "\n") {
		return errors.New("包含非法字符(\\t \\r \\n)")
	}
	return nil
}

func (baseW *BaseWriter) txtWriteField(sVal string, cellType int, stripString bool) error {
	baseW.IDataWriter.BeginCell()
	defer baseW.IDataWriter.EndCell()

	switch cellType {
	case base.BoolTableType, base.Vector2TableType, base.Vector3TableType:
		baseW.IDataWriter.WriteString(strings.Trim(sVal, " "))
		return nil
	case base.ByteTableType:
		iVal, err := base.String2Byte(strings.Trim(sVal, " "))
		baseW.IDataWriter.WriteInt(iVal)
		return err
	case base.IntTableType:
		fVal, err := base.String2Int(strings.Trim(sVal, " "))
		baseW.IDataWriter.WriteInt(fVal)
		return err
	case base.FloatTableType:
		fVal, err := base.String2Float(strings.Trim(sVal, " "))
		baseW.IDataWriter.WriteFloat(fVal)
		return err
	case base.StringTableType:
		if stripString == true {
			sVal = strings.Trim(sVal, " ")
		}
		baseW.IDataWriter.WriteString(sVal)
		return nil
	default:
		return errors.New("类型必须是 'int','bool','float','string','vector2', 'vector3'")
	}
}

func (baseW *BaseWriter) xcWriteField(sVal string, cellType int, stripString bool) error {
	switch cellType {
	case base.IntTableType:
		iVal, err := base.String2Int(strings.Trim(sVal, " "))
		if err != nil {
			return err
		}
		baseW.IDataWriter.WriteInt(iVal)
		return nil
	case base.BoolTableType:
		bVal, err := base.String2Bool(strings.Trim(sVal, " "))
		if err != nil {
			return err
		}
		baseW.IDataWriter.WriteBool(bVal)
		return nil
	case base.FloatTableType:
		fVal, err := base.String2Float(strings.Trim(sVal, " "))
		if err != nil {
			return err
		}
		baseW.IDataWriter.WriteFloat(fVal)
		return nil
	case base.StringTableType:
		if stripString == true {
			sVal = strings.Trim(sVal, " ")
		}
		baseW.IDataWriter.WriteString(sVal)
		return nil
	case base.ByteTableType:
		iVal, err := base.String2Byte(strings.Trim(sVal, " "))
		if err != nil {
			return err
		}
		baseW.IDataWriter.WriteInt(iVal)
		return nil
	case base.Vector2TableType:
		v2Val, err := base.String2Vector2(strings.Trim(sVal, " "))
		if err != nil {
			return err
		}
		baseW.IDataWriter.WriteVector2(v2Val)
		return nil
	case base.Vector3TableType:
		v3Val, err := base.String2Vector3(strings.Trim(sVal, " "))
		if err != nil {
			return err
		}
		baseW.IDataWriter.WriteVector3(v3Val)
		return nil
	default:
		return errors.New("类型必须是 'int','bool','float','string','vector2', 'vector3'")
	}
}

func (baseW *BaseWriter) writeField(cell *xlsx.Cell, cellType int) error {
	err := baseW.checkHasInvalidCharSymbol(cell.Value)
	if err != nil {
		return err
	}
	err = baseW.fmtWriterFunc(cell.Value, cellType, baseW.StripString)
	return err
}

func isEmptyRow(row *xlsx.Row) bool {
	if row == nil {
		return true
	}
	if len(row.Cells) <= 0 {
		return true
	}
	if len(strings.Trim(row.Cells[0].Value, " ")) <= 0 {
		return true
	}
	return false
}

func (baseW *BaseWriter) WriteHead(sheet *xlsx.Sheet, onNeedWrite func(col int) bool) error {

	for j, row := range sheet.Rows {
		if j == base.BelongLine {
			continue
		}
		//忽略表头
		if j >= base.IgnoreLine {
			break
		}
		baseW.IDataWriter.BeginRow()
		for k, cell := range row.Cells {

			if k >= baseW.TableHeader.ColNum {
				break
			}
			if onNeedWrite(k) == false {
				continue
			}
			err := baseW.fmtWriterFunc(cell.Value, base.StringTableType, true)

			if err != nil {
				errMsg := fmt.Sprintf("表头错误：错误位置(第%d行,第%d列) 错误描述: %s", j+1, base.FormatXlsxCol(k), err.Error())
				return errors.New(errMsg)
			}
		}
		baseW.IDataWriter.EndRow()
	}
	return nil
}

func (baseW *BaseWriter) write(sheets []*xlsx.Sheet, onNeedWrite func(col int) bool) error {
	err := baseW.IDataWriter.BeginWrite(baseW.OutFilePath)
	defer baseW.IDataWriter.EndWrite()
	if err != nil {
		return err
	}
	err = baseW.WriteHead(sheets[0], onNeedWrite)
	if err != nil {
		return err
	}
	for i, sheet := range sheets {
		if sheet == nil {
			break
		}
		for j, row := range sheet.Rows {
			//忽略表头
			if j < base.IgnoreLine {
				continue
			}

			if isEmptyRow(row) {
				continue
			}
			baseW.IDataWriter.BeginRow()

			for k, cell := range row.Cells {
				if k >= baseW.TableHeader.ColNum {
					break
				}
				if onNeedWrite(k) == false {

					continue
				}
				err := baseW.writeField(cell, baseW.TableHeader.Types[k])
				if err != nil {
					errMsg := fmt.Sprintf("表数据错误：错误位置(第%d标签,第%d行,第%s列) 错误描述: %s", i+1, j+1, base.FormatXlsxCol(k), err.Error())
					return errors.New(errMsg)
				}
			}
			//有可能需要补齐
			for k := len(row.Cells); k < baseW.TableHeader.ColNum; k++ {
				if onNeedWrite(k) == false {

					continue
				}
				err := baseW.fmtWriterFunc("", baseW.TableHeader.Types[k], false)
				if err != nil {
					errMsg := fmt.Sprintf("表数据错误：错误位置(第%d标签,第%d行,第%s列) 错误描述: %s", i+1, j+1, base.FormatXlsxCol(k), err.Error())
					return errors.New(errMsg)
				}
			}
			baseW.IDataWriter.EndRow()
		}
	}
	return nil
}
