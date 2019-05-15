package gao_serial

import (
	"bytes"
	"compress/flate"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"unsafe"
)

func StringBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return
}

func BytesString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//DeepCopy 深度拷贝对象
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func DoFlateCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := flate.NewWriter(&in, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(src)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func DoFlateUnCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r := flate.NewReader(b)
	_, err := io.Copy(&out, r)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

const InvalidString = "<invalid convert>"

func ConvString(jsonItf interface{}) string {
	switch jsonItf.(type) {
	case float64:
		tFloat := jsonItf.(float64)
		tInt := int64(tFloat)
		if float64(tInt) == tFloat {
			return strconv.FormatInt(tInt, 10)
		} else {
			return fmt.Sprint(tFloat)
		}
	case string:
		return jsonItf.(string)
	}
	return InvalidString
}

func ConvInt64(jsonItf interface{}) int64 {
	switch jsonItf.(type) {
	case float64:
		return int64(jsonItf.(float64))
	case string:
		if jsonItf.(string) == "" {
			return 0
		}
		tInt, _ := strconv.ParseInt(jsonItf.(string), 10, 64)
		return tInt
	case bool:
		tBool := jsonItf.(bool)
		if tBool {
			return 1
		} else {
			return 0
		}
	}
	return 0
}
func StrToMap(s string) (map[string]string, error) {
	m := make(map[string]string)
	if len(s) == 0 {
		return m, nil
	}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, fmt.Errorf("%s cannot to Unmarshal to map[string]string", s)
	}
	return m, nil
}

func GetCurrentDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
