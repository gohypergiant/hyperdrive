PROJECT = hyperdrive
SHELL := /bin/bash

clean:
	find . | grep -E "(__pycache__|\.pyc|\.pyo)" | xargs rm -rf
	rm -rf .executor
	rm -rf _jobs

image:
	docker build --target cpu-local -t ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-localstable -f docker/Dockerfile.main .

up:
	docker run -it --rm -p 8888:8888 --name hyperdrive-cpu-local cpu-local
