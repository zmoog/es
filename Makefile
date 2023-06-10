NOW := $(shell date +%FT%T%z)

tidy:
	go mod tidy

run:
	echo "{\"@timestamp\": \"$(NOW)\", \"name\": \"Maurizio Branca\"}" | go run main.go docs index -

