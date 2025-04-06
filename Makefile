OUTPUT_DIR?=out
OUTPUT?=envoy_config.json

run:
	@mkdir -p ${OUTPUT_DIR}
	 PILOT_ENABLE_RDS_CACHE=false go run main.go -output=${OUTPUT_DIR}/${OUTPUT}
