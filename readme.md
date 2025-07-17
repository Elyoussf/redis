<!-- Redis server compatible with redis

You can use redis clients for test 

Folder : RESP a complete Serializer and deserializer of Resp protocol which is a binary safe protocol and length prefixed (you can use any character with in it , starts with the full length not a delimiter based protocol )

High level overview of resp :
redis protocol based on request response schema there is a client that sends commands to the server these commands are piped as array with the first element the command and teh rest is args 

Most used (implemented here) commands :
SET , GET , HSET, HGET , HGETALL

Since resp is a length prefixed protocol any message from client start with *[l] with l is a number represents how many bulks or whitespace separated strings are written in fact it's not whitespace but completely new like but it's added automatically you as a user you just write inline so the delimiter is CRLF : \r\n

for example GET name hamza 
is translated to
*3
$3
GET
$4
name
$5
hamza

* means array
$ means bulk (string but not really a string )
+ means the real string 

So on 

writer folder is where the send mechanism is abstracted with a struct and bufio 

The whole protocol is based on TCP so its generally byte-stream based protocol 

commands handlers is where the five commands mentioned above are written ping , set , get .... and saved within a map to access them on the fly 

server folder contains the whole server : TCP linstener , handlers caller etc 

TEST :
pre-requisite : install a redis client in normal cases all of them should work but i tested just the debian based one so far but i guess they can work all 


Open two terminals :

in the first terminal run :
redis-cli ping 
redis-cli set <key> <value>
redis-cli hset <hash> <key> <value>

in the second teminal just run :

go run main.go 



no data persistence implemented yet :
i don;t want to use AOF like the real redis any idea ?? -->
# Redis (From Scratch)

A minimal Redis-compatible server written in Go.

Implements core Redis commands using the RESP (REdis Serialization Protocol).  
You can use any standard Redis client (like `redis-cli`) to test and interact with it.

---

## ✨ Features

- ✅ RESP protocol (serializer/deserializer) implemented from scratch
- ✅ TCP-based byte stream handling
- ✅ Command support:
  - `PING`
  - `SET`
  - `GET`
  - `HSET`
  - `HGET`
  - `HGETALL`
- ✅ Works with real Redis clients
- 🚫 No persistence yet (AOF not implemented)

---

## 📁 Folder Structure

```bash
.
├── main.go         # Entry point for server
├── server/         # TCP server setup and connection handling
├── resp/           # RESP protocol (serialization/deserialization)
├── commands/       # Implementation of supported commands
├── writer/         # Response writer abstraction using bufio

## 🔌 RESP Protocol Overview

RESP is a **binary-safe**, **length-prefixed** protocol (not delimiter-based).

### Key Concepts:

- `*` denotes an **array**
- `$` denotes a **bulk string**
- `+` denotes a **simple string**

For example, the command:

GET name hamza

Is translated into RESP as:

*3\r\n
$3\r\n
GET\r\n
$4\r\n
name\r\n
$5\r\n
hamza\r\n


This means:

- `*3` → 3 elements in the array
- `$3 GET` → the command
- `$4 name`, `$5 hamza` → the arguments

Note: While the protocol is CRLF-based (`\r\n`), clients like `redis-cli` abstract that away—you can write inline commands normally.

---

## 🧪 Testing the Server

### Prerequisite

Install a Redis client (`redis-cli`).  
Any client should work, but this was tested with the Debian package version.

### Run Instructions

Open two terminals:

**Terminal 1 – Start the server**

```bash
go run main.go
Terminal 2 – Test commands with redis-cli

redis-cli ping
redis-cli set mykey myvalue
redis-cli get mykey
redis-cli hset myhash field1 value1
redis-cli hget myhash field1
redis-cli hgetall myhash
🧠 How It Works

    resp/: Parses and constructs RESP messages (serializer and deserializer).

    commands/: Implements supported commands and stores them in a command map.

    writer/: Wraps bufio to abstract the TCP response sending mechanism.

    server/: Bootstraps the TCP listener and routes requests to handlers.

🚧 Data Persistence

Persistence is not implemented yet (AOF is not used like in real Redis).
Looking for a lightweight alternative—open to ideas and contributions!