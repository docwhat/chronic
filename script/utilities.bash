#!/bin/bash
#
# Utility functions for pretty output

cd "$(dirname "$0")" || exit 10
cd .. || exit 10

sent_header=f
colors=$(($(tput colors 2>/dev/null || :) + 0))
export PATH="${PWD}/bin:${PATH}"

function colorize() {
  if ((colors >= 8)); then
    tput bold
    tput setaf 4
  fi
  cat -
  if ((colors >= 8)); then
    tput sgr0
  fi
}

function header() {
  if [ "$sent_header" = t ]; then
    echo
  fi
  echo "$*" | colorize
  echo '----------------------------------------------------------------' | colorize
  sent_header=t
}

function join_by() {
  local IFS="$1"
  shift
  echo "$*"
}

function bash_is_new_enough() {
  bash -c 'if ((BASH_VERSINFO[0] < 5)); then exit 2; fi'
}

function assert_bash_is_new_enough() {
  if ! bash_is_new_enough; then
    fatal "BASH version 5 or newer is required. Try running ./script/bootstrap"
  fi
}

function fatal() {
  local msg="$*"
  echo "[FATAL] $msg" 1>&2
  exit 4
}

# vim: set ft=sh :
