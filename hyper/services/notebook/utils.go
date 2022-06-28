package notebook

import (
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"strings"

	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
)

func GetNotebookName(manifestPath string) string {
	return strings.ToLower(manifest.GetName(manifestPath))
}

func GetNotebookImageOptions(flavor string) types.ImageOptions {

	// TODO: We could just make this a map, and even read it in from a remote endpoint
	switch flavor {
	case "pytorch":
		return types.ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-pytorchstable",
			RepoTag: "hyperdrive-jupyter-cpu-pytorch",
			Profile: "pytorch-cpu",
		}
	case "dev":
		return types.ImageOptions{
			Image:   "cpu-local:latest",
			RepoTag: "cpu-local:latest",
			Profile: "dev",
		}
	default:
		return types.ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-basestable",
			RepoTag: "hyperdrive-jupyter-cpu-local",
			Profile: "minimal",
		}
	}
}
