.PHONY: server serverloop client test testsh tcpdump

port = 4243
host = localhost
telnet_timeout = 1s

server:
	nc -l -N $(host) $(port)

serverloop:
	while true;do echo "\n----- New server";nc -l -N $(host) $(port);done

client:
	go run . --timeout=$(telnet_timeout) $(host) $(port)

clientbreakstdout:
	for i in $$(seq 10);do sleep 1; printf "$$i\n";test "$$i" = "5" && break; done | go run . --timeout=$(telnet_timeout) $(host) $(port)

test:
	go test -v -count=1 -race -timeout=2m -run .

testsh:
	./test.sh

tcpdump:
	sudo tcpdump -i lo port $(port)
