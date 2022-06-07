from hypertrain.flavor.sklearn import SklearnModelHandler
from hypertrain.utilities import zip_study


def import_trained_model(filename="", flavor="sklearn", shape=0):
    if flavor == "sklearn":
        sklearn_model = SklearnModelHandler._load_model(filename)
        onnx_model = SklearnModelHandler._convert(sklearn_model, shape)
        save_filename = "imported_" + flavor + "_model"
        SklearnModelHandler._save(onnx_model, save_filename)
        zip_study(save_filename)
