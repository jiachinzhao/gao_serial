package gao_serial

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGaoSerialOpen(t *testing.T) {
	ports := strings.Split(os.Getenv("PORTS"), ",")
	commands := []string{
		"ATZ\r",
		"AT\r",
		"AT+CSQ\r",
		"AT+CGSN\r",
		"AT+CREG?\r",
		"AT+CMGF=1\r",
	}
	for _, port := range ports {
		gao := NewGaoSerial(2 * time.Second)
		fmt.Printf("start to open port: %s\n", port)
		if err := gao.Open(port, 115200); err != nil {
			fmt.Printf("open port: %s error: %s\n", port, err.Error())
			continue
		}
		for _, cmd := range commands {
			fmt.Printf("write cmd: %s\n", cmd)
			if _, err := gao.Write([]byte(cmd)); err != nil {
				fmt.Printf("write %s error: %s\n", cmd, err.Error())
				if _, ok := err.(ErrPortBlock); ok {
					break
				}
			}
			start := time.Now()
			timeout := time.Second
			if strings.Contains(cmd, "test") {
				time.Sleep(2 * time.Second)
			}
			bs, err := gao.Read(timeout)
			fmt.Printf("read cost time: %v\n", time.Since(start))
			if err != nil {
				fmt.Printf("read error: %s\n", err.Error())
				if _, ok := err.(ErrPortBlock); ok {
					break
				}
			}
			fmt.Printf("read content: %s\n", string(bs))
		}
	}
}
