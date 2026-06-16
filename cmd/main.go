package cmd

import (
	"context"

	"github.com/ashintv/Zeu/internal/ai"
	"github.com/ashintv/Zeu/internal/logger"
	"github.com/ashintv/Zeu/internal/types"
)

func main() {
	Ai := ai.NewAI()
	logger.Info("CREATED AI", Ai)

	logger.Info("Invoking ... ")

	messages := []types.Coversation{{
		Role:    "user",
		Content: "Hello how can u help me",
	}}

	err , resChan := Ai.Invoke(context.Background(), ai.WithMessages(messages))

	if err != nil{
		logger.Error(err)
		return

	}

	for res := range resChan{
		logger.Info(res)
	}

}
