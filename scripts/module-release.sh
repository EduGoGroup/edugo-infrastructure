#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

usage() {
  cat <<USAGE
Uso:
  scripts/module-release.sh prepare <module> <version>
  scripts/module-release.sh notes <module> <version>
  scripts/module-release.sh github <module> <version>

Ejemplos:
  scripts/module-release.sh prepare postgres v0.62.0
  scripts/module-release.sh notes tools/mock-generator v0.1.3
  scripts/module-release.sh github mongodb v0.54.0
USAGE
}

fail() {
  echo "Error: $*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "No se encontró '$1' en PATH"
}

validate_version() {
  local version="$1"
  [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail "VERSION debe tener formato vX.Y.Z"
}

resolve_module() {
  MODULE_PATH="$1"
  MODULE_DIR="$ROOT_DIR/$MODULE_PATH"
  CHANGELOG_FILE="$MODULE_DIR/CHANGELOG.md"

  [[ -d "$MODULE_DIR" ]] || fail "No existe el módulo '$MODULE_PATH'"
  [[ -f "$CHANGELOG_FILE" ]] || fail "No existe $CHANGELOG_FILE"
  grep -q '^## \[Unreleased\]' "$CHANGELOG_FILE" || fail "CHANGELOG sin sección Unreleased en $MODULE_PATH"
}

extract_section() {
  local section="$1"
  awk -v section="$section" '
    $0 ~ "^## \\[" section "\\]" {flag=1; next}
    flag && /^## \[/ {exit}
    flag {print}
  ' "$CHANGELOG_FILE"
}

prepare_release() {
  local version="$1"
  local version_number="${version#v}"
  local today
  today="$(date +%F)"

  grep -q "^## \[$version_number\]" "$CHANGELOG_FILE" && fail "La versión $version_number ya existe en $CHANGELOG_FILE"

  local unreleased
  unreleased="$(extract_section "Unreleased")"
  [[ -n "$(printf '%s' "$unreleased" | tr -d '[:space:]')" ]] || fail "La sección Unreleased está vacía en $CHANGELOG_FILE"

  local tmp
  tmp="$(mktemp)"

  awk -v version="$version_number" -v today="$today" '
    {
      print
      if (!inserted && $0 ~ /^## \[Unreleased\]/) {
        print ""
        print "## [" version "] - " today
        inserted=1
      }
    }
  ' "$CHANGELOG_FILE" > "$tmp"

  mv "$tmp" "$CHANGELOG_FILE"

  echo "CHANGELOG actualizado: $CHANGELOG_FILE"
  echo "Tag sugerido: $MODULE_PATH/$version"
}

print_notes() {
  local version="$1"
  local version_number="${version#v}"
  local notes
  notes="$(extract_section "$version_number")"
  [[ -n "$(printf '%s' "$notes" | tr -d '[:space:]')" ]] || fail "No se encontró la sección $version_number en $CHANGELOG_FILE"
  printf '%s\n' "$notes"
}

create_github_release() {
  local version="$1"
  local tag="$MODULE_PATH/$version"
  local tmp_notes
  tmp_notes="$(mktemp)"

  require_cmd gh
  print_notes "$version" > "$tmp_notes"

  gh release create "$tag" \
    --verify-tag \
    --title "$MODULE_PATH $version" \
    --notes-file "$tmp_notes"

  rm -f "$tmp_notes"
}

main() {
  local command="${1:-}"
  local module="${2:-}"
  local version="${3:-}"

  [[ -n "$command" && -n "$module" && -n "$version" ]] || {
    usage
    exit 1
  }

  validate_version "$version"
  resolve_module "$module"

  case "$command" in
    prepare)
      prepare_release "$version"
      ;;
    notes)
      print_notes "$version"
      ;;
    github)
      create_github_release "$version"
      ;;
    *)
      usage
      exit 1
      ;;
  esac
}

main "$@"
