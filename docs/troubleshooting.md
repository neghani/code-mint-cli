# Troubleshooting

## Browser did not open

Copy the printed login URL and open it manually in your browser.

## Callback port blocked

Allow local loopback traffic to `127.0.0.1` and retry login.

## Token revoked or expired

Run `codemint auth login` again to issue a new token.
