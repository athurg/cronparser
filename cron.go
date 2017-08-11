package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//解析Crontab文件内容
func parseCrontab(content []byte) ([]Schedule, error) {
	lineReader := bufio.NewReader(bytes.NewBuffer(content))

	schedules := []Schedule{}
	for {
		line, err := lineReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("读取行错误: %s", err)
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if len(line) == 0 {
			continue
		}

		strs := strings.SplitN(line, " ", 6)
		if len(strs) != 6 {
			continue
		}

		schedule, err := parseLine(strs)
		if err != nil {
			return nil, fmt.Errorf("解析行出错: %s", err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

//解析Cron配置行
func parseLine(strs []string) (schedule Schedule, err error) {
	//解析分钟
	schedule.Minutes, err = parseWord(strs[0])
	if err != nil {
		return
	}

	//解析小时
	schedule.Hours, err = parseWord(strs[1])
	if err != nil {
		return
	}

	//解析日
	schedule.Days, err = parseWord(strs[2])
	if err != nil {
		return
	}

	//解析月
	schedule.Months, err = parseWord(strs[3])
	if err != nil {
		return
	}

	//解析星期
	schedule.Weekdays, err = parseWord(strs[4])
	if err != nil {
		return
	}

	//调用可选的命令文本修正方法
	schedule.Command = strs[5]
	return
}

//（递归）解析Crontab里的配置值，目前支持以下格式：
// 1. 数字
// 2. “*”，直接返回一个长度为零的slice，error也为nil
// 3. “,”分割的数字列表
// 4. “-”分割的数字范围
// 5. “,”、“-”混用（优先解析逗号，然后解析连字符）
//TODO: 支持“/”分割的数字
func parseWord(word string) ([]int, error) {
	if word == "*" {
		return []int{}, nil
	}

	result := make([]int, 0)

	//遇到逗号，则分割成组，递归调用
	if strings.Index(word, ",") != -1 {
		for _, str := range strings.Split(word, ",") {
			values, err := parseWord(str)
			if err != nil {
				return nil, err
			}
			result = append(result, values...)
		}
		return result, nil
	}

	//遇到连字符，分别解析范围的起止点，然后计算该范围的数字组
	if strings.Index(word, "-") != -1 {
		ranges := strings.Split(word, "-")
		if len(ranges) != 2 {
			return nil, fmt.Errorf("Invalid range %s", word)
		}

		minValues, err := parseWord(ranges[0])
		if err != nil {
			return nil, err
		}

		maxValues, err := parseWord(ranges[1])
		if err != nil {
			return nil, err
		}

		if maxValues[0] < minValues[0] {
			return nil, fmt.Errorf("Invalid range %s", word)
		}

		for value := minValues[0]; value <= maxValues[0]; value++ {
			result = append(result, value)
		}

		return result,nil
	}

	//解析纯数字
	value, err := strconv.Atoi(word)
	if err != nil {
		return nil, err
	}
	result = append(result, value)

	return result, nil
}
