package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/larspensjo/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	//"time"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release.bin release.go
// CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o release.app release.go
// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o release.exe release.go

var usage = `Usage: %s [options] 
Options are:
    -f configuration file
`

func infoExit(info string) {
	fmt.Print(info)
	os.Exit(0)
}

func main() {
	//根据CPU数量设置多核运行
	runtime.GOMAXPROCS(runtime.NumCPU())

	//初始工作对象
	w := Work{
		Params: map[string]string{
			"host":  "127.0.0.1",
			"port":  "9999",
			"token": "",
			"dir":   ".",
		},
	}
	w.Init()
	w.Http()
}

type Work struct {
	Params map[string]string
}

func (w *Work) Init() {
	//获取配置文件
	var configFile string
	flag.Usage = func() {
		infoExit(fmt.Sprintf(usage, os.Args[0]))
	}
	flag.StringVar(&configFile, "f", "config.ini", "configuration file")
	flag.Parse()
	if len(configFile) <= 0 {
		infoExit(fmt.Sprintf(usage, os.Args[0]))
	} else if _, err := os.Stat(configFile); err != nil && os.IsNotExist(err) {
		infoExit(fmt.Sprintf("%s configFile not exist", configFile))
	}
	//解析配置文件
	cfg, err := config.ReadDefault(configFile)
	if err != nil {
		infoExit(fmt.Sprintf("%s configFile parse fail: %s", configFile, err.Error()))
	}
	if cfg.HasSection("app") {
		section, err := cfg.SectionOptions("app")
		if err == nil {
			for _, v := range section {
				param, err := cfg.String("app", v)
				if err == nil {
					if v=="dir" {
						w.Params[v] = strings.TrimRight(param, "/")
					}else{
						w.Params[v] = param
					}
				}
			}
		}
	}
}

//http线程
func (w *Work) Http() {
	//获取配置
	log.Println("[http] start ", w.Params["host"], w.Params["port"])
	listen := (w.Params["host"] + ":" + w.Params["port"])

	//路由
	http.HandleFunc("/debug/", w.HttpDebug)
	http.HandleFunc("/release/", w.HttpRelease)

	//开始监听
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatalln("[http] ListenAndServe: ", err)
		return
	}
}

func (w *Work) HttpDebug(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(w.Params["token"]) > 0 {
		if !(len(req.Form["token"]) > 0 && req.Form["token"][0] == w.Params["token"]) {
			resp.Write(retrunJson("[debug] 权限错误", false, nil))
			return
		}
	}

	resp.Write(retrunJson("debug", true, w.Params))
}

func (w *Work) HttpRelease(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(w.Params["token"]) > 0 {
		if !(len(req.Form["token"]) > 0 && req.Form["token"][0] == w.Params["token"]) {
			resp.Write(retrunJson("[release] 权限错误", false, nil))
			return
		}
	}

	var project string = ""
	if len(req.Form["project"]) > 0 && len(req.Form["project"][0]) > 0 {
		project = req.Form["project"][0]
	}
	if len(project) < 1 {
		resp.Write(retrunJson("[release] 参数错误", false, nil))
		return
	}
	items := strings.Split(project, ".")
	var dir string

	//目录形式一
	var dirSlice []string = []string{}
	dirSlice = append(dirSlice, w.Params["dir"])
	//log.Println(items);
	if len(items) > 2 {
		dirSlice = append(dirSlice, strings.Join(items[1:], "."))
	}
	dirSlice = append(dirSlice, items[0])
	dir = strings.Join(dirSlice, "/")
	//log.Println(dir)

	//目录形式二
	s, err := os.Stat(dir + "/.git")
	if err != nil || !s.IsDir() {
		dirFail := dir

		var dirSlice []string = []string{}
		var len int = len(items)
		dirSlice = append(dirSlice, w.Params["dir"])
		dirSlice = append(dirSlice, strings.Join(items[len-2:], "."))
		dirSlice = append(dirSlice, strings.Join(items[:len-2], "."))		
		dir = strings.Join(dirSlice, "/")
		//log.Println(dir)

		s, err = os.Stat(dir + "/.git")
		if err != nil || !s.IsDir() {
			resp.Write(retrunJson("[release] 目录不存在", false, []string{dirFail, dir}))
			return
		}
	}

	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	cmdLine := "cd " + dir + " && git checkout release && git pull origin release && git push release release && git show -2 --name-status"
	cmd := exec.Command("bash", "-c", cmdLine)
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s", opBytes)
	resp.Write(retrunJson(dir, true, strings.Split(fmt.Sprintf("%s", opBytes), "\n")))
}

//http线程返回结果结构函数
func retrunJson(msg string, status bool, data interface{}) []byte {
	if data == nil {
		data = struct{}{}
	}
	b, err := json.Marshal(Result{status, msg, data})
	if err != nil {
		log.Println("[retrunJson] Marshal", err)
	}
	return b
}

//http线程返回结果结构
type Result struct {
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}
