# KVS - Redis-like key-value storage in Go

Available commands:
- SET <key> <value>
- GET <key>
- DELETE <key>

Can be used with `redis-cli` client

## Starting KVS
To run kvs, you need go v1.25.4 installed
To start KVS, write following in command line in repository root

```
go run .
```

Also you can specify port to run your app on. Default port is 8080

```
go run . -port=7123
```
