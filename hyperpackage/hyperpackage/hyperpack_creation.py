import os
import torch
from datetime import datetime
from hyperpackage.flavor.pytorch import torch_onnx_export
from hyperpackage.utilities import (
    generate_folder_name,
    write_json,
    write_yaml,
    zip_study,
)

SUPPORTED_MODEL_FLAVORS = ["automl"]


def create_hyperpack(
    trained_model=None, model_flavor: str = None, num_train_columns: int = 0
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
        num_train_columns: the number of columns in the training dataset used to
                           train the pretrained model
    """
    curr_dir = os.getcwd()

    print("*** Verifying trained_model and model_flavor args ***")
    verify_args(model=trained_model, flavor=model_flavor)

    print("*** Loading the trained model **")
    loaded_model = load_trained_model(model=trained_model)

    # checking the num_train_columns arg. We did not check this arg in the
    # verify_args func because if an automl model is passed in and the user
    # doesn't pass in the number of training dataset columns, we
    # can get the number of training dataset columns for the user
    if model_flavor == "automl" and num_train_columns == 0:
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
    )

    print("*** Creating _study.json ***")
    create_study_json(hyperpack_path=hyperpack_path, best_trial=best_trial_folder_name)

    print("*** Creating study.yaml ***")
    create_study_yaml(current_dir=curr_dir, name=model_flavor)

    print("*** Zipping hyperpack folder to {}.zip ***".format(hyperpack_path))
    zip_study(folder_path=hyperpack_path)

    print("*** Hyperpack created! ***")


def verify_args(model, flavor: str):
    """Verifies the model and flavor args that are passed to the create_hyperpack
       function
    Args:
        model: pretrained model. Can be either a string, which is a path
               to a pickled model of type <class 'neural_network.network.Network'>,
               or to an object in memory of type
               <class 'neural_network.network.Network'>
        flavor: library/package used to build pretrained model
    """
    supported_flavors = "\n".join(map(str, SUPPORTED_MODEL_FLAVORS))
    if model is None:
        raise TypeError("You must pass in a trained model.")
    elif isinstance(model, str):
        if not os.path.exists(model):
            raise FileNotFoundError("No file could be found at {}.".format(model))

    if flavor is None:
        raise TypeError(
            "You must specify a model flavor. Supported model flavors are:\n{}".format(
                supported_flavors
            )
        )
    elif flavor not in SUPPORTED_MODEL_FLAVORS:
        raise TypeError(
            "You have selected a model flavor that is currently not supported. Supported model flavors are:\n{}".format(
                supported_flavors
            )
        )


def load_trained_model(model):
    """Loads the pretrained model
    Args:
        model: pretrained model. Can be either a string, which is a path
               to a pickled model of type <class 'neural_network.network.Network'>,
               or to an object in memory of type
               <class 'neural_network.network.Network'>
    Returns: the loaded model
    """
    if isinstance(model, str):
        try:
            the_model = torch.load(model)
        except Exception:
            print("Error while attempting to load torch model.")
    elif str(type(model)) == "<class 'neural_network.network.Network'>":
        the_model = model
    else:
        raise TypeError("The model type you have passed in is currently not supported.")

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


def save_best_model_to_onnx(model, flavor: str, save_path: str, num_cols: int):
    """Saves the model to ONNX format
    Args:
        model: pretrained model object
        flavor: library/package used to build pretrained model
        save_path: path where ONNX model will be saved
        num_cols: number of training dataset columns
    """
    if flavor == "automl":
        torch_onnx_export(model=model, trial_path=save_path, train_shape_cols=num_cols)


def create_study_json(hyperpack_path: str, best_trial: str):
    """Creates the "_study.json" file
    Args:
        hyperpack_path: path to hyperpack folder
        best_trial: name of best trial
    """
    created_time = datetime.now().strftime("%Y-%m-%d %H:%M")
    study_json_dict = {"best_trial": best_trial, "created_at": created_time}
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
