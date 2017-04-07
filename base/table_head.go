package base

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
)

//  exported comment
const (
	IntTableType     = 0
	BoolTableType    = 1
	FloatTableType   = 2
	StringTableType  = 3
	ByteTableType    = 4
	Vector2TableType = 5
	Vector3TableType = 6
)

type Vector2 struct {
	X float32
	Y float32
}

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

//TableHead .....
type TableHead struct {
	ColNum          int
	KeyCol          int
	KeyType         int
	KeyName         string
	ClientBelongNum int
	ServerBelongNum int
	Comments        []string
	Belongs         []string
	Types           []int
	Fields          []string
}

func NewTableHead() *TableHead {
	taleHead := TableHead{}
	return &taleHead
}

func CheckBelongClient(belong string) bool {
	if belong == "A" || belong == "K" || belong == "C" {
		return true
	}
	return false
}

func CheckBelongServer(belong string) bool {
	if belong == "A" || belong == "K" || belong == "S" {
		return true
	}
	return false
}

func checkBelongsIsLegal(str string) (string, error) {
	ret := strings.Trim(str, " ")
	ret = strings.ToUpper(ret)
	switch ret {
	case "A", "S", "C", "N", "K":
		return ret, nil
	default:
		return ret, errors.New("必须是 'A','S','C','N','K'")
	}
}

func TableTypeToString(tableType int) (string, error) {
	switch tableType {
	case IntTableType:
		return "int", nil
	case BoolTableType:
		return "bool", nil
	case FloatTableType:
		return "float", nil
	case StringTableType:
		return "string", nil
	case Vector2TableType:
		return "Vector2", nil
	case Vector3TableType:
		return "Vector3", nil
	case ByteTableType:
		return "byte", nil
	}
	return "nil", errors.New("内部错误:Not support TableType")
}

func getTableTypeByString(tableType string) (int, error) {
	switch tableType {
	case "int":
		return IntTableType, nil
	case "bool":
		return BoolTableType, nil
	case "float":
		return FloatTableType, nil
	case "string":
		return StringTableType, nil
	case "vector2":
		return Vector2TableType, nil
	case "vector3":
		return Vector3TableType, nil
	case "byte":
		return ByteTableType, nil
	}
	return -1, errors.New("类型必须是 'int','bool','float','string','vector2', 'vector3'")
}

func formatTypes(str string) (string, int, error) {
	ret := strings.Trim(str, " ")
	ret = strings.ToLower(ret)
	tableType, err := getTableTypeByString(ret)
	return ret, tableType, err

}

func formatField(str string) (string, error) {
	ret := strings.Trim(str, " ")
	if len(ret) <= 0 {
		return ret, errors.New("字段名不能为空")
	}
	if strings.HasPrefix(ret, "table@") {
		ret = ret[len("table@"):]
	}
	ret = strings.ToUpper(ret[:1]) + ret[1:]
	return ret, nil
}

func (tableHead *TableHead) CheckFieldIsUnique() error {
	return nil
	fieldMap := make(map[string]int, tableHead.ColNum)
	for i, field := range tableHead.Fields {
		if tableHead.Belongs[i] == "N" {
			continue
		}
		field = strings.ToLower(field)
		if _, ok := fieldMap[field]; ok {
			errMsg := fmt.Sprintf("字段名 %s 在第%d列 和第%d列重复(请注意大小写,不区分大小)", field, fieldMap[field], i)
			return errors.New(errMsg)
		} else {
			fieldMap[field] = i
		}
	}
	return nil
}

func (tableHead *TableHead) parseKey() {
	tableHead.KeyCol = -1
	for i, belong := range tableHead.Belongs {
		if belong == "K" {
			tableHead.KeyCol = i
			break
		}
	}
	if tableHead.KeyCol == -1 {
		for i, belong := range tableHead.Fields {
			if belong == "Id" {
				tableHead.KeyCol = i
				break
			}
		}
	}

	if tableHead.KeyCol == -1 {
		tableHead.KeyCol = 0
	}
	if tableHead.KeyCol != -1 {
		tableHead.KeyType = tableHead.Types[tableHead.KeyCol]
		tableHead.KeyName = tableHead.Fields[tableHead.KeyCol]
	}
}

func (taleHead *TableHead) parseBelongCount() (int, int) {
	clientRet := 0
	serverRet := 0
	for i := 0; i < len(taleHead.Belongs); i++ {
		if i >= taleHead.ColNum {
			break
		}
		if CheckBelongClient(taleHead.Belongs[i]) {
			clientRet = clientRet + 1
		}
		if CheckBelongServer(taleHead.Belongs[i]) {
			serverRet = serverRet + 1
		}
	}
	return clientRet, serverRet
}

//// 解析表头....
func (taleHead *TableHead) Parse(sheet *xlsx.Sheet) error {
	if len(sheet.Rows) < IgnoreLine {
		return errors.New("表头错误: 表头必须是4行的")
	}
	taleHead.ColNum = 0
	taleHead.Belongs = make([]string, len(sheet.Rows[BelongLine].Cells))
	//先解析belong
	for i, cell := range sheet.Rows[BelongLine].Cells {

		str := strings.Trim(cell.Value, " ")
		if len(str) <= 0 {
			break
		}
		str, err := checkBelongsIsLegal(str)
		if err != nil {
			errMsg := fmt.Sprintf("表头错误: 错误位置(第%d行 第%d列), 错误描述:%s", BelongLine, i, err.Error())
			return errors.New(errMsg)
		}
		taleHead.Belongs[i] = str
		taleHead.ColNum = taleHead.ColNum + 1
	}
	taleHead.ClientBelongNum, taleHead.ServerBelongNum = taleHead.parseBelongCount()
	taleHead.Types = make([]int, taleHead.ColNum)
	taleHead.Fields = make([]string, taleHead.ColNum)
	taleHead.Comments = make([]string, taleHead.ColNum)
	//解析注释
	for i, cell := range sheet.Rows[CommentLine].Cells {
		if i >= taleHead.ColNum {
			break
		}
		if taleHead.Belongs[i] == "N" {
			continue
		}
		taleHead.Comments[i] = strings.Trim(cell.Value, " ")
	}

	//检测类型是否错误
	for i, cell := range sheet.Rows[TypeLine].Cells {
		if i >= taleHead.ColNum {
			break
		}
		if taleHead.Belongs[i] == "N" {
			continue
		}
		cVal, tableType, err := formatTypes(cell.Value)
		cell.Value = cVal
		if err != nil {
			errMsg := fmt.Sprintf("表头错误: 错误位置(第%d行 第%d列), 错误描述:%s", TypeLine, i, err.Error())
			return errors.New(errMsg)
		}
		taleHead.Types[i] = tableType
	}

	//检测字段是否合法
	for i, cell := range sheet.Rows[FiledNameLine].Cells {
		if i >= taleHead.ColNum {
			break
		}
		if taleHead.Belongs[i] == "N" {
			continue
		}

		str, err := formatField(cell.Value)
		if err != nil {
			errMsg := fmt.Sprintf("表头错误: 错误位置(第%d行 第%d列), 错误描述:%s", FiledNameLine, i, err.Error())
			return errors.New(errMsg)
		}
		taleHead.Fields[i] = str
	}

	taleHead.parseKey()

	err := taleHead.CheckFieldIsUnique()
	if err != nil {
		errMsg := fmt.Sprintf("表头错误: 错误描述:%s", err.Error())
		return errors.New(errMsg)
	}
	return err
}
