package main

import (
	"io"
	"os"
)

type EnCoder struct {
	r io.Reader
}

func NewEnCoder(r io.Reader) *EnCoder {
	return &EnCoder{r:r}
}

func(d *EnCoder) Encode() {
	var data = make([]byte, 3)
	for {
		n, err := d.r.Read(data)
		if err == io.EOF {
			return
		}
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		d.handle(data[:n])
	}
}

func (d *EnCoder) handle(data []byte) {
	if len(data) == 0 {
		return
	}
	if len(data) > 3 {
		os.Stderr.WriteString("len(data)>3")
		os.Exit(1)
	}
	if len(data) == 3 {
		// 取出第一个字节的高6位
		t := (0xff >> 2) & (data[0]>>2)
		d.output(Mapping[t])
		// 第一个字节的低2位，第二字节的高4位
		t = ((3 & data[0]) << 4) | (0x0f & (data[1]>>4))
		d.output(Mapping[t])
		// 第二个字节的低4位，第三个字节的高2位
		t = (0x0f & data[1]) << 2 | ((data[2] >> 6) & 3)
		d.output(Mapping[t])
		// 第三个字节的低6位
		t = (data[2]<<2) >> 2
		d.output(Mapping[t])
	} else if len(data) == 2 {
		// 取出第一个字节的高6位
		t := (0xff >> 2) & (data[0]>>2)
		d.output(Mapping[t])
		// 第一个字节的低2位，第二字节的高4位
		t = ((3 & data[0]) << 4) | (0x0f & (data[1]>>4))
		d.output(Mapping[t])
		// 第二个字节的低4位
		t = ((data[1]<<4)>>4) << 2
		d.output(Mapping[t],'=')
	} else {
		// 取出第一个字节的高6位
		t := (0xff >> 2) & (data[0]>>2)
		d.output(Mapping[t])
		// 第一个字节的低2位
		t = ((3 & data[0]) << 4) << 4
		d.output(Mapping[t], '=','=')
	}
}

func (d *EnCoder) output(b ...byte) {
	os.Stdout.Write(b)
}
