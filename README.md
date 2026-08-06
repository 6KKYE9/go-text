# go-text

纯标准库实现的文本处理小工具，零第三方依赖。处理日志、代码、笔记时常用的几件事。

## 功能

- 大小写转换：转大写 / 转小写 / 首字母大写
- 行去重：去掉重复行，保持首次出现顺序
- 统计：行数 / 词数 / 字符数（中文按字算，不是按字节）
- 排序：按行排序，可配合 `-u` 去重

文件参数可省略，省略时从标准输入读，所以能和管道配合。

## 用法

```bash
# 大小写
go run . case upper notes.txt
echo "hello" | go run . case title

# 去重
go run . dedup log.txt

# 统计
go run . stat essay.txt

# 排序 + 去重
go run . sort -u names.txt
```

## 编译

```bash
go build -o go-text .
```

依赖：仅 Go 标准库（`bufio`、`strings`、`sort`、`unicode/utf8`、`os`、`io`、`fmt`）。
