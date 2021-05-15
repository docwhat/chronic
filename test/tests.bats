#!/bin/bash
# shellcheck disable=SC2154,SC2120

load test_helper/bats-assert/load.bash
load test_helper/bats-support/load.bash

setup_file() {
  export random_version="v${RANDOM}"
  go build -a -o "${BATS_TMPDIR:?}/chronic" -ldflags="-X 'main.version=${random_version:?}'" ./
  go build
  if ((BASH_VERSINFO[0] < 4)); then
    fail "*** bash version 4+ is required ***"
  fi
}

function run_chronic {
  run "${BATS_TMPDIR:?}/chronic" "$@"
}

## Help Tests
#
function it_shows_help_successfully { #@test
  run_chronic
  assert_success
}

function help_shows_usage { #@test
  run_chronic
  assert_line --index 0 "Usage: chronic <command> [args]..."
}

function help_shows_the_version { #@test
  run_chronic
  assert_line --index 1 "Version: ${random_version:?}"
}

function help_shows_a_description { #@test
  run_chronic
  # Bug in bats; empty lines are stripped.
  assert_line --index 2 --regexp 'Chronic runs'
}

## Success Tests
#
function it_discards_stdout_on_success { #@test
  run_chronic bash -c 'echo stdout; exit 0'
  assert_equal "$output" ""
}

function it_discards_stderr_on_success { #@test
  run_chronic bash -c 'echo stderr 1>&2; exit 0'
  assert_equal "$output" ""
}

function it_returns_success_on_success { #@test
  run_chronic bash -c 'exit 0'
  assert_success
}

## Failure Tests
#
function it_returns_the_correct_failure_on_failure { #@test
  local -ri ec=$((random % 64 + 2))
  run_chronic bash -c "exit $ec"
  assert_failure "$ec"
}

function it_shows_the_command_on_failure { #@test
  local -r str="q${RANDOM}q"
  run_chronic bash -c "echo '${str}'; exit 9"
  assert_line "[\`bash\` \`-c\` \`echo '${str}'; exit 9\`]"
}

function it_displays_a_stdout_header_on_failure { #@test
  run_chronic bash -c 'echo stdout; exit 9'
  assert_line "**** stdout ****"
}

function it_displays_standard_out_on_failure { #@test
  run_chronic bash -c 'echo -e "to be or not to be\nthat is the question"; exit 9'
  assert_line "stdout: to be or not to be"
  assert_line "stdout: that is the question"
}

function it_displays_a_stderr_header_on_failure { #@test
  run_chronic bash -c 'echo stderr 1>&2; exit 9'
  assert_line "**** stderr ****"
}

function it_displays_standard_error_on_failure { #@test
  run_chronic bash -c 'echo -e "In Xanadu did Kubla Khan\nA stately pleasure-dome decree" 1>&2; exit 9'
  assert_line "stderr: In Xanadu did Kubla Khan"
  assert_line "stderr: A stately pleasure-dome decree"
}

# EOF
