SUPPORTED_MODEL_FLAVORS = ["automl"]


def create_hyperpack(trained_model=None, model_flavor: str = None):
    if trained_model is None:
        raise ValueError("You must pass in a trained model.")
    elif model_flavor is None:
        print("You must specify a model flavor. Supported model flavors are:\n")
        print("\n".join(map(str, SUPPORTED_MODEL_FLAVORS)))
        raise ValueError("Please specify a model flavor.")
    print("ahoy environs!")
