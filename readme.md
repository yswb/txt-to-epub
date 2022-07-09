### 简介
把 txt 文件（仅限UTF-8编码）转换成 epub 格式文件，基于 [go-epub](https://github.com/bmaupin/go-epub) 

### 使用

默认以`^第\S+章\s*.+`匹配每章的标题
```shell
./txt-to-epub -s test.txt
```

全部参数：

* `-s` 指定txt源文件
* `-a`指定作者
* `-reg`指定每章标题的正则

```shell
./txt-to-epub -s test.txt -a name -reg '^第\d+章\s*.+'
```

