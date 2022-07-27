import os
from pathlib import Path
from hyperpackage.flavor.pytorch import torch_onnx_export

SUPPORTED_MODEL_FLAVORS = ["automl"]


def create_hyperpack(trained_model=None, model_flavor: str = None):
    verify_args(model=trained_model, flavor=model_flavor)
    hyperpack_path = make_hyperpack_path(name=model_flavor)
    try:
        os.makedirs(hyperpack_path, exist_ok=False)
    except FileExistsError:
        i = 1
        while os.path.exists(f"{hyperpack_path}_{str(i)}"):
            i += 1
        hyperpack_path = hyperpack_path + "_" + str(i)
        os.makedirs(hyperpack_path, exist_ok=False)
    torch_onnx_export(model=trained_model, hyperpack_dir=hyperpack_path)
    print("ahoy environs!")


def verify_args(model, flavor: str):
    supported_flavors = "\n".join(map(str, SUPPORTED_MODEL_FLAVORS))
    if model is None:
        raise TypeError("You must pass in a trained model.")
    elif flavor is None:
        raise TypeError(
            "You must specify a model flavor. Supported model flavors are:\n{}".format(
                supported_flavors
            )
        )
    elif flavor not in SUPPORTED_MODEL_FLAVORS:
        raise TypeError(
            "You have selected a model flavor that is currently not supported. Supported model flavors are:\n{}".format(
                supported_flavors
            )
        )


def make_hyperpack_path(name: str) -> str:
    home_dir = str(Path.home())
    hyperpack_folder_name = name + ".hyperpack"
    path = os.path.join(home_dir, hyperpack_folder_name)
    return path
