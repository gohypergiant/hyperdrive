import joblib
import onnx
import skl2onnx
from skl2onnx.common.data_types import FloatTensorType


class SklearnModelHandler:
    @classmethod
    def _convert(cls, trained_model, train_shape):
        initial_types = [("X", FloatTensorType([None, train_shape]))]
        return skl2onnx.convert_sklearn(
            model=trained_model, initial_types=initial_types
        )

    @classmethod
    def _save(cls, onnx_model, local_artifact_path):
        onnx.save(onnx_model, local_artifact_path)

    @classmethod
    def _load_model(cls, pickle_file):
        with open(pickle_file, "rb") as scaler_file:
            sklearn_pre = joblib.load(scaler_file)
        return sklearn_pre
