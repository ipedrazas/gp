package remote

import (
	"context"
	"strings"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

func PullCMD(image string, output string) []string {
	command := []string{
		"oras",
		"pull",
		image,
		"-o",
		output,
	}

	return command
}

func FetchOCI(ociArtifactUri, output string) (string, error) {
	ctx := context.Background()

	repoUri := strings.ReplaceAll(ociArtifactUri, "oci://", "")
	repo, err := remote.NewRepository(repoUri)
	if err != nil {
		return "", err
	}

	dst, err := file.New(output)
	if err != nil {
		return "", err
	}

	copyOptions := oras.DefaultCopyOptions
	desc, err := oras.Copy(ctx, repo, repo.Reference.Reference, dst, repo.Reference.Reference, copyOptions)
	if err != nil {
		return "", err
	}

	return desc.Digest.String(), nil
}
