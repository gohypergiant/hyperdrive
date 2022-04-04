import pandas as pd
import pytest
from sklearn import datasets

from hypertrain.datamodel.hyperparameter_study import HyperparameterStudy


@pytest.fixture(scope="session")
def hyperparameter_search_file():
    return "tests/data/rfc-nested-tugg-speedman.yaml"


@pytest.fixture(scope="session")
def my_hp_study(hyperparameter_search_file):
    hp_study = HyperparameterStudy()
    hp_study._load_yaml_file_as_search_dictionary(hyperparameter_search_file)
    hp_study._get_study_metadata()
    hp_study._create_optuna_study()

    return hp_study


@pytest.fixture(scope="session")
def features_df():
    data = datasets.load_iris()
    features = pd.DataFrame(data["data"])

    return features


@pytest.fixture(scope="session")
def target_df():
    data = datasets.load_iris()
    target = pd.DataFrame(data["target"])

    return target


@pytest.fixture(scope="session")
def features_df_lasso():
    data = pd.read_csv("tests/data/run_rating.csv")
    features = data.drop("runner_rating", axis=1)

    return features


@pytest.fixture(scope="session")
def target_df_lasso():
    data = pd.read_csv("tests/data/run_rating.csv")
    target = data[["runner_rating"]]

    return target
