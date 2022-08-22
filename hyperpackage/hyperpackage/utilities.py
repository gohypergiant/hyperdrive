import json
import shutil
import yaml


def generate_folder_name(
    trial_id: int = 0,
    name: str = None,
    format_precision: str = "06",
    suffix: str = "trial",
) -> str:
    """Saves pytorch or automl model to ONNX format
    Args:
        trial_id: integer id of trial
        name: string to be used as part of folder name
        format_precision: number of digits to use for trial_id
        suffix: string to be used as trailing part of folder name
    """
    prefix = format(trial_id, format_precision)
    if name is not None:
        folder_name = f"{prefix}-{name}-{suffix}"
    else:
        folder_name = f"{prefix}-{suffix}"

    return folder_name


def write_hyperpack_to_s3(
    hyperpack_path="",
    access_key_id=None,
    secret_access_key=None,
    session_token=None,
    s3_bucket_path=None,
):
    if hyperpack_path == "":
        raise ValueError("Please pass in the path to your hyperpack zip file.")
    if (
        access_key_id in [None, ""]
        or secret_access_key in [None, ""]
        or session_token in [None, ""]
    ):
        raise ValueError(
            "Please pass in the following AWS S3 credentials: access_key_id, secret_access_key, and session_token."
        )
    if s3_bucket_path in [None, ""]:
        raise ValueError("Please pass in a S3 bucket path.")
    print("ciao mondo!")


def write_json(dictionary, json_file_path):
    """Writes object to JSON format
    Args:
        dictionary: python dict to be written to JSON
        json_file_path: save path of JSON object
    """
    with open(json_file_path, "w") as json_file:
        json_file.write(json.dumps(dictionary))


def write_yaml(dictionary: dict, yaml_file_path: str):
    """Writes object to YAML
    Args:
        dictionary: python dict to be written to YAML
        yaml_file_path: save path of YAML object
    """
    with open(yaml_file_path, "w") as stream:
        yaml.dump(dictionary, stream)


def zip_study(folder_path):
    """Zips folder to create final hyperpack
    Args:
        folder_path: path to dir to be zipped
    """
    shutil.make_archive(folder_path, "zip", folder_path)
