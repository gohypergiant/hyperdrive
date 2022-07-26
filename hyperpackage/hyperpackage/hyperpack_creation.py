def create_hyperpack(trained_model=None, model_flavor: str = None):
    if trained_model is None:
        raise ValueError("You must pass in a trained model.")
    elif model_flavor is None:
        raise ValueError("You must specify a model flavor.")
    print("ahoy environs!")
