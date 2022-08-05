# hyperpackage

A package to create hyperpacks.

## Before You Begin (Requirements)

Your machine should have the following software installed:

1. Python >= 3.9
2. Docker
   -Get Docker Here: https://docs.docker.com/get-docker/
3. Hyper CLI
   -See installation instructions below

## Installation Instructions - Hyper CLI

1. Start Docker Desktop

2. In a terminal/bash session, git clone the `hyperdrive` repository:

```bash
git clone git@github.com:gohypergiant/hyperdrive.git
```

3. In a terminal/bash session, in the `hyperdrive` repository from step 2, navigate to the `hyper` subdirectory. Here is a view of the `hyperdrive` directory structure:

```
hyperdrive
├──dataclient
├──docker
├──executor
├──hyper        # NAVIGATE TO THIS DIRECTORY
├──hyperpackage
├──hypertrain
```

4. In the `hyper` subdirectory, run this command to build the binary:

```bash
make build
```

5. Put the binary in your path by running the following command:

```bash
  chmod +x hyper && sudo mv hyper /usr/local/bin/hyper
```

## Installation Instructions - hyperpackage

1. If you're working in a JupyterLab or Jupyter notebook, run this magic in a cell:
```bash
%pip install -e git+https://github.com/gohypergiant/hyperdrive.git#egg=hyperpackage\&subdirectory=hyperpackage
```

2. If you'd prefer to install from a terminal/bash session, please run this:
```bash
pip install -e git+https://github.com/gohypergiant/hyperdrive.git#egg=hyperpackage\&subdirectory=hyperpackage
```
NOTE: if you install via terminal/bash session, a "src" folder with the hyperpackage source code will be created in the directory from which you ran the pip install command.

## Usage

1. In a JupyterLab/Notebook or python session, import the "create_hyperpack" function:

```
from hyperpackage.hyperpack_creation import create_hyperpack
```

2. Currently, the only supported model flavor is automl, aka the "research-EfficientAutoML" library. The "create_hyperpack" function can accept either the model object itself, or a string path to a model saved via the torch.save() function. We'll go over both situations, starting with the model object in memory.

3. Assuming that you've run an automl study and have obtained outputs via this command:

```
output = model.fit(x=features, y=target)
```

You can then call the "create_hyperpackage" function by passing in the pretrained automl model from the output object (specifically, output["model"]), like so:

```
create_hyperpack(trained_model=output["model"], model_flavor="automl")
```
 
4. To call "create_hyperpackage" with a string path to a saved model, e.g., "/Users/hanswilsdorf/saved_model":

```
create_hyperpack(trained_model="/Users/hanswilsdorf/saved_model", model_flavor="automl")
```

5. Successful execution of the "create_hyperpack" function will create the following artifacts IN THE CURRENT DIRECTORY from which you called the "create_hyperpack" function:

```
automl.hyperpack
automl.hyperpack.zip
study.yaml
```

6. To run the prediction server, in a terminal/bash session, run the following Hyper CLI command in the same directory that the automl.hyperpack.zip and study.yaml artifacts are located:
``` bash
hyper pack run
```

NOTE: If you want to run the prediction server from a different directory, you'll need to move BOTH the automl.hyperpack.zip AND study.yaml files.

7. From the printed output of the "hyper pack run" command, please make note of the following 2 items:

Server port number: look for the message "Hyperpackage now running via Docker Container 9672b19c0b on port [XXXXX]". You'll need the port number [XXXXX].

Fast App API key: look for the message "Fast App API key is: [FAST_API_KEY]"

8. Retrieve a prediction by running the following command in a terminal/bash session. The input data that you pass with the -d flag should be an array of values that is an appropriate shape (i.e., you should be passing 5 values if your model is expecting 5 values). You'll need both the server port number and Fast App API key from the previous step:

``` bash
curl -X 'POST' \
  http://127.0.0.1:[SERVER_PORT_NUMBER_HERE]/predict \
  -H 'x-api-key: [FAST_APP_API_KEY_HERE]' \
   -d '[ARRAY_OF_VALUES_HERE]'
```
