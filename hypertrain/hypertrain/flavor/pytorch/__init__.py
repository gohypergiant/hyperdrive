import torch


class TorchModelHandler:
    @classmethod
    def _export(cls, trained_model, local_artifact_path, train_shape):
        initial_types = torch.randn(1, train_shape)
        torch.onnx.export(
            trained_model, initial_types, local_artifact_path, train_shape
        )
