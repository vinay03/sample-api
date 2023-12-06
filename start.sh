#!/bin/bash
END=2
for i in $(seq 1 $END)
do
	bash -c "go run main.go 8081"
done

# ReplicaNumber=1 go run main.go 8081
# ReplicaNumber=2 go run main.go 8082
# ReplicaNumber=3 go run main.go 8083