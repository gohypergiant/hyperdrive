# hyperdrive

A repository for code that runs hyperparameter experiments.

## Before You Begin (Requirements)

Your machine should have the following software installed:

1. Python >= 3.9
2. Docker
   -Get Docker Here: https://docs.docker.com/get-docker/

## Contributing

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
