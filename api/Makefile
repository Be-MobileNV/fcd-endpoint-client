# Child Makefile
# Always pass OUTPUT_DIR=../output to your targets

OPENAPI_OUTPUT_DIR = $(OUTPUT_DIR)/openapi
$(shell mkdir -p $(OPENAPI_OUTPUT_DIR))

version:
	cat openapi.yaml | docker run -i karlkfi/yq -r '.info.version'

lint:
	docker run --rm -i -v $(CURDIR):/tmp stoplight/spectral lint "/tmp/openapi.yaml" --ruleset=/tmp/.spectral.yaml

package:
	tar -czvf $(OPENAPI_OUTPUT_DIR)/openapi.tar.gz openapi.yaml

redoc-serve:
	redoc-cli serve openapi.yaml -p 8060