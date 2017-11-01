#!/bin/bash
counter=0
while [ $counter -le 7 ]
do
nohup ./lib/python/prediction_server "500$counter" > ./data/logs_go/prediction_server.log 2>&1 &
((counter++))
done
echo All servers are up
# nohup /usr/local/go/bin/go run lib/go/apps/balancer/balance.go -bind 0.0.0.0:8081 -balance "0.0.0.0:5000,0.0.0.0:5001,0.0.0.0:5002,0.0.0.0:5003,0.0.0.0:5004,0.0.0.0:5005,0.0.0.0:5006,0.0.0.0:5007" > /dev/null 2>&1 &
