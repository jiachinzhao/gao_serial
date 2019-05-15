package gao_serial

import (
	"bytes"
	"fmt"
	"time"
)

var (
	//AT check a equipment whether has at response
	AT = "AT\r"
	//ATZ init equipment configure
	ATZ       = "ATZ\r"
	ECHOCLOSE = "ATE0\r"
	ECHOOPEN  = "ATE1\r"
	//IME is international mobile equipment identity
	IME = "AT+CGSN\r"
	//CSQ is signal intensity of equipment
	CSQ = "AT+CSQ\r"
	//NETWORK is mobile equipment network situation
	NETWORK = "AT+CREG?\r"
)

func CheckAT(gs *GaoSerial) error {
	bs, err := gs.WriteAndRead([]byte(AT), time.Second/10, time.Second/10)
	if err != nil {
		return fmt.Errorf("checkAT err: %s", err.Error())
	}
	if bytes.Contains(bs, []byte("OK")) {
		return nil
	}
	return fmt.Errorf("cmd: %s,response not OK", AT)
}
func CheckATZ(gs *GaoSerial) error {
	bs, err := gs.WriteAndRead([]byte(ATZ), time.Second/10, time.Second/10)
	if err != nil {
		return fmt.Errorf("checkATZ err: %s", err.Error())
	}
	if bytes.Contains(bs, []byte("OK")) {
		return nil
	}
	return fmt.Errorf("cmd: %s,response not OK", ATZ)
}

//SetEchoClose close equipment echo
func SetEchoClose(gs *GaoSerial) error {
	bs, err := gs.WriteAndRead([]byte(ECHOCLOSE), time.Second/10, time.Second/10)
	if err != nil {
		return fmt.Errorf("set echo close err: %s", err.Error())
	}
	if bytes.Contains(bs, []byte("OK")) {
		return nil
	}
	return fmt.Errorf("cmd: %s,response not OK", ECHOCLOSE)
}

//GetIMEI get IME from equipment
func GetIMEI(gs *GaoSerial) (string, error) {
	//关闭回显， 保证读取imei时第一行为imei号码
	if err := SetEchoClose(gs); err != nil {
		_ = SetEchoClose(gs)
	}
	lines, err := gs.WriteAndReadLines([]byte(IME), time.Second/10, time.Second/10)
	if err != nil {
		return "", fmt.Errorf("get imei err: %s", err.Error())
	}
	if len(lines) == 0 {
		return "", fmt.Errorf("get imei err: response is empty")
	}
	return string(lines[0]), nil
}

//GetCSQ get csq  from equipment
func GetCSQ(gs *GaoSerial) (string, error) {
	if err := SetEchoClose(gs); err != nil {
		_ = SetEchoClose(gs)
	}
	//TODO

	return "", nil
}

//GetNetWork return equipment network register situation
func GetNetWork(gs *GaoSerial) (string, error) {
	if err := SetEchoClose(gs); err != nil {
		_ = SetEchoClose(gs)
	}
	//TODO

	return "", nil
}
