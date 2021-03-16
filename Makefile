
all: producer consumer

producer: P/*.go
	go build P/producer.go 

consumer: C/*.go
	go build C/consumer.go 
