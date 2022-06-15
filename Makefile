PROJECT = hyperdrive
SHELL := /bin/bash

clean:
	find . | grep -E "(__pycache__|\.pyc|\.pyo)" | xargs rm -rf
	rm -rf .executor
	rm -rf _jobs

image:
	docker build --target cpu-local -t ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-localstable -f docker/Dockerfile.main .

image-base:
	./docker/gpu/create-gpu-dockerfile.sh
	docker build --target gpu-base -t ghcr.io/gohypergiant/hyperdrive-jupyter:gpu-chungus -f docker/Dockerfile.base .
	./docker/gpu/cleanup-gpu-dockerfile.sh

up:
	docker run -it --rm -p 8888:8888 --name hyperdrive-cpu-local cpu-local
