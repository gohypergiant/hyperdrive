import os
import pwd
import shutil
import warnings

from pandas import DataFrame
from sklearn.model_selection import train_test_split

from ..datamodel.hyperparameter_study import HyperparameterStudy
from ..exceptions import HyperparameterStudyError
from ..utilities import load_yaml_file


class HyperparameterStudyController:
    @classmethod
    def create_hyperparameter_study(
        self,
        hyperparameter_search_file: str,
        features: DataFrame,
        target: DataFrame,
        study_path: str = "my_study",
    ):
        """Creates a Hyperparameter experiment.

        Parameters
        ----------
        hyperparameter_search_file : str
            Path of the YAML file containing the hyperparameter search information.
        features : pd.DataFrame
            Features dataframe.
        target : pd.DataFrame
            Target dataframe.

        Returns
        -------
        Object of class HyperparameterStudy with specific attributes.
        """
        yaml_dictionary = load_yaml_file(hyperparameter_search_file)

        if "n_trials" not in yaml_dictionary.keys():
            warnings.warn(
                "The number of trials (n_trials) is not specified in the YAML. Therefore, n_trials will be set to 100."
            )

        supported_model_flavors = [
            "sklearn",
            "lightgbm",
            "xgboost",
            "tensorflow",
            "pytorch",
        ]

        if "model_flavor" not in yaml_dictionary.keys():
            raise HyperparameterStudyError(
                f"You must specify a model_flavor in the YAML. Supported types of model_flavor include: {supported_model_flavors}"
            )
        else:
            model_flavor = yaml_dictionary["model_flavor"]

        if model_flavor not in supported_model_flavors:
            raise HyperparameterStudyError(
                f"The model flavor '{model_flavor}' is not currently supported. Please use one of the following: {supported_model_flavors}"
            )

        if model_flavor == "xgboost":
            self.verify_xgboost_yaml_params(yaml_dictionary)
            features = self.prep_xgboost_features(features)

        hp_search_exp = HyperparameterStudy()
        hp_search_exp._load_yaml_file_as_search_dictionary(hyperparameter_search_file)
        hp_search_exp._get_study_metadata()
        hp_search_exp._create_optuna_study()
        hp_search_exp._get_trial_names()

        if "test_size" in yaml_dictionary.keys():
            test_size = yaml_dictionary["test_size"]
        else:
            test_size = 0.25

        if "random_state" in yaml_dictionary.keys():
            random_state = yaml_dictionary["random_state"]
        else:
            random_state = 0

        (train_features, test_features, train_target, test_target,) = train_test_split(
            features, target, test_size=test_size, random_state=random_state
        )

        hp_search_exp.train_features = train_features
        hp_search_exp.train_target = train_target
        hp_search_exp.test_features = test_features
        hp_search_exp.test_target = test_target
        hp_search_exp.train_shape = train_features.shape

        current_user = pwd.getpwuid(os.geteuid()).pw_name
        if current_user == "jovyan":
            hp_search_exp.my_study_path = (
                f"/home/jovyan/_jobs/{study_path}/{study_path}.hyperpack"
            )
        else:
            hp_search_exp.my_study_path = (
                os.getcwd() + f"/{study_path}/{study_path}.hyperpack"
            )

        try:
            os.makedirs(hp_search_exp.my_study_path, exist_ok=False)
        except FileExistsError:
            i = 1
            while os.path.exists(f"{hp_search_exp.my_study_path}_{str(i)}"):
                i += 1
            hp_search_exp.my_study_path = hp_search_exp.my_study_path + "_" + str(i)
            os.makedirs(hp_search_exp.my_study_path, exist_ok=False)

        shutil.copyfile(
            hyperparameter_search_file, f"{hp_search_exp.my_study_path}/_hyperpack.yaml"
        )

        return hp_search_exp

    @classmethod
    def verify_xgboost_yaml_params(cls, yaml_dictionary: dict):
        """Verifies XGBoost YAML parameters

        Parameters
        ----------
        yaml_dictionary : dict
            Python dictionary containing the hyperparameter search information.

        Returns
        ----------
        None

        Notes regarding verification of XGBoost YAML parameters
        ----------
        For XGBoost, currently ONNX only supports "xgboost.XGBClassifier" and
            "xgboost.XGBRegressor".

        For XGBoost, when using the "xgboost.XGBClassifier", ONNX currently does
            not support the "gblinear" booster.
        """

        onnx_supported_xgboost_models = [
            "xgboost.XGBClassifier",
            "xgboost.XGBRegressor",
        ]

        if not all(
            model_type in onnx_supported_xgboost_models
            for model_type in yaml_dictionary["models"].keys()
        ):
            raise ValueError(
                "The hyperparameter yaml file contains xgboost models that "
                "are not supported by ONNX. At this time, ONNX only supports "
                "xgboost.XGBClassifer and xgboost.XGBRegressor models."
            )

        if "xgboost.XGBClassifier" in yaml_dictionary["models"].keys():
            if (
                "booster" in yaml_dictionary["models"]["xgboost.XGBClassifier"].keys()
                and "gblinear"
                in yaml_dictionary["models"]["xgboost.XGBClassifier"]["booster"]
            ):
                raise ValueError(
                    "ONNX does not currently support the usage of a "
                    "xgboost.XGBClassifier model with a gblinear booster. "
                    "Please modify the hyperparameter yaml file."
                )

    @classmethod
    def prep_xgboost_features(cls, features: DataFrame):
        """Verifies XGBoost YAML parameters

        Parameters
        ----------
        features : pd.DataFrame
            Features dataframe.

        Returns
        ----------
        A numpy.ndarray of the features.
        """
        if isinstance(features, DataFrame):
            features_column_names = features.columns.values
            warnings.warn(
                f"Features DataFrame automatically converted to ndarray with these columns: {features_column_names}"
            )
            features = features.values
        return features
