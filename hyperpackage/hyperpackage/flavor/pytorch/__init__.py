import os
import torch


def torch_onnx_export(
    model, hyperpack_dir: str, train_shape: int = 30,
):
    local_artifact_path = os.path.join(hyperpack_dir, "trained_model")
    initial_types = torch.randn(1, train_shape)
    torch.onnx.export(model, initial_types, local_artifact_path, train_shape)
