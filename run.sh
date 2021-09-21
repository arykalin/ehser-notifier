#!/bin/bash
while :
do
 echo "running notifier" `date`
 ./esher-notifier &>> esher-notifier.log
 sleep 300
done
