.PHONY: build privacy-scan
build:
	./scripts/build.sh

privacy-scan:
	@! grep -RInE --exclude-dir=.git --exclude='go.sum' 'LPA:1\$$[^[:space:]]+\$$[^[:space:]]+|[0-9]{30,40}|IMSI[[:space:]]*[:=][[:space:]]*[0-9]{14,16}' .
