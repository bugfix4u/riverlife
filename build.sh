#!/bin/bash

export RLCODE=$PWD

cd ${RLCODE}/scripts

case "$1" in
  compile)
    sudo -E ./rl_compile ${@:2}
    ;;
  build)
    sudo -E ./rl_build ${@:2}
    ;;
  package)
    sudo -E ./rl_package
    ;;
  install)
    sudo -E ./rl_install
    ;;
  run)
    sudo -E ./rl_run ${@:2}
    ;;
  *)
    echo "Usage: $0 {compile|build|run|install|package}"
esac