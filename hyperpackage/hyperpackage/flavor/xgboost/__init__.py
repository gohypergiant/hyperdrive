import onnx
import onnxmltools
import os

from skl2onnx.common.data_types import FloatTensorType


def xgboost_onnx_save(model, trial_path: str, train_shape_cols: int):
    model_path = os.path.join(trial_path, "trained_model")
    initial_types = [("X", FloatTensorType([None, train_shape_cols]))]
    onnx_model = onnxmltools.convert.convert_xgboost(
        model, initial_types=initial_types)
    onnx.save(onnx_model, model_path)
