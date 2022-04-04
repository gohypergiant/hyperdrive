# Hyperdrive CLI

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

To use a local Firefly server, first create the server

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

