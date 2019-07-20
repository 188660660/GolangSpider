package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)



type Parameter struct {
	Name, Url, Local string
	Start, End       int
}

func main() {
	// 获取用户输入参数
	param := InitInput()
	// 设置参数默认值
	param = InitValue(param)
	// 开始工作
	StartWork(param)
}

// 工作函数 分发处理
func StartWork(p Parameter) {
	var ch = make(chan struct{})
	for i := p.Start; i <= p.End; i++ {
		url := p.Url + strconv.Itoa((i-1)*50)
		local := p.Local + strconv.Itoa(i) + ".txt"
		dir := p.Local
		go SpiderRun(local,dir,url, i, ch)
	}
	for i:=p.Start;i<=p.End ;i++  {
		fmt.Printf("Waiting for page %d \n", i)
		<-ch
		fmt.Printf("Page %d was successfully written \n", i)
	}
}

// 运行蜘蛛
func SpiderRun(local,dir,url string, i int, ch chan struct{}) {
	content, err := HttpGet(url)
	if err != nil {
		fmt.Printf("Read the %d page error info: %v \n", i, err)
		ch <- struct{}{}
		return
	}
	err = SaveFile(content,local,dir)
	if err != nil {
		fmt.Printf("insert the %d page error info: %v \n", i, err)
		ch <- struct{}{}
		return
	}
	ch <- struct{}{}
}

// 初始化参数值
func InitInput() (param Parameter) {
	fmt.Printf("Please enter local save address")
	_, _ = fmt.Scanln(&param.Local)
	fmt.Printf("Please enter Baidu Post Bar Address")
	_, _ = fmt.Scanln(&param.Url)
	fmt.Printf("Please enter the start page")
	_, _ = fmt.Scanln(&param.Start)
	fmt.Printf("Please enter the end page")
	_, _ = fmt.Scanln(&param.End)
	return
}

func InitValue(param Parameter) Parameter {
	// 设置默认值 做测试使用
	if strconv.Itoa(param.Start) == "0" {
		param.Start = 1
	}
	if strconv.Itoa(param.End) == "0" {
		param.End = 15
	}
	if param.Url == "" {
		param.Url = "https://tieba.baidu.com/f?kw=golang&ie=utf-8&pn="
	}
	if param.Local == "" {
		param.Local = "E:/SpiderResult/"
	}
	return param
}

// 通过http.get获取目标网址内容
func HttpGet(path string) (content []byte, err error) {
	result, err := http.Get(path)
	if err != nil {
		return
	}
	defer result.Body.Close()

	content, err = ioutil.ReadAll(result.Body)
	if err != nil {
		return
	}
	return
}

// 将读取内容保存至本地
func SaveFile(content []byte, path,dir string) (err error) {
	exist, _ := PathExists(dir)
	if exist == false {
		err = os.Mkdir(dir, os.ModePerm)
	}
	err = ioutil.WriteFile(path, content, 0666)
	return
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}