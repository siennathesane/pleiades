#!/bin/sh

#
# Copyright (c) 2022 Sienna Lloyd
#
# Licensed under the PolyForm Strict License 1.0.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License here:
#  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
#

wget -q -O mage-linux.tar.gz $(curl -s https://api.github.com/repos/magefile/mage/releases/latest | jq -r '.assets[] | select(.name | contains("Linux-64bit")) | .browser_download_url')

tar zxvf mage-linux.tar.gz

mv mage /usr/local/bin/mage

mage -version