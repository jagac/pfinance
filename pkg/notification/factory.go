package notification

import"github.com/jagac/pfinance/pkg/config"

func NotifierFactory(notifierType string, config *config.GlobalConfig) Notifier {
	switch notifierType {
	case "email":
		return NewEmailNotifier(config)
	case "matrix":
		return NewMatrixNotifier(config)
	default:
		return nil
	}
}
