package base

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func FileIsExist(filePath string) bool {
	var exist = true
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func String2Byte(sVal string) (int, error) {
	if len(sVal) <= 0 {
		return 0, nil
	}
	iv, err := String2Int(sVal)
	if err != nil {
		return 0, err
	}
	if iv < 0 && iv > 255 {
		return 0, errors.New("Byte 必须为[0-255]范围内")
	}
	return iv, nil
}

func String2Int(sVal string) (int, error) {
	if len(sVal) <= 0 {
		return 0, nil
	}
	f, err := strconv.ParseFloat(sVal, 64)
	if err != nil {
		return -1, err
	}
	return int(f), nil
}

func String2Bool(sVal string) (bool, error) {
	lowerStr := strings.ToLower(sVal)
	if lowerStr == "false" || lowerStr == "0" || len(lowerStr) == 0 {
		return false, nil
	}
	return true, nil
}

func String2Float(sVal string) (float32, error) {
	if len(sVal) <= 0 {
		return 0, nil
	}
	f, err := strconv.ParseFloat(sVal, 64)
	if err != nil {
		return float32(math.NaN()), err
	}
	return float32(f), nil
}

func String2Vector2(sVal string) (Vector2, error) {

	vec2 := Vector2{}
	vec2.X = 0
	vec2.Y = 0
	if len(sVal) > 0 {
		strList := strings.Split(sVal, ",")
		if len(strList) >= 1 {
			x, err := String2Float(strList[0])
			if err != nil {
				return vec2, errors.New("无法将string转成Vector2")
			}
			vec2.X = x
		}

		if len(strList) >= 2 {
			y, err := String2Float(strList[1])
			if err != nil {
				return vec2, errors.New("无法将string转成Vector2")
			}
			vec2.Y = y
		}
	}

	return vec2, nil
}

func String2Vector3(sVal string) (Vector3, error) {
	vec3 := Vector3{}
	vec3.X = 0
	vec3.Y = 0
	vec3.Z = 0
	if len(sVal) > 0 {
		strList := strings.Split(sVal, ",")
		if len(strList) >= 1 {
			x, err := String2Float(strList[0])
			if err != nil {
				return vec3, errors.New("无法将string转成Vector3")
			}
			vec3.X = x
		}

		if len(strList) >= 2 {
			y, err := String2Float(strList[1])
			if err != nil {
				return vec3, errors.New("无法将string转成Vector3")
			}
			vec3.Y = y
		}
		if len(strList) >= 3 {
			z, err := String2Float(strList[2])
			if err != nil {
				return vec3, errors.New("无法将string转成Vector3")
			}
			vec3.Z = z
		}
	}
	return vec3, nil
}

func FormatXlsxCol(i int) string {
	ret := ""
	i = i + 1
	for i != 0 {
		j := i % 26
		ret = ret + string(j+'A'-1)
		i = i / 26
	}
	return ret
}

func FormatClassName(srcClassName string) string {
	srcClassName = strings.ToUpper(srcClassName[0:1]) + srcClassName[1:]
	ret := ""
	for i := 0; i < len(srcClassName); {
		if srcClassName[i:i+1] == "_" {
			if i+1 < len(srcClassName) {
				ret = ret + strings.ToUpper(srcClassName[i+1:i+2])
			}
			i = i + 2
		} else {
			ret = ret + srcClassName[i:i+1]
			i = i + 1
		}
	}
	return ret
}

func FormatGetPropName(propName string) string {
	atIdx := strings.Index(propName, "@")
	if atIdx != -1 {
		ret := propName[atIdx+1:]
		ret = strings.ToLower(ret[:1]) + ret[1:]
		return ret
	}

	return strings.ToLower(propName[:1]) + propName[1:]
}

func FormatPrivatePropName(propName string) string {
	atIdx := strings.Index(propName, "@")
	if atIdx != -1 {
		ret := propName[atIdx+1:]
		return ret
	}
	return propName
}

func FormatDirName(dirPath string) string {
	dirPath = strings.Replace(dirPath, "//", "/", -1)
	dirPath = strings.Replace(dirPath, "\\", "/", -1)
	dirPath = strings.Replace(dirPath, "\\\\", "/", -1)
	dirs := strings.Split(dirPath, "/")
	ret := ""
	for i := 0; i < len(dirs); i++ {
		if len(dirs[i]) > 0 {
			ret = ret + FormatClassName(dirs[i])
		}

		if i != len(dirs)-1 {
			ret = ret + "/"
		}

	}

	return ret
}

func CheckFileNameIsLegal(fileName string) bool {
	reg := regexp.MustCompile("[a-z]+$|[a-z][a-z_]+[a-z]$")
	strList := reg.FindAllString(fileName, -1)
	return strList != nil && len(strList) > 0
}

func GetFileName(filePath string) string {
	basename := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	return basename[:len(basename)-len(ext)]
}
