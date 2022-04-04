from ..meta.trial import TrialMeta
from ..datamodel.trial import Trial


class TrialController:
    @classmethod
    def _create_trial(
        cls,
        trained_model,
        metadata=None,
        metrics=None,
        hyperparameters=None,
        preprocessor=None,
        preprocessor_shape=None,
        model_flavor=None,
        train_shape=None,
        train_features=None,
        train_target=None,
        trial_id=None,
        trial_name=None,
        my_study_path=None,
    ):
        trial = Trial()

        trial.log_hyperparameters(hyperparameters)
        trial.log_metadata(metadata)
        trial.log_metrics(metrics)
        trial.log_model(trained_model)
        trial.log_preprocessor(preprocessor, preprocessor_shape)
        trial.preprocessor = preprocessor
        trial.model_flavor = model_flavor
        trial.train_shape = train_shape
        trial.train_features = train_features
        trial.train_target = train_target
        trial.trial_id = trial_id
        trial.trial_name = trial_name
        trial.my_study_path = my_study_path

        trial_manifest = TrialMeta.create(
            metadata=trial.metadata,
            metrics=trial.metrics,
            hyperparameters=trial.hyperparameters,
            trained_model=trial.trained_model,
            preprocessor=trial.preprocessor,
            preprocessor_shape=trial.preprocessor_shape,
            train_features=trial.train_features,
            train_target=trial.train_target,
            trial_id=trial.trial_id,
            trial_name=trial.trial_name,
            my_study_path=trial.my_study_path,
        )
        trial.extend_from_dict(trial_manifest)
        trial.write_artifacts()
        return trial
