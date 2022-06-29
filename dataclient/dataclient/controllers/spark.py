from functools import wraps
from dataclasses import dataclass

import os
import pyspark


def spark_decorator(func, spark):
    @wraps(func)
    def wrapper(self, *args, **kwargs):

        args_list = list(args)
        for i in range(len(args)):
            if func.__code__.co_varnames[i + 1] in ["path", "paths"]:
                if isinstance(args_list[i], str):
                    args_list[i] = spark.root_path + args_list[i]
                elif isinstance(args_list[i], list):
                    args_list[i] = [spark.root_path + path for path in args_list[i]]
                else:
                    raise TypeError(
                        "The path or paths must be either a string or list."
                    )
            if func.__code__.co_varnames[i + 1] == "key":
                if args_list[i] == "path":
                    args_list[i + 1] = spark.root_path + args_list[i + 1]
        args = tuple(args_list)
        if "path" in kwargs:
            kwargs["path"] = spark.root_path + kwargs["path"]
        if "paths" in kwargs:
            kwargs["paths"] = spark.root_path + kwargs["paths"]
        if "keys" in kwargs:
            if kwargs["keys"] == spark.root_path:
                kwargs["value"] = spark.root_path + kwargs["value"]

        return func(self, *args, **kwargs)

    return wrapper


@dataclass
class SparkController:
    def __init__(self, dataclient, decorate=True):
        if decorate:
            self.decorate_spark_methods(pyspark.sql.readwriter.DataFrameReader)
            self.decorate_spark_methods(pyspark.sql.readwriter.DataFrameWriter)
        self.dataclient = dataclient
        self.configure_spark()

    def decorate_spark_methods(self, spark_class):
        readwrite_methods = [
            "csv",
            "json",
            "load",
            "option",
            "orc",
            "parquet",
            "save",
            "text",
        ]
        method_list = [
            func
            for func in dir(spark_class)
            if callable(getattr(spark_class, func)) and func in readwrite_methods
        ]

        for method in method_list:
            decotared_method = spark_decorator(getattr(spark_class, method), self)
            setattr(spark_class, method, decotared_method)

    def configure_spark(self):
        self.root_path = "s3a://" + volume_name + "/"

        access_key = self.dataclient.access_key
        secret_key = self.dataclient.secret_access_key
        volume_region = "s3-" + self.dataclient.volume_region + ".amazonaws.com"

        self.dataclient.spark._sc._jsc.hadoopConfiguration().set(
            "fs.s3a.access.key", access_key
        )
        self.dataclient.spark._sc._jsc.hadoopConfiguration().set(
            "fs.s3a.secret.key", secret_key
        )
        self.dataclient.spark._sc._jsc.hadoopConfiguration().set(
            "fs.s3a.endpoint", volume_region
        )

    def set_local(self):
        self.root_path = os.getcwd()
