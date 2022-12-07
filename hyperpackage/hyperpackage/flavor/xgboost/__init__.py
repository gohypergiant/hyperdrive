import onnx
import onnxmltools
import os

from skl2onnx.common.data_types import FloatTensorType


def xgboost_onnx_save(model, trial_path: str, train_shape_cols: int):
    print("xgboost to onnx")
