import json
import shutil
import yaml


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


def write_json(dictionary, json_file_path):
    with open(json_file_path, "w") as json_file:
        json_file.write(json.dumps(dictionary))


def write_yaml(dictionary: dict, yaml_file_path: str):
    with open(yaml_file_path, "w") as stream:
        yaml.dump(dictionary, stream)


def zip_study(folder_path):
    shutil.make_archive(folder_path, "zip", folder_path)
