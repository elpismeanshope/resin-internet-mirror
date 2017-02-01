all: docker

Dockerfile: Dockerfile.template
	echo "# AUTOGENERATED BY MAKE - DO NOT MODIFY MANUALLY" > $@
	sed 's/%%RESIN_MACHINE_NAME%%/amd64/;s/kiwix-server-arm/kiwix-linux-x86_64/' \
		$< >> $@

docker: Dockerfile
	docker build  -t resin-internet-mirror .

.PHONY: all docker
