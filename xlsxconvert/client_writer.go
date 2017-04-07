package xlsxconvert

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

import "fmt"

type ClientWriter struct {
	BaseWriter
}

func NewClientWriter(inFilePath string, outFilePath string, outputFmt string) *ClientWriter {
	writer := ClientWriter{}
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

func (writer *ClientWriter) onNeedWriter(col int) bool {
	if writer.TableHeader.Belongs[col] == "A" || writer.TableHeader.Belongs[col] == "C" || writer.TableHeader.Belongs[col] == "K" {
		return true
	}
	return false
}

func (writer *ClientWriter) Write(sheets []*xlsx.Sheet) error {
	if len(sheets) <= 0 {
		errMsg := fmt.Sprintf("表数据为空 表明:%s", writer.InFilePath)
		return errors.New(errMsg)
	}
	err := writer.TableHeader.Parse(sheets[0])
	if err != nil {
		return err
	}
	if writer.TableHeader.ClientBelongNum <= 0 {
		fmt.Println(filepath.Base(writer.InFilePath) + " Ingore")
		return nil
	}
	err = writer.write(sheets, writer.onNeedWriter)
	if err != nil {
		return err
	}
	return nil
}
