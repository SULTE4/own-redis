# own-redis

A lightweight Redis-like in-memory key-value store implemented in Go.  
This project is designed for learning purposes, focusing on networking, concurrency, and data storage principles.

---

---

## Abstract

In this project, implemented a custom key-value database. Such databases are crucial for efficiently storing and retrieving data using unique keys. This application operates entirely in memory for fast access, using the UDP protocol for communication between clients and the server.  

Supported operations:

- `SET` – store a key-value pair  
- `GET` – retrieve a value by key  
- `PING` – check server health  

---

## Context

Key-value stores are widely used when storing unstructured data without explicit relationships. Examples:

- Redis  
- Memcached  

They provide constant-time `O(1)` operations for storing and retrieving values.

---

## Project Status

- ✅ `SET` command implemented with optional `PX` expiry  
- ✅ `GET` command implemented with expiry check  
- ✅ `PING` implemented  

### Future Plans

1. Refactor code structure for easier reuse in other projects  
2. Improve synchronization and concurrency safety  
3. Write unit tests for all commands  

---

## Installation

Clone the repository:

```bash
git clone https://github.com/SULTE4/own-redis.git
cd own-redis
```

Build the project:

```bash
go build -o own-redis
```

Run the server (default port 8080):

```bash
./own-redis --port 8080
```

---

## Usage Examples

Use `nc` (netcat) or any UDP client:

### PING

```bash
$ nc -u 0.0.0.0 8080
PING
PONG
```

### SET

```bash
$ nc -u 0.0.0.0 8080
SET Foo Bar
OK

SET Foo Baz px 10000
OK
```

- `PX milliseconds` sets the expiration time for the key

### GET

```bash
$ nc -u 0.0.0.0 8080
GET Foo
Baz

# After expiry
GET Foo
(nil)
```

Invalid commands or wrong number of arguments return an error:

```bash
SET KEYVAL
(error) ERR wrong number of arguments for 'SET' command
```

---

## Project Structure

```
own-redis/
├── main.go        # Entry point, UDP server loop
├── redis.go       # Core Redis-like logic (SET, GET, PING)
├── go.mod         # Module configuration
└── README.md      # Documentation
```

