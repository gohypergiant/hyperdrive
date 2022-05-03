package notebook

import (
	"fmt"
	"os"
	"strings"

	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
)

func GetNotebookName(manifestPath string) string {
	if manifestPath == "local" {
		cwdPath, _ := os.Getwd()
		cwdName := strings.Replace(cwdPath, "/", "-", -1)
		return fmt.Sprintf("firefly-jupyter-%s", cwdName)
	}
	return manifest.GetName(manifestPath)
}

type ImageOptions struct {
	Image   string
	RepoTag string
	Profile string
}

func GetNotebookImageOptions(flavor string) ImageOptions {

	// TODO: We could just make this a map, and even read it in from a remote endpoint
	switch flavor {
	case "huggingface-pytorch":
		return ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-huggingface-pytorchstable",
			RepoTag: "hypergiant-jupyter-cpu-huggingface-pytorch",
			Profile: "hugging-face-pytorch-cpu",
		}
	case "huggingface-tensorflow":
		return ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-huggingface-tensorflowstable",
			RepoTag: "hypergiant-jupyter-cpu-huggingface-tensorflow",
			Profile: "hugging-face-tensorflow-cpu",
		}
	case "pytorch":
		return ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-pytorchstable",
			RepoTag: "hypergiant-jupyter-cpu-pytorch",
			Profile: "pytorch-cpu",
		}
	case "spark":
		return ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-sparkstable",
			RepoTag: "hypergiant-jupyter-cpu-spark",
			Profile: "spark-cpu",
		}
	case "tensorflow":
		return ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-tensorflowstable",
			RepoTag: "hypergiant-jupyter-cpu-tensorflow",
			Profile: "tensorflow-cpu",
		}
	case "xgboost":
		return ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-xgbooststable",
			RepoTag: "hypergiant-jupyter-cpu-xgboost",
			Profile: "xgboost-cpu",
		}
	case "dev":
		return ImageOptions{
			Image:   "cpu-local:latest",
			RepoTag: "cpu-local:latest",
			Profile: "dev",
		}
	default:
		return ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-localstable",
			RepoTag: "hypergiant-jupyter-cpu-local",
			Profile: "minimal",
		}
	}
}
