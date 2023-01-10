#!/bin/sh
# Inspired by https://github.com/nodejs/node/blob/21fb98e2bf50648a31a1ef0fa23dd0dbe4002e69/tools/macos-firewall.sh
# Script that adds rules to Mac OS X Socket Firewall to avoid
# popups asking to accept incoming network connections when
# running tests.
SFW="/usr/libexec/ApplicationFirewall/socketfilterfw"
SCRIPTSDIR="$(dirname "$0")"
SCRIPTSDIR="$(cd "$SCRIPTSDIR" && pwd)"
ROOTDIR="$(cd "$SCRIPTSDIR/.." && pwd)"
BUILDDIR="$SCRIPTSDIR/../build"
# Using cd and pwd here so that the path used for socketfilterfw does not
# contain a '..', which seems to cause the rules to be incorrectly added
# and they are not removed when this script is re-run. Instead the new
# rules are simply appended. By using pwd we can get the full path
# without '..' and things work as expected.
BUILDDIR="$(cd "$BUILDDIR" && pwd)"
TESTONE="$BUILDDIR/test.test"
TESTTWO="$ROOTDIR/app/test/test.test"

add_and_unblock () {
  if [ -e "$1" ]
  then
    echo Processing "$1"
    $SFW --remove "$1" >/dev/null
    $SFW --add "$1"
    $SFW --unblock "$1"
  fi
}

if [ -f $SFW ];
then
#   add_and_unblock "$TESTONE"
  add_and_unblock "$TESTTWO"
else
  echo "SocketFirewall not found in location: $SFW"
fi
