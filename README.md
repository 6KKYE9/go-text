# go-text

处理文本的几个常用操作：大小写转换、行去重、统计、排序。

功能：

- 大小写转换：转大写 / 转小写 / 首字母大写
- 行去重：去掉重复行，保留首次出现顺序
- 统计：行数 / 词数 / 字符数（中文按字算，不是按字节）
- 排序：按行排序，可配合 `-u` 去重
- 替换：逐行把某个子串换成另一个（`replace <旧> <新>`）

文件参数省略时从标准输入读，所以能和管道配合；`case`/`dedup`/`sort`/`replace` 都支持 `-o <文件>` 把结果写到文件。

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
