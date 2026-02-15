# Enterprise

## Proxy and certificates

The CLI respects standard `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY` environment variables.

For private CAs, ensure your system trust store includes the required certificate chain.

## Security

- Tokens are stored in OS secure storage on supported platforms.
- CLI output redacts bearer tokens in common error paths.
