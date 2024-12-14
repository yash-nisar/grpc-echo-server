# gRPC Go Streaming APIs

The primary purpose of this readme is for testing grpc streaming scenarios with envoy. This was written while testing gRPC streaming APIs with CNAG.

## Build and push image to docker registry

To build and push a mutliplatform image:

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t yourrepo.azurecr.io/external/grpcecho:e2e \
  --push \
  .
```

## Testing all scenarios

1. Unary Echo (simplest, just one request-response):

```bash
# Regular command
grpcurl -d '{"message": "hello"}' -insecure 172.18.255.203:443 grpc.examples.echo.Echo/UnaryEcho

# Verbose mode for debugging
grpcurl -v -d '{"message": "hello"}' -insecure 172.18.255.203:443 grpc.examples.echo.Echo/UnaryEcho
```

2. Server Streaming (server will send responses for 10 seconds):

```bash
# Regular command
grpcurl -d '{"message": "hello"}' -insecure 172.18.255.203:443 grpc.examples.echo.Echo/ServerStreamingEcho

# Verbose mode
grpcurl -v -d '{"message": "hello"}' -insecure 172.18.255.203:443 grpc.examples.echo.Echo/ServerStreamingEcho
```

3. Client Streaming:

```bash
# Interactive mode - type messages manually, press Ctrl+D (Unix) or Ctrl+Z (Windows) when done
grpcurl -d @ -insecure 172.18.255.203:443 grpc.examples.echo.Echo/ClientStreamingEcho
{"message": "message 1"}
{"message": "message 2"}
{"message": "message 3"}
# Press Ctrl+D or Ctrl+Z here

# Automated mode (Unix/Linux/MacOS)
(echo '{"message": "automated 1"}'; 
 echo '{"message": "automated 2"}'; 
 echo '{"message": "automated 3"}') | grpcurl -d @ -insecure 172.18.255.203:443 grpc.examples.echo.Echo/ClientStreamingEcho

# Automated with delays (more realistic)
(echo '{"message": "automated 1"}'; sleep 1; 
 echo '{"message": "automated 2"}'; sleep 1; 
 echo '{"message": "automated 3"}') | grpcurl -d @ -insecure 172.18.255.203:443 grpc.examples.echo.Echo/ClientStreamingEcho
```

4. Bidirectional Streaming:

```bash
# Interactive mode - type messages and see responses in real-time
grpcurl -d @ -insecure 172.18.255.203:443 grpc.examples.echo.Echo/BidirectionalStreamingEcho
{"message": "bi-di message 1"}
{"message": "bi-di message 2"}
{"message": "bi-di message 3"}
# Press Ctrl+D or Ctrl+Z when done

# Automated mode (Unix/Linux/MacOS)
(echo '{"message": "bi-di 1"}'; 
 echo '{"message": "bi-di 2"}'; 
 echo '{"message": "bi-di 3"}') | grpcurl -d @ -insecure 172.18.255.203:443 grpc.examples.echo.Echo/BidirectionalStreamingEcho

# Automated with delays (better for observing bi-directional behavior)
(echo '{"message": "bi-di 1"}'; sleep 2; 
 echo '{"message": "bi-di 2"}'; sleep 2; 
 echo '{"message": "bi-di 3"}'; sleep 2) | grpcurl -d @ -insecure 172.18.255.203:443 grpc.examples.echo.Echo/BidirectionalStreamingEcho
```

## Utility Commands

```bash
# List all available services
grpcurl -insecure 172.18.255.203:443 list

# Describe the Echo service
grpcurl -insecure 172.18.255.203:443 describe grpc.examples.echo.Echo

# Describe a specific method
grpcurl -insecure 172.18.255.203:443 describe grpc.examples.echo.Echo.UnaryEcho
```
