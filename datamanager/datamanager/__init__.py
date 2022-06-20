from dataclasses import dataclass, field
import os

from .controllers.dataset import DatasetController
from .controllers.file import FileController
from .filesystem import FileSystem


class DataManagerError(Exception):
    pass


@dataclass
class DataManager(
    DatasetController,
    FileController,
):
    storage_provider: str = "local"
    volume_name: str = "/tmp"
    volume_region: str = None
    access_key: str = field(repr=False, default=None)
    secret_access_key: str = field(repr=False, default=None)
    session_token: str = field(repr=False, default=None)
    expiration: str = field(repr=False, default=None)

    def __post_init__(self):
        if self.storage_provider == "aws":
            if self.access_key is None:
                self.access_key = os.environ.get("AWS_ACCESS_KEY_ID", self.access_key)
            if self.secret_access_key is None:
                self.secret_access_key = os.environ.get(
                    "AWS_SECRET_ACCESS_KEY", self.secret_access_key
                )
            if self.session_token is None:
                self.session_token = os.environ.get("SESSION_TOKEN", self.session_token)

            if self.access_key in [None, ""] or self.secret_access_key in [None, ""]:
                raise DataManagerError(
                    "No valid credential profile configured for storage provider: "
                    f"{self.storage_provider}"
                )

        self.filesystem = FileSystem(
            storage_provider=self.storage_provider,
            volume_name=self.volume_name,
            volume_region=self.volume_region,
            access_key=self.access_key,
            secret_access_key=self.secret_access_key,
            session_token=self.session_token,
            expiration=self.expiration,
        )


__version__ = "0.0.1"
