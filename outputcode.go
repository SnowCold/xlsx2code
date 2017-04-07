package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/sachinkung/xlsx2code/base"
	"github.com/sachinkung/xlsx2code/code_gen"
	"strings"
	"sync"
	"time"

	"github.com/tealeg/xlsx"
)

func outputCode(config base.Config, excelFileName string, extendDir string, generateDir string) error {
	fileBaseName := base.GetFileName(excelFileName)
	if base.CheckFileNameIsLegal(fileBaseName) == false {
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

	tableHeader := base.TableHead{}
	err = tableHeader.Parse(xlFile.Sheets[0])
	if err != nil {
		return err
	}

	if tableHeader.ClientBelongNum <= 0 {
		return nil
	}

	os.MkdirAll(extendDir, 0666)
	os.MkdirAll(generateDir, 0666)

	csharpCodeGen := code_gen.CSharpCodeGen{}
	err = csharpCodeGen.Run(&tableHeader, config, fileBaseName, extendDir, generateDir)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var inDir = flag.String("i", "", "(输入目录或文件) input a directory or file")
	var outDir = flag.String("o", "", "(输出目录) Output to directory")
	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}
	config := base.NewConfig()
	config.Load("config.json")
	wg := sync.WaitGroup{}
	workFun := func(inFilePath string, extendDir string, generateDir string) {
		err := outputCode(*config, inFilePath, extendDir, generateDir)
		wg.Done()
		if err != nil {
			errMsg := fmt.Sprintf("%s %s", inFilePath, err.Error())
			log.Panicln(errors.New(errMsg))
		}
		fmt.Println(filepath.Base(inFilePath) + " Success")
	}

	start := time.Now()
	f, _ := os.Stat(*inDir)
	if f.IsDir() {

		filepath.Walk(*inDir, func(filePath string, info os.FileInfo, err error) error {
			if err != nil { //忽略错误
				return err
			}
			if info.IsDir() { // 忽略目录
				return nil
			}
			fileName := base.GetFileName(filePath)
			//忽略不需要输出的文件
			if config.CheckNeedOutputCode(fileName) == false {
				return nil
			}

			if !strings.HasSuffix(filePath, ".xlsx") || strings.HasPrefix(fileName, "~$") {
				return nil
			}
			fileDir := filepath.Dir(filePath)

			relDir, _ := filepath.Rel(*inDir, fileDir)
			relDir = base.FormatDirName(relDir)
			extendDir := *outDir + "/Extend/" + relDir + "/"
			generateDir := *outDir + "/Generate/" + relDir + "/"
			wg.Add(1)
			go workFun(filePath, extendDir, generateDir)
			return nil
		})

	} else {

		fileDir := filepath.Dir(*inDir)
		fileDir = base.FormatDirName(fileDir)
		extendDir := *outDir + "/Extend/" + fileDir + "/"
		generateDir := *outDir + "/Generate/" + fileDir + "/"

		wg.Add(1)
		go workFun(*inDir, extendDir, generateDir)
	}
	wg.Wait()
	fmt.Printf("Success cost time: (%s) \n", time.Since(start))

}
