#!/bin/bash

mkdir -p \
  /home/jovyan/.executor

cp -R /tmp/repo/data/notebooks /home/jovyan/.executor/notebooks

start-notebook.sh \
  --NotebookApp.token=${NB_API_KEY} \
  --NotebookApp.password=${NB_PASSWORD} \
  --NotebookApp.allow_origin='*' \
  --NotebookApp.base_url=${NB_PREFIX}
