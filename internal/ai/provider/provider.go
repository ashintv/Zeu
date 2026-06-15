package provider

import (
	"context"

	"github.com/ashintv/Zeu/internal/types"
)


type Provider interface {
	Info() types.ProviderInfo
	Process(ctx context.Context, req *types.AiRequest, streamCh chan<- types.AiResponse) (err error)
	Default() *types.DefaultOptions
}
