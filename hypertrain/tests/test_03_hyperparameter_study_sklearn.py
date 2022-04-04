import numpy as np
from hypertrain.controllers.hyperparameter_study import (
    HyperparameterStudyController as hec,
)


class TestHyperparameterStudySklearn:
    def test_run_hyperparameter_search_xgboost(self, features_df, target_df):
        study = hec.create_hyperparameter_study(
            "tests/data/xgboost-simple-jack.yaml", features_df, target_df
        )
        study.run_hyperparameter_search()

        assert len(study.optuna_study.trials) == 7

        best_trial_trained_model = study.optuna_study.best_trial.user_attrs[
            "trained_model"
        ]
        assert (
            study.optuna_study.best_params.get("algorithm") == "xgboost.XGBClassifier"
        )
        assert len(best_trial_trained_model.predict(study.test_features)) == 38

    def test_run_hyperparameter_search_lasso(self, features_df_lasso, target_df_lasso):
        study = hec.create_hyperparameter_study(
            "tests/data/lasso-alpa-chino.yaml", features_df_lasso, target_df_lasso
        )

        study.run_hyperparameter_search()

        assert len(study.optuna_study.trials) == 10

        best_trial_trained_model = study.optuna_study.best_trial.user_attrs[
            "trained_model"
        ]
        predictions = best_trial_trained_model.predict(study.test_features)
        assert (
            study.optuna_study.best_params.get("algorithm")
            == "sklearn.linear_model.Lasso"
        )
        assert len(predictions) == 57
        assert predictions.dtype.type == np.float64
        assert True
