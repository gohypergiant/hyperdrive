# Hyperdrive CLI

## Installation Instructions

### Installing From Source

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
├──dataclient
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


### Usage

#### Prerequisites

1.  Make sure docker is running

#### Using Hyper with Hypertrain

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

## Commands

### `hyper remoteStatus` : start a status endpoint for polling

Summons the status endpoint, if the specified port is unavailable the command will fail. Ping `localhost:3001/status` to receive the status in JSON, example:

```json
{
  "message": "Pulling docker image"
}
```

If the `message` is an empty string, or the statusFile does not exist, The response will be a [204](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/204).

**Optional flags:**

- `--port`: _(Default: `3001`)_ Specify the custom port to launch the endpoint with
- `--statusFile`: _(Default: `./statusfile.json`)_ Specify the location of the status file.

### `hyper remoteStatus update "<message>"` : Update the status file

Updates the status file with the supplied `message` in JSON format.

**Optional flags:**

- `--statusFile`: _(Default: `./statusfile.json`)_ Specify the location of the status file.

## Remote

To use with a remote Firefly target, add a `.hyperdrive` file to your `$HOME` directory with contents that look like this:
```json
{
   "remotes":{
      "dev":{
         "url":"http://localhost:8090/hub/api",
         "username":"firefly",
         "hub_token":"TOKEN"
      }
   }
}
```

You should then be able to use the CLI with the `--remote=<REMOTE NAME>` flag. The above example will provide a remote named `dev`


The `test-data` directory in this repo contains a `my_study.yaml` file which is the firefly manifest. From that directory, you should be able to run
```bash
# If a manifest path isn't provided, it will look for a file name study.yaml by default
> hyper jupyter start --remote=dev --manifestPath=./my_study.yaml
```
Which will start a notebook server instance on the remote using the `study_name` field in the manifest (`log_reg_health_tracker`)

Once you have a running notebook server instance, you should be able to start a training session on it by
```bash
# Currently all this does is upload the manifest and data to the _jobs folder
> hyper train --remote=dev --manifestPath=./my_study.yaml
```

## Local

To use a local jupyter notebook server, first create the server

```bash
# If a manifest path isn't provided, it will look for a file name study.yaml by default
> hyper jupyter
```

Once you have a running notebook server instance, you should be able to start a training session on it by
```bash
> hyper train --manifestPath=./my_study.yaml
```

The study will be scheduled to be executed on the server, to fetch the hyperpackage from the training session run
```bash
> hyper train fetch --manifestPath=./my_study.yaml
```


Note: To use a local Firefly server for training, it is necessary to create the notebook server instance and execute the traning session from within the same git project.

