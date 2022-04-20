PROJECT = hyperdrive
SHELL := /bin/bash

clean:
	find . | grep -E "(__pycache__|\.pyc|\.pyo)" | xargs rm -rf

image:
	docker build --target cpu-local -t cpu-local -f docker/Dockerfile.main .

up:
	docker run -p 8888:8888 cpu-local
