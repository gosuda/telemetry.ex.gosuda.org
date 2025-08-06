package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

var _index_html = []byte(
	`<!DOCTYPE html>
<html>
<head>
	<title>GoSuda Telemetry Service</title>
</head>
<body>
	<h1>GoSuda Telemetry Service</h1>
	<hr/>
	<p>Welcome to the GoSuda Telemetry Service!</p>
	<p>Public APIs:</p>
	<ul>
		<li>GET <a href="/healthz">/healthz</a> - Check the health of the service</li>
		<li>GET <a href="/idz">/idz</a> - Generate a new randflake ID</li>
	</ul>
</body>
</html>`,
)

func IndexHandler(is types.InternalServiceProvider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(_index_html)
	}
}
