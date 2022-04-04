import onnx
import onnxmltools
from skl2onnx.common.data_types import FloatTensorType
from onnxconverter_common.onnx_ex import DEFAULT_OPSET_NUMBER


class XGBoostModelHandler:
    @classmethod
    def _convert(cls, trained_model, train_shape):
        initial_types = [("X", FloatTensorType([None, train_shape]))]
        onnx_model = onnxmltools.convert.convert_xgboost(
            trained_model,
            initial_types=initial_types,
            target_opset=DEFAULT_OPSET_NUMBER,
        )
        return onnx_model

    @classmethod
    def _save(cls, onnx_model, local_artifact_path):
        onnx.save(onnx_model, local_artifact_path)
