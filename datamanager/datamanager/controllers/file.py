class FileController:
    def get_file(self, remoteFileKey, localFilePath=None):
        """Gets a file from a Data Repo and saves it locally.

        Parameters
        ----------
        remoteFileKey : str
            The name of the file in the Data Repo.
        localFilePath : str
            The path to the location to write the file locally.

        Returns
        -------
        bool
            Returns True on success.
        """

        if localFilePath is None:
            localFilePath = f"/tmp/{remoteFileKey}"
        return self.filesystem.write_file_from_remote(remoteFileKey, localFilePath)

    def get_file_info(self, file):
        """Gets file information.

        Read more about Pyarrow FileInfo:
        https://arrow.apache.org/docs/python/generated/pyarrow.fs.FileInfo.html

        Parameters
        ----------
        file : str
            Refers to an arbitrary file in a Data Repo. Can be nested.

        Returns
        -------
        dictionary
            File information.
        """

        return self.filesystem.get_file_info(file)

    def get_object(self, file_key, path=None):
        """Gets an object from the Data Repo and save it locally.

        Parameters
        ----------
        file_key : str
            Refers to an object in a Data Repo.

        path : str
            Refers to path to save the object locally.

        Returns
        -------
        bool
            Returns True on success.
        """
        if path is None:
            path = f"/tmp/{file_key}"
        return self.filesystem.write_file_from_remote(file_key, path)

    def write_object(self, obj=None, file_key=None, metadata=None):
        """Writes an object to the Data Repo.

        Parameters
        ----------
        obj : str, optional
            Refers to the path of the object (default is None).
        file_key : str, optional
            Refers to name of artifact to be written on Data Repo. (default is None).
        metadata : dictionary, optional
            The metadata of the artifact (default is None).

        Returns
        -------
        bool
            Returns True on success.
        """
        return self.filesystem.write_file_to_remote(obj, file_key, metadata)
