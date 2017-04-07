package xlsxconvert

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

type ServerWriter struct {
	BaseWriter
}

func NewServerWriter(inFilePath string, outFilePath string, outputFmt string) *ServerWriter {
	writer := ServerWriter{}
	writer.InFilePath = inFilePath
	writer.OutputFmt = outputFmt
	writer.OutFilePath = outFilePath
	basename := filepath.Base(outFilePath)
	if strings.Contains(basename, "language") {
		writer.StripString = false
	} else {
		writer.StripString = true
	}
	writer.BaseInit()
	return &writer
}

func (writer *ServerWriter) onNeedWriter(col int) bool {
	if writer.TableHeader.Belongs[col] == "A" || writer.TableHeader.Belongs[col] == "S" || writer.TableHeader.Belongs[col] == "K" {
		return true
	}
	return false
}

func (writer *ServerWriter) Write(sheets []*xlsx.Sheet) error {
	if len(sheets) <= 0 {
		errMsg := fmt.Sprintf("表数据为空 表明:%s", writer.OutFilePath)
		return errors.New(errMsg)
	}

	err := writer.TableHeader.Parse(sheets[0])
	if err != nil {
		return err
	}
	if writer.TableHeader.ServerBelongNum <= 0 {
		fmt.Println(filepath.Base(writer.InFilePath) + " Ingore")
		return nil
	}
	err = writer.write(sheets, writer.onNeedWriter)
	if err != nil {
		return err
	}
	return nil
}
