# qrpc
A small distributed queue messaging service backed by grpc.

I've created the project mainly as an exercise because I wanted to play with [grcp](http://www.grpc.io/).
This is not super secure queue service, rather a prototype implemented in go.

### Note
For the purpose of the project I also "stole" some ideas :)
QRPC stores data using Peter's awesome [diskv](https://github.com/peterbourgon/diskv) library. Also multi-master architecture is something what I really like in [disque](https://github.com/antirez/disque), so I've tried to implement similar mechanism.

### Install & Run
```sh
$ go get github.com/kuba--/qrpc/qrpc
$ qrpc --help
usage: ./qrpc [flags] [cluster peer(s) ...]
flags:
  -port int
  	port to listen on (default 9033)
  -data string
  	directory in which queue data is stored (default "/tmp/qrpc")
  -cache uint
  	max. cache size (bytes) before an item is evicted. (default 1048576)
  -interval duration
  	cluster watch timer interval. (default 1s)
  -timeout duration
  	cluster request (gossip) timeout. (default 3s)

# start standalone server :9091
$ qrpc -data /tmp/qrpc-9091 -port 9091

# ...join the cluster
$ qrpc -data /tmp/qrpc-9092 -port 9092 127.0.0.1:9091
$ qrpc -data /tmp/qrpc-9093 -port 9093 127.0.0.1:9092
```
