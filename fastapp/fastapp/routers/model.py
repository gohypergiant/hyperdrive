from fastapi import APIRouter, Request

from fastapp.controllers.model import batch, info, predict
from fastapp.services.utils import check_api_key

router = APIRouter()


@router.post("/predict")
async def default_predict(request: Request):
    fast_key_header = request.headers.get("x-api-key")
    check_api_key(fast_key_header)
    body = await request.json()
    return predict(body)


@router.post("/batch")
async def default_batch(request: Request):
    fast_key_header = request.headers.get("x-api-key")
    check_api_key(fast_key_header)
    body = await request.json()
    return batch(body)


@router.get("/info")
def default_info() -> dict:
    return info()


@router.post("/{model_id}/predict")
async def model_predict(model_id: str, request: Request):
    fast_key_header = request.headers.get("x-api-key")
    check_api_key(fast_key_header)
    body = await request.json()
    return predict(body, model_id)


@router.post("/{model_id}/batch")
async def model_batch(model_id: str, request: Request):
    fast_key_header = request.headers.get("x-api-key")
    check_api_key(fast_key_header)
    body = await request.json()
    return batch(body, model_id)


@router.get("/{model_id}/info")
def model_info(model_id: str) -> dict:
    return info(model_id=model_id)
