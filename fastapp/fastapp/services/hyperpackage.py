import boto3
import json
import zipfile
from os import listdir, path

from fastapp.services.utils import model_slug_info

path_map = {}
hyperpackage_path = "/hyperpackage"
hyperpackage_upload_file = "upload/hyperpack.zip"
hyperpackage_s3_file = "hyperpack_s3.txt"
hyperpackage_s3_path = "/" + hyperpackage_s3_file

def extract_hyperpackage():
    with zipfile.ZipFile(hyperpackage_upload_file, "r") as zip_file:
      zip_file.extractall(hyperpackage_path)

    path_map.clear()

    for item in listdir(hyperpackage_path):
        item_path = path.join(hyperpackage_path, item)
        if not path.isdir(item_path):
            continue
        slug_info = model_slug_info(item)
        path_map[slug_info["name"]] = item_path
        path_map[slug_info["id"]] = item_path
        path_map[slug_info["trimmed_id"]] = item_path

if path.exists(hyperpackage_s3_path):
    session = boto3.Session()
    client = session.client("s3")

    with open(hyperpackage_s3_path, "r") as f:
        lines = f.read().splitlines()
        bucket = lines[0]
        s3_hyperpack_path = lines[1]
        client.download_file(bucket, s3_hyperpack_path, hyperpackage_upload_file)

        extract_hyperpackage()

def model_path(slug: str) -> str:
    return path_map[slug]


def get_study_info() -> dict:
    try:
        study_info_path = path.join(hyperpackage_path, "_study.json")
        study_info_file = open(study_info_path)
        return json.load(study_info_file)
    except Exception as e:
        print(e)
        return "";

