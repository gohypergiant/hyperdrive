from fastapp.services.model import get_default_model_id
from fastapp.services.model import info as svc_model_info
from fastapp.services.model import predict as svc_model_predict

default_model_id = get_default_model_id()


def batch(body, model_id=default_model_id):
    result = svc_model_predict(body, model_id=model_id)
    return {"predictions": list(result[0])}


def predict(body, model_id=default_model_id):
    result = svc_model_predict([body], model_id=model_id)[0]
    return {"prediction": list(result)[0]}


def info(model_id=default_model_id) -> dict:
    return svc_model_info(model_id=model_id)
