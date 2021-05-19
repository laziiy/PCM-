package main

import (
	"fmt"
	//"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	var ip [3]string
	var path_selected_num [3]int64
	var path0_reward [1000]float64
	var path1_reward [1000]float64
	var path2_reward [1000]float64
	var sum_reward float64
	var min_reward float64
	min_reward = 99999
	var accumulate_regret int
	i := 0
	validaddr := "192.168.1.203"
	//	j_count := 0

	//获取本地IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//检查IP地址，判断是否为回环地址,将ip地址存入ip[i]中
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			//			j_count++
			if ipnet.IP.To4() != nil {
				if ipnet.IP.String() != validaddr{
					ip[i] = ipnet.IP.String()
					fmt.Println(ip[i])
					i++}
			}

		}
	}
	//初始化3条路径
	for i = 0; i <= 999; i++ {
		//		duration := time.Duration(1) * time.Second
		time.Sleep(500 * time.Millisecond) // 1s
		t1 := time.Now().UnixNano() / 1e6  // unix时间戳，毫秒
		t11 := time.Now()
		//		fmt.Println(i % 3)
		//		localip := net.ParseIP(ip[i%3])
		localip := net.ParseIP(ip[0]) //ipv4/v6 ，返回地址
		//	fmt.Println(i)
		lAddr := &net.UDPAddr{IP: localip}

		socket, err := net.DialUDP("udp", lAddr, &net.UDPAddr{
			IP:   net.IPv4(192, 168, 0, 100),
			Port: 20001,
		})
		if err != nil {
			fmt.Println("error,err:", err)
			return
		} // 创建套接字，与服务器连接

		defer socket.Close()

		sendData := []byte("hello")

		_, err = socket.Write(sendData)
		if err != nil {
			fmt.Println("error,err:", err)
			return
		} //向服务器发送数据
		//		if i == 1 {
		data := make([]byte, 4096)
		n, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("error,err:", err)
			return
		} //接收到的数据
		f, err := os.OpenFile(`get.txt`, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660) //存接收数据
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.WriteString(string(data[:n]))
		f.WriteString("\r\n") //存数据
		
		fmt.Printf("SIP:%v,DIP:%v,Bytes:%v\n", localip, remoteAddr, n)
		//		}
		fmt.Println("文件接收完毕")
		/*if i == 0 {
			time.Sleep(40 * time.Millisecond)
		}*/
		// if i <= 300 {
		// 	if i%4 == 2 {
		// 		time.Sleep(100 * time.Millisecond)
		// 	} else {
		// 		time.Sleep(200 * time.Millisecond)
		// 	}
		// }
		// if i > 300 {
		// 	if i%4 == 1 {
		// 		time.Sleep(100 * time.Millisecond)
		// 	} else {
		// 		time.Sleep(200 * time.Millisecond)
		// 	}
		// }
		// if i <= 300 {
		// 	if i%4 != 2 {
		// 		time.Sleep(200 * time.Millisecond)
		// 	}
		// }
		// if i > 300 {
		// 	if i%4 != 1 {
		// 		time.Sleep(200 * time.Millisecond)
		// 	}
		// }
		/*ret := rand.NormFloat64()*100 + 100
		ret1 := time.Duration(ret) * time.Millisecond
		time.Sleep(ret1)*/
		t2 := time.Now().UnixNano() / 1e6
		t22 := time.Now()
		fmt.Println(t22.Sub(t11))
		f, err = os.OpenFile(`test.txt`, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.WriteString(t22.Sub(t11).String())
		f.WriteString("\r\n") //存延时
		f, err = os.OpenFile(`path.txt`, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(strconv.Itoa(i % 3))
		f.WriteString("\r\n") //存路径
		switch i % 3 {
		case 0:
			path0_reward[path_selected_num[0]+1] = float64(t2 - t1)
			path_selected_num[0] += 1
			sum_reward += path0_reward[path_selected_num[0]]
			if min_reward > path0_reward[path_selected_num[0]] {
				min_reward = path0_reward[path_selected_num[0]]
			}
		case 1:
			path1_reward[path_selected_num[1]+1] = float64(t2 - t1)
			path_selected_num[1] += 1
			sum_reward += path1_reward[path_selected_num[1]]
			if min_reward > path1_reward[path_selected_num[1]] {
				min_reward = path1_reward[path_selected_num[1]]
			}
		case 2:
			path2_reward[path_selected_num[2]+1] = float64(t2 - t1)
			path_selected_num[2] += 1
			sum_reward += path2_reward[path_selected_num[2]]
			if min_reward > path2_reward[path_selected_num[2]] {
				min_reward = path2_reward[path_selected_num[2]]
			}
		default:
			fmt.Println(err)
		}
		fmt.Println(i%3)
		accumulate_regret += int(0.8 * (sum_reward - float64(i)*min_reward))
		fmt.Println(accumulate_regret)
		fmt.Println(i)
		f, err = os.OpenFile(`regret.txt`, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	    if err != nil {
		    panic(err)
		}
		defer f.Close()
		f.WriteString(strconv.Itoa(accumulate_regret))
		f.WriteString("\r\n")
	}
	// fmt.Println("jishu1", path_selected_num[0])
	// fmt.Println("jishu2", path_selected_num[1])
	// fmt.Println("jishu3", path_selected_num[2])
	// fmt.Println("lujing1", path0_reward[0])
	// fmt.Println("lujing2", path1_reward[0])
	// fmt.Println("lujing3", path2_reward[0])
	// fmt.Println("lujing4", path3_reward[0])
}
