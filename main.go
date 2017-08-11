package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var filename, fromDate, toDate, listenAddr string
	flag.StringVar(&filename, "f", "", "要解析的文件路径")
	flag.StringVar(&fromDate, "from", time.Now().Format("2006-01-02"), "开始日期")
	flag.StringVar(&toDate, "to", time.Now().Format("2006-01-02"), "结束日期")
	flag.StringVar(&listenAddr, "l", "", "获取实时解析crontab -l输出的监听的地址")

	flag.Parse()

	//命令行一次性解析模式
	if filename != "" {
		fmt.Println("分析", filename, "中的计划任务从", fromDate, "到", toDate, "期间的计划任务分布")
		from, _ := time.Parse("2016-01-02", fromDate)
		to, _ := time.Parse("2016-01-02", toDate)
		textDump(filename, from, to)
		return
	}

	//后台WEB服务实时解析模式
	if listenAddr != "" {
		http.HandleFunc("/", httpHandle)
		err := http.ListenAndServe(listenAddr, nil)
		if err!=nil {
			fmt.Println(err)
		}
		return
	}

	//Usage
	flag.Usage()
	fmt.Println("请提供要解析的文件路径或者监听地址")
	fmt.Println("  >> 如果提供了解析文件，则将其解析输出")
	fmt.Println("  >> 如果提供的监听地址，则在该地址监听请求，每次请求会重新获取crontab -l的内容并解析成HTML")
}

func httpHandle(w http.ResponseWriter, r *http.Request) {
	startTime, _ := time.Parse("2006-01-02", r.FormValue("start"))
	endTime, _ := time.Parse("2006-01-02", r.FormValue("end"))

	cmd := exec.Command("crontab", "-l")
	content, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(w, "执行命令`crontab -l`失败:%s", err)
		return
	}

	//解析文件
	schedules, err := parseCrontab(content)
	if err != nil {
		fmt.Println(err)
		return
	}

	showall := r.FormValue("showall")
	if showall != "" {
		startTime = time.Time{}
		endTime = time.Time{}
	}

	scheduleMap, commands := schedulesToMinutesCommandMap(schedules, startTime, endTime)

	templateData := map[string]interface{}{}
	templateData["end"] = endTime
	templateData["start"] = startTime
	templateData["showall"] = showall
	templateData["commands"] = commands
	templateData["scheduleMap"] = scheduleMap

	tmpl, err := template.New("name").Parse(htmlTemplate)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	err = tmpl.Execute(w, templateData)
}

func textDump(filename string, from, to time.Time) {
	//读取文件
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("打开文件%s失败:%s\n", filename, err)
		return
	}

	//解析文件
	schedules, err := parseCrontab(content)
	if err != nil {
		fmt.Println(err)
		return
	}

	scheduleMap, _ := schedulesToMinutesCommandMap(schedules, from, to)

	fmt.Println(strings.Repeat("=", 73))
	fmt.Printf("分钟 =>")
	for minute := 0; minute < 60; minute++ {
		if minute%10 == 0 {
			fmt.Printf("%d", minute)
		} else {
			fmt.Printf("_")
		}
	}

	fmt.Println("|")
	for hour := 0; hour < 24; hour++ {
		fmt.Printf("%02d:00 ", hour)
		for minute := 0; minute < 60; minute++ {
			if minute%10 == 0 {
				fmt.Print("|")
			}
			commands := scheduleMap[hour][minute]
			if len(commands) > 0 {
				fmt.Printf("%d", len(commands))
			} else {
				fmt.Printf("_")
			}
		}
		fmt.Println("|")
	}
	fmt.Println(strings.Repeat("=", 73))
}

//去掉whenever生成的Cron命令前后缀
func wheneverCommandTrim(command string) string {
	command = strings.TrimSuffix(command, "\n")
	command = strings.TrimSuffix(command, "\r")
	command = strings.TrimPrefix(command, "/bin/bash -l -c 'cd /app && RAILS_ENV=production bundle exec ")
	command = strings.TrimPrefix(command, "/bin/bash -l -c 'cd /app && script/")
	command = strings.TrimSuffix(command, " >> /app/log/cron_log.log 2>&1'")
	command = strings.TrimSuffix(command, " --silent")

	command = strings.Replace(command, "'\\'", "", -1)
	return command
}
