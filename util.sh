#!/bin/sh

set -o errexit
set -o pipefail
set -o nounset

case $1 in

	help)
		cat ./util.sh
		;;

	format)
		git ls-files | grep '\.go' | xargs -I {} gofmt -w=true '{}'
		;;

	gitignore)
		cat > .gitignore << EndOfGitIgnore
/.gitignore
.DS_Store
._*
.*.swp
.*.swo
EndOfGitIgnore
		;;

	build)
		go get github.com/setekhid/jormungand/cmd/jormsrv
		go get github.com/setekhid/jormungand/cmd/jormcli
		;;
esac
