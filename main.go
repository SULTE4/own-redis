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

var SetError = "(error) ERR wrong number of arguments for 'SET' command"

var ExpireMinLimit int = 500

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

	buf := make([]byte, 2048)

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		responseText := handleMessage(*redis, buf[:n])
		if _, err := conn.WriteTo([]byte(responseText+"\n"), addr); err != nil {
			slog.Error(err.Error())
		}
		slog.Info("response handled")
	}
}

func handleMessage(redis ORedis, buf []byte) string {
	args := strings.Fields(string(buf))
	argCount := len(args)
	if argCount == 0 {
		return ""
	}

	command := strings.ToLower(args[0])
	switch command {
	case "ping":
		return "PONG"
	case "set":
		if argCount < 3 {
			return SetError
		}
		key := args[1]
		duration := time.Second * math.MaxInt32
		valueCount := argCount

		if argCount >= 5 {
			if args[argCount-2] == "px" {
				valueCount = argCount - 2
				tmp, err := strconv.Atoi(args[argCount-1])
				if err != nil || tmp <= ExpireMinLimit {
					return "Invalid time duration"
				}
				duration = time.Millisecond * time.Duration(tmp)
				slog.Info("duration integrating", slog.String("time (ms):", duration.String()))
			}
		}

		var value strings.Builder
		for i, el := range args[2:valueCount] {
			value.WriteString(el)
			if argCount-1 != i {
				value.WriteString(" ")
			}
		}

		nowTime := time.Now()
		expireTime := nowTime.Add(duration)

		v := Value{data: value.String(), expireTime: expireTime}

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
		currentTime := time.Now()
		if value.expireTime.Before(currentTime) {
			slog.Info("key time expired")
			return "(nil)"
		}
		return value.data
	default:
		return "Invalid operation"
	}
}
