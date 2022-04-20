from glob import glob
from pathlib import Path
from hddataclient import DataRepo
import os
import papermill
import yaml


class WalkerMethods:
    def walk_job_tree(self):
        job_tree = glob(f"{self.job_dir}/*")
        job_names = [job.split("/")[-1] for job in job_tree]
        for job_name in job_names:
            files = glob(f"{self.job_dir}/{job_name}/*")
            started_status = any(["STARTED" in file for file in files])
            completed_status = any(["COMPLETED" in file for file in files])
            if not started_status and not completed_status:
                self.set_status(job_name=job_name, status="started")
                return job_name

    def run_next_job(self, job_name):
        executor_notebook_paths = {
            "test": "/home/jovyan/.executor/notebooks/executor-test.ipynb",
            "automl": "/home/jovyan/.executor/notebooks/executor-automl.ipynb",
            "low-code": "/home/jovyan/.executor/notebooks/executor-low-code.ipynb",
        }
        study_yaml_path = f"/home/jovyan/_jobs/{job_name}/_study.yaml"
        study_definition = self._parse_study_yaml(study_yaml_path)

        executor_notebook_type = study_definition["notebook_type"]
        executor_input_notebook_path = executor_notebook_paths[executor_notebook_type]
        executor_output_notebook_path = executor_input_notebook_path.replace(
            ".executor/notebooks", f"_jobs/{job_name}"
        )

        features_path = study_definition["features_source"]
        target_path = study_definition["target_source"]

        try:
            papermill.execute_notebook(
                input_path=executor_input_notebook_path,
                output_path=executor_output_notebook_path,
                parameters={
                    "features": features_path,
                    "target": target_path,
                    "job_name": job_name,
                    "study_yaml": study_yaml_path,
                },
            )
        except papermill.exceptions.PapermillExecutionError:
            pass
        self.set_status(job_name=job_name, status="completed")
        return True

    def set_status(self, job_name, status):
        paths = {
            "started": f"{self.job_dir}/{job_name}/STARTED",
            "completed": f"{self.job_dir}/{job_name}/COMPLETED",
        }
        for status, path in paths.items():
            Path(path).unlink(missing_ok=True)
        Path(paths[status]).touch()

    def _parse_study_yaml(self, study_yaml_path):
        with open(study_yaml_path) as fh:
            my_study = yaml.safe_load(fh)
        if "automl" in my_study.keys() and my_study["automl"]:
            notebook_type = "automl"
        elif "test" in my_study.keys() and my_study["test"]:
            notebook_type = "test"
        else:
            notebook_type = "low-code"

        features_source = my_study["training"]["data"]["features"]["source"]
        target_source = my_study["training"]["data"]["target"]["source"]
        join_id = my_study["training"]["data"]["join_id"]
        target_response_variable = my_study["training"]["data"]["target"][
            "response_variable"
        ]

        study_definition = {
            "notebook_type": notebook_type,
            "features_source": features_source,
            "target_source": target_source,
            "join_id": join_id,
            "target_response_variable": target_response_variable,
        }

        return study_definition

    def _prepare_data(self, job_name, study_definition):
        datarepo = DataRepo(
            description="Study as Data Repo",
            volume_name=f"/home/jovyan/_jobs/{job_name}/",
        )
        features_df = datarepo.load_dataset(study_definition["features_source"])
        target_df = datarepo.load_dataset(study_definition["target_source"])
        join_id = study_definition["join_id"]
        return features_df.merge(target_df, left_on=join_id, right_on=join_id)
