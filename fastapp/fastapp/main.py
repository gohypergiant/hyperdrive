import logging
import os
from logging.config import dictConfig

from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles

from fastapp.logging.config import LogConfig
from fastapp.routers import model

dictConfig(LogConfig().dict())
logger = logging.getLogger("hyperpack-wrapper")

logger.info("Dummy Info")
logger.error("Dummy Error")
logger.debug("Dummy Debug")
logger.warning("Dummy Warning")

app = FastAPI(title="mlsdk-fastapp")
app.include_router(model.router)

@app.on_event("startup")
def show_fast_app_api_key():
    fast_key = os.environ.get("FASTKEY")
    if fast_key is None:
        logger.error("Fast App API key from environment is None.")
        raise TypeError("Fast App API key from environment is None.")
    fast_key_msg = "Fast App API key is: {api_key}".format(api_key=fast_key)
    logger.info(fast_key_msg)

@app.get("/status", status_code=200)
def root() -> dict:
    return {"status": "ok"}

app.mount("/", StaticFiles(directory="fastapp/ui", html=True), name="ui")
