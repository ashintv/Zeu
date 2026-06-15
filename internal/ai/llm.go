package ai

import (
	"context"

	"github.com/ashintv/Zeu/internal/ai/provider"
	"github.com/ashintv/Zeu/internal/logger"
	"github.com/ashintv/Zeu/internal/types"
)


type LLm struct{
	Provider provider.Provider
	logger *logger.Logger
} 


func (llm *LLm) Invoke(ctx context.Context, req *types.AiRequest) (error , <-chan types.AiResponse){
	streamCh := make(chan types.AiResponse)
	
	logger.Info("channel created invoking agent" , llm.Provider.Info())
	err := llm.Provider.Process(ctx ,  req , streamCh)
	
	return err, streamCh
 
}