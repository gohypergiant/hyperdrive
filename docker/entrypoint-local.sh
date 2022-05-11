#!/bin/bash

# Launch Executor Daemon
ipython -m executor

start-notebook.sh \
  --NotebookApp.token='' \
  --NotebookApp.password='' \
  --NotebookApp.allow_origin='*' \
  --NotebookApp.base_url=${NB_PREFIX}
