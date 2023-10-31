// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package extloadtest

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/event-kit/go/event_kit_api"
	extension_kit "github.com/steadybit/extension-kit"
	"github.com/steadybit/extension-kit/exthttp"
	"net/http"
)

func RegisterEventListenerHandlers() {
	exthttp.RegisterHttpHandler("/events/log", handle())
}

func handle() func(w http.ResponseWriter, r *http.Request, body []byte) {
	return func(w http.ResponseWriter, r *http.Request, body []byte) {
		var event event_kit_api.EventRequestBody
		err := json.Unmarshal(body, &event)
		if err != nil {
			exthttp.WriteError(w, extension_kit.ToError("Failed to decode event request body", err))
			return
		}

		log.Info().Str("event", event.EventName).Msg("Received event.")
		exthttp.WriteBody(w, "{}")
	}
}
