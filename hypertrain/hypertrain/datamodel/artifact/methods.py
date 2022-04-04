import os
from pathlib import Path
from ...exceptions import HyperparameterStudyError
from ...utilities import generate_file_key, generate_folder_name


class ArtifactMethods:
    def generate_artifact_file_key(self, name=None):
        """Generates the file key for the Artifact of the Run.

        Parameters
        ----------

        name: str

        Returns
        -------
        A formatted string in the Project structure for the Data Repo.

        """
        if name is None:
            name = Path(self.local_artifact_path).name

        file_key_params = {
            "trial_id": self.trial_id,
            "name": name,
        }

        self.file_key = generate_file_key(**file_key_params)


class MachineLearningModelMethods:
    def save_to_workspace(self):
        """Uses the model_flavor attribute of the Experiment to choose the appropriate
        trained model and transform the artifact type before saving the data.

        Parameters
        ----------
        artifact: hdsdk.datamodel.Artifact object
            An Artifact object

        """
        trial_folder_name = generate_folder_name(
            self.trial.trial_id, self.trial.trial_name
        )
        self.local_artifact_path = (
            f"{self.trial.my_study_path}/{trial_folder_name}/" + self.artifact_type
        )
        self.generate_artifact_file_key()

        if (
            self.trial.model_flavor in ("sklearn", "lightgbm")
            or self.artifact_type == "preprocessor"
        ):
            from ...flavor.sklearn import SklearnModelHandler

            if self.artifact_type == "trained_model":
                model = self.artifact
                shape = self.trial.train_shape
                if type(self.trial.train_shape) is not int:
                    shape = self.trial.train_shape[1]
            else:
                model = self.trial.preprocessor
                shape = self.trial.preprocessor_shape
            onnx_model = SklearnModelHandler._convert(model, shape)
            SklearnModelHandler._save(onnx_model, self.local_artifact_path)
        elif self.trial.model_flavor == "xgboost":
            from ...flavor.xgboost import XGBoostModelHandler

            shape = self.trial.train_shape
            if type(self.trial.train_shape) is not int:
                shape = self.trial.train_shape[1]
            onnx_model = XGBoostModelHandler._convert(self.artifact, shape)
            XGBoostModelHandler._save(onnx_model, self.local_artifact_path)
        elif self.trial.model_flavor == "tensorflow":
            from ...flavor.tensorflow import TensorflowModelHandler

            shape = self.trial.train_shape
            if type(self.trial.train_shape) is not int:
                shape = self.trial.train_shape[1]
            onnx_model = TensorflowModelHandler._convert(self.artifact, shape)
            TensorflowModelHandler._save(onnx_model, self.local_artifact_path)

        elif self.trial.model_flavor == "pytorch":
            from ...flavor.pytorch import TorchModelHandler

            shape = self.trial.train_shape
            if type(self.trial.train_shape) is not int:
                shape = self.trial.train_shape[1]
            TorchModelHandler._export(
                self.artifact,
                self.local_artifact_path,
                shape,
            )
        else:
            raise HyperparameterStudyError("Other model flavors not yet implemented.")

        self.metadata = {"Content-Type": "application/octet-stream"}


class PythonFileMethods:
    def save_to_workspace(self):
        """Sets the path of the notebook which satisfies the artifact type for the
        Run.

        Parameters
        ----------
        artifact: hdsdk.datamodel.Artifact object
            An Artifact object

        """
        if self.artifact_type == "data_exploration_file":
            self.artifact_type = "data_exploration_notebook"
            self.local_artifact_path = os.path.realpath(
                self.run.experiment.data_exploration_file
            )
        elif self.artifact_type == "model_training_file":
            self.artifact_type = "model_training_notebook"
            self.local_artifact_path = os.path.realpath(self.run.model_training_file)
        else:
            self.artifact_type = "data_preparation_notebook"
            self.local_artifact_path = os.path.realpath(
                self.run.experiment.data_preparation_file
            )

        self.generate_artifact_file_key()
        self.remove_local_artifact = False
        self.metadata = {"Content-Type": "application/x-ipynb+json"}
