# qrpc
A small distributed queue messaging service backed by grpc.

I've created  *the learning project* mostly because I wanted to play with [grcp](http://www.grpc.io/).
This is not super secure queue service, rather a grpc example implemented in go.

### Note
For the purpose of the project I also "stole" some ideas :)
QRPC stores data using Peter's awesome [diskv](https://github.com/peterbourgon/diskv) library. Also multi-master architecture is something what I really like in [disque](https://github.com/antirez/disque), so I've tried to implement similar mechanism.

