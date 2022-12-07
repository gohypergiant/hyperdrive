import onnx
import os
import tf2onnx


def tensorflow_onnx_export(model, trial_path: str):
    """Saves tensorflow model to ONNX format
    Args:
        model: pretrained model
        trial_path: dir path to the trial
        train_shape_cols: the number of columns in the training dataset used to
                          train the pretrained model
        ml_task : the type of machine learning task
    """
    onnx_model, _ = tf2onnx.convert.from_keras(model)
    model_path = os.path.join(trial_path, "trained_model")
    onnx.save(onnx_model, model_path)
