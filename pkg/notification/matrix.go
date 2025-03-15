package notification

import (
	"bytes"
	"context"
	"github.com/jagac/pfinance/pkg/config"
	"github.com/jagac/pfinance/pkg/rand"

	"encoding/json"
	"fmt"
	"net/http"
)

type MatrixNotifier struct {
	config *config.GlobalConfig
}

func NewMatrixNotifier(config *config.GlobalConfig) *MatrixNotifier {
	return &MatrixNotifier{
		config: config,
	}
}

func (m *MatrixNotifier) Send(ctx context.Context, to, subject, body string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

		txnID := rand.GenerateTransactionID()
		newMatrixUrl := fmt.Sprintf("https://matrix.intra.kulvej.net/_matrix/client/r0/rooms/%s/send/m.room.message/%s",
			"!XlrscsEPQGKNCWMlKf:matrix.intra.kulvej.net", txnID)

		newPayload := map[string]any{
			"msgtype": "m.text",
			"body":    fmt.Sprintf("[New message] %s\nSubject: %s\n%s", to, subject, body),
			"format":  "org.matrix.custom.html",
			"formatted_body": fmt.Sprintf(
				"New message to <b>%s</b><br>Subject: <b>%s</b><br><pre><code>%s</code></pre>",
				to, subject, body,
			),
		}

		newJsonPayload, err := json.Marshal(newPayload)
		if err != nil {
			return err
		}

		newReq, err := http.NewRequestWithContext(ctx, "PUT", newMatrixUrl, bytes.NewBuffer(newJsonPayload))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		newReq.Header.Set("Content-Type", "application/json")
		newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.config.MatrixAccessToken))
		newClient := &http.Client{}

		resp, err := newClient.Do(newReq)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.Status != http.StatusText(http.StatusOK) {
			return fmt.Errorf("failed to execute request: %w", err)
		}

		return nil
	}
}
