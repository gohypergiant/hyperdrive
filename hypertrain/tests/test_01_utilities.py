import pandas as pd
import pytest

from hypertrain.utilities import (
    generate_ddl,
    generate_folder_name,
    get_column_names_and_types,
    load_yaml_file,
)


@pytest.fixture(scope="session")
def utility_test_df():
    yield pd.read_csv("tests/data/ht_agg.csv")


class TestUtilites:
    def test_generate_ddl(self, utility_test_df):

        column_names_and_types = get_column_names_and_types(utility_test_df)
        my_df_ddl = generate_ddl(column_names_and_types)

        assert (
            my_df_ddl == "_id text, "
            "avg(BMI) float, "
            "avg(active_heartrate) float, "
            "avg(resting_heartrate) float, "
            "avg(VO2_max) float"
        )

    def test_generate_folder_name(self):
        id_23 = 23
        name_23 = "air"
        folder_name_23 = generate_folder_name(trial_id=id_23, name=name_23)
        assert folder_name_23 == "000023-air-trial"

        id_45 = 45
        folder_name_45 = generate_folder_name(trial_id=id_45)
        assert folder_name_45 == "000045-trial"

    def test_get_column_names_and_types(self, utility_test_df):

        column_names_and_types = get_column_names_and_types(utility_test_df)
        assert column_names_and_types == [
            ("_id", "text"),
            ("avg(BMI)", "float"),
            ("avg(active_heartrate)", "float"),
            ("avg(resting_heartrate)", "float"),
            ("avg(VO2_max)", "float"),
        ]

    def test_load_yaml(self):
        assert (
            load_yaml_file("tests/data/rfc-nested-tugg-speedman.yaml")["study_name"]
            == "RandomForestClassifier_nested"
        )
