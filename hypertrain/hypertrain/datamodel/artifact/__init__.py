from dataclasses import dataclass
from .__doc__ import ArtifactDoc
from .methods import (
    ArtifactMethods,
    MachineLearningModelMethods,
    PythonFileMethods,
)
from ..base import Base


@dataclass
class Artifact(Base, ArtifactDoc, ArtifactMethods):
    id = None
    run = None
    run_id = None
    artifact = None
    artifact_type = None
    file_key = None
    input_signature = None
    output_signature = None
    mutable = tuple()

    def __repr__(self):
        return f"<Artifact {self.file_key}>"


@dataclass
class MachineLearningModel(Artifact, MachineLearningModelMethods):
    id = None
    run = None
    run_id = None
    artifact = None
    artifact_type = None
    file_key = None
    input_signature = None
    output_signature = None
    mutable = tuple()

    def __repr__(self):
        return f"<MachineLearningModel artifact: {self.file_key}>"


@dataclass
class Manifest(Artifact):
    id = None
    run = None
    trial_id = None
    artifact = None
    artifact_type = None
    file_key = None
    input_signature = None
    output_signature = None
    mutable = tuple()

    def __repr__(self):
        return f"<Manifest artifact: {self.file_key}>"


@dataclass
class Preprocessor(MachineLearningModel):
    id = None
    run = None
    run_id = None
    artifact = None
    artifact_type = None
    file_key = None
    input_signature = None
    output_signature = None
    mutable = tuple()

    def __repr__(self):
        return f"<Preprocessor artifact: {self.file_key}>"


@dataclass
class PythonFile(Artifact, PythonFileMethods):
    id = None
    run = None
    run_id = None
    artifact = None
    artifact_type = None
    file_key = None
    input_signature = None
    output_signature = None
    mutable = tuple()

    def __repr__(self):
        return f"<PythonFile artifact: {self.file_key}>"


@dataclass
class TrainedModel(MachineLearningModel):
    id = None
    trial = None
    trial_id = None
    artifact = None
    artifact_type = None
    file_key = None
    input_signature = None
    output_signature = None
    mutable = tuple()

    def __repr__(self):
        return f"<TrainedModel artifact: {self.file_key}>"


__doc__ = ArtifactDoc.__doc__
