package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	go func() { //新开一个goroutine用来接收消息
		io.Copy(os.Stdout, conn) // 通过io.Copy来操作I/O，从标准输入中读取数据并写入TCP连接，以及从TCP连接中读取数据并写入标准输出中
		log.Println("done")
		done <- struct{}{} // 新开的goroutine通过一个channel done和main goroutine通信
	}()

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
