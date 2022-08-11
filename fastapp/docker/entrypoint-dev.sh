#!/bin/bash
conda run --no-capture-output -n fast-app uvicorn --reload --port $PORT --host 0.0.0.0 fastapp.main:app 