from dataclasses import dataclass, field
from mimetypes import init
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
    init_spark: bool = False

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

        if self.init_spark:
            import pyspark

            builder = pyspark.sql.SparkSession.builder.config(
                "spark.jars.packages",
                (
                    "org.apache.hadoop:hadoop-aws:3.2.0,"
                    "com.amazonaws:aws-java-sdk-bundle:1.12.119"
                ),
            )

            if self.suppress_stage:
                builder.config("spark.ui.showConsoleProgress", False)

            spark = builder.getOrCreate()
            spark.sparkContext.setLogLevel("ERROR")
            spark._sc._jsc.hadoopConfiguration().set(
                "fs.s3a.impl", "org.apache.hadoop.fs.s3a.S3AFileSystem"
            )
            from .controllers.spark import SparkController

            self.spark = spark
            self.spark_controller = SparkController(self, decorate=self.decorate)


__version__ = "0.0.1"
