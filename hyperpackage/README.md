# hyperpackage

A package to create hyperpacks.

A hyperpack is a machine learning, model deployment-ready artifact (i.e., a zipped/bundled set of files) that can be used with Hyperdrive to serve real-time predictions.

## Before You Begin (Requirements)

Your machine should have the following software installed:

1. Python >= 3.9
2. Docker
   - Get Docker Here: https://docs.docker.com/get-docker/
3. Hyper CLI
   - [See installation instructions](../hyper/README.md#installation)


## Installation Instructions - hyperpackage

1. If you're working in a JupyterLab or Jupyter notebook, run this magic in a cell:
```bash
    %pip install git+https://github.com/gohypergiant/hyperdrive.git#egg=hyperpackage\&subdirectory=hyperpackage
```

2. OR if you'd prefer to install from a terminal/bash session, you can run this:
```bash
    pip install git+https://github.com/gohypergiant/hyperdrive.git#egg=hyperpackage\&subdirectory=hyperpackage
```

## Usage

1. In a JupyterLab/Notebook or python session, import the `create_hyperpack` function:

```
    from hyperpackage.hyperpack_creation import create_hyperpack
```

2. Currently, the only supported model flavor is automl, aka the [research-EfficientAutoML](https://github.com/gohypergiant/research-EfficientAutoML) library. The `create_hyperpack` function can accept either the automl model object itself, or a string path to a model saved via the `torch.save()` function from `pytorch`. We'll go over both situations, starting with the model object in memory.

3. Model object in memory - assuming that you've run an automl study and have obtained outputs via this command:

```
    output = model.fit(x=features, y=target)
```

    You can then call the `create_hyperpackage` function by passing in  
    the pretrained automl model from the output object (specifically,  
    output["model"]), like so:

```
    create_hyperpack(trained_model=output["model"], model_flavor="automl", ml_task="binary_classification")
```
NOTE: You MUST pass in `ml_task`, which refers to the machine learning task type of the model. Available options are "regression", "binary_classification", and "multi_class_classification".

4. String path to model - to call `create_hyperpackage` with a string path to a saved model, e.g., "/Users/hanswilsdorf/saved_model":

```
    create_hyperpack(trained_model="/Users/hanswilsdorf/saved_model", model_flavor="automl", ml_task="binary_classification")
```
NOTE: You MUST pass in `ml_task`, which refers to the machine learning task type of the model. Available options are "regression", "binary_classification", and "multi_class_classification".

5. Successful execution of the `create_hyperpack` function will create the following artifacts IN THE CURRENT DIRECTORY from which you called the `create_hyperpack` function:

```
    automl.hyperpack
    automl.hyperpack.zip
    study.yaml
```

6. Before running the prediction server, start Docker Desktop.

7. To run the prediction server, in a terminal/bash session, run the following Hyper CLI command in the same directory that the `automl.hyperpack.zip` and `study.yaml` artifacts are located:

``` bash
    hyper pack run
```

    NOTE: If you want to run the prediction server from a different  
    directory, you'll need to move BOTH the automl.hyperpack.zip  
    AND study.yaml files into that directory.

8. From the printed output of the `hyper pack run` command, please make note of the following 2 items:

    Server port number: look for the message "Hyperpackage now running via Docker Container 9672b19c0b on port [XXXXX]". You'll need the port number [XXXXX].

    Fast App API key: look for the message "Fast App API key is: [FAST_API_KEY]"

9. Retrieve a prediction by running the following `curl` command in a terminal/bash session. The input data that you pass with the -d flag should be an array of values that is an appropriate shape (i.e., you should be passing 5 values if your model is expecting 5 inputs). You'll need both the server port number and Fast App API key from the previous step:

``` bash
    curl -X 'POST' \
    http://127.0.0.1:[SERVER_PORT_NUMBER_HERE]/predict \
    -H 'x-api-key: [FAST_APP_API_KEY_HERE]' \
    -d '[ARRAY_OF_VALUES_HERE]'
```

## Writing hyperpacks to S3

Hyperpackage provides a utility function to write your hyperpack to an Amazon S3 bucket. Here are the steps:

1. In a JupyterLab/Notebook or python session, import the `write_hyperpack_to_s3` function:

```
    from hyperpackage.utilities import write_hyperpack_to_s3
```

2. Call the `write_hyperpack_to_s3` function. All args are of string type. The `hyperpack_file` arg is either the hyperpack file name (if you're calling this function from the same directory where your hyperpack file is located), or a path to a hyperpack file. The `access_key_id`, `secret_access_key`, and `session_token` args are your S3 credentials (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN, respectively). The `s3_bucket` arg is the name of the target S3 bucket where you want the hyperpack written to:

```
    write_hyperpack_to_s3(
        hyperpack_file="my_hyperpack.zip",
        access_key_id=AWS_ACCESS_KEY_ID,
        secret_access_key=AWS_SECRET_ACCESS_KEY,
        session_token=AWS_SESSION_TOKEN,
        s3_bucket="my_beautiful_bucket"
    )
```

## Hyperpack Schema

The hyperpack, when unzipped, has the following structure/schema:

```
   my_hyperpack.hyperpack.zip
   ├── _hyperpack.yaml            # contains information about the model and training parameters 
   ├── _study.json                # contains details of best trial
   ├── 000001-friendly-trial
   │   ├── trained_model          # model in ONNX format
   │   ├── _trial.json            # contains details of best trial
   │   └── ...                    # optional additional contents can include an .ipynb notebook from the run
   ├── 000002-unwieldy-trial
       ├── ....

```

An example of a compressed and uncompressed hyperpack are available in the examples folder of this repository
