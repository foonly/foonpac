# Foonpac

A declarative package wrapper for pacman.

## Configuration

Foonpac uses a configuration file and group files to manage your packages.

### Main Config

Located at `~/.config/foonpac/config.yaml`.

```yaml
hosts:
  my-desktop:
    - base
    - desktop
    - dev
  my-laptop:
    - base
    - laptop
```

### Group Files

Located at `~/.config/foonpac/groups/<group_name>.txt`. These are plain text files with one package per line.

Example `base.txt`:

```text
linux
linux-firmware
base
vim
git
```

## Usage

### Sync

Find declared packages that are not installed and installed packages that are not declared.

```bash
foonpac sync
```

### List

List packages based on their management status.

- `foonpac list unmanaged`: Lists installed native packages not in any group file for this host.
- `foonpac list managed`: Lists installed native packages that are in the group files.
- `foonpac list missing`: Lists packages in group files that are not installed.
- `foonpac list local`: Lists installed packages not in official repositories (AUR/manual).
- `foonpac list dependencies`: Lists installed dependencies that are required by managed packages. (Requires `pacman-contrib` for `pactree`).
