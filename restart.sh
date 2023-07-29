#!/bin/bash

echo "Usage: "
echo "  [fishing-gateway script]"
echo ""
echo "Available Commands:"
echo "  restart"
echo "  stop"
echo ""
echo "Input Commands: "


RestartServer() {
    git pull
    go build -o fishing-gateway
    supervisorctl restart fishing-gateway
}

StopServer() {
    supervisorctl stop fishing-gateway
}

while :
do
  read cmd
  case $cmd in
    restart) RestartServer
      echo "重启服务完成;"
      break
    ;;
    stop) StopServer
      echo "关闭服务完成;"
      break
    ;;
    *) echo "cmd not found Re-enter;"
      continue
    ;;
  esac
done