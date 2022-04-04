Hyperparameter Tuning
=====================

Hyperparameters can be tuned with Hyperdrive using the following method:

1. Get or Create a Model using the ``hyperdrive`` interface.
2. Specify the hyperparameters and the values they can assume using a YAML file. 
3. Get or create a Hyperparameter Experiment using the ``hyperdrive`` interface.
4. Use the Hyperparameter Experiment object to tune the hyperparameters.

Get or Create a Model
---------------------
A Model should be created or loaded::

    my_model = hyperdrive.get_or_create_model(model_name, model_type)

A unique string should be passed to the ``model_name`` argument. The argument 
``model_type`` refers to the type of machine learning model being trained
(e.g. ``classification``, ``regression``or ``clustering``). Custom values (e.g.
``image_classification`` or ``natural_language_processing``) can also be provided.

Specifying the Hyperparameters
------------------------------

The hyperparameters can be specified using a YAML file. The hyperparameter search file 
has the following syntax::

    experiment_name: <EXPERIMENT_NAME>
    metric: <METRIC>
    direction: <DIRECTION>
    n_trials: <NUMBER_OF_TRIALS_TO_BE_RUN>
    models:
        <MODELNAME1>:
            <HYPERPARAMETER_1>: <NUMPY_COMMAND>
            <HYPERPARAMETER_2>: 
                - <VALUE_1>
                - <VALUE_2>
        <MODELNAME2>:
            <HYPERPARAMETER_3>: 
                distribution: <DISTRIBUTION>
                low: <LOW>
                high: <HIGH>
            <HYPERPARAMETER_4>: <VALUE>

Experiment names should be unique. 
``metric`` is a string referring to the fully qualified namespace name of an sklearn 
metric, for example, ``"sklearn.metrics.mean_squared_error"`` or 
``"sklearn.metrics.accuracy_score"``.
``direction`` is a string describing the direction in which the evaluating metric should
be optimized, for example, ``"minimize"`` or ``"maximize"``.

Numpy commands can be used to define a list of values that can be assumed by a 
hyperparameter. The commands ``numpy.linspace``, ``numpy.logspace``, ``numpy.arange`` 
and ``numpy.geomspace`` are supported.

For numeric hyperparameters, the ``low`` and ``high`` keys can be used to specify the 
minimum and maximum values that the hyperparameter can assume. The ``distribution`` key
can be used to specify the numeric distribution of the hyperparameter.

Get or Create a Hyperparameter Experiment
-----------------------------------------

Hyperparameter Experiments are a subclass of Experiments used for hyperparameter tuning.

Hyperparameter Experiments can be created using the ``hyperdrive`` interface::

    hp_exp = hyperdrive.get_or_create_hyperparameter_experiment(
        hyperparameter_search_file,
        features,
        target,
    )

If the hyperparameter experiment already exists it will be loaded rather than created.

An Experiment can also be extended to a Hyperparameter Experiment in order to support
hyperparameter tuning::

    hp_exp = hyperdrive.extend_to_hyperparameter_experiment(
        experiment,
        hyperparameter_search_file,
        features,
        target,
    )

The features and target use a 75-25 split for the train and test data. The models are 
trained on the train data and the metric score is evaluated on the test data. A ratio
from 0 to 1 can be provided to the optional ``test_size`` argument of both the methods
to change the split used.

Tune the Hyperparameters
------------------------

The hyperparameter tuning can be initiated using the following method::

    hp_exp.run_hyperparameter_search(n_trials)

If the optional ``n_trials`` argument is not provided, the number of trials specified
in the hyperparameter search file is used instead.

Each trial is logged as a Hyperdrive Run. The executed runs can be viewed by using the
``list_runs`` method::

    hp_exp.list_runs()

The Pandas DataFrame returned by this method specifies the models trained in each Run,
the values chosen for each hyperparameter and the metric score obtained by the model on
the test data. 

The trained models are logged as artifacts. The Run corresponding to the best metric 
score can be determined as with a standard Hyperdrive Experiment (see 
:ref:`standard-hyperdrive-experiment-label` for details) and the ``get_artifact`` 
method can be used to download the trained model Artifact from the Data Repo.

