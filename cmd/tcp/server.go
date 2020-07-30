package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

//用Go语言进行Socket编程的框架
func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

//broadcaster 用于记录聊天室用户，并进行消息广播
//1.新用户进来；2.用户普通消息；3.用户离开
func broadcaster() {

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	//1.新用户进来，构建该用户的实例
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	//2.当前在一个新的 goroutine 中，用来进行读操作，因此需要开一个 goroutine 用于写操作
	//读写 goroutine 之间可以通过 channel 进行通信
	go sendMessage(conn, user.MessageChannel)

	//3.给当前用户发送欢迎信息，给所有用户告知新用户到来
	user.MessageChannel <- "Welcome, " + user.String()
	msg := Message{
		OwnerID: user.ID,
		Content: "user: `" + strconv.Itoa(user.ID) + "` has enter",
	}
	messageChannel <- msg

	//4.记录到全局用户列表中，避免用锁
	enteringChannel <- user

	//5.循环读取用户输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		messageChannel <- msg
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}

	//6.用户离开
	leavingChannel <- user
	msg.Content = "user: `" + strconv.Itoa(user.ID) + "` has left"
	messageChannel <- msg
}

type User struct {
	ID             int         //用户唯一标识，通过GenUserID函数生成
	Addr           string      //用户的IP地址接口
	EnterAt        time.Time   //用户进入的时间
	MessageChannel chan string //当前用户发送信息的通道
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func (u *User) String() string {
	return u.Addr + ", UID:" + strconv.Itoa(u.ID) + ", Enter At:" +
		u.EnterAt.Format("2006-01-02 15:04:05+8000")
}

// 给用户发送的消息
type Message struct {
	OwnerID int
	Content string
}

var (
	// 新用户到来，通过该 channel 进行登记
	enteringChannel = make(chan *User)
	// 用户离开，通过该 channel 进行登记
	leavingChannel = make(chan *User)
	// 广播专用的用户普通消息 channel，缓冲是尽可能避免出现异常情况堵塞
	messageChannel = make(chan Message, 8)
)

// 生成用户 ID
var (
	globalID int
	idLocker sync.Mutex
)

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalID++
	return globalID
}
