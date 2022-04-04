import numpy as np
import pytest
import yaml
from optuna import Study
from sklearn.linear_model import ElasticNet
from xgboost.sklearn import XGBClassifier

from hypertrain.controllers.hyperparameter_study import (
    HyperparameterStudyController as hec,
)
from hypertrain.datamodel.hyperparameter_study import HyperparameterStudy


class TestHyperparameterStudy:
    def test_create_optuna_study(self, my_hp_study):
        my_hp_study._create_optuna_study()
        assert isinstance(my_hp_study.optuna_study, Study)
        assert len(my_hp_study.optuna_study.trials) == 0

    def test_suggest_hyperparameter_value(self, my_hp_study):
        trial = my_hp_study.optuna_study.ask()
        attribute_dict = {"low": 0, "high": 1}
        value = my_hp_study._suggest_hyperparameter_value(
            trial, "l1_ratio", attribute_dict
        )
        assert type(value) in [float, int]

    def test_build_hyperparameter_dictionary(self, my_hp_study):
        trial = my_hp_study.optuna_study.ask()
        hyperparameter_dict = {
            "alpha": "np.logspace(0,1,11)",
            "l1_ratio": {"low": 0, "high": 1},
            "selection": ["random", "cyclic"],
        }
        hyperparams_for_current_model = my_hp_study._build_hyperparameter_dictionary(
            trial, hyperparameter_dict
        )

        assert type(hyperparams_for_current_model) is dict
        assert set(["alpha", "l1_ratio", "selection"]) == set(
            hyperparams_for_current_model.keys()
        )

    def test_call_numpy(self, my_hp_study):

        assert np.array_equal(
            my_hp_study._call_numpy("np.linspace(2.0, 3.0, num=5, endpoint=False)"),
            np.linspace(2.0, 3.0, num=5, endpoint=False),
        )
        assert np.array_equal(
            my_hp_study._call_numpy("np.linspace(2.0, 3.0, num=5, retstep=True)")[0],
            np.linspace(2.0, 3.0, num=5, retstep=True)[0],
        )
        assert np.array_equal(my_hp_study._call_numpy("np.arange(3.0)"), np.arange(3.0))
        assert np.array_equal(
            my_hp_study._call_numpy("np.logspace(2.0, 3.0, num=4, base=2.0)"),
            np.logspace(2.0, 3.0, num=4, base=2.0),
        )
        assert np.array_equal(
            my_hp_study._call_numpy("np.logspace(0.1, 1, 10, endpoint=True)"),
            np.logspace(0.1, 1, 10, endpoint=True),
        )
        assert np.array_equal(
            my_hp_study._call_numpy("np.geomspace(1, 1000, num=3, endpoint=False)"),
            np.geomspace(1, 1000, num=3, endpoint=False),
        )

    def test_discrete_handler(self, my_hp_study):
        trial = my_hp_study.optuna_study.ask()
        attribute_dict = {
            "distribution": "discrete-uniform",
            "low": 10,
            "high": 25,
            "step": 5,
        }

        assert my_hp_study._discrete_handler(
            trial, "n_estimators_discrete", attribute_dict
        ) in [10, 15, 20, 25]

    def test_float_handler(self, my_hp_study):
        trial = my_hp_study.optuna_study.ask()

        attribute_dict1 = {
            "distribution": "uniform",
            "low": 0.001,
            "high": 0.25,
        }
        attribute_dict2 = {
            "distribution": "log-uniform",
            "low": 0.01,
            "high": 1,
        }

        float_value1 = my_hp_study._float_handler(
            trial, "learning_rate", attribute_dict1
        )
        float_value2 = my_hp_study._float_handler(trial, "C", attribute_dict2)

        assert (
            attribute_dict1["low"] <= float_value1
            and float_value1 <= attribute_dict1["high"]
        )
        assert (
            attribute_dict2["low"] <= float_value2
            and float_value2 <= attribute_dict2["high"]
        )
        assert type(float_value1) == float

    def test_get_study_metadata(self):
        yaml_file = "tests/data/xgboost-simple-jack.yaml"

        hp_study = HyperparameterStudy()
        hp_study._load_yaml_file_as_search_dictionary(yaml_file)
        hp_study._get_study_metadata()

        assert hp_study.study_name == "xgboost_test_3"
        assert hp_study.direction == "maximize"
        assert hp_study.metric == "sklearn.metrics.accuracy_score"
        assert "xgboost.XGBClassifier" in hp_study.models.keys()

    def test_get_trial_names(self, my_hp_study):

        my_hp_study._get_trial_names()
        assert isinstance(my_hp_study.trial_names, list)
        assert len(my_hp_study.trial_names) == 7

    def test_handle_numpy_value(self, my_hp_study):
        trial = my_hp_study.optuna_study.ask()

        assert my_hp_study._handle_numpy_value(
            trial, "alpha", "np.logspace(0,1,11)"
        ) in np.logspace(0, 1, 11)

        assert my_hp_study._handle_numpy_value(
            trial, "max_depth_np", "np.linspace(5,20,2)"
        ) in np.linspace(5, 20, 2)

    def test_handle_explicit_values(self, my_hp_study):
        trial = my_hp_study.optuna_study.ask()

        assert my_hp_study._handle_explicit_values(
            trial, "criterion", ["mse", "mae"]
        ) in ["mse", "mae"]

        assert my_hp_study._handle_explicit_values(
            trial, "n_estimators_rf", [2, 4, 6, 8]
        ) in [2, 4, 6, 8]

        assert (
            my_hp_study._handle_explicit_values(trial, "max_depth_rf_classifier", 5)
            == 5
        )

    def test_instantiate_sklearn_object(self, my_hp_study):
        object_name = "sklearn.linear_model.ElasticNet"
        sklearn_object = my_hp_study._instantiate_sklearn_object(object_name)

        assert isinstance(sklearn_object, ElasticNet)

    def test_instantiate_xgboost_object(self, my_hp_study):
        yaml_file = "tests/data/xgboost-simple-jack.yaml"

        hp_study = HyperparameterStudy()
        hp_study._load_yaml_file_as_search_dictionary(yaml_file)
        hp_study._get_study_metadata()
        hp_study._create_optuna_study()

        object_name = "xgboost.XGBClassifier"
        xgboost_object = hp_study._instantiate_xgboost_object(object_name)

        assert isinstance(xgboost_object, XGBClassifier)

    def test_int_handler(self, my_hp_study):
        trial = my_hp_study.optuna_study.ask()

        attribute_dict1 = {"distribution": "int", "low": 5, "high": 7}
        attribute_dict2 = {"distribution": "int-log-uniform", "low": 1, "high": 3}

        assert my_hp_study._int_handler(
            trial, "max_depth_rf_regressor", attribute_dict1
        ) in [5, 6, 7]
        assert my_hp_study._int_handler(
            trial, "alpha_elastic_net", attribute_dict2
        ) in [1, 2, 3]

    def test_nested_handler(self, my_hp_study):
        nested_search_dict = my_hp_study.search_dictionary["models"][
            "sklearn.ensemble.RandomForestClassifier"
        ]["n_estimators"]

        assert len(nested_search_dict) == 2

        nested_dict_item = next(
            (
                entry
                for entry in nested_search_dict
                if entry["distribution"] == "int-uniform"
            ),
            None,
        )

        assert nested_dict_item["low"] == 80
        assert nested_dict_item["high"] == 120
        assert len(nested_dict_item["criterion"]) == 2

    def test_create_hyperparameter_study(self, features_df, target_df):
        study = hec.create_hyperparameter_study(
            "tests/data/xgboost-simple-jack.yaml", features_df, target_df
        )

        assert isinstance(study, HyperparameterStudy)
        assert set(
            [
                "study_name",
                "direction",
                "metric",
                "model_flavor",
                "n_trials",
                "random_state",
                "data",
                "models",
            ]
        ).issubset(set(study.search_dictionary.keys()))
        assert isinstance(study.optuna_study, Study)
        assert study.n_trials == 7
        assert study.model_flavor == "xgboost"
        assert study.direction == "maximize"
        assert len(study.train_features) == 112
        assert len(study.train_target) == 112
        assert len(study.test_features) == 38
        assert len(study.test_target) == 38

    def test_objective(self, features_df, target_df):
        study = hec.create_hyperparameter_study(
            "tests/data/xgboost-simple-jack.yaml", features_df, target_df
        )
        trial = study.optuna_study.ask()

        assert type(study._objective(trial)) is np.float64
        n = len(study.optuna_study.trials)
        assert "trained_model" in study.optuna_study.trials[n - 1].user_attrs

    def test_parse_np_argument(self, my_hp_study):
        assert my_hp_study._parse_np_argument("True")
        assert not my_hp_study._parse_np_argument("False")
        assert my_hp_study._parse_np_argument("None") is None
        assert my_hp_study._parse_np_argument("int") == int
        assert my_hp_study._parse_np_argument("float") == float
        assert my_hp_study._parse_np_argument("22.5") == 22.5
        assert my_hp_study._parse_np_argument("5") == 5

    def test_load_yaml_file_as_search_dictionary(self, my_hp_study):
        assert type(my_hp_study.search_dictionary) is dict
        assert set(
            ["study_name", "direction", "n_trials", "metric", "models"]
        ).issubset(set(my_hp_study.search_dictionary.keys()))

    def test_prep_xgboost_features(self, features_df):
        features = hec.prep_xgboost_features(features_df)

        assert isinstance(features, np.ndarray)
        assert features.shape == (150, 4)

    def test_verify_xgboost_yaml_params(self):
        with open("tests/data/xgboost-not-supp-models.yaml", "r") as stream:
            xgboost_not_supp_yaml = yaml.safe_load(stream)
        with pytest.raises(ValueError) as e1:
            hec.verify_xgboost_yaml_params(xgboost_not_supp_yaml)
        assert str(e1.type) == "<class 'ValueError'>"
        assert (
            str(e1.value)
            == "The hyperparameter yaml file contains xgboost models that "
            "are not supported by ONNX. At this time, ONNX only supports "
            "xgboost.XGBClassifer and xgboost.XGBRegressor models."
        )

        with open("tests/data/xgboost-classif-gblinear.yaml", "r") as stream:
            xgboost_classif_gblinear_yaml = yaml.safe_load(stream)
        with pytest.raises(ValueError) as e2:
            hec.verify_xgboost_yaml_params(xgboost_classif_gblinear_yaml)
        assert str(e2.type) == "<class 'ValueError'>"
        assert (
            str(e2.value) == "ONNX does not currently support the usage of a "
            "xgboost.XGBClassifier model with a gblinear booster. "
            "Please modify the hyperparameter yaml file."
        )
