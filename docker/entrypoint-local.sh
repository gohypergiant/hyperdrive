#!/bin/bash

# mkdir -p \
#   /home/jovyan/.executor \
#   /home/jovyan/_jobs \
#   /home/jovyan/_jobs/auto-ml-demo/data \
#   /home/jovyan/_jobs/low-code-demo/data \
#   /home/jovyan/_jobs/executor-test/data

# cp -R /tmp/repo/data/notebooks /home/jovyan/.executor/notebooks
# cp -R /tmp/repo/data/sample-data/. /home/jovyan/_jobs/auto-ml-demo/data
# cp -R /tmp/repo/data/yaml/auto-ml.yaml /home/jovyan/_jobs/auto-ml-demo/_study.yaml
# cp -R /tmp/repo/data/sample-data/. /home/jovyan/_jobs/low-code-demo/data
# cp -R /tmp/repo/data/yaml/low-code.yaml /home/jovyan/_jobs/low-code-demo/_study.yaml
# cp -R /tmp/repo/data/sample-data/. /home/jovyan/_jobs/executor-test/data
# cp -R /tmp/repo/data/yaml/test.yaml /home/jovyan/_jobs/executor-test/_study.yaml

# touch /home/jovyan/_jobs/low-code-demo/COMPLETED
# touch /home/jovyan/_jobs/auto-ml-demo/COMPLETED

# # Launch Executor Daemon
# ipython -m executor

start-notebook.sh \
  --NotebookApp.token='' \
  --NotebookApp.password='' \
  --NotebookApp.allow_origin='*' \
  --NotebookApp.base_url=${NB_PREFIX}
