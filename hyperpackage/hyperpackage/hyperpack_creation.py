import torch

SUPPORTED_MODEL_FLAVORS = ["automl"]


def create_hyperpack(trained_model=None, model_flavor: str = None):
    if trained_model is None:
        raise TypeError("You must pass in a trained model.")
    elif model_flavor is None:
        supported_flavors = "\n".join(map(str, SUPPORTED_MODEL_FLAVORS))
        raise TypeError(
            "You must specify a model flavor. Supported model flavors are:\n{}".format(
                supported_flavors
            )
        )
    torch_export(trained_model=trained_model)
    print("ahoy environs!")


def torch_export(
    trained_model,
    local_artifact_path="/home/jovyan/automl_trained_model",
    train_shape=30,
):
    initial_types = torch.randn(1, train_shape)
    torch.onnx.export(trained_model, initial_types, local_artifact_path, train_shape)
