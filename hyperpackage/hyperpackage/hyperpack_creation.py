SUPPORTED_MODEL_FLAVORS = ["automl"]


def create_hyperpack(trained_model=None, model_flavor: str = None):
    if trained_model is None:
        raise TypeError("You must pass in a trained model.")
    elif model_flavor is None:
        supported_flavors = "\n".join(map(str, SUPPORTED_MODEL_FLAVORS))
        raise TypeError(
            "You must specify a model flavor. Supported model flavors are:\n{}".format(
                supported_flavors
            )
        )
    print("ahoy environs!")
