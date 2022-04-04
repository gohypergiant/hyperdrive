FROM python:3.9-buster AS hypertrain

ARG HYPERDRIVE_REGISTRY_TOKEN

RUN pip install git+https://$HYPERDRIVE_REGISTRY_TOKEN@github.com/gohypergiant/mlsdk-hypertrain.git@MLSDK-6 --quiet

ADD docker/requirements.txt .
RUN pip install -r requirements.txt

RUN useradd -ms /bin/bash mlsdk

USER mlsdk
WORKDIR /home/mlsdk

# use this ENTRYPOINT + CMD later after implementing __main__.py
# ENTRYPOINT ["python"]
# CMD ["-m", "hypertrain"]

ENTRYPOINT ["top", "-b"]
CMD ["-c"]
