package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"work/model"
)

func ReadFile(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GenerateKeyPair() (string, string) {
	genkey, _ := exec.Command("/bin/bash", "-c", "wg genkey").Output()
	pubkey, _ := exec.Command("/bin/bash", "-c", fmt.Sprintf(`echo "%s" | wg pubkey`, string(genkey))).Output()
	return strings.Replace(string(genkey), "\n", "", -1), strings.Replace(string(pubkey), "\n", "", -1)
}

func GetClientPort() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
	}
	defer listen.Close()
	return strconv.Itoa(listen.Addr().(*net.TCPAddr).Port)
}

func GetIP() string {
	ipv4, _ := exec.Command("/bin/sh", "-c", "curl -4 ip.sb").Output()
	return string(ipv4)
}

// GenerateRandInt 生成随机端口
func GenerateRandInt(min, max int) int {
	rand.Seed(time.Now().UnixMilli())
	return rand.Intn(max-min) + min
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

func GenerateIp(subnet string) string {
	rand.Seed(time.Now().Unix())
	ss := strings.FieldsFunc(subnet, SplitFunc)
	switch ss[4] {
	case "24":
		return fmt.Sprintf("%s.%s.%s.%d", ss[0], ss[1], ss[2], rand.Intn(255))
	case "16":
		return fmt.Sprintf("%s.%s.%d.%d", ss[0], ss[1], rand.Intn(255), rand.Intn(255))
	case "8":
		return fmt.Sprintf("%s.%d.%d.%d", ss[0], rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}
	return ""
}

func SplitFunc(r rune) bool {
	return r == '.' || r == '/'
}

// AppendWriteFile 追加保存文件
func AppendWriteFile(path string, content string) error {
	var fileByte = []byte(content)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(fileByte)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

// FileForEach 遍历文件夹 获取文件
func FileForEach(fileFullPath string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(fileFullPath)
	if err != nil {
		return nil, err
	}
	var myFile []fs.FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		myFile = append(myFile, file)
	}
	return myFile, nil
}

// ReadConfigs 读取每个wg配置的配置
func ReadConfigs(path string) (ws *model.ConfigObjConfig, err error) {
	var serverClientConfig []byte
	serverClientConfig, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(serverClientConfig, &ws)
	if err != nil {
		return nil, err
	}
	return
}

// WriteFile 保存文件
func WriteFile(path string, content string) error {
	var fileByte = []byte(content)
	err := ioutil.WriteFile(path, fileByte, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetBetweenStr 截取字符串
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

// RunCommand 执行linux命令
func RunCommand(arg ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/sh", arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// RandString 生成随机字符串
func RandString(len int) string {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().Unix()))
	bs := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bs[i] = byte(b)
	}
	return string(bs)
}

// RandStringlowercase 生成随机小写字符串
func RandStringlowercase(len int) string {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().Unix()))
	bs := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bs[i] = byte(b)
	}
	return strings.ToLower(string(bs))
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

func SaveJoinServerConfig(filepath string, co *model.ConfigObjConfig) (err error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = AppendWriteFile(filepath, string(file)+BuildAppendWCC(co)+"\n")
	if err != nil {
		return err
	}
	return nil
}

func GenerateClientConfiguration(co *model.ConfigObjConfig) (string, error) {
	var ws = &model.ConfigObjConfig{}
	serverConfig, err := ioutil.ReadFile("/etc/wgctl/server/" + co.JoinServerId)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(serverConfig, &ws)
	if err != nil {
		return "", err
	}
	return BuildWCS(co) + "\n\n" + BuildWCC(ws, co.Endpoint), nil
}

// BuildWCS wireguard服务端配置文件
func BuildWCS(wc *model.ConfigObjConfig) string {
	var sb strings.Builder
	sb.Write([]byte(WCS))
	var envMap = make(map[string]interface{})
	envMap["PrivateKey"] = wc.PrivateKey
	envMap["ListenPort"] = wc.ListenPort
	envMap["Address"] = wc.Address
	envMap["DNS"] = wc.DNS
	envMap["MTU"] = wc.MTU
	return FromTemplateContent(sb.String(), envMap)
}

// BuildWCC wireguard客户端配置文件
func BuildWCC(wc *model.ConfigObjConfig, PublicIp string) string {
	var sb strings.Builder
	sb.Write([]byte(WCC))
	var envMap = make(map[string]interface{})
	envMap["PublicKey"] = wc.PublicKey
	envMap["AllowedIPs"] = wc.AllowedIPs
	envMap["Endpoint"] = PublicIp
	envMap["PersistentKeepalive"] = wc.PersistentKeepalive
	return FromTemplateContent(sb.String(), envMap)
}

// BuildAppendWCC 保存最新的wireguard服务端的客户端配置文件
func BuildAppendWCC(wc *model.ConfigObjConfig) string {
	var sb strings.Builder
	sb.Write([]byte(APPENDSERVERCONFIG))
	var envMap = make(map[string]interface{})
	envMap["PublicKey"] = wc.PublicKey
	envMap["AllowedIPs"] = wc.Address + "/32"
	return FromTemplateContent(sb.String(), envMap)
}

// BuildAppendWCS 保存最新的wireguard服务端配置文件
func BuildAppendWCS(wc *model.ConfigObjConfig) string {
	var sb strings.Builder
	sb.Write([]byte(APPENDSERVERCONFIGS))
	var envMap = make(map[string]interface{})
	envMap["PrivateKey"] = wc.PrivateKey
	envMap["ListenPort"] = wc.ListenPort
	return FromTemplateContent(sb.String(), envMap)
}

// BuildServerConfigTemplate 初始化默认的服务端配置
func BuildServerConfigTemplate(wc *model.ConfigObjConfig) string {
	var sb strings.Builder
	sb.Write([]byte(SERVERCONFIGTEMPLATE))
	var envMap = make(map[string]interface{})
	envMap["Name"] = wc.PrivateKey
	envMap["Port"] = wc.ListenPort
	envMap["PrivateKey"] = wc.PrivateKey
	envMap["PublicKey"] = wc.PublicKey
	envMap["Address"] = wc.Address
	envMap["DNS"] = wc.DNS
	envMap["MTU"] = wc.MTU
	envMap["AllowedIPs"] = wc.AllowedIPs
	envMap["PersistentKeepalive"] = wc.PersistentKeepalive
	return FromTemplateContent(sb.String(), envMap)
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
