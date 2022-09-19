"""
Utility functions for the following:
- generating a hyperpack
- unpacking a hyperpack
- retrieving files from the unzipped/unpacked hyperpack
"""

import json
import os
import shutil
import zipfile

import boto3
from icecream import ic
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
    hyperpack_file: str = "",
    access_key_id: str = None,
    secret_access_key: str = None,
    session_token: str = None,
    s3_bucket: str = None,
):
    """Writes hyperpack to a S3 bucket
    Args:
        hyperpack_file: name of, or path to, hyperpack file
        access_key_id: AWS_ACCESS_KEY_ID
        secret_access_key: AWS_SECRET_ACCESS_KEY
        session_token: AWS_SESSION_TOKEN
        s3_bucket: name of S3 bucket
    """
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
        raise ValueError("Please pass in a S3 bucket name.")

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
        f"*** COMPLETED: The {s3_object_key} hyperpack was written to the {s3_bucket} S3 bucket ***"
        )


def write_json(dictionary, json_file_path):
    """Writes object to JSON format
    Args:
        dictionary: pythson dict to be written to JSON
        json_file_path: save path of JSON object
    """
    with open(json_file_path, "w", encoding='utf-8') as json_file:
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


# Unpacking functions
def search_by_filename_substring(file_directory: str = None,
                                filename_substring: str = '.zip',
                                verbose: bool = False):
    """Searches for a hyperpack.zip file in the current directory
    :file_directory: file directory to search
    :filename_substring: substring the filename contains
    :verbose: True to print the filepath found; False to suppress print
    """
    directory = file_directory if file_directory is not None else './'
    for root, _, files in os.walk(directory):
        for file in files:
            if filename_substring in file:
                if verbose:
                    print (f"{root}{file}")
                return f"{root}{file}"
            else:
                pass


def unzip_study(zip_filepath: str = None, unzip_to: str = './'):
    """Unzips the zipped hyperpack to specified directory
    :param unzip_to: file directory to write the unzipped files to
    :return: unzip_to
    """
    with zipfile.ZipFile(zip_filepath, 'r') as zip_ref:
        zip_ref.extractall(unzip_to)
    return unzip_to


class UnpackHyperpack:
    """Unpack Hyperpack files & retrieve its contents
    """

    def __init__(self, hyperpack_filedir: str = None, verbose: bool = False):
        """
        :hyperpack_filedir: file directory where hyperpack is located
        :verbose: print if True, else suppress printing
        :return: None
        """
        self.hyperpack_filedir = hyperpack_filedir or './'
        self.hp_filepath = search_by_filename_substring(self.hyperpack_filedir, 'hyperpack.zip')
        self.hp_directory = unzip_study(self.hp_filepath)
        self.trial_names = self.get_trial_names()
        self.onnx_filepaths = self.get_onnx_files()

        self.verbose = verbose
        if self.verbose: 
            ic(self.hp_filepath)
            ic(self.hp_directory)
            
    def get_trial_names(self):
        """Retrieves the trial name(s) from the unzipped hyperpack directory
        :return: list of trial name(s) from hyperpacked optuna study
        """
        files = os.listdir(self.hp_directory)
        studies = [f for f in files if '_study.json' in f]

        trial_names = []
        for study in studies:
            with open(study, 'r', encoding='utf-8') as json_file:
                trial_meta = json.load(json_file)
                trial_names.append(trial_meta.get('best_trial'))
        if self.verbose:
            ic(trial_names)
        return trial_names

    def get_onnx_files(self):
        """retrieve trained models from the respective trials
        :return: list of onnx filepaths
        """
        trial_names = self.get_trial_names()
        onnx_filepaths = []
        for trial in trial_names:
            trial_dir = f"{self.hp_directory}/{trial}"
            for file in trial_dir:
                if '.onnx' in file or 'trained_model' in file:
                    onnx_filepaths.append(f"{trial_dir}/{file}")
        if self.verbose:
            ic(onnx_filepaths)
        return onnx_filepaths
