import torch


def torch_onnx_export(
    trained_model,
    local_artifact_path="/home/jovyan/automl_trained_model",
    train_shape=30,
):
    initial_types = torch.randn(1, train_shape)
    torch.onnx.export(trained_model, initial_types, local_artifact_path, train_shape)
