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
	case "pytorch":
		return ImageOptions{
			Image:   "ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-pytorchstable",
			RepoTag: "hyperdrive-jupyter-cpu-pytorch",
			Profile: "pytorch-cpu",
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
			RepoTag: "hyperdrive-jupyter-cpu-local",
			Profile: "minimal",
		}
	}
}
