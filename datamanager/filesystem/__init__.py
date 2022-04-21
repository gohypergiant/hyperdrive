from pathlib import Path

from datetime import datetime
import os
import pyarrow
import pyarrow.csv as csv
import pyarrow.parquet as parquet
import numpy as np
from pandas import DataFrame, Series
from ..exceptions import FileSystemError
from ..utilities import infer_target_format


class FileSystem:
    """Defines the FileSystem class, can be either local or aws.

    Attributes
    ----------
    local : bool
        The flag that indicates if the FileSystem is local (default is False).
    bucket : str
        The path of the root folder to the FileSystem or the S3 bucket name.
    storage_provider : str
        The string that indicates the FileSystem type.
    access_key : str
        The access_key with permission to access the S3 bucket, if the storage_provider
        is 'aws' (default is None).
    secret_access_key : str
        The secret_access_key with permission to access the S3 bucket,
        if the storage_provider is 'aws' (default is None).
    region : str
        The region where the S3 bucket is located, if the storage_provider
        is 'aws' (default is None).
    account : str
        The account id of AWS, if the storage_provider is 'aws' (default is None).

    Methods
    -------
    create_dir(directory: str)
        Deletes a directory of the FileSystem.
    delete_dir(directory: str)
        Deletes a directory of the FileSystem.
    delete_dir_contents(directory: str)
        Delete the content of a directory of the FileSystem.
    delete_file(file: str)
        Delete a file of the FileSystem.
    get_file_info(file: str)
        Get file information.
    get_object(file: str)
        Get an object from the FileSystem.
    list(subpath: str = "", recursive: bool = False, exclude_projects:
    bool = True)
        List the content of a bucket of the FileSystem.
    write_file_to_remote(file: str, file_key: str, metadata: dictionary)
        Writes a file to the bucket.
    write_file_from_remote(file_key: str, path: str)
        Write a file locally from the remote bucket.
    write_dataset(dataframe: pandas.DataFrame, file_key: str)
        Writes a dataset to the bucket.
    """

    file_types = [
        "target does not exist",
        "unknown",
        "file",
        "directory",
    ]

    def __init__(
        self,
        storage_provider,
        volume_name,
        volume_account=None,
        volume_region=None,
        access_key=None,
        secret_access_key=None,
        session_token=None,
        expiration=None,
    ):

        self.bucket = volume_name
        self.storage_provider = storage_provider
        self.access_key = access_key
        self.secret_access_key = secret_access_key
        self.region = volume_region
        self.account = volume_account
        self.session_token = session_token
        self.expiration = expiration
        if self.expiration is not None:
            self.expiration = datetime.strptime(expiration[:-5], "%Y-%m-%dT%H:%M:%S")

        self.local = False
        if storage_provider == "local":
            self.fs = pyarrow.fs.LocalFileSystem()
            self.local = True
        elif storage_provider == "aws":
            access_key = self.access_key
            if (
                access_key is None
                and (access_key := os.getenv("AWS_ACCESS_KEY_ID", self.access_key))
                is None
            ):
                raise FileSystemError(
                    "Configuring an S3 FileSystem requires an access key passed "
                    "or configured as the environment variable AWS_ACCESS_KEY_ID"
                )
            secret_key = self.secret_access_key
            if (
                secret_key is None
                and (
                    secret_key := os.getenv(
                        "AWS_SECRET_ACCESS_KEY", self.secret_access_key
                    )
                )
                is None
            ):
                raise FileSystemError(
                    "Configuring an S3 FileSystem requires a secret access key passed "
                    "or configured as the environment variable AWS_SECRET_ACCESS_KEY"
                )
            self.fs = pyarrow.fs.S3FileSystem(
                access_key=access_key,
                secret_key=secret_access_key,
                session_token=self.session_token,
                region=self.region,
            )

    def create_dir(self, directory):
        """Create a directory in the FileSystem.

        Parameters
        ----------
        directory : str
            The name of the directory.
        """
        self.fs.create_dir(directory)
        return True

    def delete_dir(self, directory):
        """Deletes a directory of the FileSystem.

        Parameters
        ----------
        directory : str
            The name of the directory.
        """
        try:
            self.fs.delete_dir(f"{self.bucket}/{directory}")
        except (OSError, FileNotFoundError) as err:
            if "code 100" in str(err) or type(err) is FileNotFoundError:
                pass
            else:
                raise err
        return True

    def delete_dir_contents(self, directory):
        """Delete the content of a directory of the FileSystem.

        Parameters
        ----------
        directory : str
            the name of the directory.
        """
        self.fs.delete_dir_contents(f"{self.bucket}/{directory}")
        return True

    def delete_file(self, file):
        """Delete a file of the FileSystem.

        Parameters
        ----------
        file : str
            The name of the file.
        """
        try:
            self.fs.delete_file(f"{self.bucket}/{file}")
        except OSError:
            pass
        return True

    def get_file_info(self, file):
        """Get file information.

        Read more about Pyarrow FileInfo:
        https://arrow.apache.org/docs/python/generated/pyarrow.fs.FileInfo.html

        Parameters
        ----------
        file : str
            The name of the file.

        Returns
        -------
        dictionary
            File information.
        """
        file_info = self.fs.get_file_info(f"{self.bucket}/{file}")

        return {
            "base_name": file_info.base_name,
            "extension": file_info.extension,
            "is_file": file_info.is_file,
            "mtime": file_info.mtime,
            "mtime_ns": file_info.mtime_ns,
            "path": file_info.path,
            "size": file_info.size,
            "type": FileSystem.file_types[file_info.type],
        }

    def get_object(self, file):
        """Get an object from the FileSystem.

        Parameters
        ----------
        file : str
            The file name.

        Returns
        -------
        pyarrow.NativeFile
            The object.
        """
        return self.fs.open_input_file(f"{self.bucket}/{file}")

    def list(self, subpath="", recursive=False, exclude_projects=True):
        """List the content of a bucket of the FileSystem.

        Parameters
        ----------
        subpath : str, optional
            The subpath to the folder to list (default is "").
        recursive : bool, optional
            The flag that indicates if should list content of subfolders (default is
            False).
        exclude_projects : bool, optional
            The flag that indicates if shouldn't list projects (default is True).

        Returns
        -------
        pandas.DataFrame
            Contents of the bucket.

        Raises
        ------
        FileSystemError
            Raises error for the following error codes:
                "core 15" : Wrong credientials configured.
                "code 100" : Path doesn't exist in bucket or wrong region was
                configured.
        """
        try:
            if self.local:
                Path(f"{self.bucket}/{subpath}").mkdir(parents=True, exist_ok=True)
            files = self.fs.get_file_info(
                pyarrow.fs.FileSelector(f"{self.bucket}/{subpath}", recursive=recursive)
            )
        except OSError as err:
            if "code 100" in str(err):
                raise FileSystemError(
                    "Either the path does not exist in your bucket or the DataRepo "
                    "has been configured with the wrong region."
                )
            if "code 15" in str(err):
                raise FileSystemError(
                    "The DataRepo has been configured with the wrong credentials."
                )
        files_df = DataFrame(
            [
                {
                    "base_name": file.base_name,
                    "mtime": file.mtime,
                    "type": FileSystem.file_types[file.type],
                    "extension": file.extension,
                    "size": file.size,
                    "handle": file,
                }
                for file in files
            ]
        )
        if files_df.empty:
            columns = ["base_name", "mtime", "type", "extension", "size", "handle"]
            files_df = DataFrame(columns=columns)
        if exclude_projects:
            files_df = files_df[files_df.base_name != "Projects"]
        return files_df

    def write_file_to_remote(self, file, file_key, metadata):
        """Writes a file to the bucket.

        Parameters
        ----------
        file : str
            The name of the file.
        file_key : str
            The name of the file to be written in the bucket.
        metadata : dictionary
            The metadata of the object.
        """
        if self.local:
            Path(self.bucket + f"/{file_key}").parents[0].mkdir(
                parents=True, exist_ok=True
            )
        with self.fs.open_output_stream(
            self.bucket + f"/{file_key}", compression=None, metadata=metadata
        ) as fw, open(file, "rb") as local_file:
            fw.write(local_file.read())

    def write_file_from_remote(self, file_key, path):
        """Write a file locally from the remote bucket.

        Parameters
        ----------
        file_key : str
            The name of the file in the bucket.
        path : str
            The path to the location to write the file locally.
        """
        file_stream = self.fs.open_input_stream(self.bucket + f"/{file_key}")
        with open(path, "wb") as write_path:
            write_path.write(file_stream.read())

    def write_dataset(self, data_object, file_key, target_format="parquet"):
        """Writes a dataset to the bucket.

        Parameters
        ----------
        data_object : {numpy.ndarray, pandas.DataFrame, pandas.Series, torch.Tensor}
            The data object to be written to the bucket.
        file_key : str
            The name of the data object in the bucket.
        target_format : str {parquet, csv, json}
            Format of the target object to be written.
        """
        target_format, file_key = infer_target_format(
            file_key=file_key, target_format=target_format
        )

        if type(data_object) is not DataFrame:
            if type(data_object) not in (np.ndarray, Series):
                data_object = data_object.numpy()
            data_object = DataFrame(data_object)

        with self.fs.open_output_stream(f"{self.bucket}/{file_key}") as fw:
            if target_format == "json":
                data_object.to_json(fw, orient="records")
                return

            table = pyarrow.Table.from_pandas(data_object)
            if target_format == "parquet":
                parquet.write_table(table, fw)
            elif target_format == "csv":
                csv.write_csv(table, fw)
