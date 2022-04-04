class ArtifactDoc:
    __doc__ = (
        """
    A helper class for loading Study Artifacts.

    Attributes
    ----------
    id: str
        The id of the experiment (default is None).
    run_id: str
        The id of the parent run (default is None).
    artifact: varies
        The artifact object in memory (default is None).
    artifact_type: {"data_exploration_file", "model_training_file", """
        """"data_preparation_file", "trained_model", "preprocessor","baseline"}
        The type of artifact to be saved (default is None).
    file_key: str
        The key in object storage for the artifact (default is None).
    input_signature: str
        input signature for a feature vector passed to the model (default is None).
    output_signature: str
        output signature for a target value returned by the model (default is None).
    """
    )
