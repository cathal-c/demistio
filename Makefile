BIN=demistio
BIN_DIR?=bin
OUTPUT_DIR?=out
OUTPUT?=envoy_config.json

run:
	@mkdir -p ${OUTPUT_DIR}
	@ENABLE_DELIMITED_STATS_TAG_REGEX=false \
		PILOT_ENABLE_RDS_CACHE=false \
		PILOT_ENABLE_CDS_CACHE=false \
		go run main.go -output=${OUTPUT_DIR}/${OUTPUT} -debug

build:
	@mkdir -p ${BIN_DIR}
	go build -v -o ${BIN_DIR}/${BIN} main.go

test:
	@go test ./...

lint:
	@golangci-lint run
