package csvproc

import (
	"io"
	"strconv"
	"encoding/csv"
	"math/rand"
	"encoding/hex"
	"encoding/binary"
	"compress/zlib"
	"bytes"
	"encoding/base64"
	"encoding/json"
)


type File struct {
	Headers []string
	Data [][]float32
}


func Load(r io.Reader) (*File, error) {
	cr := csv.NewReader(r)
	headers, err := cr.Read()
	if err != nil {
		return nil, err
	}

	var ret File
	for _,v := range headers {
		ret.Headers = append(ret.Headers, v)
	}
	var line = 2
	for {
		row,err := cr.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
		var floatrow = make([]float32, len(ret.Headers))	
		for col,v := range row {
			f,err := strconv.ParseFloat(v,32)
			if err != nil {
				return nil, &csv.ParseError{Line:line, Column:col, Err:err}
			}
			floatrow[col] = float32(f)
		}
		line++
		ret.Data = append(ret.Data, floatrow)
	}
	return &ret, nil
}

func (f *File) Store(w io.Writer) error {
	cw := csv.NewWriter(w)	
	err := cw.Write(f.Headers)
	if err != nil {
		return err
	}

	for _,row := range f.Data {
		var rowStrings = make([]string, len(row))
		for k,v := range row {
			rowStrings[k] = strconv.FormatFloat(float64(v), 'G', -1, 32)
		}
		err = cw.Write(rowStrings)
		if err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}


func genHeaderName(nameLen int) (string) {
	bs := make([]byte, nameLen/2)
	for i := range bs {
		bs[i] = byte(rand.Intn(256))
	}
	return hex.EncodeToString(bs)
}

func Generate(rows, cols int) *File {
	var ret File
	for i:=0;i<cols;i++ {
		ret.Headers = append(ret.Headers,genHeaderName(8))		
	}
	for i:=0;i<rows;i++ {
		var row []float32
		for j:=0;j<cols;j++ {
			row = append(row, float32(rand.Intn(11)-5 + 2048))
		}
		ret.Data = append(ret.Data, row)
	}
	return &ret
}

func (f *File) ExtractWaves() []Wave {
	var ret = make([]Wave,len(f.Headers))
	for i := range ret {
		ret[i].Name = f.Headers[i]
	}

	for _,row := range f.Data {
		for i,v := range row {
			ret[i].Data = append(ret[i].Data, v)
		}
	}
	return ret
}

type Wave struct {
	Name string
	Data []float32
}

type encWave struct {
	Name string
	Value string
}

func (w *Wave) Store(wr io.Writer) error {
		var buf bytes.Buffer
		b64enc := base64.NewEncoder(base64.URLEncoding,&buf)
		zlibEnc := zlib.NewWriter(b64enc)
		err := binary.Write(zlibEnc, binary.LittleEndian, w.Data)
		if err != nil {
			return err
		}
		zlibEnc.Close()
		b64enc.Close()
		js := json.NewEncoder(wr)
		return js.Encode(encWave{Name: w.Name, Value: buf.String()})
}
