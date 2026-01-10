package verify

import (
	"fmt"

	"github.com/zevwings/workflow/internal/config"
	"github.com/zevwings/workflow/internal/prompt"
)

// VerifyLogConfig verifies log configuration
func VerifyLogConfig(logConfig *config.LogConfig) {
	if logConfig.Level == "" {
		return
	}
	msg := prompt.GetMessage()
	msg.Info("Log Configuration")
	msg.Info("%s", fmt.Sprintf("Log Level: %s", logConfig.Level))
	msg.Break()
}
