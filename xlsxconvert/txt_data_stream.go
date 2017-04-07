package xlsxconvert

import "os"
import "fmt"
import "github.com/sachinkung/xlsx2code/base"

const (
	kTxtDelimiter = "\t"
)

type TxtDataStreamWriter struct {
	fd        *os.File
	isNewLine bool
}

func (w *TxtDataStreamWriter) BeginWrite(filePath string) error {
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
	w.isNewLine = false
	//write BOM
	w.fd.WriteString(string('\uFEFF'))
	return nil
}

func (w *TxtDataStreamWriter) BeginRow() {
	w.isNewLine = true
}

func (w *TxtDataStreamWriter) BeginCell() {
	if w.isNewLine == true {
		w.isNewLine = false
	} else {
		w.fd.WriteString(kTxtDelimiter)
	}
}

func (w *TxtDataStreamWriter) WriteInt(iVal int) {
	w.fd.WriteString(fmt.Sprintf("%d", iVal))
}

func (w *TxtDataStreamWriter) WriteBool(bVal bool) {

	if bVal {
		w.fd.WriteString("1")
	} else {
		w.fd.WriteString("0")
	}

}
func (w *TxtDataStreamWriter) WriteFloat(fVal float32) {
	w.fd.WriteString(fmt.Sprintf("%f", fVal))

}

func (w *TxtDataStreamWriter) WriteString(sVal string) {
	w.fd.WriteString(sVal)

}

func (w *TxtDataStreamWriter) WriteVector2(vec2 base.Vector2) {
	w.fd.WriteString(fmt.Sprintf("%d,%d", vec2.X, vec2.Y))

}

func (w *TxtDataStreamWriter) WriteVector3(vec3 base.Vector3) {
	w.fd.WriteString(fmt.Sprintf("%d,%d,%d", vec3.X, vec3.Y, vec3.Z))

}

func (w *TxtDataStreamWriter) EndCell() {

}
func (w *TxtDataStreamWriter) EndRow() {
	w.fd.WriteString("\r\n")
}

func (w *TxtDataStreamWriter) EndWrite() {
	if w.fd != nil {
		w.fd.Close()
	}
}
