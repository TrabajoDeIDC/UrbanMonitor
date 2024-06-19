# Nix file for creating development dependencies for Urban Server project

let
	pkgs = import <nixpkgs> {};
	nixpkgsUnstable = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-unstable";
	pkgsUnstable = import nixpkgsUnstable {};
	stablePkgs = with pkgs; [
		go
	];
	unstablePkgs = with pkgsUnstable; [
		vscode
	];

in
	pkgs.mkShell {
		packages = stablePkgs ++ unstablePkgs;
	}
