package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/mbarbita/golib-cfgutils"
)

func checkStatus() {
	i := 1
	for {
		res, err := http.Get("http://example.com/info.txt")

		if err == nil {
			// log.Fatal(err)
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				// log.Fatal(err)
				log.Println("body err:", err)
			}
			fmt.Printf("body: %s", body)
			fmt.Println("byte 0,1:", body[0], body[1])
			// fmt.Println("byte 0,1:", body...)
		}

		if err != nil {
			log.Println("get err:", err)

		}

		log.Println("Sleeping...", i)
		time.Sleep(30 * time.Second)
		i++
	}
}

func echo() {
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
				fmt.Println("conn read err:", err)
				for j := 10; j > 0; j-- {
					fmt.Printf("conn reconnecting in: %2v\r", j)
					time.Sleep(1 * time.Second)
				}
				fmt.Println()
				break
			}
			fmt.Printf("conn read status: %v", status)
			time.Sleep(30 * time.Second)
			// ...
		}
	}
}

func main() {

	fmt.Println("Starting...")

	cfgMap := cfgutils.ReadCfgFile("cfg.ini", false)

	cmd := cfgMap["cmd"]
	fields := strings.Split(strings.TrimSpace(cfgMap["arg"]), " ")
	fmt.Println("command:", cmd)
	fmt.Println("args:")
	for _, s := range fields {
		fmt.Print(s + " ")
	}
	fmt.Println()
	go echo()
	i := 1
	for {
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
