import onnx
import os
import tf2onnx


def tensorflow_onnx_export(model, trial_path: str):
    onnx_model, _ = tf2onnx.convert.from_keras(model)
    model_path = os.path.join(trial_path, "trained_model")
    onnx.save(onnx_model, model_path)
