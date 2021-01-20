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
  TARGET_VERSION_CORE="$(getVersionCore)"
}

semverRegex='\([0-9]\+\.[0-9]\+\.[0-9]\+\)\(-\?\)\([0-9a-zA-Z.]\+\)\?'

function getTargetVersion {
  # parse target version from CHANGELOG
  _version=$(sed -n 's/^## '"$semverRegex"' (Upcoming)$/\1\2\3/p' CHANGELOG.md)
  if [ -z $_version ]; then
   echo "Target version not found in changelog, exiting"
   exit 1
  fi

  echo $_version

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
  sed -i "s/const Version =.*/const Version = \"${TARGET_VERSION_CORE}\"/" version/version.go
  ## Set pre-release version to empty string
  sed -i "s/const VersionPrerelease =.*/const VersionPrerelease = \"\"/" version/version.go
}

function commitChanges {
  git add CHANGELOG.md
  modifyVersionFiles
  git add version/version.go

	## Enable GPG Signing on commits
	if [ "$CI" == true ]; then
    git commit --gpg-sign="${GPG_KEY_ID}" -m "v${TARGET_VERSION} [skip ci]"
    git tag -a -m "v${TARGET_VERSION}" -s -u "${GPG_KEY_ID}" "v${TARGET_VERSION}"
    #git push origin "${CIRCLE_BRANCH}"
  else
    printf "Skipping GPG signature on non CI releases...\n"
    git commit -m "v${TARGET_VERSION} [skip ci]"
    git tag -a -m "v${TARGET_VERSION}" "v${TARGET_VERSION}"
    git push origin
  fi

  #git push origin "v${TARGET_VERSION}"
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
