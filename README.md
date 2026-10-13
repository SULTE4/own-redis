<div align="center">

# own-redis

**`SET Foo Bar px 10000` over UDP. `GET Foo`. Ten seconds later: `(nil)`.**

A tiny Redis-like in-memory key-value store in Go, built to learn networking, concurrency<br>
and storage basics. One file, no dependencies, talks plain text over UDP.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![Protocol](https://img.shields.io/badge/protocol-UDP-orange)
![Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen)
![Lines](https://img.shields.io/badge/main.go-~130%20lines-informational)

</div>

---

## Why

Redis and Memcached look like magic until you write a small one yourself. At its core, a key-value store is a map, a lock and a network loop, plus expiry. This project is that core and nothing else, so the whole thing fits on one screen.

## How it works

```
 nc -u localhost 8080
        │  "SET Foo Bar px 10000"
        ▼
 ┌─────────────────────────────┐
 │ UDP loop (net.ListenPacket) │  one datagram = one command
 └──────────────┬──────────────┘
                ▼
   parse: strings.Fields, command is case-insensitive
                │
                ▼
   map[string]{data, expireTime}   writes guarded by sync.RWMutex
                │
                ▼
   reply datagram: OK │ value │ (nil) │ PONG │ error
```

Expiry is lazy: an expired key is reported as `(nil)` when you read it.

## Quick start

```bash
git clone https://github.com/SULTE4/own-redis.git
cd own-redis
go build -o own-redis
./own-redis --port 8080
```

In another terminal:

```bash
$ nc -u 127.0.0.1 8080
PING
PONG
SET Foo Bar
OK
GET Foo
Bar
SET Foo Baz px 10000
OK
GET Foo
Baz
```

Ten seconds later:

```bash
GET Foo
(nil)
```

## Commands

| Command | Reply | Notes |
|---|---|---|
| `PING` | `PONG` | Health check |
| `SET key value...` | `OK` | Values can contain spaces: `SET greeting hello world` |
| `SET key value px <ms>` | `OK` | Expires after `<ms>` milliseconds. Must be more than 500 |
| `GET key` | value or `(nil)` | `(nil)` for missing or expired keys |

| Mistake | Reply |
|---|---|
| `SET KEYVAL` | `(error) ERR wrong number of arguments for 'SET' command` |
| `SET k v px 100` | `Invalid time duration` |
| `DEL Foo` (anything else) | `Invalid operation` |

## Options

| Flag | Default | Meaning |
|---|---|---|
| `--port` | `8080` | UDP port to listen on |

## Roadmap

- [ ] `DEL`, `EXISTS`, `TTL` and an active expiry sweep
- [ ] Lock reads too, and pass the store by pointer
- [ ] Trim the trailing space from returned values
- [ ] Unit tests for every command
- [ ] Split `main.go` into a reusable package
