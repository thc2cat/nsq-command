
# NSQ cli producer and consumer for spreading data to multiple hosts commands

Purpose of this project is to build a tool able to send data to multiple hosts (in a serialized way or multicast way), and execute a local command with received data using [NSQ](https://nsq.io/) realtime distributed messaging platform.

## nsq server

Use docker for a rapid setup.

```shell
# cd contrib
# docker compose -d up 
# # will launch nsqlookupd, nsqd, and nsqadmin
```

Now you should reflet your __NSQ_SERVER__ environment variable to YourNSQServer:4150

Verify the web interface of nsq instance [http://YourNSQServer:4171](http://YourNSQServer:4171)

## Build producer and consumer

```shell
# make ( will do a go build )
```

## Test

```shell
$ # launch consumer 
$ ./C.exe
$ # launch producer
$ ./P.exe -s "this host" -m "just a test"
$ # default output text has been received from consumer 
NSQ stamp[10:47:32], from[this host], data[just a test]
$ # you can use command mode where %s is replace with text
$ ./C.exe -c "echo %s"
$ # output => just a test
```

So now, you can spread data to multiple hosts.

* use topic to separate usage
* use channel to differentiate data usage
  * same channel mean one consumer out of many will have data
  * different channel mean each consumer will have a copy

### Exemple

#### Spreading ip address to group of hosts for blocking ip address

* ./P.exe -C "fw_deny_ip" -m "192.168.99.0/24"
* on Host1 (using ipfw) $ ./C.exe -C "fw_deny_ip" -cmd "ipfw table 1 add %s"
* on Host2 (using ipset) $ ./C.exe -C "fw_deny_ip" -cmd "ipset add blacklist %s timeout 0"

#### Additional tips

Verbose mode (-v) show details about server, topic, channel and nsq communication.

nsqd options --max-msg-size can be used to raise maximum size of message.
