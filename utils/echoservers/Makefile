run:
	@go run main.go
certs:
	@rm -rf $(CURDIR)/echoserver.key $(CURDIR)/echoserver.crt
	# generate the key
	@openssl genrsa \
		-out cert.key \
		8192
	# generate the certificate
	@openssl req \
		-subj "/C=SG/ST=SG/L=Kepler12b/O=JediInc/OU=CouncilOfElders/CN=localhost" \
		-new \
		-x509 \
		-sha256 \
		-key cert.key \
		-days 1 \
		-out cert.crt
