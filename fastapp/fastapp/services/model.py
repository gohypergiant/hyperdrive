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
    if study_info:
        return model_slug_info(study_info["best_trial"])["id"]
    else:
        return study_info


def predict(input_data, model_id: str):
    study_info = get_study_info()
    ml_task = study_info["ml_task"]
    model_flavor = study_info["model_flavor"]

    trained_model_path = path.join(model_path(model_id), "trained_model")
    model = ONNXModel(trained_model_path)

    try:
        result = model.predict(input_data=np.array(input_data, dtype=np.float32))
        if model_flavor == "automl":
            if ml_task == "binary_classification":
                result = int(expit(result).round())
            elif ml_task == "multi_class_classification":
                result = log_softmax(result).argmax().item()
            else:
                result = result[0].item()
        elif model_flavor == "tensorflow":
            if ml_task == "binary_classification":
                result = int(result[0].argmax())
            else:
                result = float(result[0])
        elif model_flavor == "xgboost":
            if ml_task == "binary_classification":
                result = int(result[0])
            else:
                result = float(result[0])
        else:
            raise TypeError(
                "The '{}' model flavor is currently not supported.".format(model_flavor)
            )
        return result
    except ValueError as err:
        logging.error(err)


def batch_predict(input_data, model_id: str):
    study_info = get_study_info()
    ml_task = study_info["ml_task"]

    trained_model_path = path.join(model_path(model_id), "trained_model")
    model = ONNXModel(trained_model_path)

    try:
        if ml_task == "binary_classification":
            predictions = [
                model.predict(input_data=np.array([input], dtype=np.float32))
                for input in input_data
            ]
            results = [int(expit(pred).round()) for pred in predictions]
        elif ml_task == "multi_class_classification":
            predictions = [
                model.predict(input_data=np.array([input], dtype=np.float32))
                for input in input_data
            ]
            results = [log_softmax(pred).argmax().item() for pred in predictions]
        else:
            predictions = [
                model.predict(input_data=np.array([input], dtype=np.float32))
                for input in input_data
            ]
            results = [pred[0].item() for pred in predictions]
        return results
    except ValueError as err:
        logging.error(err)


class ONNXModel:
    def __init__(self, path):

        self.session = InferenceSession(path)
        self.input_name = self.session.get_inputs()[0].name
        self.label_name = self.session.get_outputs()[0].name

    def predict(self, input_data):
        return self.session.run([self.label_name], {self.input_name: input_data})
