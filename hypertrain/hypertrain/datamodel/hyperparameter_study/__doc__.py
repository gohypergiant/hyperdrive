class HyperparameterStudyDoc:
    __doc__ = """
    A container for hyperparameter tuning experiments.

    Attributes
    ----------
    id: str, optional
        The id of the experiment (default is None).
    model: hdsdk.Model, optional
        The Model containing this experiment (default is None).
    model_id: str, optional
        The id of the Model to which this experiment belongs (default is None).
    experiment_name: str, optional
        The name of the experiment (default is None).
    model_flavor: str {pytorch, sklearn, tensorflow, xgboost}, optional
        The machine learning framework used in this experiment (default
        is None).
    best_run_id: str, optional
        Id for the best run of the experiment (default is None).
    preprocessor: Sklearn object, optional
        The trained sklearn preprocessor for data preprocessing (default is None).
    preprocessor_shape: int , optional
        The dimensions of the input to the preprocessor (default is None).
    data_exploration_file: str, optional
        Path of the eda file (default is None).
    data_preparation_file: str, optional
        Path of the data engineering file (default is None).
    model_training_file: str, optional
        Path of the model training file (default is None).
    optuna_study: optuna.study.study.Study, optional
        The Optuna study associated with this experiment (default is None).
    direction: str {maximize, minimize}, optional
        The direction in which the metric should be optimized (default is None).
    pruner: optuna.pruners._base.BasePruner, optional
        A pruner object that determines early stopping (default is None).
    sampler: optuna.samplers._base.BaseSampler, optional
        A sampler object that specifies the sampling algorithm to be used (default is
        None).
    n_trials: int, optional
        Number of trials to be run (default is None).
    metric: str, optional
        The metric used to evaluate the models trained in each trial (default is None).
    models: dict, optional
        Python dictionary specifying a list of models which can be tuned for this
        experiment (default is None).
    train_shape: int, optional
        The dimensions of an input vector (default is None).
    search_dictionary: dict, optional
        Python dictionary specifying the hyperparameters to be tuned and the values
        they can assume (default is None).
    """
