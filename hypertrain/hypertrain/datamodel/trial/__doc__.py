class TrialDoc:
    __doc__ = """
    A container for specific runs of a model training study.


    Attributes
    ----------
    id: str
        The id of the run.
    model_experiment_id: str
        The model experiment id.
    metadata: dictionary, optional
        Metadata to be stored in the run (default is None).
    hyperparameters: dictionary, optional
        The hyperparameters chosen for the run (default is None).
    metrics: dictionary, optional
        Performance metrics for the training run.
    trained_model: Model object, optional
        Trained model to which the run belongs (default is None).
    baseline: str
        A dictionary of the expected value and standard deviation
        for features and target to be monitored by Hyperdrive (default is None).
    preprocessor: Sklearn object, optional
        The trained sklearn preprocessor for data preprocessor (default is None).
    preprocessor_shape: int, optional
        The number of inputs (default is None).
    data_exploration_file: str
        Path of the eda file (default is None).
    data_preparation_file: str
        Path of the data engineering file (default is None).
    model_training_file: str
        Path of the model training file (default is None).
    """
