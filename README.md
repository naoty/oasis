# oasis

A HTTP proxy building docker containers for each commits

## Usage

```sh
% oasis start \
  --proxy feature.example.com \
  --container-host "$(docker-machine ip dev)" \
  --repository github.com/naoty/sample_rails_app
```

When the proxy receives a request to `http://feature.example.com`, it will do below things.

1. Clone `https://github.com/naoty/sample_rails_app`.
2. Checkout the repo to `feature`.
3. Build a docker container based on `Dockerfile` at the repo.
4. Redirect the request to the docker container.

## Installation

```sh
% go get github.com/naoty/oasis
```
