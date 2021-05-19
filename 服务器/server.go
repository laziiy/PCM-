package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	/**
	  创建监听的地址，并且指定udp协议
	*/
	udp_addr, err := net.ResolveUDPAddr("udp", "172.27.153.114:4328")
	if err != nil {
		fmt.Println("获取监听地址失败,错误原因: ", err)
		return
	}

	/**
	  创建数据通信socket
	*/
	conn, err := net.ListenUDP("udp", udp_addr)
	if err != nil {
		fmt.Println("开启UDP监听失败,错误原因: ", err)
		return
	}
	defer conn.Close()

	for i := 0; ; i++ {
		fmt.Println("开启监听...")
		buf := make([]byte, 1024)

		/**
		  通过ReadFromUDP可以读取数据，可以返回如下三个参数:
		      dataLength:
		          数据的长度
		      raddr:
		          远程的客户端地址
		      err:
		          错误信息
		*/
		dataLength, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("获取客户端传递数据失败,错误原因: ", err)
			return
		}
		fmt.Println("获取到客户端的数据为: ", string(buf[:dataLength]))

		/**
		  写回数据
		*/
		fmt.Println(i)
		// if i%4 == 0 {
		// 	f, err := os.Open("c.txt")

		// 	if err != nil {
		// 		fmt.Println("os.Open err = ", err)
		// 		return
		// 	}
		// 	defer f.Close()
		// 	for {
		// 		buf1 := make([]byte, 1024*4)
		// 		n, err := f.Read(buf1)
		// 		fmt.Println(n)
		// 		if err != nil {
		// 			if err == io.EOF {
		// 				fmt.Println("文件发送完毕")
		// 				break
		// 			} else {
		// 				fmt.Println("f.Read err = ", err)
		// 			}
		// 		}
		// 		fmt.Println(raddr)
		// 		time.Sleep(800 * time.Millisecond)
		// 		conn.WriteToUDP([]byte(strings.ToUpper(string(buf1[:n]))), raddr)
		// 	}
		// 	// conn.WriteToUDP([]byte("你好"), raddr)
		// }
		// if i%4 == 1 {
		// 	f, err := os.Open("c.txt")

		// 	if err != nil {
		// 		fmt.Println("os.Open err = ", err)
		// 		return
		// 	}
		// 	defer f.Close()
		// 	for {
		// 		buf1 := make([]byte, 1024*4)
		// 		n, err := f.Read(buf1)
		// 		fmt.Println(n)
		// 		if err != nil {
		// 			if err == io.EOF {
		// 				fmt.Println("文件发送完毕")
		// 				break
		// 			} else {
		// 				fmt.Println("f.Read err = ", err)
		// 			}
		// 		}
		// 		fmt.Println(raddr)
		// 		time.Sleep(800 * time.Millisecond)
		// 		conn.WriteToUDP([]byte(strings.ToUpper(string(buf1[:n]))), raddr)
		// 	}
		// }
		// if i%4 == 2 {
		// 	f, err := os.Open("c.txt")

		// 	if err != nil {
		// 		fmt.Println("os.Open err = ", err)
		// 		return
		// 	}
		// 	defer f.Close()
		// 	for {
		// 		buf1 := make([]byte, 1024*4)
		// 		n, err := f.Read(buf1)
		// 		fmt.Println(n)
		// 		if err != nil {
		// 			if err == io.EOF {
		// 				fmt.Println("文件发送完毕")
		// 				break
		// 			} else {
		// 				fmt.Println("f.Read err = ", err)
		// 			}
		// 		}
		// 		fmt.Println(raddr)
		// 		time.Sleep(100 * time.Millisecond)
		// 		conn.WriteToUDP([]byte(strings.ToUpper(string(buf1[:n]))), raddr)
		// 	}
		// }
		// if i%4 == 3 {
		f, err := os.Open("c.txt")

		if err != nil {
			fmt.Println("os.Open err = ", err)
			return
		}
		defer f.Close()
		for {
			buf1 := make([]byte, 1024*4)
			n, err := f.Read(buf1)
			fmt.Println(n)
			if err != nil {
				if err == io.EOF {
					fmt.Println("文件发送完毕")
					break
				} else {
					fmt.Println("f.Read err = ", err)
				}
			}
			fmt.Println(raddr)
			time.Sleep(800 * time.Millisecond)
			conn.WriteToUDP([]byte(strings.ToUpper(string(buf1[:n]))), raddr)
		}
		// }

	}
}
