# Introduction
Light weight http rate limiting proxy.

The proxy will perform rate limiting based on the rules defined in the configuration [file](config.yml). If no rule match the request, it'll proxy transparently.

# Launch in development mode
```sh
make dev
```

# HTTP Request
Perform a http request using HTTPie
```sh
http :8080
```

# Trigger rate limiting
```sh
seq 10 | xargs -P 4 -I '{}' http post :8080/health
```

# Credits
* https://github.com/didip/tollbooth

# TODO
* Comment exported types and functions
* Write tests
