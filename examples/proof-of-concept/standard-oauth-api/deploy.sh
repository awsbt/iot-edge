#!/usr/bin/env bash
set -e

#
# Copyright 2022 ForgeRock AS
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

INITIAL_DIR=$(PWD)
FORGEOPS_DIR=$(PWD)/../../../deployments/forgeops
CUSTOM_OVERLAY_DIR=$(PWD)/forgeops/overlay

cd $FORGEOPS_DIR
./deploy.sh $CUSTOM_OVERLAY_DIR 6KZjOxJU1xHGWHI0hrQT24Fn

cd $INITIAL_DIR
./deploy-ig.sh $FORGEOPS_DIR/tmp/forgeops/bin