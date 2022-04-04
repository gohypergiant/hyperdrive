from dataclasses import dataclass
from .__doc__ import HyperparameterStudyDoc
from ...controllers.hyperparameter_study_tuning import (
    HyperparameterStudyTuningController,
)
from ...controllers.trial import TrialController


@dataclass
class HyperparameterStudy(
    HyperparameterStudyDoc, HyperparameterStudyTuningController, TrialController,
):
    id = None
    experiment_name = None
    model_flavor = None
    best_run_id = None
    preprocessor = None
    optuna_study = None
    direction = None
    pruner = None
    sampler = None
    n_trials = 100
    metric = None
    train_shape = None
    search_dictionary = None
    my_study_path = None

    mutable = tuple()

    def __repr__(self):
        return f"<HyperparameterExperiment {self.experiment_name}>"


__doc__ = HyperparameterStudyDoc.__doc__
