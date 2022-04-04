import onnx
import tensorflow as tf
import tf2onnx.convert

OPSET = 10


class TensorflowModelHandler:
    @classmethod
    def _convert(cls, trained_model, train_shape):
        onnx_model, _ = tf2onnx.convert.from_keras(trained_model, opset=OPSET)
        return onnx_model

    @classmethod
    def _save(cls, onnx_model, local_artifact_path):
        onnx.save(onnx_model, local_artifact_path)

    @classmethod
    def _load_model(cls, json_file, h5_file):
        with open(json_file) as model_architecture:
            trained_model = tf.keras.models.model_from_json(model_architecture.read())
        trained_model.load_weights(h5_file)
        return trained_model
