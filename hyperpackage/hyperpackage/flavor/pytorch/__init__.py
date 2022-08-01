import os
import torch


def torch_onnx_export(
    model, trial_path: str, train_shape: int = 30,
):
    local_artifact_path = os.path.join(trial_path, "trained_model")
    initial_types = torch.randn(1, train_shape)
    torch.onnx.export(model, initial_types, local_artifact_path, train_shape)
