import json
from os import listdir, path

from fastapp.services.utils import model_slug_info

path_map = {}
hyperpackage_path = "/hyperpackage"
for item in listdir(hyperpackage_path):
    item_path = path.join(hyperpackage_path, item)
    if not path.isdir(item_path):
        continue
    item_id = item.split("-")[0]
    trimmed_item_id = item_id.lstrip("0")
    slug_info = model_slug_info(item)
    path_map[slug_info["name"]] = item_path
    path_map[slug_info["id"]] = item_path
    path_map[slug_info["trimmed_id"]] = item_path


def model_path(slug: str) -> str:
    return path_map[slug]


def get_study_info() -> dict:
    study_info_path = path.join(hyperpackage_path, "_study.json")
    study_info_file = open(study_info_path)
    return json.load(study_info_file)
