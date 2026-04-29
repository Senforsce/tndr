# Installation

## go install (global)

With Go 1.24 or greater installed, run:

```bash
go install github.com/senforsce/tndr/cmd/t1@latest
```

This installs t1 into your path.

## go install (as tool)

To install t1 locally in your project, run:

```bash
go get -tool github.com/senforsce/tndr/cmd/t1@latest
```

:::info 
This uses the [tool directive](https://tip.golang.org/doc/modules/managing-dependencies#tools) feature of Go added in v1.24. 

To run t1 once installed, use `go tool t1` instead of `t1`.
:::

## GitHub binaries

Download the latest release from https://github.com/senforsce/tndr/releases/latest

## Nix

tndr provides a Nix flake with an exported package containing the binary at https://github.com/senforsce/tndr/blob/main/flake.nix

```bash
nix run github:senforsce/tndr
```

tndr also provides a development shell which includes all of the tools required to build tndr, e.g. go, gopls etc. but not tndr itself.

```bash
nix develop github:senforsce/tndr
```

To install in your Nix Flake:

This flake exposes an overlay, so you can add it to your own Flake and/or NixOS system.

```nix
{
  inputs = {
    ...
    tndr.url = "github:senforsce/tndr";
    ...
  };
  outputs = inputs@{
    ...
  }:

  # For NixOS configuration:
  {
    # Add the overlay,
    nixpkgs.overlays = [
      inputs.tndr.overlays.default
    ];
    # and install the package
    environment.systemPackages = with pkgs; [
      tndr
    ];
  };

  # For a flake project:
  let
    forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
      inherit system;
      pkgs = import nixpkgs { inherit system; };
    });
    tndr = system: inputs.tndr.packages.${system}.tndr;
  in {
    packages = forAllSystems ({ pkgs, system }: {
      myNewPackage = pkgs.buildGoModule {
        ...
        preBuild = ''
          ${tndr system}/bin/t1 generate
        '';
      };
    });

    devShell = forAllSystems ({ pkgs, system }:
      pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          (tndr system)
        ];
      };
  });
}
```

## Docker

A Docker container is pushed on each release to https://github.com/senforsce/tndr/pkgs/container/tndr

Pull the latest version with:

```bash
docker pull ghcr.io/senforsce/tndr:latest
```

To use the container, mount the source code of your application into the `/app` directory, set the working directory to the same directory and run `t1 generate`, e.g. in a Linux or Mac shell, you can generate code for the current directory with:

```bash
docker run -v `pwd`:/app -w=/app ghcr.io/senforsce/tndr:latest generate
```

If you want to build templates using a multi-stage Docker build, you can use the `tndr` image as a base image.

Here's an example multi-stage Dockerfile. Note that in the `generate-stage` the source code is copied into the container, and the `t1 generate` command is run. The `build-stage` then copies the generated code into the container and builds the application.

The permissions of the source code are set to a user with a UID of 65532, which is the UID of the `nonroot` user in the `ghcr.io/senforsce/tndr:latest` image.

Note also the use of the `RUN ["t1", "generate"]` command instead of the common `RUN t1 generate` command. This is because the tndr Docker container does not contain a shell environment to keep its size minimal, so the command must be ran in the ["exec" form](https://docs.docker.com/reference/dockerfile/#shell-and-exec-form).

```Dockerfile
# Fetch
FROM golang:latest AS fetch-stage
COPY go.mod go.sum /app
WORKDIR /app
RUN go mod download

# Generate
FROM ghcr.io/senforsce/tndr:latest AS generate-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["t1", "generate"]

# Build
FROM golang:latest AS build-stage
COPY --from=generate-stage /app /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app

# Test
FROM build-stage AS test-stage
RUN go test -v ./...

# Deploy
FROM gcr.io/distroless/base-debian12 AS deploy-stage
WORKDIR /
COPY --from=build-stage /app/app /app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
```
