# hyperdrive
A repository for code that runs hyperparameter experiments.

## Before You Begin (Requirements)
Your machine should have the following software installed:
1. Python >= 3.9
2. Docker
    -Get Docker Here: https://docs.docker.com/get-docker/

## Installation Instructions

### Using the CLI (local)
1. Start Docker Desktop

1. In a terminal/bash session, git clone the `hyperdrive` repository:
```bash
git clone git@github.com:gohypergiant/hyperdrive.git
```

1. In a terminal/bash session, in the `hyperdrive` repository from step 2, navigate to the `hyper` subdirectory. Here is a view of the `hyperdrive` directory structure:
```
hyperdrive
├──LICENSE
├──Makefile
├──README.md
├──datamanager
├──docker
├──executor
├──hyper        # NAVIGATE TO THIS DIRECTORY
├──hypertrain
```

1. In the `hyper` subdirectory, run this command to build the binary:
```bash
make build
```

1. If you want to put the binary in your path, run this:
```bash
  sudo mv hyper /usr/local/bin/hyper && \
  sudo chmod +x /usr/local/bin/hyper
```

1. In a terminal/bash session, create a project folder. Put your data into this project folder. Here is an example:
```
threat_detection
├──data
    ├──label_data.csv
    ├──object_data.json
```

1. Create a local jupyter server by running this command in a terminal/bash session:
```bash
hyper jupyter
```

1. When Step 6 is finished, you will see a message telling you which port Jupyter Lab is running on. Make note of the port number. The message will look something like this (with the exception of the port number):
```bash
Jupyter Lab Now Running via Docker Container 5433bddb4d on port 50491
```

1. Open a web browser (e.g., Chrome, Firefox, MS Edge)

1. In the address bar, navigate to the Jupyter Lab session by typing in the following into the address bar. We will be using the port number from Step 7:
```
localhost:[PORT_NUMBER_FROM_STEP_7]
```
