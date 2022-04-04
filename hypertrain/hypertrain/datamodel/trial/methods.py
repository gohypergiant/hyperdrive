from ..artifact import Manifest, Preprocessor, TrainedModel


class TrialMethods:
    def write_artifacts(self):
        manifest_artifact_definition = {
            "trial": self,
            "artifact_type": "manifest",
        }

        self.manifest_artifact = Manifest.create_from_dict(manifest_artifact_definition)
        self.manifest_artifact.generate_artifact_file_key(name="_trial.json")

        if self.preprocessor is not None:
            preprocessor_artifact_definition = {
                "run": self,
                "artifact": self.preprocessor,
                "artifact_type": "preprocessor",
            }
            self.preprocessor_artifact = Preprocessor.create_from_dict(
                preprocessor_artifact_definition
            )
            self.preprocessor_artifact.save_to_workspace()
            self.preprocessor_artifact.write_to_datarepo()
            self.preprocessor_artifact.write_metadata()

        if self.trained_model is not None:
            trained_model_artifact_definition = {
                "trial": self,
                "artifact": self.trained_model,
                "artifact_type": "trained_model",
                "model_flavor": self.model_flavor,
                "train_shape": self.train_shape,
            }
            self.trained_model_artifact = TrainedModel.create_from_dict(
                trained_model_artifact_definition
            )
            self.trained_model_artifact.save_to_workspace()
