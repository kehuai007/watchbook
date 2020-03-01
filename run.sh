#!/bin/sh
name='watchbook'
process=`ps -ef | grep $name | grep -v grep | grep -v PPID|awk '{ print $2}'`
for i in $process
	do
	echo "kill the $name process $i"
	kill -9 $i
	done
rm -f nohup.out
nohup ./$name &
pidof $name
