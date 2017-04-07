package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sachinkung/xlsx2code/base"
	"github.com/sachinkung/xlsx2code/xlsxconvert"

	"github.com/tealeg/xlsx"
)

func generateCSVFromXLSXFile(excelFileName string, outFilePath string, fmt string, target string) error {
	if base.CheckFileNameIsLegal(base.GetFileName(excelFileName)) == false {
		return errors.New("xlsx命名非法: 只能包含 'a-z, _' 全部小写，字母开头结尾  例如: skill_effect.xlsx")
	}

	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		return err
	}

	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("This XLSX file contains no sheets.")

	}

	//只保留第一个标签页
	if !strings.HasSuffix(excelFileName, "language.xlsx") {
		for i := 1; i < len(xlFile.Sheets); i++ {
			xlFile.Sheets[i] = nil
		}
	}

	if target == "client" {
		clientW := xlsxconvert.NewClientWriter(excelFileName, outFilePath, fmt)
		err = clientW.Write(xlFile.Sheets)
		return err
	} else {
		clientS := xlsxconvert.NewServerWriter(excelFileName, outFilePath, fmt)
		err = clientS.Write(xlFile.Sheets)
		return err
	}
}

func main() {
	var inDir = flag.String("i", "", "(输入目录或文件) input a directory or file")
	var outDir = flag.String("o", "", "(输出目录) Output to directory")
	var target = flag.String("t", "client", "(客户端表或者服务器表) Option: client|server|")
	var format = flag.String("f", "txt", "(选择格式) Select txt format or xc format")
	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}

	start := time.Now()
	f, _ := os.Stat(*inDir)

	if f.IsDir() {
		wg := sync.WaitGroup{}
		workFun := func(inFilePath string, outFilePath string) {
			err := generateCSVFromXLSXFile(inFilePath, outFilePath, *format, *target)
			wg.Done()
			if err != nil {
				errMsg := fmt.Sprintf("%s %s", inFilePath, err.Error())
				log.Panicln(errors.New(errMsg))
			}
			fmt.Println(filepath.Base(inFilePath) + " Success")
		}
		filepath.Walk(*inDir, func(filePath string, info os.FileInfo, err error) error {
			if err != nil { //忽略错误
				return err
			}
			if info.IsDir() { // 忽略目录
				return nil
			}
			fileName := base.GetFileName(filePath)
			if !strings.HasSuffix(filePath, ".xlsx") || strings.HasPrefix(fileName, "~$") {
				return nil
			}
			fileDir := filepath.Dir(filePath)

			outFileDir := *outDir + "/"
			if *target == "server" {
				relDir, _ := filepath.Rel(*inDir, fileDir)
				outFileDir = *outDir + "/" + relDir + "/"
			}
			os.MkdirAll(outFileDir, 0666)
			outFilePath := outFileDir + fileName + "." + *format

			wg.Add(1)
			go workFun(filePath, outFilePath)
			return nil
		})

		wg.Wait()
	} else {

		fileName := base.GetFileName(*inDir)
		outFilePath := *outDir + "/" + fileName + "." + *format
		os.MkdirAll(*outDir, 0666)
		err := generateCSVFromXLSXFile(*inDir, outFilePath, *format, *target)
		if err != nil {
			errMsg := fmt.Sprintf("%s %s", *inDir, err.Error())
			log.Panicln(errors.New(errMsg))
		}
	}

	fmt.Printf("Success cost time: (%s) \n", time.Since(start))

}
