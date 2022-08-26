import json
import logging
from os import path

import numpy as np
from onnx.backend.test.case.node.softmax import softmax
from onnxruntime import InferenceSession

from fastapp.services.hyperpackage import get_study_info, model_path
from fastapp.services.utils import model_slug_info


def info(model_id: str) -> dict:
    model_info_path = path.join(model_path(model_id), "_trial.json")
    model_info_file = open(model_info_path)
    return json.load(model_info_file)


def get_default_model_id() -> str:
    study_info = get_study_info()
    return model_slug_info(study_info["best_trial"])["id"]


def predict(input_data, model_id: str):
    trained_model_path = path.join(model_path(model_id), "trained_model")
    model = ONNXModel(trained_model_path)
    try:
        result = model.predict(input_data=np.array(input_data, dtype=np.float32))
        result = softmax(result).argmax().item()
        return np.array(result).tolist()
    except ValueError as err:
        logging.error(err)


class ONNXModel:
    def __init__(self, path):

        self.session = InferenceSession(path)
        self.input_name = self.session.get_inputs()[0].name
        self.label_name = self.session.get_outputs()[0].name

    def predict(self, input_data):
        return self.session.run([self.label_name], {self.input_name: input_data})
