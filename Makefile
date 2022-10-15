default: 
	./scripts/make.sh		

# Server Development makefile
run-dev:
	./scripts/make.sh
	./kv-middleware --log_level=debug --config_path=deployment/config/development.yaml

build-binary:
	./scripts/make.sh

# Build image with tag latest
build-image:
	./scripts/make.sh build-image

# Infra Development makefile
compose-up:
	./scripts/make.sh compose-up
	# sleep to wait docker to start
	sleep 10 
	soda create -a
	./scripts/make.sh migrate-up
	
compose-down:
	./scripts/make.sh migrate-down
	./scripts/make.sh compose-down