# hyperdrive
A repository for code that runs hyperparameter experiments.

## Before You Begin (Requirements)
Your machine should have the following software installed:
1. Python >= 3.9
2. Docker
    -Get Docker Here: https://docs.docker.com/get-docker/

##  Contributing

See CONTRIBUTORS.md for details on requirements.

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
hyper jupyter --browser
```

1. Train the Model
```bash
hyper train
```

1. Retrieve the hyperpack
```bash
cp ../_jobs/threat_detection/threat_detection.hyperpack.zip .
```

1. Run the prediction server
```bash
hyper hyperpackage run
```

1. Submit a prediction
```
curl -X 'POST' \
  http://127.0.0.1:PORT_OF_YOUR_SERVER_HERE/predict \
  -H 'accept: application/json' \
   -d '[ 22.6734323 , 133.87953978,  71.25881828,  24.74618134]'
```
