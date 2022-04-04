from dataclasses import dataclass
from .__doc__ import TrialDoc
from .methods import TrialMethods
from ..base import Base
from ...controllers.trial_logging import TrialLoggingController


@dataclass
class Trial(Base, TrialDoc, TrialMethods, TrialLoggingController):
    id = None
    model_experiment_id = None
    metadata = None
    hyperparameters = None
    metrics = None
    trained_model = None
    baseline = None
    preprocessor = None
    preprocessor_shape = None
    dataset = None
    data_exploration_file = None
    data_preparation_file = None
    model_training_file = None

    mutable = (
        "preprocessor",
        "preprocessor_shape",
        "dataset",
        "data_exploration_file",
        "data_preparation_file",
        "model_training_file",
    )


__doc__ = TrialDoc.__doc__
