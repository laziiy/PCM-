package main

import (
	"fmt"
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
	for i = 0; i <= 2; i++ {
		//		duration := time.Duration(1) * time.Second
		time.Sleep(500 * time.Millisecond) // 1s
		t1 := time.Now().UnixNano() / 1e6  // unix时间戳，毫秒
		t11 := time.Now()
		//		localip := net.ParseIP(ip[i%4])
		localip := net.ParseIP(ip[i]) //ipv4/v6 ，返回地址
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

		//fmt.Printf("send:%v,SIP:%v,DIP:%v,Bytes:%v\n", string(data[:n]), localip, remoteAddr, n)
		fmt.Printf("SIP:%v,DIP:%v,Bytes:%v\n", localip, remoteAddr, n)

		fmt.Println("文件接收完毕")
		/*
		if i == 2 {
			time.Sleep(100 * time.Millisecond)
		} else if i == 1 {
			time.Sleep(150 * time.Millisecond)
		} else {
			time.Sleep(200 * time.Millisecond)
		}
		if i == 0 {
			time.Sleep(40 * time.Millisecond)
		}*/
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
		f.WriteString(strconv.Itoa(i))
		f.WriteString("\r\n") //存路径
		if i == 0 {
			path0_reward[0] = float64(t2 - t1) //路径i的第一次TCT
			path_selected_num[0] = 1           //路径i选择一次
		}
		if i == 1 {
			path1_reward[0] = float64(t2 - t1)
			path_selected_num[1] = 1
		}
		if i == 2 {
			path2_reward[0] = float64(t2 - t1)
			path_selected_num[2] = 1
		}
	}

	//UCB/greddy/polling
	var path_sum_reward float64
	var path_average_reward [3]float64
	var path_selected int
	var j int64
	var sum_reward float64
	var min_reward float64
	min_reward = 99999
	var accumulate_regret int

	for i = 3; i <= 999; i++ { //for start
		for j = 0; j < path_selected_num[0]; j++ {
			path_sum_reward += path0_reward[j]
		}
		path_average_reward[0] = path_sum_reward / float64(path_selected_num[0]) //路径i当前的平均收益

		path_sum_reward = 0

		for j = 0; j <= (path_selected_num[1] - 1); j++ { //默认第一次选择路径0
			path_sum_reward += path1_reward[j]
		}
		path_average_reward[1] = path_sum_reward / float64(path_selected_num[1]) //路径i当前的平均收益

		path_sum_reward = 0

		for j = 0; j <= (path_selected_num[2] - 1); j++ {
			path_sum_reward += path2_reward[j]
		}
		path_average_reward[2] = path_sum_reward / float64(path_selected_num[2]) //路径i当前的平均收益
		path_sum_reward = 0

		path_selected = 0
		minReward := path_average_reward[0]
		for k := 0; k < 3; k++{
			if path_average_reward[k] < minReward { //选择平均收益最高的路径
				minReward = path_average_reward[k]
				path_selected = k
			}
		}

		fmt.Println(path_selected)
		time.Sleep(500 * time.Millisecond)
		localip := net.ParseIP(ip[path_selected])
		t11 := time.Now()
		t1 := time.Now().UnixNano() / 1e6
		fmt.Println(i)
		lAddr := &net.UDPAddr{IP: localip}

		socket, err := net.DialUDP("udp", lAddr, &net.UDPAddr{
			IP:   net.IPv4(192, 168, 0, 100),
			Port: 20001,
		})
		if err != nil {
			fmt.Println("error,err:", err)
			return
		}

		defer socket.Close()

		sendData := []byte("hello")

		_, err = socket.Write(sendData)
		if err != nil {
			fmt.Println("error,err:", err)
			return
		}
		data := make([]byte, 4096)
		n, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("error,err:", err)
			return
		}
		f, err := os.OpenFile(`get.txt`, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660) //存接收数据
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.WriteString(string(data[:n]))
		f.WriteString("\r\n") //存数据
		fmt.Printf("SIP:%v,DIP:%v,Bytes:%v\n", localip, remoteAddr, n)
		/*
		if i <= 300 {
			if path_selected == 2 {
				time.Sleep(100 * time.Millisecond)
			} else {
				time.Sleep(200 * time.Millisecond)
			}
		}
		if i > 300 {
			if path_selected == 1 {
				time.Sleep(100 * time.Millisecond)
			} else {
				time.Sleep(200 * time.Millisecond)
			}
		}*/
		//		fmt.Println("文件接收完毕")
		//fmt.Printf("send:%v,SIP:%v,DIP:%v,Bytes:%v\n", string(data[:n]), localip, remoteAddr, n)
		t2 := time.Now().UnixNano() / 1e6
		t22 := time.Now()
		fmt.Println(t22.Sub(t11))
		f, err = os.OpenFile(`test.txt`, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(t22.Sub(t11).String())
		f.WriteString("\r\n")
		f, err = os.OpenFile(`path.txt`, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(strconv.Itoa(path_selected))
		f.WriteString("\r\n")
		switch path_selected {
		case 0:
			path0_reward[path_selected_num[0]] = float64(t2 - t1)
			sum_reward += float64(t2 - t1)
			path_selected_num[0] += 1
			if min_reward > float64(t2-t1) {
				min_reward = float64(t2 - t1)
			}
		case 1:
			path1_reward[path_selected_num[1]] = float64(t2 - t1)
			sum_reward += float64(t2 - t1)
			path_selected_num[1] += 1
			if min_reward > float64(t2-t1) {
				min_reward = float64(t2 - t1)
			}
		case 2:
			path2_reward[path_selected_num[2]] = float64(t2 - t1)
			sum_reward += float64(t2 - t1)
			path_selected_num[2] += 1
			if min_reward > float64(t2-t1) {
				min_reward = float64(t2 - t1)
			}
		default:
			fmt.Println(err)
		}
		accumulate_regret += int(0.8 * (sum_reward - float64(i-3)*min_reward))
		fmt.Println(accumulate_regret)
		fmt.Println(i)
		f, err = os.OpenFile(`regret.txt`, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		if err != nil {
			panic(err)
	    }
		defer f.Close()
		f.WriteString(strconv.Itoa(accumulate_regret))
		f.WriteString("\r\n")
		// if path_selected == 0 {
		// 	path0_reward[path_selected_num[0]] = float64(t2 - t1)
		// 	path_selected_num[0] += 1
		// }
		// if path_selected == 1 {
		// 	path1_reward[path_selected_num[1]] = float64(t2 - t1)
		// 	path_selected_num[1] += 1
		// }
		// if path_selected == 2 {
		// 	path2_reward[path_selected_num[2]] = float64(t2 - t1)
		// 	path_selected_num[2] += 1
		// }
		// if path_selected == 3 {
		// 	path3_reward[path_selected_num[3]] = float64(t2 - t1)
		// 	path_selected_num[3] += 1
		// }

	}

}
