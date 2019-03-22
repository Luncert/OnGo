package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	cr         = [4]rune{'-', '\\', '|', '/'}
	concurrent = 16
	url        string // = "https://self-billing-receiver.cfapps.sap.hana.ondemand.com/api/idoc"
	body       []byte
	client     *http.Client
	logBuf     = &MutexBuffer{}
)

func main() {
	// init
	if os.Args == nil || len(os.Args) == 1 {
		Usage()
	}
	url = os.Args[1]
	reqNum, err := strconv.Atoi(os.Args[2])
	if nil != err {
		panic(err)
	}
	if len(os.Args) >= 4 {
		concurrent, err = strconv.Atoi(os.Args[3])
		if nil != err {
			panic(err)
		}
	}
	body, err = ioutil.ReadFile("body.xml")
	if nil != err {
		panic(err)
	}

	// make Client
	client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(25 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	// start
	taskQueue := make(chan bool, concurrent)
	wait := make(chan bool)
	go Dispatch(taskQueue, wait, reqNum)

	for i := 0; i < reqNum; i++ {
		taskQueue <- true
	}
	close(taskQueue)
	<-wait
}

// Usage ...
func Usage() {
	fmt.Println("usage: url reqNum [concurrent]")
	os.Exit(0)
}

// Dispatch ...
func Dispatch(taskQueue <-chan bool, signal chan<- bool, total int) {
	processing, finished, failed, crI := 0, 0, 0, 0
	wait := make(chan bool, concurrent)

	process := func() {
		fmt.Printf("%c hitting %.2f%%\r", cr[crI], float32(finished)*float32(100)/float32(total))
		crI = (crI + 1) % 4
	}
	waitOne := func() {
		if !<-wait {
			failed++
		}
		processing--
		finished++
		process()
	}

	start := time.Now()
	process()
	for range taskQueue {
		go makeRequest(wait)
		processing++
		if processing == concurrent {
			waitOne()
		}
	}
	for processing > 0 {
		waitOne()
	}
	totalTime := time.Since(start).Seconds()

	fmt.Printf("total time: %.6f sec\n", totalTime)
	fmt.Printf("qps: %.6f\n", float64(total)/totalTime)
	fmt.Printf("total: %d\n", total)
	fmt.Printf("failed: %d\n", failed)
	ioutil.WriteFile("failed.log", logBuf.b.Bytes(), 0777)
	signal <- true
}

func makeRequest(ch chan<- bool) {
	rep, err := client.Post(url, "application/xml", bytes.NewBuffer(body))
	successed := err == nil && rep.StatusCode == 200
	if !successed {
		var detail string
		if rep != nil {
			detail = fmt.Sprintf("[%v] %v\n", rep.StatusCode, err)
		} else {
			detail = fmt.Sprintf("[no response] %v\n", err)
		}
		logBuf.Write([]byte(detail))
	}
	ch <- successed
}

// MutexBuffer ...
type MutexBuffer struct {
	b  bytes.Buffer
	rw sync.RWMutex
}

func (mb *MutexBuffer) Read(p []byte) (n int, err error) {
	mb.rw.RLock()
	defer mb.rw.RUnlock()
	return mb.b.Read(p)
}

func (mb *MutexBuffer) Write(p []byte) (n int, err error) {
	mb.rw.Lock()
	defer mb.rw.Unlock()
	return mb.b.Write(p)
}
