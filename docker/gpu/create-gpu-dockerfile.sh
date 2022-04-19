#!/bin/bash
cat docker/gpu/Dockerfile.gpu.base > docker/Dockerfile.base
curl https://raw.githubusercontent.com/jupyter/docker-stacks/master/base-notebook/Dockerfile | grep -v ROOT_CONTAINER >> docker/Dockerfile.base
curl https://raw.githubusercontent.com/jupyter/docker-stacks/master/minimal-notebook/Dockerfile | grep -v BASE_CONTAINER >> docker/Dockerfile.base
curl https://raw.githubusercontent.com/jupyter/docker-stacks/master/scipy-notebook/Dockerfile | grep -v BASE_CONTAINER >> docker/Dockerfile.base
sed 's/cpu/gpu/g' docker/Dockerfile.base | grep -v jupyter/scipy-notebook >> docker/Dockerfile.base
perl -pi -e 's/tini \\//g' docker/Dockerfile.base
curl -O https://raw.githubusercontent.com/jupyter/docker-stacks/master/base-notebook/fix-permissions
curl -O https://raw.githubusercontent.com/jupyter/docker-stacks/master/base-notebook/jupyter_notebook_config.py
curl -O https://raw.githubusercontent.com/jupyter/docker-stacks/master/base-notebook/jupyter_server_config.py
curl -O https://raw.githubusercontent.com/jupyter/docker-stacks/master/base-notebook/start-notebook.sh
curl -O https://raw.githubusercontent.com/jupyter/docker-stacks/master/base-notebook/start-singleuser.sh
curl -O https://raw.githubusercontent.com/jupyter/docker-stacks/master/base-notebook/start.sh
