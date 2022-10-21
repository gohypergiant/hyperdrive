import onnx
import os
import tf2onnx


def tensorflow_onnx_export(model, trial_path: str):
    onnx_model, _ = tf2onnx.convert.from_keras(model)
    onnx.save(onnx_model, "trained_model.onnx")
