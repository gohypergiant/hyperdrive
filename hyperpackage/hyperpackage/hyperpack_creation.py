import os
import tensorflow as tf
import torch

from datetime import datetime
from xgboost import XGBClassifier

from hyperpackage.flavor.pytorch import torch_onnx_export
from hyperpackage.flavor.tensorflow import tensorflow_onnx_export
from hyperpackage.flavor.xgboost import xgboost_onnx_save
from hyperpackage.utilities import (
    generate_folder_name,
    write_json,
    write_yaml,
    zip_study,
)

SUPPORTED_MODEL_FLAVORS = ["automl", "tensorflow", "xgboost"]

SUPPORTED_ML_TASKS = [
    "regression",
    "binary_classification",
    "multi_class_classification",
    "univariate_timeseries",
]


def create_hyperpack(
    trained_model=None,
    model_flavor: str = None,
    ml_task: str = None,
    num_train_columns: int = 0,
):
    """Entrypoint function for the hyperpackage package. The user will call this
       function
    Args:
        trained_model: pretrained model. Can be either a string, which is a path
                       to a pickled model of type <class 'neural_network.network.Network'>,
                       or to an object in memory of type <class 'neural_network.network.Network'>
        model_flavor: library/package used to build pretrained model. Currently,
                      we only support "automl", but examples of other flavors
                      that will eventually be supported include "sklearn",
                      "pytorch", and "xgboost"
        ml_task: the type of machine learning task, e.g., regression,
                 binary_classification, or multi_class_classification
        num_train_columns: the number of columns in the training dataset used to
                           train the pretrained model
    """
    curr_dir = os.getcwd()

    print("*** Verifying trained_model and model_flavor args ***")
    verify_args(model=trained_model, flavor=model_flavor, task=ml_task)

    print("*** Loading the trained model ***")
    loaded_model = load_trained_model(model=trained_model, flavor=model_flavor)

    # checking the num_train_columns arg. We did not check this arg in the
    # verify_args func because if an automl model is passed in and the user
    # doesn't pass in the number of training dataset columns, we
    # can get the number of training dataset columns for the user
    if model_flavor == "automl" and num_train_columns == 0:
        if ml_task == "univariate_timeseries":
            num_train_columns = loaded_model[1].fc1.in_features
        else:
            num_train_columns = loaded_model.fc1.in_features
    elif num_train_columns == 0:
        raise ValueError(
            "Please pass in the number of columns in the training dataset."
        )

    # make the hyperpack directory (e.g., automl.hyperpack)
    hyperpack_path = make_hyperpack_path(base_dir=curr_dir, name=model_flavor)
    try:
        os.makedirs(hyperpack_path, exist_ok=False)
    except FileExistsError:
        i = 1
        while os.path.exists(f"{hyperpack_path}_{str(i)}"):
            i += 1
        hyperpack_path = hyperpack_path + "_" + str(i)
        os.makedirs(hyperpack_path, exist_ok=False)

    # make the directory for the best trial
    best_trial_name = "adventurous"
    best_trial_folder_name = generate_folder_name(name=best_trial_name)
    best_trial_path = os.path.join(hyperpack_path, best_trial_folder_name)
    os.makedirs(best_trial_path, exist_ok=True)

    print("*** Saving the model to ONNX format ***")
    save_best_model_to_onnx(
        model=loaded_model,
        flavor=model_flavor,
        save_path=best_trial_path,
        num_cols=num_train_columns,
        ml_task=ml_task,
    )

    print("*** Creating _study.json ***")
    create_study_json(
        hyperpack_path=hyperpack_path,
        best_trial=best_trial_folder_name,
        ml_task=ml_task,
        flavor=model_flavor
    )

    print("*** Creating study.yaml ***")
    create_study_yaml(current_dir=curr_dir, name=model_flavor)

    print("*** Zipping hyperpack folder to {}.zip ***".format(hyperpack_path))
    zip_study(folder_path=hyperpack_path)

    print("*** Hyperpack created! ***")


def verify_args(model, flavor: str, task: str):
    """Verifies the model and flavor args that are passed to the create_hyperpack
       function
    Args:
        model: pretrained model. Can be either a string or an object
        flavor: library/package used to build pretrained model
        task: type of machine learning task, e.g., regression, binary_classification,
              multi_class_classification
    """
    supported_flavors = "\n".join(map(str, SUPPORTED_MODEL_FLAVORS))
    supported_tasks = "\n".join(map(str, SUPPORTED_ML_TASKS))
    if model is None:
        raise TypeError("You must pass in a trained model.")
    elif isinstance(model, str):
        if not os.path.exists(model):
            raise FileNotFoundError(
                "No file could be found at {}.".format(model))

    if flavor is None:
        raise TypeError(
            "You must specify a model flavor. Supported model flavors are:\n{}".format(
                supported_flavors
            )
        )
    elif flavor not in SUPPORTED_MODEL_FLAVORS:
        raise TypeError(
            "You have specified a model flavor, {}, that is currently not supported. Supported model flavors are:\n{}".format(
                flavor, supported_flavors
            )
        )

    if task is None:
        raise TypeError(
            "You must specify a machine learning task. Supported tasks are:\n{}".format(
                supported_tasks
            )
        )
    elif task not in SUPPORTED_ML_TASKS:
        raise TypeError(
            "You have specified a task, {}, that is currently not supported. Supported tasks are:\n{}".format(
                task, supported_tasks
            )
        )


