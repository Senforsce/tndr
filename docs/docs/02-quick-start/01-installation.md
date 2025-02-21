# Installation

## go install

With Go 1.20 or greater installed, run:

```sh
go install github.com/senforsce/tndr/cmd/t1@latest
```

## Github binaries

Download the latest release from https://github.com/senforsce/tndr/releases/latest

## Nix

t1 provides a Nix flake with an exported package containing the binary at https://github.com/senforsce/tndr/blob/main/flake.nix

```sh
nix run github:senforsce/t1
```

t1 also provides a development shell which includes a Neovim configuration setup to use the t1 autocompletion features.

```sh
nix develop github:senforsce/t1
```

To install in your Nix Flake:

This flake exposes an overlay, so you can add it to your own Flake and/or NixOS system.

```nix
{
  inputs = {
    ...
    t1.url = "github:senforsce/t1";
    ...
  };
  outputs = inputs@{
    ...
  }:

  # For NixOS configuration:
  {
    # Add the overlay,
    nixpkgs.overlays = [
      inputs.t1.overlays.default
    ];
    # and install the package
    environment.systemPackages = with pkgs; [
      t1
    ];
  };

  # For a flake project:
  let
    forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
      inherit system;
      pkgs = import nixpkgs { inherit system; };
    });
    t1 = system: inputs.t1.packages.${system}.t1;
  in {
    packages = forAllSystems ({ pkgs, system }): {
      myNewPackage = pkgs.buildGoModule {
        ...
        preBuild = ''
          ${t1 system}/bin/t1 generate
        '';
      };
    };

    devShell = forAllSystems ({ pkgs, system }):
      pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          (t1 system)
        ];
      };
  };
}
```
