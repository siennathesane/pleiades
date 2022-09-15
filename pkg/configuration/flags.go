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
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

func appendFlags[T any](name string, fs *pflag.FlagSet) {
	tHost := reflect.TypeOf((*T)(nil)).Elem()

	if tHost.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%s must be a struct", tHost.Name()))
	}

	for i := 0; i < tHost.NumField(); i++ {
		field := tHost.Field(i)

		// get the flag
		flag, ok := field.Tag.Lookup("flag")
		if !ok {
			panic(fmt.Sprintf("the `flag:\"%s\"` struct tag must be present", strings.ToLower(field.Name)))
		}

		if name != "" {
			flag = fmt.Sprintf("%s-%s", name, flag)
		}

		// the default value
		def, ok := field.Tag.Lookup("default")
		if !ok {
			panic(fmt.Sprintf("the `default:\"%s\"` struct tag must be present", strings.ToLower(field.Name)))
		}

		// the default value
		usage, ok := field.Tag.Lookup("usage")
		if !ok {
			panic(fmt.Sprintf("the `usage:\"%s\"` struct tag must be present", strings.ToLower(field.Name)))
		}

		switch field.Type.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(def)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to bool", def))
			}
			fs.Bool(flag, b, usage)
			break
		case reflect.Int:
			i, err := strconv.Atoi(def)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int", def))
			}
			fs.Int(flag, i, usage)
			break
		case reflect.Int8:
			i, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int8", def))
			}
			fs.Int8(flag, int8(i), usage)
			break
		case reflect.Int16:
			i, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int16", def))
			}
			fs.Int16(flag, int16(i), usage)
			break
		case reflect.Int32:
			i, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int32", def))
			}
			fs.Int32(flag, int32(i), usage)
			break
		case reflect.Int64:
			i, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int64", def))
			}
			fs.Int64(flag, int64(i), usage)
			break
		case reflect.Uint:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint", def))
			}
			fs.Uint(flag, uint(u), usage)
			break
		case reflect.Uint8:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint8", def))
			}
			fs.Uint8(flag, uint8(u), usage)
			break
		case reflect.Uint16:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint16", def))
			}
			fs.Uint16(flag, uint16(u), usage)
			break
		case reflect.Uint32:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint32", def))
			}
			fs.Uint32(flag, uint32(u), usage)
			break
		case reflect.Uint64:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint64", def))
			}
			fs.Uint64(flag, uint64(u), usage)
			break
		case reflect.Float32:
			f, err := strconv.ParseFloat(def, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to float32", def))
			}
			fs.Float32(flag, float32(f), usage)
			break
		case reflect.Float64:
			f, err := strconv.ParseFloat(def, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to float64", def))
			}
			fs.Float64(flag, float64(f), usage)
			break
		case reflect.String:
			fs.String(flag, def, usage)
			break
		default:
			panic(fmt.Sprintf("%s is an unsupported type", field.Type.Kind()))
		}
	}
}

func ToFlagSet[T any](name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ExitOnError)

	tHost := reflect.TypeOf((*T)(nil)).Elem()

	if tHost.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%s must be a struct", tHost.Name()))
	}

	for i := 0; i < tHost.NumField(); i++ {
		field := tHost.Field(i)

		// get the flag
		flag, ok := field.Tag.Lookup("flag")
		if !ok {
			panic(fmt.Sprintf("the `flag:\"%s\"` struct tag must be present", strings.ToLower(field.Name)))
		}

		if name != "" {
			flag = fmt.Sprintf("%s-%s", name, flag)
		}

		// the default value
		def, ok := field.Tag.Lookup("default")
		if !ok {
			panic(fmt.Sprintf("the `default:\"%s\"` struct tag must be present", strings.ToLower(field.Name)))
		}

		// the default value
		usage, ok := field.Tag.Lookup("usage")
		if !ok {
			panic(fmt.Sprintf("the `usage:\"%s\"` struct tag must be present", strings.ToLower(field.Name)))
		}

		switch field.Type.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(def)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to bool", def))
			}
			fs.Bool(flag, b, usage)
			break
		case reflect.Int:
			i, err := strconv.Atoi(def)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int", def))
			}
			fs.Int(flag, i, usage)
			break
		case reflect.Int8:
			i, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int8", def))
			}
			fs.Int8(flag, int8(i), usage)
			break
		case reflect.Int16:
			i, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int16", def))
			}
			fs.Int16(flag, int16(i), usage)
			break
		case reflect.Int32:
			i, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int32", def))
			}
			fs.Int32(flag, int32(i), usage)
			break
		case reflect.Int64:
			i, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to int64", def))
			}
			fs.Int64(flag, int64(i), usage)
			break
		case reflect.Uint:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint", def))
			}
			fs.Uint(flag, uint(u), usage)
			break
		case reflect.Uint8:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint8", def))
			}
			fs.Uint8(flag, uint8(u), usage)
			break
		case reflect.Uint16:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint16", def))
			}
			fs.Uint16(flag, uint16(u), usage)
			break
		case reflect.Uint32:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint32", def))
			}
			fs.Uint32(flag, uint32(u), usage)
			break
		case reflect.Uint64:
			u, err := strconv.ParseUint(def, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to uint64", def))
			}
			fs.Uint64(flag, uint64(u), usage)
			break
		case reflect.Float32:
			f, err := strconv.ParseFloat(def, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to float32", def))
			}
			fs.Float32(flag, float32(f), usage)
			break
		case reflect.Float64:
			f, err := strconv.ParseFloat(def, 64)
			if err != nil {
				panic(fmt.Sprintf("can't parse %s to float64", def))
			}
			fs.Float64(flag, float64(f), usage)
			break
		case reflect.String:
			fs.String(flag, def, usage)
			break
		default:
			panic(fmt.Sprintf("%s is an unsupported type", field.Type.Kind()))
		}
	}

	return fs
}
