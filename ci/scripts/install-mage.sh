#!/bin/sh

#
# Copyright (c) 2022 Sienna Lloyd
#
# Licensed under the PolyForm Strict License 1.0.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License here:
#  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
#

wget https://github.com/magefile/mage/releases/download/v${MAGE_VERSION}/mage_${MAGE_VERSION}_Linux-64bit.tar.gz

tar zxvf mage_${MAGE_VERSION}_Linux-64bit.tar.gz

mv mage /usr/local/bin/mage
