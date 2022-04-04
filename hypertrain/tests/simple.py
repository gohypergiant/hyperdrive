import pandas as pd
import yaml
from sklearn.preprocessing import LabelEncoder
from hddataclient import DataRepo
from hypertrain.controllers.hyperparameter_study import (
    HyperparameterStudyController as hec,
)

features = "tests/data/ht_agg.csv"
target = "tests/data/user_data.csv"
study_yaml = "tests/data/rfc-nested-tugg-speedman.yaml"


with open(study_yaml) as fh:
    my_study_yaml = yaml.safe_load(fh)


X_df = pd.read_csv(features)
y_df = pd.read_csv(target)


drop_cols = list(X_df.dtypes.index[X_df.dtypes == "object"])
X = X_df.drop(drop_cols, axis=1)

response_variable = my_study_yaml["data"]["response_variable"]
y_resp = y_df[response_variable]
le = LabelEncoder()
y = le.fit_transform(y_resp)


hyp_study_rf_nested = hec.create_hyperparameter_study(study_yaml, X, y)
hyp_study_rf_nested.run_hyperparameter_search()
