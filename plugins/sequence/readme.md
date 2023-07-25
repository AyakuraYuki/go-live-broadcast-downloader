ID生成工具

然后方法中：

```
id := sequence.ID()
fmt.Println(id)
```

为避免 ecs 机器部署场景下，同一业务不同进程可能的 id 冲突，可以自定义机器 id，环境中设置 `SEQUENCE_SN` 环境变量，int64 值

注意事项：

为兼容现有业务，生成的id从70开头，共计19位