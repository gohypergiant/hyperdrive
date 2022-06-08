import os

from hypertrain.flavor.sklearn import SklearnModelHandler
from hypertrain.utilities import zip_study


def import_trained_model(filename="", flavor="sklearn", shape=0):
    save_filename = "imported_" + flavor + "_model"
    save_dir = f"/home/jovyan/_jobs/{save_filename}/"
    os.makedirs(save_dir, exist_ok=True)

    if flavor == "sklearn":
        sklearn_model = SklearnModelHandler._load_model(filename)
        onnx_model = SklearnModelHandler._convert(sklearn_model, shape)
        save_path = save_dir + save_filename
        SklearnModelHandler._save(onnx_model, save_path)
        zip_study(save_dir)
