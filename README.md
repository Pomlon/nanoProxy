# nanoProxy

NanoProxy sits on your server and proxies traffic only to defined interfaces/endpoints.

The aim here is to be dead simple, and used in situations where stuff like HAProxy and Traefik are way overkill.

## How to use
config.yaml needs to be sitting next to the executable.

```
listenInterface: localhost
listenPort: 8080
endpoints:
  "test1": "https://serverOrInterface/enpoint1"
  "test2": "https://serverOrInterface/endpoint2"
```
In the above the Proxy will listen on 127.0.0.1:8080 and when it receives traffic to http://localhost:8080/test1 
it will display contents of https://serverOrInterface/enpoint1 and so on. If no rule matches it will return an error message.

If you wish to make yours a service I recommend https://nssm.cc/ on Windows, or systemd on Linux
