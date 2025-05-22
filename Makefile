.PHONY: openapi
openapi:
	@./scripts/generate_openapi.sh

.PHONY: protobuf
protobuf:
	@./scripts/generate_protobuf.sh
