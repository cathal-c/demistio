OUTPUT_DIR?=out
OUTPUT?=envoy_config.json

run:
	@mkdir -p ${OUTPUT_DIR}
	go run main.go -output=${OUTPUT_DIR}/${OUTPUT}
