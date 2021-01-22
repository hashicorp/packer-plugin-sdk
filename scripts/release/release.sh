#!/bin/bash

set -e
# release.sh will:
# 1. Modify changelog
# 2. Modify version in version/version.go
# 3. Commit and push changes
# 4. Create a Git tag
# 5. Push Git tag

### Script shamelessly taken form hashicorp/terraform-plugin-sdk
function pleaseUseGNUsed {
    echo "Please install GNU sed to your PATH as 'sed'."
    exit 1
}

function gpgKeyCheck {
  if [ -z "${GPG_KEY_ID}" ]; then
    printf "A valid GPG_KEY_ID is needed to sign the release...exiting\n"
		exit 1
  fi
}

function init {
  sed --version > /dev/null || pleaseUseGNUsed

  ## Enable GPG Check
  if [ "$CI" == true ]; then
    gpgKeyCheck
  fi

  DATE=`date '+%B %d, %Y'`
  START_DIR=`pwd`

  TARGET_VERSION="$(getTargetVersion)"

  if [ -z "${TARGET_VERSION}" ]; then
   printf "Target version not found in changelog, exiting\n"
   exit 1
  fi

  TARGET_VERSION_CORE="$(getVersionCore)"
}

semverRegex='\([0-9]\+\.[0-9]\+\.[0-9]\+\)\(-\?\)\([0-9a-zA-Z.]\+\)\?'

function getTargetVersion {
  # parse target version from CHANGELOG
  sed -n 's/^## '"$semverRegex"' (Upcoming)$/\1\2\3/p' CHANGELOG.md
}

function getVersionCore {
    # extract major.minor.patch version, e.g. 1.2.3
    echo "${TARGET_VERSION}" | sed -n 's/'"$semverRegex"'/\1/p'
}

function modifyChangelog {
  sed -i "s/$TARGET_VERSION (Upcoming)$/$TARGET_VERSION ($DATE)/" CHANGELOG.md
}

function changelogLinks {
  ./scripts/release/changelog_links.sh
}

function changelogMain {
  printf "Modifying Changelog..."
  modifyChangelog
  printf "ok!\n"
}

function modifyVersionFiles {
  printf "Modifying version files..."
  sed -i "s/var Version =.*/var Version = \"${TARGET_VERSION_CORE}\"/" version/version.go
  ## Set pre-release version to empty string
  sed -i "s/var VersionPrerelease =.*/var VersionPrerelease = \"\"/" version/version.go
}

function commitChanges {
  git add CHANGELOG.md
  modifyVersionFiles
  git add version/version.go

	if [ "$CI" == true ]; then
    git commit --gpg-sign="${GPG_KEY_ID}" -m "v${TARGET_VERSION} [skip ci]"
    git tag -a -m "v${TARGET_VERSION}" -s -u "${GPG_KEY_ID}" "v${TARGET_VERSION}"
    git push origin "${CIRCLE_BRANCH}"
  else
    printf "Skipping GPG signature on non CI releases...\n"
    git commit -m "v${TARGET_VERSION} [skip ci]"
    git tag -a -m "v${TARGET_VERSION}" "v${TARGET_VERSION}"
    git push origin main
  fi

  git push origin "v${TARGET_VERSION}"
}

function commitMain {
  printf "Committing Changes..."
  commitChanges
  printf "ok!\n"
}

function main {
  init
  changelogMain
  commitMain
}

main
