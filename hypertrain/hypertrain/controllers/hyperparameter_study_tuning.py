import logging
import random
from uuid import uuid4

import numpy as np
import optuna

from ..exceptions import HyperparameterStudyError
from ..meta.interface.manifest import ManifestInterface
from ..utilities import (
    generate_folder_name,
    load_yaml_file,
    write_json_to_local,
    zip_study,
)
from .trial_log_callback import TrialLogCallback


class HyperparameterStudyTuningController:
    def _objective(self, trial):
        """Defines the objective function used by Optuna's optimize method. Defines
        the model, trains it on the training data and evaluates it on the validation
        data.

        Parameters
        ----------
        trial : optuna.trial._trial.Trial
            A single execution of the objective function which is instantiated each
            time the objective function is called.

        Returns
        -------
        metric : float
            Evaluation metric of the model on the validation data.

        """
        model_names = list(self.models.keys())
        algo_name = trial.suggest_categorical("algorithm", model_names)

        hyperparameter_dict = self.models[algo_name]
        hyperparams_for_current_model = self._build_hyperparameter_dictionary(
            trial, hyperparameter_dict
        )

        if self.model_flavor == "xgboost":
            model = self._instantiate_xgboost_object(
                algo_name, hyperparams_for_current_model
            )
        else:
            model = self._instantiate_sklearn_object(
                algo_name, hyperparams_for_current_model
            )
        model.fit(self.train_features, self.train_target)

        trial.set_user_attr("trained_model", model)

        metric = self._instantiate_sklearn_object(self.metric, is_model=False)

        y_pred = model.predict(self.test_features)
        return metric(self.test_target, y_pred)

    def run_hyperparameter_search(self):
        """Runs the hyperparameter search."""

        self.optuna_study.optimize(
            self._objective, n_trials=self.n_trials, callbacks=[TrialLogCallback(self)],
        )

        best_trial_id = self.optuna_study.best_trial.number
        best_trial_name = self.trial_names[best_trial_id]

        best_trial = generate_folder_name(trial_id=best_trial_id, name=best_trial_name,)

        study_manifest = {"best_trial": best_trial}

        ManifestInterface.write(
            manifest=study_manifest, kind="study", my_study_path=self.my_study_path,
        )

        self.write_summary_log(best_trial_id=best_trial)

        zip_study(self.my_study_path)

        return None

    def write_summary_log(self, best_trial_id):
        summary_string = f"""

----- [ SUMMARY: HYPERPARAMETER SEARCH ] -----
Study Name: {self.study_name}
Model Flavor: {self.model_flavor}
Number of Trials: {self.n_trials}
Metric: {self.metric}
Direction: {self.direction}
Best Trial Id: {best_trial_id}
Best Trial Metric Score: {self.optuna_study.best_trial.value}
Best Trial Hyperparameters: {self.optuna_study.best_trial.params}
----- [ END: HYPERPARAMETER SEARCH ] -----
            """
        print(summary_string)

        logging.basicConfig(
            filename=f"{self.my_study_path}/study_summary.json",
            encoding="utf-8",
            level=logging.INFO,
        )

        logging.info(summary_string)

        summary_log = {
            "study_name": self.study_name,
            "model_flavor": self.model_flavor,
            "n_trials": self.n_trials,
            "metric": self.metric,
            "direction": self.direction,
            "best_trial_id": best_trial_id,
            "best_trial_metric_score": self.optuna_study.best_trial.value,
            "best_trial_hyperparameters": self.optuna_study.best_trial.params,
        }

        write_json_to_local(summary_log, f"{self.my_study_path}/study_summary.json")

    def _suggest_hyperparameter_value(self, trial, hyperparameter, attribute_dict):
        """Suggests value(s) for the hyperparameter using the details in
        the attribute dictionary.
        The attribute must belong to the following distributions ("float", "int",
        "uniform", "discrete-uniform", "log-uniform", "int-uniform", "int-log-uniform").

        Parameters
        ----------
        trial : optuna.trial._trial.Trial
            A single execution of the objective function.
        hyperparameter : str
            The hyperparameter to which the attributes belong.
        attribute_dict : Dict
            Dictionary specifying the attributes of the hyperparameter.

        Returns
        -------
        float/int
            Suggested value for the attribute.

        """
        if "distribution" in attribute_dict.keys():
            distribution = attribute_dict["distribution"]
        else:
            attribute_dict["distribution"] = "float"
            distribution = "float"

        supported_distributions = (
            "float",
            "int",
            "uniform",
            "log-uniform",
            "discrete-uniform",
            "int-uniform",
            "int-log-uniform",
        )

        if distribution not in supported_distributions:
            raise HyperparameterStudyError(
                "Hyperdrive Hyperparameter "
                "Experiment definition using yaml only supports the following "
                f"distributions: {supported_distributions}"
            )

        if distribution in ("float", "uniform", "log-uniform"):
            return self._float_handler(trial, hyperparameter, attribute_dict)
        elif distribution in ("int", "int-uniform", "int-log-uniform"):
            return self._int_handler(trial, hyperparameter, attribute_dict)
        elif distribution in ("discrete-uniform",):
            return self._discrete_handler(trial, hyperparameter, attribute_dict)

    def _build_hyperparameter_dictionary(self, trial, hyperparameter_dict):
        """Builds the hyperparameter dictionary for a model.

        Parameters
        ----------
        trial : optuna.trial._trial.Trial
            A single execution of the objective function.
        hyperparameter_dict : Dict
            Dictionary containing the hyperparameters and details describing the range
            of values they can assume.

        Returns
        -------
        hyperparams_for_current_model : Dict
            Dictionary containing the suggested values for hyperparameters.

        """
        hyperparams_for_current_model = {}
        for hyperparameter, attributes in hyperparameter_dict.items():
            attribute = hyperparameter_dict[hyperparameter]
            if type(attribute) is dict:
                hyperparams_for_current_model[
                    hyperparameter
                ] = self._suggest_hyperparameter_value(trial, hyperparameter, attribute)
            elif type(attribute) is str and (
                "np." in attribute or "numpy." in attribute
            ):
                hyperparams_for_current_model[
                    hyperparameter
                ] = self._handle_numpy_value(trial, hyperparameter, attribute)
            elif type(attribute) is list and type(attribute[0]) is dict:
                hyperparams_for_current_model = self._nested_handler(
                    trial, hyperparameter, attribute, hyperparams_for_current_model
                )
            else:
                hyperparams_for_current_model[
                    hyperparameter
                ] = self._handle_explicit_values(trial, hyperparameter, attribute)

        return hyperparams_for_current_model

    def _call_numpy(self, numpy_command):
        """Calls the Numpy command and returns the results.

        Parameters
        ----------
        numpy_command : str
            The Numpy command to be called. The commands np.linspace,
            np.logspace, np.arange and np.geomspace are supported.

        Returns
        -------
        np_values : list
            List of values returned by the Numpy command.

        """
        idx1 = numpy_command.find(".")
        idx2 = numpy_command.find("(")

        np_args = numpy_command[idx2 + 1 : -1].split(",")
        function_name = numpy_command[idx1 + 1 : idx2]

        args = []
        kwargs = {}
        for argument in np_args:
            if "=" in argument:
                arg_name, arg_value = argument.split("=")
                kwargs[arg_name.strip()] = self._parse_np_argument(arg_value.strip())
            else:
                args.append(self._parse_np_argument(argument.strip()))

        np_values = getattr(__import__("numpy"), function_name)(*args, **kwargs)
        return np_values

    def _create_optuna_study(self):
        """Creates an Optuna study and stores it in an attribute."""
        self.optuna_study = optuna.create_study(
            study_name=self.study_name, direction=self.direction,
        )

    def _discrete_handler(self, trial, hyperparameter, attribute_dict):
        """Returns a value in the discrete uniform distribution for the
        hyperparameter according to the attributes specified in the
        attribute dictionary.

        Parameters
        ----------
        trial : optuna.trial._trial.Trial
            A single execution of the objective function.
        hyperparameter : str
            The hyperparameter to which the attributes belong.
        attribute_dict : Dict
            Dictionary specifying the attributes for the hyperparameter.

        Returns
        -------
        float
        """

        trial_dict = {}
        trial_dict["name"] = hyperparameter
        trial_dict["low"] = attribute_dict["low"]
        trial_dict["high"] = attribute_dict["high"]

        if "step" in attribute_dict:
            trial_dict["q"] = attribute_dict["step"]
        return trial.suggest_discrete_uniform(**trial_dict)

    def _float_handler(self, trial, hyperparameter, attribute_dict):
        """Returns a value for the hyperparameter according to the attributes
        specified in the attribute dictionary.
        Value belongs to one of the following distributions: {"float", "uniform",
        "log-uniform"}

        Parameters
        ----------
        trial : optuna.trial._trial.Trial
            A single execution of the objective function.
        hyperparameter : str
            The hyperparameter to which the attributes belong.
        attribute_dict : Dict
            Dictionary specifying the attributes for the hyperparameter.

        Returns
        -------
        float
        """
        trial_dict = {}
        trial_dict["name"] = hyperparameter
        trial_dict["low"] = attribute_dict["low"]
        trial_dict["high"] = attribute_dict["high"]

        distribution = attribute_dict["distribution"]

        if distribution == "log-uniform":
            if trial_dict["low"] < 0:
                raise HyperparameterStudyError(
                    "If log is True, low cannot be less than zero."
                )
            elif trial_dict["low"] == 0:
                trial_dict["low"] = 1e-7
            trial_dict["log"] = True

        return trial.suggest_float(**trial_dict)

    def _get_study_metadata(self):
        """Stores the experiment metadata as attributes."""
        self.study_name = self.search_dictionary.get("study_name")
        self.direction = self.search_dictionary.get("direction", "minimize")
        self.pruner = self.search_dictionary.get("pruner", "median")
        self.sampler = self.search_dictionary.get("sampler", "random")
        self.n_trials = self.search_dictionary.get("n_trials", 100)
        self.metric = self.search_dictionary.get("metric")
        self.models = self.search_dictionary.get("models")
        self.model_flavor = self.search_dictionary.get("model_flavor")

    def _handle_explicit_values(self, trial, hyperparameter, attributes):
        """Returns a value from a list of specified values for the
        hyperparameter. Can handle a single value.

        Parameters
        ----------
        trial : optuna.trial._trial.Trial
            A single execution of the objective function.
        hyperparameter : str
            The hyperparameter to which the attributes belong.
        attributes : List/str/int/float
            Value or list of values the hyperparameter can assume.

        Returns
        -------
        float/int/str
        """
        hyperparameter_name = hyperparameter + "*" + str(uuid4())[-6:]

        if type(attributes) is not list:
            attributes = [attributes]
        return trial.suggest_categorical(hyperparameter_name, attributes)

    def _handle_numpy_value(self, trial, hyperparameter, numpy_command):
        """Returns a value for the hyperparameter by using the Numpy command
        specified.

        Parameters
        ----------
        trial : optuna.trial._trial.Trial
            A single execution of the objective function.
        hyperparameter : str
            The hyperparameter to which the attributes belong.
        numpy_command : str
            Numpy command specifying the range of values the hyperparameter can
            assume.

        Returns
        -------
        float
        """
        hyperparameter_name = hyperparameter + "*" + str(uuid4())[-6:]

        numpy_values = self._call_numpy(numpy_command)
        return trial.suggest_categorical(hyperparameter_name, numpy_values)

    def _instantiate_sklearn_object(
        self, clf_name, hyperparams_for_current_model={}, is_model=True
    ):
        """Instantiates the sklearn object.

        Parameters
        ----------
        clf_name : str
            Name of the object to be instantiated.
        hyperparams_for_current_model :  Dict, optional
            Dictionary containing the hyperparameters to be used when instantiating
            the model (default is an empty dictionary).
        is_model : boolean, optional
            Indicates if the sklearn object to be instantiated is a model (default is
            True).

        Returns
        -------
        sklearn object
        """
        sklearn_string_length = len("sklearn.")
        idx1 = clf_name.rfind(".")
        module_name = clf_name[sklearn_string_length:idx1]
        function_name = clf_name[idx1 + 1 :]
        module = getattr(__import__("sklearn"), module_name)

        if is_model:
            return getattr(module, function_name)(**hyperparams_for_current_model)
        return getattr(module, function_name)

    def _instantiate_xgboost_object(
        self, clf_name, hyperparams_for_current_model={}, is_model=True
    ):
        """Instantiates the xgboost object.

        Parameters
        ----------
        clf_name : str
            Name of the object to be instantiated.
        hyperparams_for_current_model :  Dict, optional
            Dictionary containing the hyperparameters to be used when instantiating
            the model (default is an empty dictionary).
        is_model : boolean, optional
            Indicates if the xgboost object to be instantiated is a model (default
            is True).

        Returns
        -------
        xgboost object
        """
        xgboost_string_length = len("xgboost.")

        if clf_name.count(".") == 1:
            module = __import__("xgboost")
            function_name = clf_name[xgboost_string_length:]
        else:
            idx1 = clf_name.rfind(".")
            module_name = clf_name[xgboost_string_length:idx1]
            function_name = clf_name[idx1 + 1 :]
            module = getattr(__import__("xgboost"), module_name)

        if is_model:
            return getattr(module, function_name)(**hyperparams_for_current_model)
        return getattr(module, function_name)

    def _int_handler(self, trial, hyperparameter, attribute_dict):
        """Returns a value for the hyperparameter according to the attributes
        specified in the attribute dictionary.
        Value belongs to one of the following distributions: {"int", "int-uniform",
        "int-log-uniform"}

        Parameters
        ----------
        trial : optuna.trial._trial.Trial
            A single execution of the objective function.
        hyperparameter : str
            The hyperparameter to which the attributes belong.
        attribute_dict : Dict
            Dictionary specifying the attributes for the hyperparameter.

        Returns
        -------
        int
        """

        trial_dict = {}
        trial_dict["name"] = hyperparameter
        trial_dict["low"] = attribute_dict["low"]
        trial_dict["high"] = attribute_dict["high"]

        distribution = attribute_dict["distribution"]

        if distribution == "int-log-uniform":
            if trial_dict["low"] < 0:
                raise HyperparameterStudyError(
                    "If log is True, low cannot be less than or equal to zero."
                )
            trial_dict["log"] = True

        return trial.suggest_int(**trial_dict)

    def _nested_handler(
        self, trial, hyperparameter, nested_list, hyperparams_for_current_model,
    ):
        """Returns the hyperparameter dictionary for nested hyperparameters.

        Parameters
        ----------
        trial : optuna.trial._trial.Trial
            A single execution of the objective function.
        hyperparameter : str
            The hyperparameter to which the attributes belong.
        nested_list : list
            List containing dictionaries which specify the nested hyperparameters.
        hyperparams_for_current_model : Dict
            Dictionary containing the suggested values for hyperparameters.

        Returns
        -------
        hyperparams_for_current_model : Dict
            Updated dictionary containing suggested values for the nested
            hyperparameters.
        """
        first_attribute_values = [
            next(iter(attribute_dict))
            if type(attribute_dict) is dict
            else attribute_dict
            for attribute_dict in nested_list
        ]

        numeric_keys = set(["low", "high", "distribution", "step"])
        selected_idx = np.random.randint(len(first_attribute_values))

        if type(nested_list[selected_idx]) is not dict:
            nested_list[selected_idx] = {nested_list[selected_idx]: {}}

        if type(first_attribute_values[selected_idx]) is str and (
            "np" in first_attribute_values[selected_idx]
            or "numpy" in first_attribute_values[selected_idx]
        ):
            selected_command = first_attribute_values[selected_idx]

            selected_value = self._handle_numpy_value(
                trial, hyperparameter, selected_command
            )
            remove_keys = [selected_command]
        elif numeric_keys.isdisjoint(set(nested_list[selected_idx].keys())):
            selected_explicit_value = first_attribute_values[selected_idx]
            selected_value = self._handle_explicit_values(
                trial, hyperparameter, selected_explicit_value
            )

            remove_keys = [selected_value]
        else:
            selected_attribute_list = nested_list[selected_idx]

            selected_value = self._suggest_hyperparameter_value(
                trial, hyperparameter, selected_attribute_list
            )
            remove_keys = ["low", "high", "step", "distribution"]

        hyperparams_for_current_model[hyperparameter] = selected_value

        nested_hyperparam_dict = nested_list[selected_idx].copy()
        for key in remove_keys:
            if key in nested_hyperparam_dict:
                nested_hyperparam_dict.pop(key)

        nested_hyperparams = self._build_hyperparameter_dictionary(
            trial, nested_hyperparam_dict
        )

        hyperparams_for_current_model.update(nested_hyperparams)
        return hyperparams_for_current_model

    def _parse_np_argument(self, arg):
        """Parses the argument for a Numpy command and returns it.

        Parameters
        ----------
        arg : str
            Argument for the Numpy command.

        Returns
        -------
        int/float/boolean/type
        """
        if not arg.isalpha():
            if "." in arg:
                return float(arg)
            return int(arg)
        elif arg == "True":
            return True
        elif arg == "False":
            return False
        elif arg == "int":
            return int
        elif arg == "float":
            return float
        elif arg == "None":
            return None

    def _load_yaml_file_as_search_dictionary(self, yaml_file_name):
        """Converts the YAML hyperparameter search file into a Python dictionary
        and stores it in an attribute."""
        self.search_dictionary = load_yaml_file(yaml_file_name)

    def _get_trial_names(self):
        n_trials = self.n_trials

        names_list = [
            "abiding",
            "able",
            "able_bodied",
            "absolute",
            "abstract",
            "accommodating",
            "accomplished",
            "accordant",
            "active",
            "adamant",
            "adapted",
            "adorable",
            "adumbrative",
            "adventuresome",
            "adventurous",
            "affable",
            "affecting",
            "affectionate",
            "aggressive",
            "agreeable",
            "agreeing",
            "alien",
            "allegiant",
            "alluring",
            "amazing",
            "ambrosial",
            "amiable",
            "amicable",
            "amical",
            "amusing",
            "angelic",
            "animated",
            "appreciating",
            "approachable",
            "archetypal",
            "ardent",
            "arresting",
            "assured",
            "astonishing",
            "astounding",
            "astral",
            "athletic",
            "attached",
            "attentive",
            "attractive",
            "audacious",
            "august",
            "auspicious",
            "awe_inspiring",
            "awesome",
            "beatific",
            "beautiful",
            "beneficent",
            "beneficial",
            "benevolent",
            "benign",
            "benignant",
            "big",
            "big_hearted",
            "bland",
            "bleeding_heart",
            "blessed",
            "blest",
            "blissful",
            "blithe",
            "bloodless",
            "bold",
            "bonhomous",
            "bound",
            "boundless",
            "brawny",
            "breathless",
            "breathtaking",
            "breezeless",
            "breezy",
            "bright",
            "brilliant",
            "brisk",
            "bucolic",
            "buddy_buddy",
            "calm",
            "capable",
            "captivated",
            "captivating",
            "caring",
            "celebrated",
            "celestial",
            "changeless",
            "charitable",
            "charming",
            "cheerful",
            "cheery",
            "chipper",
            "chirpy",
            "chivalrous",
            "chummy",
            "civil",
            "classic",
            "classical",
            "clever",
            "close",
            "clubby",
            "collected",
            "commendable",
            "commiserating",
            "commiserative",
            "companionable",
            "compassionate",
            "compatible",
            "complaisant",
            "composed",
            "comprehending",
            "comradely",
            "conciliatory",
            "concordant",
            "condolatory",
            "condoling",
            "confidential",
            "confiding",
            "congenial",
            "congruous",
            "considerate",
            "consistent",
            "consonant",
            "constant",
            "consummate",
            "content",
            "contented",
            "conversable",
            "convivial",
            "cool",
            "cooperative",
            "copacetic",
            "cordial",
            "courageous",
            "courteous",
            "courtly",
            "cozy",
            "daredevil",
            "daring",
            "darling",
            "daunting",
            "dauntless",
            "dazzling",
            "dear",
            "dedicated",
            "delectable",
            "delicious",
            "delighted",
            "delightful",
            "delineative",
            "demoniac",
            "dependable",
            "depictive",
            "determined",
            "devoted",
            "distinguished",
            "distressing",
            "divine",
            "doting",
            "doughty",
            "dreadful",
            "driving",
            "dummy",
            "durable",
            "dynamic",
            "dynamite",
            "eager",
            "easeful",
            "easy",
            "ecstatic",
            "effulgent",
            "elated",
            "elevated",
            "elysian",
            "emblematic",
            "eminent",
            "empathetic",
            "empathic",
            "empyreal",
            "empyrean",
            "enchanting",
            "enduring",
            "energetic",
            "engaging",
            "enjoyable",
            "enterprising",
            "entertaining",
            "enthusiastic",
            "entire",
            "epic",
            "equable",
            "established",
            "esteemed",
            "eternal",
            "ethereal",
            "evocative",
            "exaggerated",
            "exalted",
            "exceeding",
            "excellent",
            "exciting",
            "exemplary",
            "exultant",
            "facsimile",
            "fair",
            "faithful",
            "famed",
            "familiar",
            "famous",
            "fantastic",
            "fascinating",
            "fast",
            "favorable",
            "fearless",
            "fearsome",
            "fiery",
            "fine",
            "finished",
            "fire_eating",
            "firm",
            "fit",
            "fixed",
            "flawless",
            "fond",
            "forbearing",
            "forceful",
            "forcible",
            "formidable",
            "forthcoming",
            "fresh",
            "friendly",
            "gallant",
            "game",
            "genial",
            "gentle",
            "glad",
            "gleeful",
            "glorious",
            "godlike",
            "good",
            "good_hearted",
            "good_humored",
            "good_natured",
            "gorgeous",
            "gracious",
            "grand",
            "grandiose",
            "gratified",
            "gratifying",
            "great",
            "grind",
            "gritty",
            "gutsy",
            "gutty",
            "hair_raising",
            "halcyon",
            "hale",
            "hallowed",
            "happy",
            "hard_working",
            "hardy",
            "harmonious",
            "heart_stirring",
            "heart_stopping",
            "hearty",
            "heavenly",
            "heavy",
            "heavy_duty",
            "helpful",
            "heroic",
            "high",
            "high_flown",
            "high_powered",
            "high_spirited",
            "honored",
            "hospitable",
            "huggy",
            "humane",
            "humanitarian",
            "hushed",
            "hyper",
            "hypothetical",
            "ideal",
            "illimitable",
            "illustrative",
            "illustrious",
            "imitation",
            "immeasurable",
            "immense",
            "immobile",
            "immortal",
            "immovable",
            "impavid",
            "imposing",
            "impossible",
            "impressive",
            "inactive",
            "incalculable",
            "incessant",
            "incommunicable",
            "incomparable",
            "incredible",
            "indefatigable",
            "indefinable",
            "indefinite",
            "indescribable",
            "indomitable",
            "indulgent",
            "industrious",
            "ineffable",
            "inexhaustible",
            "inexorable",
            "inexpressible",
            "infinite",
            "inflated",
            "inflexible",
            "innate",
            "inspiring",
            "intact",
            "intellectual",
            "intense",
            "intent",
            "interested",
            "intimate",
            "intimidating",
            "intoxicated",
            "intrepid",
            "intuitive",
            "irenic",
            "jolly",
            "jovial",
            "joyful",
            "joyous",
            "jubilant",
            "jumping",
            "kind",
            "kindhearted",
            "kindly",
            "kindred",
            "kinetic",
            "laughing",
            "lenient",
            "level",
            "liege",
            "light",
            "like_minded",
            "limitless",
            "lion_hearted",
            "lionhearted",
            "lively",
            "lofty",
            "lovable",
            "lovely",
            "lovey_dovey",
            "loving",
            "low_key",
            "loyal",
            "luscious",
            "lush",
            "lusty",
            "magical",
            "magnificent",
            "majestic",
            "marvelous",
            "matchless",
            "measureless",
            "mellow",
            "memorable",
            "merciful",
            "merry",
            "mighty",
            "mild",
            "mind_blowing",
            "mind_boggling",
            "miniature",
            "miraculous",
            "mirthful",
            "model",
            "motionless",
            "moving",
            "muscular",
            "mushy",
            "mystical",
            "mythological",
            "nameless",
            "neighborly",
            "nervy",
            "neutral",
            "neutralist",
            "never_failing",
            "noble",
            "nonbelligerent",
            "nonviolent",
            "notable",
            "noted",
            "obdurate",
            "obliging",
            "obscure",
            "olympian",
            "original",
            "otherworldly",
            "outgoing",
            "outrageous",
            "overjoyed",
            "overwhelming",
            "pacific",
            "pacifistic",
            "pally",
            "palsy_walsy",
            "paradigmatic",
            "partial",
            "pastoral",
            "peace_loving",
            "peaceable",
            "peaceful",
            "peerless",
            "peppy",
            "perfect",
            "perky",
            "persevering",
            "piteous",
            "pitying",
            "placatory",
            "placid",
            "playful",
            "pleasant",
            "pleased",
            "pleasing",
            "pleasurable",
            "plucky",
            "polite",
            "potent",
            "powerful",
            "preeminent",
            "presentational",
            "primordial",
            "princely",
            "prodigious",
            "propitious",
            "prototypal",
            "prototypical",
            "proud",
            "quiescent",
            "quiet",
            "quintessential",
            "radiant",
            "rapturous",
            "ravishing",
            "receptive",
            "refreshing",
            "regular",
            "reinforced",
            "relentless",
            "reliable",
            "remarkable",
            "renowned",
            "rep",
            "reposeful",
            "reposing",
            "representative",
            "resolute",
            "resplendent",
            "responsive",
            "restful",
            "right",
            "righteous",
            "rigid",
            "robust",
            "rugged",
            "rural",
            "sacred",
            "satisfied",
            "satisfying",
            "scrumptious",
            "secure",
            "sensitive",
            "seraphic",
            "serene",
            "shining",
            "shocking",
            "similar",
            "simple",
            "sinewy",
            "single_minded",
            "slow",
            "smooth",
            "snappy",
            "sociable",
            "social",
            "soft",
            "soft_shell",
            "softhearted",
            "solemn",
            "solicitous",
            "solid",
            "soothing",
            "sound",
            "sparing",
            "sparkling",
            "spartan",
            "spectacular",
            "spine_tingling",
            "spirited",
            "spiritual",
            "splendid",
            "splendiferous",
            "splendorous",
            "sprightly",
            "spry",
            "square_shooting",
            "stable",
            "staggering",
            "stalwart",
            "standard",
            "stark",
            "stately",
            "staunch",
            "steadfast",
            "steady",
            "still",
            "stormless",
            "stout",
            "stouthearted",
            "strapping",
            "strenuous",
            "striking",
            "strong",
            "stubborn",
            "stunning",
            "stupefying",
            "sturdy",
            "suave",
            "sublime",
            "substantial",
            "suitable",
            "sunny",
            "super",
            "superb",
            "superior",
            "supernal",
            "supernatural",
            "supportive",
            "supreme",
            "sure",
            "surpassing",
            "surprising",
            "sweet_tempered",
            "swell",
            "symbolic",
            "symbolical",
            "sympathetic",
            "sympathizing",
            "tenacious",
            "tender",
            "tenderhearted",
            "theoretical",
            "thick",
            "thoughtful",
            "thrilled",
            "thrilling",
            "tickled",
            "tight",
            "time_honored",
            "tireless",
            "together",
            "tough",
            "towering",
            "tranquil",
            "transcendent",
            "transcendental",
            "transcending",
            "transmundane",
            "tremendous",
            "tried_and_true",
            "triumphant",
            "true",
            "true_blue",
            "typical",
            "ultimate",
            "unafraid",
            "unanimous",
            "unbelievable",
            "unbending",
            "unblinking",
            "unbounded",
            "uncanny",
            "uncompromising",
            "unconfined",
            "unctuous",
            "undaunted",
            "understanding",
            "undisturbed",
            "unearthly",
            "unending",
            "unequalable",
            "unequalled",
            "unfaltering",
            "unflagging",
            "unflinching",
            "unintimidated",
            "unique",
            "united",
            "unlimited",
            "unmovable",
            "unparalleled",
            "unqualified",
            "unquestioning",
            "unrivalled",
            "unruffled",
            "unshrinking",
            "unspeakable",
            "unswerving",
            "untellable",
            "untiring",
            "untold",
            "untroubled",
            "unutterable",
            "unwavering",
            "unwearied",
            "unwearying",
            "unworldly",
            "unyielding",
            "up",
            "upbeat",
            "urbane",
            "usual",
            "valiant",
            "valorous",
            "vast",
            "venerable",
            "venturesome",
            "venturous",
            "very",
            "vicarious",
            "vigorous",
            "vintage",
            "visionary",
            "vital",
            "vivacious",
            "warm",
            "warmhearted",
            "waveless",
            "welcoming",
            "well_adjusted",
            "well_balanced",
            "well_built",
            "well_disposed",
            "well_founded",
            "well_known",
            "well_made",
            "well_mannered",
            "well_organized",
            "well_suited",
            "whole",
            "wholehearted",
            "windless",
            "winning",
            "wonderful",
            "wondrous",
            "yummy",
            "zippy",
        ]

        if n_trials > len(names_list):
            self.trial_names = random.choices(names_list, k=n_trials)
        else:
            self.trial_names = random.sample(names_list, k=n_trials)
