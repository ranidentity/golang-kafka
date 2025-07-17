go get github.com/confluentinc/confluent-kafka-go/kafka

[guid]::NewGuid().ToString()
docker-compose run --rm kafka /opt/kafka/bin/kafka-storage.sh format -t 37857e05-4df2-41b3-a397-09bd8df0cde6 -c /opt/kafka/config/server.properties
