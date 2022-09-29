import setuptools


def get_version(rel_path):
    with open(rel_path, "r", encoding="utf-8") as fh:
        for line in fh.readlines():
            if line.startswith("__version__"):
                return line.split('"')[1]
            else:
                raise RuntimeError("Unable to find version string.")


setuptools.setup(
    name="hyperpackage",
    version=get_version("hyperpackage/__init__.py"),
    author="@hypergiant",
    description="Hyperpackage packages",
    long_description="A package to create hyperpacks",
    long_description_content_type="text/markdown",
    url="https://github.com/gohypergiant",
    project_urls={"Bug Tracker": "https://github.com/gohypergiant/issues",},
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    include_package_data=True,
    install_requires=[
        "boto3",
        "icecream",
        "numpy",
        "onnx==1.11.0",
        "onnx2torch",
        "onnxruntime",
    ],
    extras_require={"torch": ["pytorch"],},
    packages=setuptools.find_packages(),
    python_requires=">=3.8",
)
