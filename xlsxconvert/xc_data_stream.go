package xlsxconvert

import "os"
import "fmt"
import "encoding/binary"
import "github.com/sachinkung/xlsx2code/base"

type XcDataStreamWriter struct {
	fd              *os.File
	schemeLen       int
	fieldTypeRowOff int
	fieldNameowOff  int
	row             int
	col             int
	curLen          int64
}

const (
	FIXEDHEADERLEN = 28
	FILETAGLEN     = 4
	INT0TYPE       = 0
	INT1TYPE       = 1
	INT2TYPE       = 2
	INT3TYPE       = 3
	INT4TYPE       = 4
	FLOATTYPE      = 5
	STRINGTYPE     = 6
)

func (w *XcDataStreamWriter) BeginWrite(filePath string) error {
	var err error
	if base.FileIsExist(filePath) == false {
		w.fd, err = os.Create(filePath)
	} else {
		w.fd, err = os.OpenFile(filePath, os.O_RDWR, 0666)
		w.fd.Truncate(0)
	}
	if err != nil {
		return err
	}
	return nil
}

func (w *XcDataStreamWriter) writeHead() {
	binary.Write(w.fd, binary.BigEndian, "LRGB")
	binary.Write(w.fd, binary.BigEndian, int32(0))
	binary.Write(w.fd, binary.BigEndian, int32(0))
	binary.Write(w.fd, binary.BigEndian, int32(0))
	binary.Write(w.fd, binary.BigEndian, int32(0))
	w.curLen, _ = w.fd.Seek(0, os.SEEK_CUR)
}

func (w *XcDataStreamWriter) BeginRow() {

}

func (w *XcDataStreamWriter) BeginCell() {

}

func (w *XcDataStreamWriter) writeInt1(iVal int) {

}

func (w *XcDataStreamWriter) writeInt2(iVal int) {

}

func (w *XcDataStreamWriter) writeInt3(iVal int) {

}

func (w *XcDataStreamWriter) writeInt4(iVal int) {

}

func (w *XcDataStreamWriter) writeInt(iVal int) {

}

func (w *XcDataStreamWriter) WriteInt(iVal int) {
	if iVal == 0 {
		binary.Write(w.fd, binary.BigEndian, int8(INT0TYPE))
	} else if iVal >= -128 && iVal <= 127 {
		binary.Write(w.fd, binary.BigEndian, int8(INT1TYPE))
	} else if iVal >= -32768 && iVal <= 32767 {
		binary.Write(w.fd, binary.BigEndian, int8(INT2TYPE))
	} else if iVal >= -8388608 && iVal <= 8388607 {
		binary.Write(w.fd, binary.BigEndian, int8(INT3TYPE))
	} else if iVal >= -2147483648 && iVal <= 2147483647 {
		binary.Write(w.fd, binary.BigEndian, int8(INT4TYPE))
	}

}

func (w *XcDataStreamWriter) WriteBool(bVal bool) {

	if bVal {
		w.fd.WriteString("1")
	} else {
		w.fd.WriteString("0")
	}
}
func (w *XcDataStreamWriter) WriteFloat(fVal float32) {
	binary.Write(w.fd, binary.BigEndian, float32(fVal))
}

func (w *XcDataStreamWriter) WriteString(sVal string) {
	w.fd.WriteString(sVal)
}

func (w *XcDataStreamWriter) WriteVector2(vec2 base.Vector2) {
	//binary.Write(w.fd, binary.BigEndian, int32(fVal))
}

func (w *XcDataStreamWriter) WriteVector3(vec3 base.Vector3) {
	w.fd.WriteString(fmt.Sprintf("%d,%d,%d", vec3.X, vec3.Y, vec3.Z))
}

func (w *XcDataStreamWriter) EndCell() {
	w.fd.WriteString("\r\n")
}

func (w *XcDataStreamWriter) EndRow() {
	w.fd.WriteString("\r\n")
}

func (w *XcDataStreamWriter) EndWrite() {
	if w.fd != nil {
		w.fd.Close()
	}
}
