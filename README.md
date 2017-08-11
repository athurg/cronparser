# cronparser
Crontab计划任务解析工具，主要用于可视化crontab的计划任务。

# Usage
## 命令行解析文件
```bash
crontab -l /tmp/output
./cronparser -f /tmp/output
```

效果：

![命令行模式解析效果图](https://raw.githubusercontent.com/athurg/cronparser/master/cli_example.png)

## 实时解析
```bash
./cronparser -l :8080
```

效果

![命令行模式解析效果图](https://raw.githubusercontent.com/athurg/cronparser/master/service_example.png)
