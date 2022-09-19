"""This script unpacks hyperpack onnx model artifacts that were 
originated from Hypergiant's AutoML Univariate TimeSeries.
NOTE: Since time series is comprised of multiple model artifacts,
this script currently only opens up 1 model artifact at a time. 
Therefore, if multiple model artifacts need to be unpacked, a for-loop
or iteration over multiple artifact files is required using the code herein.
"""

from icecream import ic
import onnx
from onnx2torch import convert

from utilities import UnpackHyperpack


class UnpackAutoMLHyperpack(UnpackHyperpack):
    """Unpack Hyperpack AutoML Time Series model"""

    def __init__(self, hyperpack_filedir: str = None, verbose: bool = False):
        """
        :hyperpack_filedir: file directory where a hyperpack.zip file is located
        :verbose: print if True, else suppress printing
        """
        super(UnpackAutoMLHyperpack, self).__init__(hyperpack_filedir, verbose)
        self.onnx_models = []
        self.automl_model = None

    def load_automl_onnx_model(self):
        """Load onnx model
        """
        ic(self.onnx_filepaths)
        for onnx_f in self.onnx_filepaths:
            onnx_model = onnx.load(onnx_f)
            onnx.checker.check_model(onnx_model)
            self.onnx_models.append(onnx_model)
            
        ic(type(self.onnx_models))

    def convert_onnx_to_automl_torch(self, ith_model: int = 0):
        """Convert onnx model to torch model ()
        :ith_model: if there is greater than one onnx_model in self.onnx_models, use the ith index
        """
        automl_model = self.onnx_models[ith_model]
        self.automl_model = convert(automl_model)

        if self.verbose:
            ic(self.automl_model)

        return self.automl_model
