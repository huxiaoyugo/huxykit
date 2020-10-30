## 背景
使用go sdk提供的tabwriter，更加友好的输出文本到命令行

## 安装
go get github.com/huxiaoyugo/huxykit/pw

## 使用方法

#### 文件test.txt
```
1 2 3 4
aaa bbb ccc ddd
```
#### 命令
```shell script
cat test.txt | pw
```

#### 输出
```
1    2 3 4 5
aaaa b c d d
```







