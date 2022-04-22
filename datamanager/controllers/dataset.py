from datetime import datetime
from pandas import read_csv
from pandas import read_json as pd_read_json
from pandas import read_excel as pd_read_excel
from pandas import ExcelFile
from pyarrow.lib import ArrowInvalid
from ..exceptions import DatasetError, DataRepoError, DataManagerError


class DatasetController:
    def delete_dataset(self, dataset):
        """Deletes a dataset from the default Data Repo.

        Parameters
        ----------
        dataset : str
            Refers to File Key in the Data Repo. Can be a nested value.

        Returns
        -------
        True, if successful.

        """

        file_info = self.filesystem.get_file_info(dataset)
        if file_info["type"] == "directory":
            return self.filesystem.delete_dir(dataset)
        if file_info["type"] == "file":
            return self.filesystem.delete_file(dataset)
        return True

    def list(self, subpath="", exclude_projects=True, full=False, recursive=False):
        """Lists all datasets of the default Data Repo.

        Parameters
        ----------
        subpath : str, optional
            Refers to a subpath in a Data Repo (default is "").
        exclude_projects : bool, optional
            Excluding the Project directory in results (default is True).
        full : bool
            Displaying all parameters returned by PyArrow (default is False).
        recursive : bool, optional
            The flag that indicates if should list content of subfolders (default is
            False).

        Returns
        -------
        pandas.DataFrame
            Datasets of the Data Repo.
        """

        bucket_list = self.filesystem.list(
            subpath=subpath, exclude_projects=exclude_projects, recursive=recursive
        )
        if full:
            return bucket_list
        return bucket_list[["base_name", "extension", "size"]]

    def load_dataset(
        self,
        key=None,
        path_or_paths=None,
        source_format=None,
        target_format="pandas",
        index=None,
        prepend_bucket=True,
        **kwargs,
    ):
        """Loads a dataset from the Data Repo.

        Parameters
        ----------
        key : str
            Refers to dataset.
        path_or_paths : str
            Refers to a directory of dataset files. Supports parquet only.
        source_format : {None, "csv", "json", "parquet", "xlsx"}, optional
            Refers to format of the dataset in the data repo (default is None).
        target_format : {"pandas", "numpy", "torch"}, optional
            Refers to the format in which the dataset will be loaded (default is
            "pandas").

        Returns
        -------
        pandas.DataFrame/numpy.ndarray/torch.Tensor
            Dataset of the Data Repo.
        """
        if path_or_paths is not None:
            if source_format is not None and source_format != "parquet":
                raise DataRepoError(
                    "Data Repos currently only support reading "
                    "from a directory of parquet files."
                )

            from pyarrow.parquet import ParquetDataset

            if prepend_bucket:
                path_or_paths = f"{self.filesystem.bucket}/{path_or_paths}"

            dataset = ParquetDataset(
                path_or_paths=path_or_paths, filesystem=self.filesystem.fs
            )
            table = dataset.read()

        else:
            if source_format is None:
                source_format = self.get_file_info(key)["extension"]

            if source_format not in ("csv", "json", "parquet", "xlsx"):
                raise DataRepoError(
                    "Data Repos currently only support reading "
                    "from csv, json, parquet, or xlsx."
                )

            if prepend_bucket:
                s3_key = f"{self.filesystem.bucket}/{key}"
            else:
                s3_key = key

            file_handle = self.filesystem.fs.open_input_file(s3_key)

            if source_format == "csv":
                from pyarrow.csv import read_csv

                table = read_csv(file_handle)

            elif source_format == "json":
                from pyarrow.json import read_json

                try:
                    table = read_json(file_handle)
                except ArrowInvalid:
                    self.get_file(key, "/tmp/tmp.json")
                    return pd_read_json("/tmp/tmp.json")

            elif source_format == "parquet":
                from pyarrow.parquet import ParquetDataset

                dataset = ParquetDataset(s3_key, filesystem=self.filesystem.fs)
                table = dataset.read()

            elif source_format == "xlsx":
                self.get_file(key, "/tmp/tmp.xlsx")
                excel_file = ExcelFile("/tmp/tmp.xlsx")
                if len(excel_file.sheet_names) > 1:
                    return excel_file
                else:
                    return pd_read_excel("/tmp/tmp.xlsx", **kwargs)

        df = table.to_pandas()

        if index is not None:
            df = df.set_index(index)
        if target_format == "pandas":
            return df
        elif target_format == "numpy":
            return df.to_numpy()
        elif target_format == "torch":
            try:
                from torch import tensor

                df = df.select_dtypes(include="number")

                return tensor(df.values)
            except ImportError:
                raise DataManagerError("Loading PyTorch tensors requires PyTorch.")

    def write_dataset(
        self,
        data,
        file_key: str,
        target_format: str = "parquet",
        overwrite: bool = False,
    ) -> bool:
        """Writes a dataset to the default Data Repo.

        Parameters
        ----------
        data : {numpy.ndarray, pandas.DataFrame, pandas.Series, torch.Tensor}
            Refers to dataset to be written to the Data Repo.
        file_key : str
            Refers to the name of the dataset to be written to the Data Repo.
        target_format : str
            Refers to the file format of the dataset. Can be either csv, parquet or json.
        overwrite : bool
            Boolean flag that allows for overwriting of existing dataset.

        Returns
        -------
        bool
            Returns True on success.

        Raises
        ------
        DataRepoError
            If a default Data Repo is not set.
        """
        # TODO: the overwrite method here is not working properly
        current_datasets = self.list()
        dataset_to_write = file_key + "." + target_format
        if not overwrite and dataset_to_write in current_datasets["base_name"].values:
            raise DatasetError(
                f"The dataset '{dataset_to_write}' already exists. Please set "
                "overwrite=TRUE to overwrite this dataset."
            )

        self.filesystem.write_dataset(
            data_object=data, file_key=file_key, target_format=target_format
        )
        return True

    def write_remote_dataset(self, uri, dataset, header="infer", names=None):
        """Writes a remote dataset to the default Data Repo.

        Parameters
        ----------
        uri: str
            URI of a remote tabular dataset.
        dataset : str
            Refers to the name of the dataset to be written to the Data Repo.
        header: int, list of int, default ‘infer’
            Row number(s) to use as the column names, and the start of the data.
            Default behavior is to infer the column names: if no names are passed
            the behavior is identical to header=0 and column names are inferred from
            the first line of the file, if column names are passed explicitly then
            the behavior is identical to header=None. Explicitly pass header=0 to be
            able to replace existing names. The header can be a list of integers that
            specify row locations for a multi-index on the columns e.g. [0,1,3].
            Intervening rows that are not specified will be skipped (e.g. 2 in this
            example is skipped). Note that this parameter ignores commented lines
            and empty lines if skip_blank_lines=True, so header=0 denotes the first
            line of data rather than the first line of the file.
        names : array-like, optional
            List of column names to use. If the file contains a header row, then you
            should explicitly pass header=0 to override the column names. Duplicates
            in this list are not allowed.

        Returns
        -------
        bool
            Returns True on success.

        Raises
        ------
        DataRepoError
            If a default Data Repo is not set.
        """

        if names is None:
            dataframe = read_csv(uri, header=header)
            print("Creating DataFrame for writing to DataRepo using header: ")
            print(header)
        else:
            dataframe = read_csv(uri, header=header, names=names)
            print("Creating DataFrame for writing to DataRepo using header: ")
            print(header)
            print("with names: ")
            print(names)
        return self.write_dataset(dataframe, dataset)
