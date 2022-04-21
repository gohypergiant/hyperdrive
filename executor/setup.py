import setuptools

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()


def get_version(rel_path):
    with open(rel_path, "r", encoding="utf-8") as fh:
        for line in fh.readlines():
            if line.startswith("__version__"):
                return line.split('"')[1]
        else:
            raise RuntimeError("Unable to find version string.")


setuptools.setup(
    name="executor",
    version=get_version("executor/__init__.py"),
    author="@hypergiant",
    description="Scheduler Package",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/gohypergiant/hyperdrive",
    project_urls={"Bug Tracker": "https://github.com/gohypergiant/hyperdrive/issues",},
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    include_package_data=True,
    install_requires=[
        "croniter",
        "databases",
        "jupyterlab",
        "pandas>=1.2.4",
        "papermill",
        "typer",
    ],
    packages=setuptools.find_packages(),
    python_requires=">=3.9",
)
