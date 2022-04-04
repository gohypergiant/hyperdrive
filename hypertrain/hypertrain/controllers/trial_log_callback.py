import optuna
from .trial import TrialController


class TrialLogCallback:
    def __init__(self, study):
        """Stores the Study associated with the Trial in an attribute."""
        self.study = study

    def __call__(self, study, trial):
        """Logs the Optuna trials as Trials for hyperparameter tuning experiments.

        Parameters
        ----------
        study : optuna.study.study.Study
            An Optuna study.
        trial : optuna.trial._trial.Trial
            A single execution of the objective function.

        """
        if not trial.state == optuna.trial.TrialState.PRUNED:
            old_hyperparameters = trial.params
            hyperparameters = {}
            for key, value in old_hyperparameters.items():
                if "*" in key:
                    idx = key.find("*")
                    key = key[:idx]
                hyperparameters[key] = value

            experiment_metric = self.study.metric
            idx = experiment_metric.rfind(".")
            metric_name = experiment_metric[idx + 1 :]

            model_flavor = self.study.model_flavor
            train_shape = self.study.train_shape
            train_features = self.study.train_features
            train_target = self.study.train_target
            trial_id = trial._trial_id
            trial_name = self.study.trial_names[trial_id]
            my_study_path = self.study.my_study_path

            metrics = {metric_name: trial.value}
            metadata = {
                "trial_run_time": (
                    trial.datetime_complete - trial.datetime_start
                ).microseconds
                / 1000000,
                "trial_type": "optuna",
            }
            trial_params = {
                "trained_model": trial.user_attrs.get("trained_model", None),
                "metadata": metadata,
                "metrics": metrics,
                "hyperparameters": hyperparameters,
                "model_flavor": model_flavor,
                "train_shape": train_shape,
                "train_features": train_features,
                "train_target": train_target,
                "trial_id": trial_id,
                "trial_name": trial_name,
                "my_study_path": my_study_path,
            }
            TrialController._create_trial(**trial_params)

    __doc__ = """
    Callback for Optuna. Logs the Optuna trials as Trial for hyperparameter tuning
    studies.

    Methods
    ----------
    __call__(self, trial):
            Logs the Optuna trials as Trials for hyperparameter tuning studies.
    """
