package main

import (
	"io"
	"os"
)

type Decoder struct {
	r io.Reader
	w io.Writer
}

func NewDecoder(r io.Reader, writer io.Writer) *Decoder {
	return &Decoder{r:r,w:writer}
}

func(d *Decoder) Decode() {
	var data = make([]byte, 4)
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

func (d *Decoder) handle(data []byte) {
	if len(data) == 0 {
		return
	}
	if len(data) != 4 {
		os.Stderr.WriteString("len(data)!=4")
		os.Exit(1)
	}
	for i,d := range data {
		if d != '=' {
			data[i] = byte(RMapping[d])
		}
	}
	t := data[0]<<2 | data[1]>>4
	d.output(t)
	if data[2] == '=' {
		return
	}
	t = (data[1]<<4) | data[2]>>2
	d.output(t)
	if data[3] == '=' {
		return
	}
	t = (data[2]<<6) | data[3]
	d.output(t)
}

func (d *Decoder) output(b ...byte) {
	d.w.Write(b)
}

