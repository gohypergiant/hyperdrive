import os


class TrialLoggingController:
    def log_metadata(self, metadata):
        """Logs the metadata attributes of the current Trial.

        Parameters
        ----------
        metadata : dictionary
            Metadata attributes for the current Trial.

        """
        self.metadata = metadata

    def log_metrics(self, metrics):
        """Logs the metrics attributes of the current Trial.

        Parameters
        ----------
        metrics : dictionary
            Performance metrics for the trial object.

        """
        self.metrics = metrics

    def log_model(self, model):
        """Logs the model of the current Trial.

        Parameters
        ----------
        model : trained model object
            Trained model for the Trial object.

        """
        self.trained_model = model

    def log_preprocessor(self, preprocessor, preprocessor_shape):
        """Logs the preprocessor of the current Trial.

        Parameters
        ----------
        preprocessor : preprocessor object
            Trained preprocessor for the Trial object.

        """
        self.preprocessor = preprocessor
        self.preprocessor_shape = preprocessor_shape

    def log_hyperparameters(self, hyperparameters=None):
        """Logs the hyperparameters attribute of the current Trial.

        Parameters
        ----------
        hyperparameters : dictionary
            Hyperparameters for the Trial object.

        """
        self.hyperparameters = hyperparameters

    def _process_if_not_notebook(self, experiment, filename, notebook_name):
        file = getattr(experiment, filename, None)
        if file is not None:
            if os.path.splitext(file)[-1] == ".py" and os.path.isfile(file):
                setattr(
                    experiment, filename, self._file_to_notebook(notebook_name, file)
                )
