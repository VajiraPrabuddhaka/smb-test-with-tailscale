## Run Tailscale Proxy

```shell
docker run -d -e TS_AUTH_KEY='<tailscale_auth_key>' -v ./config.yaml:/config.yaml -v tailscale-local:/.local -v tailscale-run:/var/run/tailscale -p 4789:4789 --name tailscale-proxy vajiraprabuddhaka/tailscale-proxy:latest
```