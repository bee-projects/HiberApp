#!/usr/bin/env bash
# This file:
#
#  - Deploy example Hiber Apps
#
# Usage:
#
#  LOG_LEVEL=7 ./example.sh -f /tmp/x -d (change this for your script)
#
# Based on a template by BASH3 Boilerplate v2.3.0
# http://bash3boilerplate.sh/#authors
#
# The MIT License (MIT)
# Copyright (c) 2013 Kevin van Zonneveld and contributors
# You are not obligated to bundle the LICENSE file with your b3bp projects as long
# as you leave these references intact in the header comments of your source files.


### BASH3 Boilerplate (b3bp) Header
##############################################################################

# Commandline options. This defines the usage page, and is used to parse cli
# opts & defaults from. The parsing is unforgiving so be precise in your syntax
# - A short option must be preset for every long option; but every short option
#   need not have a long option
# - `--` is respected as the separator between options and arguments
# - We do not bash-expand defaults, so setting '~/app' as a default will not resolve to ${HOME}.
#   you can use bash variables to work around this (so use ${HOME} instead)

# shellcheck disable=SC2034

read -r -d '' __usage <<-'EOF' || true # exits non-zero when EOF encountered
  -a --application [arg] Application to deploy. Default="wordpress"
  -A --action      [arg] Action to perform. Default="create"
  -l --location    [arg] Location to create the application. Default="eastus"
  -g --resource-group [arg] Name of the resource group. Default="aci-example"
  -n --name        [arg] Name of the application.
  -v               Enable verbose mode, print script as it is executed
  -d --debug       Enables debug mode
  -h --help        This page
  -n --no-color    Disable color output
  -1 --one         Do just one thing
EOF

# shellcheck disable=SC2034
read -r -d '' __helptext <<-'EOF' || true # exits non-zero when EOF encountered
This script is to download all the tables.
EOF

# shellcheck source=main.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/base.sh"


### Signal trapping and backtracing
##############################################################################

function __b3bp_cleanup_before_exit () {
  info "Cleaning up. Done"
}
trap __b3bp_cleanup_before_exit EXIT

# requires `set -o errtrace`
__b3bp_err_report() {
    local error_code
    error_code=${?}
    # shellcheck disable=SC2154
    error "Error in ${__file} in function ${1} on line ${2}"
    exit ${error_code}
}
# Uncomment the following line for always providing an error backtrace
# trap '__b3bp_err_report "${FUNCNAME:-.}" ${LINENO}' ERR


### Command-line argument switches (like -d for debugmode, -h for showing helppage)
##############################################################################

# debug mode
if [[ "${arg_d:?}" = "1" ]]; then
  set -o xtrace
  LOG_LEVEL="7"
  # Enable error backtracing
  trap '__b3bp_err_report "${FUNCNAME:-.}" ${LINENO}' ERR
fi

# verbose mode
if [[ "${arg_v:?}" = "1" ]]; then
  set -o verbose
fi

# no color mode
if [[ "${arg_n:?}" = "1" ]]; then
  NO_COLOR="true"
fi

# help mode
if [[ "${arg_h:?}" = "1" ]]; then
  # Help exists with code 1
  help "Help using ${0}"
fi


### Validation. Error out if the things required for your script are not present
##############################################################################

[[ "${LOG_LEVEL:-}" ]] || emergency "Cannot continue without LOG_LEVEL. "

if [ "${arg_A}" != "create" ] && [ "${arg_A}" != "delete" ];
then
    emergency "Action has to be 'create' or 'delete'"
fi

if [[ -z "${arg_n}" ]]; then
    __appname=${arg_a}
else
    __appname=${arg_n}
fi

# Init vars
templatefile="example/${arg_a}/azuredeploy.json"
parameterfile="example/${arg_a}/azuredeploy.parameters.json"

function create ()
{
  echo "Running create action:"
  az group create --name "$arg_g" --location "$arg_l"
  az group deployment create \
    --name "${__appname}" \
    --resource-group "$arg_g" \
    --template-file "$templatefile" \
    --parameters "$parameterfile"
}

function delete()
{
  echo "Running delete action:"
  az group deployment delete \
    --name "$__appname" \
    --resource-group "$arg_g"
}

function checkargs()
{
  ACTION="$1"
  if [ -z "$ACTION" ];
  then
    echo "Must supply action: 'create' or 'delete'"
    exit 1
  fi

  if [ "$ACTION" != "create" ] && [ "$ACTION" != "delete" ];
  then
    echo "Action has to be 'create' or 'delete'"
    exit 1
  fi
}

if [ "${arg_A}" == "create" ];
then
  create
elif [ "${arg_A}" == "delete" ]
then
  delete
fi
