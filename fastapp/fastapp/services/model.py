import json
import logging
import numpy as np
from os import path

from onnxruntime import InferenceSession
from scipy.special import expit, log_softmax

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
    study_info = get_study_info()
    ml_task = study_info["ml_task"]
    print("\n\nml_task is:", ml_task, "\n\n")
    trained_model_path = path.join(model_path(model_id), "trained_model")
    model = ONNXModel(trained_model_path)
    try:
        result = model.predict(input_data=np.array(input_data, dtype=np.float32))
        if ml_task == "binary_classification":
            result = expit(result).round().item()
            print("ONE binary result:", result)
        elif ml_task == "multi_class_classification":
            result = log_softmax(result).argmax().item()
            print("ONE multi result:", result)
        else:
            print("ONE regression result:", result)
            result = result[0].item()
        return result
    except ValueError as err:
        logging.error(err)


class ONNXModel:
    def __init__(self, path):

        self.session = InferenceSession(path)
        self.input_name = self.session.get_inputs()[0].name
        self.label_name = self.session.get_outputs()[0].name

    def predict(self, input_data):
        return self.session.run([self.label_name], {self.input_name: input_data})
