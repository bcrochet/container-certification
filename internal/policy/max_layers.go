package policy

import (
	"context"
	"fmt"

	"github.com/opdev/container-certification/internal/log"
	"github.com/opdev/knex/types"

	"github.com/go-logr/logr"
	cranev1 "github.com/google/go-containerregistry/pkg/v1"
)

const (
	acceptableLayerMax = 40
)

var _ types.Check = &MaxLayersCheck{}

// UnderLayerMaxCheck ensures that the image has less layers in its assembly than a predefined maximum.
type MaxLayersCheck struct{}

func (p *MaxLayersCheck) Validate(ctx context.Context, imgRef types.ImageReference) (bool, error) {
	layers, err := p.getDataToValidate(imgRef.ImageInfo)
	if err != nil {
		return false, fmt.Errorf("could not get image layers: %v", err)
	}

	return p.validate(ctx, layers)
}

func (p *MaxLayersCheck) getDataToValidate(image cranev1.Image) ([]cranev1.Layer, error) {
	return image.Layers()
}

func (p *MaxLayersCheck) validate(ctx context.Context, layers []cranev1.Layer) (bool, error) {
	logr.FromContextOrDiscard(ctx).V(log.DBG).Info("number of layers detected in image", "layerCount", len(layers))
	return len(layers) <= acceptableLayerMax, nil
}

func (p *MaxLayersCheck) Name() string {
	return "LayerCountAcceptable"
}

func (p *MaxLayersCheck) Metadata() types.Metadata {
	return types.Metadata{
		Description:      fmt.Sprintf("Checking if container has less than %d layers.  Too many layers within the container images can degrade container performance.", acceptableLayerMax),
		Level:            "better",
		KnowledgeBaseURL: certDocumentationURL,
		CheckURL:         certDocumentationURL,
	}
}

func (p *MaxLayersCheck) Help() types.HelpText {
	return types.HelpText{
		Message:    "Check LayerCountAcceptable encountered an error. Please review the preflight.log file for more information.",
		Suggestion: "Optimize your Dockerfile to consolidate and minimize the number of layers. Each RUN command will produce a new layer. Try combining RUN commands using && where possible.",
	}
}
