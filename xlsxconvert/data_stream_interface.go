package xlsxconvert

import "github.com/sachinkung/xlsx2code/base"

type IDataStreamWriter interface {
	BeginWrite(filePath string) error
	BeginRow()
	BeginCell()
	WriteInt(iVal int)
	WriteBool(bVal bool)
	WriteFloat(fVal float32)
	WriteString(sVal string)
	WriteVector2(vec2 base.Vector2)
	WriteVector3(vec3 base.Vector3)
	EndCell()
	EndRow()
	EndWrite()
}
