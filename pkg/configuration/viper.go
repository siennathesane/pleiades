/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package configuration

import (
	vip "github.com/spf13/viper"
)

var (
	viper *vip.Viper
)


func init() {
	viper = vip.New()
}

func Get() *vip.Viper {
	return viper
}