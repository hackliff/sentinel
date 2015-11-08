#! /bin/sh

export GO15VENDOREXPERIMENT=1

install_glide() {
  local dest=$GOPATH/src/github.com/Masterminds/glide
  mkdir -p $GOPATH/src/github.com/Masterminds
  git clone https://github.com/Masterminds/glide $GOPATH/src/github.com/Masterminds/glide
  cd $GOPATH/src/github.com/Masterminds/glide
  make bootstrap
  make install
}

deprecated() {
  local version=${1:-"0.7.0"}
  local platform=${2:-"linux-amd64"}
  local target=glide-${GLIDE_VERSION}-${platform}
  local repo="https://github.com/Masterminds/glide"

  curl -LkOs "${repo}/releases/download/${version}/${target}.tar.gz" && \
    tar xvzf ${target}.tar.gz && \
    mv ${platform}/glide $GOPATH/bin/ && \
    rm -r *${platform}*
}

install_tools() {
  go get github.com/alecthomas/gometalinter && \
    gometalinter --install --update
  go get github.com/mitchellh/gox
  go get github.com/tcnksm/ghr
}

install_deps() {
  glide install
}

# TODO detect platform
install_glide 0.7.0 linux-amd64
install_tools
install_deps
