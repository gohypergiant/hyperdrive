import os


def xgboost_onnx_save(model, trial_path: str, train_shape_cols: int):
    print("xgboost to onnx")
