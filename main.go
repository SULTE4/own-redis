package main

import (
	"flag"
	"fmt"
	"log/slog"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var port = flag.String("port", "8080", "port for redis")

type ORedis struct {
	storage map[string]Value
	mu      sync.RWMutex
}

func NewORedis() *ORedis {
	st := make(map[string]Value)
	return &ORedis{
		storage: st,
	}
}

type Value struct {
	data       string
	expireTime time.Time
}

func (v *Value) ToString() {
	fmt.Printf("\"data: %s,\t expireTime: %d\"\n", v.data, v.expireTime)
}

func main() {
	flag.Parse()
	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%s", *port))
	if err != nil {
		slog.Error(fmt.Sprintf("Occured error to create server on port: %s", *port), err.Error())
	}
	defer conn.Close()
	slog.Info("Redis server start running")

	redis := NewORedis()

	buf := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		// someCheck(buf, n)
		responseText := handleMessage(*redis, buf, n)
		fmt.Println(buf[:n])
		// slog.Info(responseText)
		t := []byte(responseText + "\n")
		if _, err := conn.WriteTo(t, addr); err != nil {
			slog.Error(err.Error())
		}
		// buf = buf[:0]
	}
}

func handleMessage(redis ORedis, buf []byte, n int) string {
	args := strings.Fields(string(buf))
	argCount := len(args) - 1

	argCount = countWord(args, n-1)

	fmt.Println("argCount:", argCount, "byte count:", n)
	someCheck(args)
	// for _, el := range args {
	// 	fmt.Printf("%s\t%T", el, el)
	// 	fmt.Println()
	// }
	// fmt.Println()
	firstKey := strings.ToLower(args[0])
	switch firstKey {
	case "ping":
		return "PONG"
	case "set":
		if argCount < 3 {
			return SetErr
		}
		key := args[1]
		if strings.ToLower(key) == "px" {
			return "Invalid key"
		}

		value := ""
		for i, el := range args[2:argCount] {
			if el == "px" {
				if i != 1 && i == argCount-2 {
					return SetErr
				}
				break
			}
			value += el + " "
		} //TODO: REMOVE LAST SPACE

		duration := time.Second * math.MaxInt32
		if args[argCount-2] == "px" {
			tmp, err := strconv.Atoi(args[argCount-1])
			if err != nil || tmp <= 500 {
				return "Invalid time duration"
			}
			duration = time.Second * time.Duration(tmp/1000)
			slog.Info("duration integrating", slog.String("time (sec):", duration.String()))
		}

		nowTime := time.Now()
		expireTime := nowTime.Add(duration)

		v := Value{data: value, expireTime: expireTime}
		// fmt.Printf("key: %s\t", key)
		// v.ToString()

		redis.mu.Lock()
		redis.storage[key] = v
		redis.mu.Unlock()
		return "OK"

	case "get":
		if argCount != 2 {
			return "(nil)"
		}
		key := args[1]
		value := redis.storage[key]
		value.ToString()
		currentTime := time.Now()
		if value.expireTime.Before(currentTime) {
			slog.Info("key time expired")
			// maps.DeleteFunc(redis.storage, func(k string, v Value) bool{
			// 	return
			// })
			return "(nil)"
		}
		return value.data
	default:
		return "invalid input"
	}
}

var SetErr = "(error) ERR wrong number of arguments for 'SET' command"

func countWord(args []string, n int) int {
	cnt := 0
	for _, el := range args {
		l := len(el)
		if n <= l {
			break
		}
		n -= l
		cnt++
	}
	return cnt
}

func someCheck(args []string) {
	for i, el := range args {
		fmt.Println(i, ":", el)
	}
}
