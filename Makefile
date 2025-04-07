OUTPUT_DIR?=out
OUTPUT?=envoy_config.json

run:
	@mkdir -p ${OUTPUT_DIR}
	@ENABLE_DELIMITED_STATS_TAG_REGEX=false \
		PILOT_ENABLE_RDS_CACHE=false \
		go run main.go -output=${OUTPUT_DIR}/${OUTPUT}

test:
	@go test ./...
