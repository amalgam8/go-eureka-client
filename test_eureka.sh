#!/bin/bash
#
# Copyright 2016 IBM Corporation
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

set -x
set -o errexit

SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

echo "Testing eureka go client.."

sudo $SCRIPTDIR/run_eureka_server.sh
sleep 10

cuurent_dir=`pwd`
cd $SCRIPTDIR/
go test -run TestCases
cd $current_dir

sleep 10

echo "eurka go client tests successful. Cleaning up.."
#$SCRIPTDIR/cleanup.sh
#sleep 5
#sudo $SCRIPTDIR/uninstall-kubernetes.sh
