package harness

import (
	"github.com/ashintv/Zeu/internal/ai"
	"github.com/ashintv/Zeu/internal/types"
)



type Agent struct {
	state   []types.Coversation
	Ai      ai.AI
	MaxIter int
	System  string
	tools []types.Tool
 
} 

func CreateAgent() *Agent {
	return &Agent{}
}

func (a *Agent) Invoke() {

}
