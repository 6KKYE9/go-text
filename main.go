// go-text 是文本处理小工具。
// 平时处理日志、代码、笔记时常要做的事：大小写转换、行去重、统计字数。
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

// readLines 从标准输入或文件逐行读出，返回不带换行的字符串切片。
func readLines(path string) ([]string, error) {
	var r *bufio.Scanner
	if path == "" || path == "-" {
		r = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = bufio.NewScanner(f)
	}
	// 默认单行上限 64K，给个更大的缓冲避免长行被截断。
	r.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lines := []string{}
	for r.Scan() {
		lines = append(lines, r.Text())
	}
	return lines, r.Err()
}

// writeLines 把结果写回标准输出（outPath 非空时写到文件）。
func writeLines(lines []string, outPath string) error {
	if outPath != "" {
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		for _, l := range lines {
			fmt.Fprintln(w, l)
		}
		return w.Flush()
	}
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, l := range lines {
		fmt.Fprintln(w, l)
	}
	return nil
}

// cmdCase 大小写转换：go-text case <upper|lower|title> [文件] [-o 输出]
func cmdCase(args []string, outPath string) {
	if len(args) < 1 {
		fmt.Println("用法: go-text case <upper|lower|title> [文件]")
		return
	}
	mode := strings.ToLower(args[0])
	path := ""
	if len(args) >= 2 {
		path = args[1]
	}
	lines, err := readLines(path)
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}
	for i, l := range lines {
		switch mode {
		case "upper":
			lines[i] = strings.ToUpper(l)
		case "lower":
			lines[i] = strings.ToLower(l)
		case "title":
			lines[i] = strings.Title(strings.ToLower(l))
		default:
			fmt.Println("未知模式:", mode)
			return
		}
	}
	if err := writeLines(lines, outPath); err != nil {
		fmt.Println("写出失败:", err)
	}
}

// cmdDedup 行去重（保留首次出现顺序）：go-text dedup [文件] [-o 输出]
func cmdDedup(args []string, outPath string) {
	path := ""
	if len(args) >= 1 {
		path = args[0]
	}
	lines, err := readLines(path)
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}
	seen := map[string]bool{}
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if seen[l] {
			continue
		}
		seen[l] = true
		out = append(out, l)
	}
	fmt.Printf("去重前 %d 行，去重后 %d 行\n", len(lines), len(out))
	if err := writeLines(out, outPath); err != nil {
		fmt.Println("写出失败:", err)
	}
}

// cmdReplace 在每行中替换子串：go-text replace <旧> <新> [文件] [-o 输出]
func cmdReplace(args []string, outPath string) {
	if len(args) < 2 {
		fmt.Println("用法: go-text replace <旧> <新> [文件]")
		return
	}
	old, nw := args[0], args[1]
	path := ""
	if len(args) >= 3 {
		path = args[2]
	}
	lines, err := readLines(path)
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}
	for i, l := range lines {
		lines[i] = strings.ReplaceAll(l, old, nw)
	}
	if err := writeLines(lines, outPath); err != nil {
		fmt.Println("写出失败:", err)
	}
}

// cmdStat 统计：行数/词数/字符数：go-text stat [文件]
func cmdStat(args []string) {
	path := ""
	if len(args) >= 1 {
		path = args[0]
	}
	// 读原始字节才能数到准确字符数。
	var data []byte
	var err error
	if path == "" || path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}
	text := string(data)
	lines := strings.Count(text, "\n")
	if len(text) > 0 && !strings.HasSuffix(text, "\n") {
		lines++ // 最后一行没有换行符也要算上
	}
	words := 0
	for _, w := range strings.Fields(text) {
		words += len(strings.Fields(w))
	}
	// utf8.RuneCountInString 才是正确的"字符数"，不是字节数。
	chars := utf8.RuneCountInString(text)
	fmt.Printf("行数: %d\n词数: %d\n字符数: %d\n", lines, wordCount(text), chars)
}

// wordCount 用 strings.Fields 按空白切词统计词数。
func wordCount(text string) int {
	return len(strings.Fields(text))
}

// cmdSort 排序行（可选去重）：go-text sort [-u] [文件] [-o 输出]
func cmdSort(args []string, outPath string) {
	uniq := false
	rest := args
	if len(args) > 0 && args[0] == "-u" {
		uniq = true
		rest = args[1:]
	}
	path := ""
	if len(rest) >= 1 {
		path = rest[0]
	}
	lines, err := readLines(path)
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}
	sort.Strings(lines)
	if uniq {
		out := lines[:0]
		for _, l := range lines {
			if len(out) == 0 || out[len(out)-1] != l {
				out = append(out, l)
			}
		}
		lines = out
	}
	if err := writeLines(lines, outPath); err != nil {
		fmt.Println("写出失败:", err)
	}
}

func usage() {
	fmt.Print(`go-text 文本处理工具

用法:
  go-text case  <upper|lower|title> [文件] [-o 输出]   大小写/首字母大写
  go-text dedup [文件] [-o 输出]                        行去重（保持顺序）
  go-text stat  [文件]                                  统计行/词/字符数
  go-text sort  [-u] [文件] [-o 输出]                   排序行，-u 同时去重
  go-text replace <旧> <新> [文件] [-o 输出]            逐行替换子串

文件可省，省略时从标准输入读（支持管道）；-o 指定输出文件
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	// 通用 -o 输出参数
	outPath := ""
	restArgs := os.Args[2:]
	for i := 0; i < len(restArgs); i++ {
		if restArgs[i] == "-o" && i+1 < len(restArgs) {
			outPath = restArgs[i+1]
			restArgs = append(restArgs[:i], restArgs[i+2:]...)
			i--
		}
	}
	switch os.Args[1] {
	case "case":
		cmdCase(restArgs, outPath)
	case "dedup":
		cmdDedup(restArgs, outPath)
	case "stat":
		cmdStat(restArgs)
	case "sort":
		cmdSort(restArgs, outPath)
	case "replace":
		cmdReplace(restArgs, outPath)
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Println("未知命令:", os.Args[1])
		usage()
	}
}
