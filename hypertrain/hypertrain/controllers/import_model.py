from hypertrain.flavor.sklearn import SklearnModelHandler


def import_trained_model(filename="", flavor="sklearn"):
    if flavor == "sklearn":
        sklearn_model = SklearnModelHandler._load_model(filename)
        shape = 8
        onnx_model = SklearnModelHandler._convert(sklearn_model, shape)
        save_filename = "imported_" + flavor + "_model"
        SklearnModelHandler._save(onnx_model, save_filename)
