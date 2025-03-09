#!/bin/bash
trap "rm server;kill 0" EXIT
pkill -f server 2>/dev/null

# 先编译 distribute_main.go
go build -o server distribute_main.go

# 运行多个实例
./server -port=8001 &
./server -port=8002 &
./server -port=8003 -api=1 &

sleep 2
echo ">>> start test"

# 发送测试请求
curl "http://localhost:10004/api?key=Tom" &
curl "http://localhost:10004/api?key=Tom" &
curl "http://localhost:10004/api?key=Tom" &

wait

