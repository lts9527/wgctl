package util

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
	"text/template"
)

// WriteFile 保存文件
func WriteFile(path string, content string) error {
	var fileByte = []byte(content)
	err := ioutil.WriteFile(path, fileByte, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetPublicIp() (string, error) {
	output, err := exec.Command("/bin/sh", "-c", "curl -4 ip.sb").Output()
	if err != nil {
		return "", err
	}
	return RemoveLineBreaks(string(output)), nil
}

func RemoveLineBreaks(str string) string {
	return strings.Replace(str, "\n", "", -1)
}

// FromTemplateContent 替换字符串里的关键字
func FromTemplateContent(templateContent string, envMap map[string]interface{}) string {
	tmpl, err := template.New("text").Parse(templateContent)
	defer func() {
		if r := recover(); r != nil {
			//logger.Error("Template parse failed:", err)
		}
	}()
	if err != nil {
		panic(1)
	}
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return string(buffer.Bytes())
}
