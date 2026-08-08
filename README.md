# go-text

复制粘贴改格式改到手酸？这玩意儿一行就搞定。

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

# 逐行替换
go run . replace "foo" "bar" notes.txt -o out.txt
```
