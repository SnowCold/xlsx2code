package code_gen

import (
	"os"
	"github.com/sachinkung/xlsx2code/base"
	"text/template"
)

type CSharpCodeGen struct {
	fileName    string
	extendDir   string
	generateDir string
	tableHead   *base.TableHead
}

func CreateCodeTemplateClass(tableHead *base.TableHead, config base.Config, fileName string) (*CodeTemplateClass, error) {
	templateClass := CodeTemplateClass{}
	templateClass.ClassName = base.FormatClassName(fileName)
	templateClass.FileBaseName = fileName
	keyType, err := base.TableTypeToString(tableHead.KeyType)
	if err != nil {
		return nil, err
	}
	templateClass.KeyType = keyType
	templateClass.KeyPropName = base.FormatGetPropName(tableHead.KeyName)
	templateClass.Members = make([]CodeTemplateMember, tableHead.ClientBelongNum)
	j := 0
	for i := 0; i < len(tableHead.Fields); i++ {
		if base.CheckBelongClient(tableHead.Belongs[i]) {
			memberClass := CodeTemplateMember{}
			memberClass.Comment = tableHead.Comments[i]
			memberClass.Name = tableHead.Fields[i]
			tableType, err := base.TableTypeToString(tableHead.Types[i])
			if err != nil {
				return nil, err
			}
			memberClass.MemberType = tableType
			memberClass.NeedVirtual = false
			memberClass.FormatPrivateName = base.FormatPrivatePropName(memberClass.Name)
			memberClass.FormatGetPropName = base.FormatGetPropName(memberClass.Name)
			templateClass.Members[j] = memberClass
			j = j + 1
		}

	}
	return &templateClass, nil
}

func GenerateTemplateToFile(tplFilePath string, outFilePath string, templateClass *CodeTemplateClass) error {
	var err error
	var fd *os.File
	if base.FileIsExist(outFilePath) == false {
		fd, err = os.Create(outFilePath)
	} else {
		fd, err = os.OpenFile(outFilePath, os.O_WRONLY, 0666)
		fd.Truncate(0)
	}
	if err != nil {
		return err
	}
	t, err := template.ParseFiles(tplFilePath)
	if err != nil {
		return err
	}
	err = t.Execute(fd, templateClass)
	if err != nil {
		return err
	}
	fd.Close()
	return nil
}

func (baseW *CSharpCodeGen) Run(tableHead *base.TableHead, config base.Config,
	fileName string, extendDir string, generateDir string) error {
	templateClass, err := CreateCodeTemplateClass(tableHead, config, fileName)
	if err != nil {
		return err
	}
	err = GenerateTemplateToFile("csharp_template/DataTableExtend.tmpl",
		extendDir+"TD"+base.FormatClassName(fileName)+"TableExtend.cs", templateClass)
	if err != nil {
		return err
	}
	err = GenerateTemplateToFile("csharp_template/DataExtend.tmpl",
		extendDir+"TD"+base.FormatClassName(fileName)+"Extend.cs", templateClass)
	if err != nil {
		return err
	}
	err = GenerateTemplateToFile("csharp_template/DataTable.tmpl",
		generateDir+"TD"+base.FormatClassName(fileName)+"Table.cs", templateClass)
	if err != nil {
		return err
	}
	err = GenerateTemplateToFile("csharp_template/Data.tmpl",
		generateDir+"TD"+base.FormatClassName(fileName)+".cs", templateClass)
	if err != nil {
		return err
	}
	return nil
}
