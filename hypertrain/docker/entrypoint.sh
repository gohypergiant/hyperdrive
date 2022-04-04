# Instantiate IPython Profiles
ipython profile create
ipython profile create daemon
INITIALIZATION_FILE=~/.ipython/profile_default/startup/initialization.py
INITIALIZATION_FILE_DAEMON=~/.ipython/profile_daemon/startup/initialization.py

# Create Base Default Profile
echo "import os" > $INITIALIZATION_FILE
echo "os.chdir('/home/jovyan')" >> $INITIALIZATION_FILE

# Copy Base Default Profile to Daemon Profile
cp $INITIALIZATION_FILE $INITIALIZATION_FILE_DAEMON

# Store Git Credentials
git config --global credential.helper store

# Prepare Dev Server
mv /tmp/mlsdk-hypertrain /home/jovyan/
rm -r /home/jovyan/work
rm /home/jovyan/requirements.txt

# Launch Notebook Server
start-notebook.sh --NotebookApp.token='' --NotebookApp.password='' --NotebookApp.allow_origin='*' --NotebookApp.base_url=${NB_PREFIX}
