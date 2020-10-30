package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"os"
	"text/tabwriter"
)

var split string
func main() {
	flag.StringVar(&split, "s"," ","分隔符")
	flag.Parse()
	NewPrettyWriter(os.Stdin, split).Write()
}


type PrettyWriter struct{
	in *bufio.Reader
	out  *tabwriter.Writer
	// 分割符
	SplitChar string
}

func NewPrettyWriter(in io.Reader, split string)*PrettyWriter {
	tabW := new(tabwriter.Writer)
	tabW.Init(os.Stdout,0,0,1,' ',0)
	return &PrettyWriter{
		in:        bufio.NewReaderSize(in,4096),
		out:       tabW,
		SplitChar: split,
	}
}

func(p *PrettyWriter) Write() {
	for {
		line, isPrefix, err := p.in.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		for isPrefix {
			var tmp = make([]byte, len(line))
			// 此处有坑，需要ReadLine返回的line和缓冲区公用内存，左右在下次ReadLine时需要先将之前的数据移走
			copy(tmp, line)
			line, isPrefix, err = p.in.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				os.Stderr.WriteString(err.Error())
				os.Exit(1)
			}
			tmp = append(tmp, line...)
			if !isPrefix {
				line = tmp
			}
		}
		//
		if p.SplitChar != "\t" {
			line = bytes.Replace(line, []byte(p.SplitChar), []byte{'\t'}, -1)
		}
		line = append(line, '\n')
		_, err = p.out.Write(line)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	}
	p.out.Flush()
}

