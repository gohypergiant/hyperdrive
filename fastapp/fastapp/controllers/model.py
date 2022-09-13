from fastapp.services.model import get_default_model_id
from fastapp.services.model import info as svc_model_info
from fastapp.services.model import predict as svc_model_predict
from fastapp.services.model import batch_predict as svc_batch_predict

default_model_id = get_default_model_id()


def batch(body, model_id=default_model_id):
    results = svc_batch_predict(body, model_id=model_id)
    return {"predictions": results}


def predict(body, model_id=default_model_id):
    result = svc_model_predict([body], model_id=model_id)
    return {"prediction": result}


def info(model_id=default_model_id) -> dict:
    return svc_model_info(model_id=model_id)
