#!/bin/sh
echo "Now we are going to list files"
ls
echo "Now we are going to print the directory we are in"
pwd
echo "Now we will print out a list of processes that are taking up too much space!"
ps -eo pid,ppid,pri,ni,cmd,%mem,%cpu --sort=-%mem | head -15

