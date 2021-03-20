.PHONY: buildwithvars installwithvars

# the real command is the following
# go build -v -ldflags "-X main.commitHash=$(git rev-parse --short HEAD) -X main.commitDate=$(git log -1 --format=%ci | awk '{ print $1 }')"
buildwithvars:
	rm -rf ./frontend; /usr/local/go/bin/go build -v -ldflags "-X main.commitHash=$$(git rev-parse --short HEAD) -X main.commitDate=$$(git log -1 --format=%ci | awk '{ print $$1 }')"

installwithvars:
	rm -rf ./frontend; go install -v -ldflags "-X main.commitHash=$$(git rev-parse --short HEAD) -X main.commitDate=$$(git log -1 --format=%ci | awk '{ print $$1 }')"
