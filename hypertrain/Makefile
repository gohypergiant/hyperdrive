PROJECT = hypertrain
PYTHON_VERSION = 3.9
DEV_PACKAGES = black flake8 ipdb isort pre-commit pytest pytest-cov
SHELL := /bin/bash

PYTEST_ARGS = -x -p no:warnings
PYTEST_COVERAGE = --cov-report term-missing --cov=${PROJECT}
PYTEST_DEBUG = -s
PYTEST_FOCUS = -k focus
SHELL=/bin/bash

DOCKER_COMPOSE_COMMAND = docker compose -f docker/docker-compose.yml

DEV_CONTAINER = $(shell docker ps --format {{.Names}} | grep hypertrain-dev)

build:
# add --no-cache to trigger full rebuild
	${DOCKER_COMPOSE_COMMAND} build --build-arg HYPERDRIVE_REGISTRY_TOKEN=${HYPERDRIVE_REGISTRY_TOKEN}

clean:
	find . | grep -E "(__pycache__|\.pyc|\.pyo|\.egg-info)" | xargs rm -rf

down:
	${DOCKER_COMPOSE_COMMAND} down -v

install:
	python setup.py install

lint:
	pre-commit run --all-files

package:
	python setup.py bdist_wheel

ps:
	${DOCKER_COMPOSE_COMMAND} ps

test:
	docker exec --workdir /home/jovyan/mlsdk-hypertrain -it ${DEV_CONTAINER} pytest ${PYTEST_COVERAGE} ${PYTEST_ARGS}

testdebug:
	docker exec --workdir /home/jovyan/mlsdk-hypertrain -it ${DEV_CONTAINER} pytest ${PYTEST_DEBUG} ${PYTEST_ARGS}

testfocus:
	pytest ${PYTEST_DEBUG} ${PYTEST_ARGS} ${PYTEST_FOCUS}

up:
	${DOCKER_COMPOSE_COMMAND} up -d

wheel:
	python setup.py bdist_wheel
