#!/bin/bash

mkdir -p \
  /home/jovyan/.executor

cp -R /tmp/repo/data/notebooks /home/jovyan/.executor/

# Launch Executor Daemon
ipython -m executor

start-notebook.sh \
  --NotebookApp.token='' \
  --NotebookApp.password='' \
  --NotebookApp.allow_origin='*' \
  --NotebookApp.base_url=${NB_PREFIX}
