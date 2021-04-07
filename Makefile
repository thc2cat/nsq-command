
all: producer consumer

producer: Producer/*.go
	go build Producer/producer.go 

consumer: Consumer/*.go
	go build Consumer/consumer.go 

clean:
	rm Consumer/Consumer.exe Producer/Producer.exe consumer producer 2>/dev/null
