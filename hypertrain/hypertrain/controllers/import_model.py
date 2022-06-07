from hypertrain.flavor.sklearn import SklearnModelHandler


def import_trained_model(flavor="dingo"):
    if flavor == "sklearn":
        sklearn_model = SklearnModelHandler._load_model("finalized_model.sav")
        shape = 8
        onnx_model = SklearnModelHandler._convert(sklearn_model, shape)
        SklearnModelHandler._save(onnx_model, "trained_model")
        print("Model flavor is:", flavor)