def load_trained_model(model, flavor: str = None):
    """Loads the pretrained model
    Args:
        model: pretrained model. For automl, this can be either a string,
               which is a path to a pickled model of type
               <class 'neural_network.network.Network'>, or to an object in
               memory of type <class 'neural_network.network.Network'>. For
               tensorflow, this can be either a string, which is a path to the
               directory for a model that has been saved in the "SavedModel"
               format (so it'll have an "assets" folder, a "variables" folder,
               and the model with a "saved_model.pb" file name), or to an object
               in memory of a type that looks like <class 'keras.engine.[MORE_STUFF]'
        flavor: library/package used to build pretrained model
    Returns: the loaded model
    """
    if flavor == "automl":
        if isinstance(model, str):
            try:
                the_model = torch.load(model)
            except Exception:
                print("Error while attempting to load torch model.")
        elif str(type(model)) == "<class 'neural_network.network.Network'>":
            the_model = model
        elif isinstance(model, dict):
            try:
                for v in model.values():
                    if str(type(v)) != "<class 'neural_network.network.Network'>":
                        raise TypeError(
                            "The dictionary of models contain a model type that is currently not supported.")
                the_model = model
            except Exception:
                print("Error while attempting to load torch model.")
        else:
            raise TypeError(
                "The model type you have passed in is currently not supported."
            )
    elif flavor == "tensorflow":
        if isinstance(model, str):
            try:
                the_model = tf.keras.models.load_model(model)
            except Exception:
                print("Error while attempting to load tensorflow model.")
        elif "keras.engine" in str(type(model)):
            the_model = model
        else:
            raise TypeError(
                "The model type you have passed in is currently not supported."
            )
    elif flavor == "xgboost":
        if isinstance(model, str):
            try:
                the_model = XGBClassifier()
                the_model.load_model(model)
            except Exception:
                print("Error while attempting to load xgboost model.")
        elif "xgboost" in str(type(model)):
            the_model = model
        else:
            raise TypeError(
                "The model type you have passed in is currently not supported."
            )

    return the_model


def make_hyperpack_path(base_dir: str, name: str) -> str:
    """Makes the directory path to the hyperpack folder
    Args:
        base_dir: string representing the base directory of the hyperpack folder
        name: string to be used as part of the name of the hyperpack folder
    Returns: path to the hyperpack folder
    """
    hyperpack_folder_name = name + ".hyperpack"
    path = os.path.join(base_dir, hyperpack_folder_name)
    return path


def save_best_model_to_onnx(model, flavor: str, save_path: str, num_cols: int, ml_task: str):
    """Saves the model to ONNX format
    Args:
        model: pretrained model object
        flavor: library/package used to build pretrained model
        save_path: path where ONNX model will be saved
        num_cols: number of training dataset columns
        ml_task: the type of machine learning task
    """
    if flavor == "automl":
        torch_onnx_export(model=model, trial_path=save_path,
                          train_shape_cols=num_cols, ml_task=ml_task)
    elif flavor == "tensorflow":
        tensorflow_onnx_export(model=model, trial_path=save_path)
    elif flavor == "xgboost":
        xgboost_onnx_save(model=model, trial_path=save_path,
                          train_shape_cols=num_cols)


def create_study_json(hyperpack_path: str, best_trial: str, ml_task: str, flavor: str):
    """Creates the "_study.json" file
    Args:
        hyperpack_path: path to hyperpack folder
        best_trial: name of best trial
        ml_task: type of machine learning task
        flavor: library/package used to build pretrained model.
                Examples include automl, pytorch, and tensorflow
    """
    created_time = datetime.now().strftime("%Y-%m-%d %H:%M")
    study_json_dict = {
        "best_trial": best_trial,
        "created_at": created_time,
        "ml_task": ml_task,
        "model_flavor": flavor
    }
    study_json_path = os.path.join(hyperpack_path, "_study.json")
    write_json(dictionary=study_json_dict, json_file_path=study_json_path)


def create_study_yaml(current_dir: str, name: str):
    """Creates the "study.yaml" file
    Args:
        current_dir: path to current directory
        name: string to be used for both project_name and study_name
    """
    study_yaml_dict = {"project_name": name, "study_name": name}
    study_yaml_path = os.path.join(current_dir, "study.yaml")
    write_yaml(dictionary=study_yaml_dict, yaml_file_path=study_yaml_path)
