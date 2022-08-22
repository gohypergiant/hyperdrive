import boto3
import json
import os
import shutil
import yaml
import zipfile


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
    hyperpack_file="",
    access_key_id=None,
    secret_access_key=None,
    session_token=None,
    s3_bucket=None,
):
    if hyperpack_file == "":
        raise ValueError(
            "Please pass in the name of, or path to, your hyperpack zip file."
        )
    elif not os.path.isfile(hyperpack_file):
        raise FileNotFoundError("No file could be found at {}".format(hyperpack_file))
    elif not zipfile.is_zipfile(hyperpack_file):
        raise TypeError("The file {} is an invalid zip file.".format(hyperpack_file))

    if (
        access_key_id in [None, ""]
        or secret_access_key in [None, ""]
        or session_token in [None, ""]
    ):
        raise ValueError(
            "Please pass in all of the following AWS S3 credentials: access_key_id, secret_access_key, and session_token."
        )

    if s3_bucket in [None, ""]:
        raise ValueError("Please pass in a S3 bucket.")

    s3 = boto3.resource(
        "s3",
        aws_access_key_id=access_key_id,
        aws_secret_access_key=secret_access_key,
        aws_session_token=session_token,
    )

    s3_object_key = hyperpack_file.rsplit("/", 1)[-1]

    try:
        s3.meta.client.upload_file(hyperpack_file, s3_bucket, s3_object_key)
    except Exception as exp:
        print("S3 upload error: ", exp)
        raise

    print(
        "*** COMPLETED: The {} hyperpack was written to the {} S3 bucket ***".format(
            s3_object_key, s3_bucket
        )
    )


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
