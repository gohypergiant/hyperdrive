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
            self.access_key = os.environ.get("AWS_ACCESS_KEY_ID")
            self.secret_access_key = os.environ.get("AWS_SECRET_ACCESS_KEY")

            if self.access_key in [None, ''] or self.secret_access_key in [None, '']:
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
