package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	cfgutils "github.com/mbarbita/golib-cfgutils"
)

func checkStatus() {
	// i := 1
	// for {
	// 	res, err := http.Get("http://example.com/info.txt")
	//
	// 	if err == nil {
	// 		// log.Fatal(err)
	// 		body, err := ioutil.ReadAll(res.Body)
	// 		res.Body.Close()
	// 		if err != nil {
	// 			// log.Fatal(err)
	// 			log.Println("body err:", err)
	// 		}
	// 		fmt.Printf("body: %s", body)
	// 		fmt.Println("byte 0,1:", body[0], body[1])
	// 		// fmt.Println("byte 0,1:", body...)
	// 	}
	//
	// 	if err != nil {
	// 		log.Println("get err:", err)
	//
	// 	}
	//
	// 	log.Println("Sleeping...", i)
	// 	time.Sleep(30 * time.Second)
	// 	i++
	// }
}

func clearscreen() {
	opSys := cfgMap["os"]
	switch opSys {
	case "windows":
		//Windows
		c := exec.Command("cmd", "/c", "cls")
		c.Stdout = os.Stdout
		c.Run()
	case "linux":
		//linux
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func echo() {
	time.Sleep(10 * time.Second)
	i := 1
	for {
	loop:
		fmt.Println("echo dial:", i)
		conn, err := net.Dial("tcp", "localhost:17700")
		if err != nil {
			// handle error
			fmt.Println("echo dial err:", err)
			// time.Sleep(10*time.Second)
			for j := 10; j > 0; j-- {
				fmt.Printf("echo dial reconnecting in: %2v\r", j)
				time.Sleep(1 * time.Second)
			}
			fmt.Println()
			i++
			goto loop
		}
		for {
			// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
			fmt.Fprintf(conn, "bzzz\n")
			status, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println()
				fmt.Println("conn read err:", err)
				for j := 10; j > 0; j-- {
					fmt.Printf("conn reconnecting in: %2v\r", j)
					time.Sleep(1 * time.Second)
				}
				fmt.Println()
				break
			}
			t := time.Now()
			fmt.Printf("%v last conn read status: %v\r", t.Format(time.Stamp), status[:len(status)-1])
			time.Sleep(30 * time.Second)
			// ...
		}
	}
}

var cfgMap = cfgutils.ReadCfgFile("cfg.ini", false)

func main() {
	clearscreen()
	if cfgMap["echo"] == "on" {
		go echo()
	}

	fmt.Println("Starting...")

	cmd := cfgMap["cmd"]
	fields := strings.Split(strings.TrimSpace(cfgMap["arg"]), " ")
	fmt.Println("command:", cmd)
	fmt.Println("args:")
	for _, s := range fields {
		fmt.Print(s + " ")
	}
	fmt.Println()

	// time.Sleep(10 * time.Second)
	// for j := 10; j > 0; j-- {
	// 	fmt.Printf("continue in: %2v\r", j)
	// 	time.Sleep(1 * time.Second)
	// }

	i := 1
	for {
		clearscreen()
		log.Println("plink connecting...", i)
		command := exec.Command(cmd, fields...)
		err := command.Start()
		if err != nil {
			// log.Fatal(err)
			log.Println("plink command start error:", err)
		}

		log.Println("waiting for plink command to finish...")
		err = command.Wait()
		log.Printf("command plink finished with error: %v\n", err)

		i++
		for j := 10; j > 0; j-- {
			fmt.Printf("reconnecting in: %2v\r", j)
			time.Sleep(1 * time.Second)
		}
		fmt.Println()
	}
}
