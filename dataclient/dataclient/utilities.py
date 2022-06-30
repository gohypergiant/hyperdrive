import json
import shutil
from datetime import date, datetime
from pathlib import Path

import yaml
from pandas.io.sql import SQLiteDatabase, SQLiteTable

from .exceptions import FileSystemError


def cron_to_datetime(cron):
    cron_split = cron.split(" ")
    return datetime(
        date.today().year,
        int(cron_split[3]),
        int(cron_split[2]),
        int(cron_split[1]),
        int(cron_split[0]),
    )


def datetime_to_cron(dt):
    return f"{dt.minute} {dt.hour} {dt.day} {dt.month} {dt.weekday()}"


def generate_ddl(column_names_and_types):
    return ", ".join([f"{name} {_type}" for name, _type in column_names_and_types])


def generate_file_key(model_name=None, experiment_name=None, trial_id=None, name=None):
    file_key = "Projects/Project - Friendly_Trial"
    if model_name is not None:
        file_key += f"/Model - {model_name}"
        if experiment_name is not None:
            file_key += f"/Experiment - {experiment_name}"
            if trial_id is not None:
                file_key += f"/Run - {trial_id}"
    if name is not None:
        file_key += f"/{name}"
    return file_key


def generate_folder_name(
    trial_id: int,
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


def get_column_names_and_types(dataframe):
    table = SQLiteTable("_", SQLiteDatabase(None), dataframe)
    column_names_and_types = [
        (
            table.frame.columns[i],
            table._sqlalchemy_type(table.frame.iloc[:, i]).__visit_name__,
        )
        for i in range(len(table.frame.columns))
    ]

    return column_names_and_types


def infer_target_format(file_key, target_format="parquet"):
    path = Path(file_key)
    if path.suffix == "":
        file_key = f"{file_key}.{target_format}"
        path = Path(file_key)
    else:
        target_format = path.suffix.replace(".", "")
    if target_format not in ("csv", "parquet", "json"):
        raise FileSystemError(
            "By design, Data Repos only support writing data "
            "using the CSV, Parquet or JSON formats."
        )
    return target_format, file_key


def load_yaml_file(yaml_file_name):
    """Converts the YAML file into a Python dictionary."""
    with open(yaml_file_name, "r") as stream:
        return yaml.safe_load(stream)


def read_json_from_local(local_artifact_path):
    with open(local_artifact_path, "r") as json_file:
        return json.loads(json_file.read())


def transform_dict_key(dictionary, old_key, new_key):
    dictionary[new_key] = dictionary[old_key]
    del dictionary[old_key]
    return dictionary


def update_yaml_file(yaml_file_name, key, value):
    """Converts the YAML file into a Python dictionary."""
    with open(yaml_file_name, "r") as stream:
        yaml_as_dict = yaml.safe_load(stream)
    yaml_as_dict[key] = value

    with open(yaml_file_name, "w") as stream:
        yaml.dump(yaml_as_dict, stream)

    return True


def write_json_to_local(dictionary, local_artifact_path):
    with open(local_artifact_path, "w") as json_file:
        json_file.write(json.dumps(dictionary))


def write_yaml_to_local(dictionary, yaml_file_name):
    with open(yaml_file_name, "w") as stream:
        yaml.dump(dictionary, stream)


def zip_study(folder_path):
    shutil.make_archive(folder_path, "zip", folder_path)
