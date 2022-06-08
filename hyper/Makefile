BINARY_NAME = hyper-binary-$(shell git rev-parse --abbrev-ref HEAD)

build:
	go build -o hyper

build-dev:
	go build -o ${BINARY_NAME}
	chmod +x ${BINARY_NAME}	
	sudo rm /usr/local/bin/hyper-dev
	sudo ln ${BINARY_NAME} /usr/local/bin/hyper-dev
	echo "${BINARY_NAME}"

path:
	chmod +x hyper && sudo mv hyper /usr/local/bin/hyper
