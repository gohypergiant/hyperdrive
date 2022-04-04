import os
from datetime import datetime
from ...utilities import read_json_from_local, write_json_to_local


class ManifestInterface:
    @classmethod
    def write(
        self, manifest, kind, my_study_path, trial_folder=None,
    ):
        if kind == "trial":
            manifest["created_at"] = datetime.now().strftime("%Y-%m-%d %H:%M")
            full_trial_dir = f"{my_study_path}/{trial_folder}"
            os.makedirs(full_trial_dir, exist_ok=True)
            write_json_to_local(manifest, f"{full_trial_dir}/_{kind}.json")
        elif kind == "study":
            best_trial = manifest["best_trial"]
            best_trial_json_path = f"{my_study_path}/{best_trial}/_trial.json"
            trial_json = read_json_from_local(best_trial_json_path)
            manifest["created_at"] = trial_json["created_at"]
            write_json_to_local(manifest, f"{my_study_path}/_{kind}.json")
