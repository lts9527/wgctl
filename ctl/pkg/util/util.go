package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
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

func CreateFolder(path string, perm os.FileMode) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		// 必须分成两步
		// 先创建文件夹
		err = os.Mkdir(path, perm)
		if err != nil {
			return err
		}
		// 再修改权限
		err = os.Chmod(path, perm)
		if err != nil {
			return err
		}
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

func Command(cmd string) (string, error) {
	c := exec.Command("/bin/bash", "-c", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return "", err
	}
	var wg sync.WaitGroup
	var res string
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
			res = fmt.Sprintf("%s", readString)
		}
	}()
	err = c.Start()
	wg.Wait()
	return res, err
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
