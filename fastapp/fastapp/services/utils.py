import logging
import os
from fastapi import HTTPException

logger = logging.getLogger("hyperpack-wrapper")


def model_slug_info(model_slug: str) -> dict:
    id = model_slug.split("-")[0]
    trimmed_id = id.lstrip("0")
    return {"name": model_slug, "id": id, "trimmed_id": trimmed_id}


def check_api_key(fast_key_header: str):
    if fast_key_header is None:
        logger.error("Fast App API key from header is None.")
        raise TypeError("Fast App API key from header is None.")

    fast_key_env_var = os.environ.get("FASTKEY")

    if fast_key_header != fast_key_env_var:
        raise HTTPException(
            status_code=401,
            detail="Fast App API keys from header vs environment don't match.",
        )
