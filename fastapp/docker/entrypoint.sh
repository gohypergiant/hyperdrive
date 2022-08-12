#!/bin/bash
conda run --no-capture-output -n fast-app uvicorn --port $PORT --host 0.0.0.0 fastapp.main:app