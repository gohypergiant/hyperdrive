import os
import torch
from datetime import datetime
from hyperpackage.flavor.pytorch import torch_onnx_export
from hyperpackage.utilities import (
    generate_folder_name,
    write_json,
    write_yaml,
    zip_study,
)

SUPPORTED_MODEL_FLAVORS = ["automl"]


def create_hyperpack(trained_model=None, model_flavor: str = None):
    curr_dir = os.getcwd()
    verify_args(model=trained_model, flavor=model_flavor)
    loaded_model = load_trained_model(model=trained_model)
    hyperpack_path = make_hyperpack_path(base_dir=curr_dir, name=model_flavor)
    try:
        os.makedirs(hyperpack_path, exist_ok=False)
    except FileExistsError:
        i = 1
        while os.path.exists(f"{hyperpack_path}_{str(i)}"):
            i += 1
        hyperpack_path = hyperpack_path + "_" + str(i)
        os.makedirs(hyperpack_path, exist_ok=False)
    best_trial_name = "adventurous"
    best_trial_folder_name = generate_folder_name(name=best_trial_name)
    best_trial_path = os.path.join(hyperpack_path, best_trial_folder_name)
    os.makedirs(best_trial_path, exist_ok=True)
    save_best_trial_model(
        model=loaded_model, flavor=model_flavor, save_path=best_trial_path
    )
    create_study_json(hyperpack_path=hyperpack_path, best_trial=best_trial_folder_name)
    create_study_yaml(current_dir=curr_dir, name=model_flavor)
    zip_study(folder_path=hyperpack_path)
    print("ahoy environs!")


def verify_args(model, flavor: str):
    supported_flavors = "\n".join(map(str, SUPPORTED_MODEL_FLAVORS))
    if model is None:
        raise TypeError("You must pass in a trained model.")
    elif isinstance(model, str):
        if not os.path.exists(model):
            raise FileNotFoundError("No file could be found at {}.".format(model))

    if flavor is None:
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


def load_trained_model(model):
    if isinstance(model, str):
        try:
            the_model = torch.load(model)
        except Exception:
            print("Error while attempting to load torch model.")
    elif str(type(model)) == "<class 'neural_network.network.Network'>":
        the_model = model
    else:
        raise TypeError("The model type you have passed in is currently not supported.")

    return the_model


def make_hyperpack_path(base_dir: str, name: str) -> str:
    hyperpack_folder_name = name + ".hyperpack"
    path = os.path.join(base_dir, hyperpack_folder_name)
    return path


def save_best_trial_model(model, flavor: str, save_path: str):
    if flavor == "automl":
        torch_onnx_export(model=model, trial_path=save_path)


def create_study_json(hyperpack_path: str, best_trial: str):
    created_time = datetime.now().strftime("%Y-%m-%d %H:%M")
    study_json_dict = {"best_trial": best_trial, "created_at": created_time}
    study_json_path = os.path.join(hyperpack_path, "_study.json")
    write_json(dictionary=study_json_dict, json_file_path=study_json_path)


def create_study_yaml(current_dir: str, name: str):
    study_yaml_dict = {"project_name": name, "study_name": name}
    study_yaml_path = os.path.join(current_dir, "study.yaml")
    write_yaml(dictionary=study_yaml_dict, yaml_file_path=study_yaml_path)
