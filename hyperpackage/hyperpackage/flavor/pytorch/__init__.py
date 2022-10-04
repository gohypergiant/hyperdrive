import os
import torch


def torch_onnx_export(
    model, trial_path: str, train_shape_cols: int, ml_task : str
):
    """Saves pytorch or automl model to ONNX format
    Args:
        model: pretrained model
        trial_path: dir path to the trial
        train_shape_cols: the number of columns in the training dataset used to
                          train the pretrained model
        ml_task : the type of machine learning task
    """
    if ml_task == "univariate_timeseries":
        for key, m in model.items():
            model_path = os.path.join(trial_path, "trained_model_" + str(key))
            initial_types = torch.randn(1, train_shape_cols)
            torch.onnx.export(m, initial_types, model_path, train_shape_cols)
    else:
        model_path = os.path.join(trial_path, "trained_model")
        initial_types = torch.randn(1, train_shape_cols)
        torch.onnx.export(model, initial_types, model_path, train_shape_cols)
