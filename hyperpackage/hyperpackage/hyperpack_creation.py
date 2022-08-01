import json
import os
import shutil
import yaml
from datetime import datetime
from hyperpackage.flavor.pytorch import torch_onnx_export

SUPPORTED_MODEL_FLAVORS = ["automl"]


def create_hyperpack(trained_model=None, model_flavor: str = None):
    curr_dir = os.getcwd()
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
    best_trial_name = "adventurous"
    best_trial_dir_name = generate_folder_name(name=best_trial_name)
    best_trial_path = os.path.join(hyperpack_path, best_trial_dir_name)
    os.makedirs(best_trial_path, exist_ok=True)
    torch_onnx_export(model=trained_model, trial_path=best_trial_path)
    created_time = datetime.now().strftime("%Y-%m-%d %H:%M")
    study_json_dict = {"best_trial": best_trial_dir_name, "created_at": created_time}
    write_json(study_json_dict, best_trial_path)
    zip_study(hyperpack_path)
    study_yaml_dict = {"project_name": model_flavor, "study_name": model_flavor}
    study_yaml_path = os.path.join(curr_dir, "study.yaml")
    write_yaml(study_yaml_dict, study_yaml_path)
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
    curr_dir = os.getcwd()
    hyperpack_folder_name = name + ".hyperpack"
    path = os.path.join(curr_dir, hyperpack_folder_name)
    return path


def write_json(dictionary, local_artifact_path):
    with open(local_artifact_path, "w") as json_file:
        json_file.write(json.dumps(dictionary))


def write_yaml(dictionary: dict, yaml_file_path: str):
    with open(yaml_file_path, "w") as stream:
        yaml.dump(dictionary, stream)


def zip_study(folder_path):
    root_dir = os.path.dirname(folder_path)
    base_dir = os.path.basename(folder_path)
    shutil.make_archive(folder_path, "zip", root_dir, base_dir)


def generate_folder_name(
    trial_id: int = 0,
    name: str = None,
    format_precision: str = "06",
    suffix: str = "trial",
):
    prefix = format(trial_id, format_precision)
    if name is not None:
        folder_name = f"{prefix}-{name}-{suffix}"
    else:
        folder_name = f"{prefix}-{suffix}"

    return folder_name
