#!/bin/sh

ps -ef|grep /usr/bin/fluentd|grep -v grep

if [ $? -eq 0 ]
then
  echo "fluentd is good!"
else
  echo "fluentd is bad!"
  exit 1
fi