import hypertrain


class TestImport:
    def test_hypertrain_version(self):
        version = hypertrain.__version__
        assert version == "0.1.0"
