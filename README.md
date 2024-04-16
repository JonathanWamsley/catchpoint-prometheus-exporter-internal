# catchpoint-prometheus-exporter
Exports [Catchpoint](https://www.catchpoint.com) monitoring metrics for nodes via HTTP for Prometheus consumption.

## Configuration
### Command line flags
The exporter can be configured through the following command line flags:
```
  -h, --help                Show context-sensitive help.
      --web.listen-address=":8080"
                            Address on which to expose metrics and web interface.
      --web.telemetry-path="/metrics"
                            Path under which to expose metrics.
      --bearer-token=""     The bearer token for authentication with Catchpoint.
      --node-ids=""         Comma-separated list of node IDs to include in the metrics.
      --request-delay=1     Delay between API requests in seconds to manage rate limiting.
      --port="8080"         The port to bind the HTTP server.
      --version             Show application version.
      --log.level=info      Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt   Output format of log messages. One of: [logfmt, json]
```

Example usage:
```sh
./catchpoint_exporter --bearer-token="YourBearerToken" --node-ids="123,456" --port="9091"
```

### Environment Variables
The exporter can also be configured using environment variables:

| Name                                      | Description                                        |
|-------------------------------------------|----------------------------------------------------|
| CATCHPOINT_EXPORTER_BEARER_TOKEN          | The bearer token for authentication with Catchpoint.|
| CATCHPOINT_EXPORTER_NODE_IDS              | Comma-separated list of node IDs to include in the metrics. |
| CATCHPOINT_EXPORTER_PORT                  | The port to bind the HTTP server.                  |
| CATCHPOINT_EXPORTER_REQUEST_DELAY         | Delay between API requests in seconds.             |
| CATCHPOINT_EXPORTER_WEB_TELEMETRY_PATH    | Path under which to expose metrics.                |

Example usage:
```sh
CATCHPOINT_EXPORTER_BEARER_TOKEN="YourBearerToken" \
CATCHPOINT_EXPORTER_NODE_IDS="123,456" \
CATCHPOINT_EXPORTER_PORT="9091" \
./catchpoint_exporter
```
