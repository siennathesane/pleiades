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
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
)

func TestSchema(t *testing.T) {
	suite.Run(t, new(schemaTestSuite))
}

type schemaTestSuite struct {
	suite.Suite
}

func (t *schemaTestSuite) TestGetFlagSet() {
	dsFs := ToFlagSet[Datastore]("datastore")
	hostFs := ToFlagSet[HostConfig]("host")

	fs := pflag.NewFlagSet("root", pflag.ExitOnError)
	fs.AddFlagSet(dsFs)
	fs.AddFlagSet(hostFs)

	t.Require().NotNil(fs)
	fs.PrintDefaults()
}
