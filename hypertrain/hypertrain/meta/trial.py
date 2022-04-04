from pandas import DataFrame, Series
from numpy import ndarray
from ..exceptions import HyperparameterStudyError
from ..utilities import generate_ddl, generate_folder_name, get_column_names_and_types
from .interface.manifest import ManifestInterface

try:
    from torch import Tensor
except ImportError:
    pass


class TrialMeta:
    @classmethod
    def create(
        cls,
        metadata,
        metrics,
        hyperparameters,
        trained_model,
        preprocessor,
        preprocessor_shape,
        train_features,
        train_target,
        trial_id,
        trial_name,
        my_study_path,
    ):
        if hyperparameters is not None:
            for key, value in hyperparameters.items():
                if "numpy" in str(type(value)):
                    hyperparameters[key] = hyperparameters[key].item()
        if metrics is not None:
            for key, value in metrics.items():
                if "numpy" in str(type(value)):
                    metrics[key] = metrics[key].item()
        if metadata is not None:
            for key, value in metadata.items():
                if "numpy" in str(type(value)):
                    metadata[key] = metadata[key].item()

        trial_manifest = {
            "metrics": metrics,
            "metadata": metadata,
            "hyperparameters": hyperparameters,
        }

        feature_data_type = type(train_features)
        target_data_type = type(train_target)

        if feature_data_type in (DataFrame, Series):
            input_signature = cls._signature_creator_pandas(train_features)
        elif feature_data_type is ndarray:
            input_signature = cls._signature_creator_ndarray(train_features)
        elif feature_data_type is Tensor:
            input_signature = cls._signature_creator_tensor(train_features)
        else:
            raise HyperparameterStudyError(
                "Input signature inference not available for "
                f"data of type {feature_data_type}"
            )

        if target_data_type in (DataFrame, Series):
            output_signature = cls._signature_creator_pandas(train_target)
        elif target_data_type is ndarray:
            output_signature = cls._signature_creator_ndarray(train_target)
        elif target_data_type is Tensor:
            output_signature = cls._signature_creator_tensor(train_target)
        else:
            raise HyperparameterStudyError(
                "Output signature inference not available for "
                f"data of type {target_data_type}"
            )

        trial_manifest["input_signature"] = input_signature
        trial_manifest["output_signature"] = output_signature

        trial_folder = generate_folder_name(trial_id=trial_id, name=trial_name)

        ManifestInterface.write(
            manifest=trial_manifest,
            kind="trial",
            my_study_path=my_study_path,
            trial_folder=trial_folder,
        )

        trial_manifest["trained_model"] = trained_model

        return trial_manifest

    @classmethod
    def _signature_creator_ndarray(cls, ndarray):
        shape = ndarray.shape[1:]
        type = ndarray.dtype.type.__name__
        if shape == ():
            shape = [1]
        return f"ndarray: {type} {shape}"

    @classmethod
    def _signature_creator_pandas(cls, dataframe):
        column_names_and_types = get_column_names_and_types(DataFrame(dataframe))
        return generate_ddl(column_names_and_types)

    @classmethod
    def _signature_creator_tensor(cls, tensor):
        return str(tensor.size()[1:])
