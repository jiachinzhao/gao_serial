package gao_serial

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gofrs/flock"
	"github.com/jiachinzhao/serial"
	"github.com/pkg/errors"
)

type ErrPortBlock struct {
	port string
}

func (e ErrPortBlock) Error() string {
	return fmt.Sprintf("ErrPortBlock: port: %s block", e.port)
}

type GaoSerial struct {
	port *serial.Port
	//port.Read and Write may forever block cause Go internal implements when read and write return EAGAIN
	//blockTimeout is set to close port to break this block
	//usually we set to 1s or more, it must higher than serial.Port.ReadTimeout to distinguish whether is block or just no data return when timeout
	blockTimeout time.Duration
	Com          string
	Bandrate     int
	readBuf      []byte
}

var readBufSize = 64

func NewGaoSerial(blockTimeout time.Duration) *GaoSerial {
	return &GaoSerial{blockTimeout: blockTimeout, readBuf: make([]byte, readBufSize)}
}
func (gs *GaoSerial) Open(com string, bandrate int) error {
	if com == "" {
		return errors.Errorf("com name is empty")
	}
	fileLock := flock.New(com)
	locked, err := fileLock.TryLock()
	if err != nil {
		return errors.Wrapf(err, "lock com: %s", com)
	}
	if !locked {
		return fmt.Errorf("com: %s already been used by other process", com)
	}
	_ = fileLock.Unlock()
	port, err := serial.OpenPort(&serial.Config{Name: com, Baud: bandrate})
	if err != nil {
		return errors.WithStack(err)
	}
	gs.port = port
	return nil
}
func (gs *GaoSerial) Close() error {
	if gs.port != nil {
		return gs.port.Close()
	}
	return nil
}

type result struct {
	n   int
	err error
}

func (r result) String() string {
	return fmt.Sprintf("result: n: %d, err: %s", r.n, r.err)
}

//read  close port when read block after block timeout
//it will return EOF after timeout or return at least one byte data
func (gs *GaoSerial) read() (int, error) {
	rCh := make(chan *result)
	go func() {
		result := &result{}
		result.n, result.err = gs.port.Read(gs.readBuf)
		rCh <- result
	}()
	select {
	case r := <-rCh:
		return r.n, r.err
	case <-time.After(gs.blockTimeout):
		_ = gs.Close()
		_ = <-rCh
		return 0, ErrPortBlock{port: gs.Com}
	}
}
func (gs *GaoSerial) SetReadTimeout(timeout time.Duration) error {
	if err := gs.port.SetReadTimeout(timeout); err != nil {
		return errors.Errorf("port SetReadTimeout err: %s", err.Error())
	}
	return nil
}

//Read read data from port, use timeout, or encounter error return
func (gs *GaoSerial) Read(timeout time.Duration) ([]byte, error) {
	var readBytes []byte
	//actually set to timeout, timeout max must <= 25s
	timer := time.NewTimer(timeout)
	readCnt := 0
	gaoTimeout := func(timeout time.Duration) error {
		readCnt++
		if readCnt == 2 {
			timeout = 10 * time.Millisecond
		}
		if readCnt > 2 {
			return nil
		}
		if err := gs.SetReadTimeout(timeout); err != nil {
			return errors.Errorf("setreadtimeout: %d error: %s", timeout, err.Error())
		}
		return nil
	}
	defer gs.port.Flush()
	for {
		select {
		//TODO consider err encounter before timeout, need to stop timer to drain timer.C
		case <-timer.C:
			return readBytes, nil
		default:
			if err := gaoTimeout(timeout); err != nil {
				return nil, err
			}
			n, err := gs.read()
			if err != nil {
				if err != io.EOF {
					return readBytes, err
				}
				return readBytes, nil
			}
			readBytes = append(readBytes, gs.readBuf[:n]...)
		}
	}
}

//write close port when write block after block timeout
func (gs *GaoSerial) write(buf []byte) (int, error) {
	rCh := make(chan *result)
	go func() {
		result := &result{}
		result.n, result.err = gs.port.Write(buf)
		rCh <- result
	}()
	select {
	case r := <-rCh:
		return r.n, r.err
	case <-time.After(gs.blockTimeout):
		_ = gs.Close()
		_ = <-rCh
		return 0, ErrPortBlock{port: gs.Com}
	}
}

//Write write data to port
func (gs *GaoSerial) Write(b []byte) (int, error) {
	return gs.write(b)
}

//WriteAndRead write data and receive response
//wrInterval is time interval between write and read
//readTimeout is time read data
func (gs *GaoSerial) WriteAndRead(b []byte, wrInterval, readTimeout time.Duration) ([]byte, error) {
	if _, err := gs.Write(b); err != nil {
		return nil, err
	}
	time.Sleep(wrInterval)
	bs, err := gs.Read(readTimeout)
	return bs, err
}

//WriteAndReadLines return receive response split by lines and filter empty lines
func (gs *GaoSerial) WriteAndReadLines(b []byte, wrInterval, readTimeout time.Duration) ([][]byte, error) {
	bs, err := gs.WriteAndRead(b, wrInterval, readTimeout)
	bsLines := bytes.Split(bs, []byte("\r\n"))
	var lineNoEmpty int
	for _, line := range bsLines {
		if len(line) > 0 {
			fmt.Println("write and read: ", string(line))
			bsLines[lineNoEmpty] = line
			lineNoEmpty++
		}
	}
	return bsLines[:lineNoEmpty], err
}
