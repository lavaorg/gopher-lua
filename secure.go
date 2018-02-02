/*
 * Copyright(c) 2018 Larry Rau, all rights reserved.
 *
 * LICENSE: see end of file.
 */

package lua

var secureBaseFuncs = map[string]LGFunction{
	"error":    baseError,
	"tonumber": baseToNumber,
	"tostring": baseToString,
	"type":     baseType,
	"next":     baseNext,
	"unpack":   baseUnpack,
	// loadlib
	"module":  loModule,
	"require": loRequire,
}

var secureOsFuncs = map[string]LGFunction{
	"clock":    osClock,
	"difftime": osDiffTime,
	"date":     osDate,
	"time":     osTime,
}

func OpenSecureBase(L *LState) int {
	global := L.Get(GlobalsIndex).(*LTable)
	L.SetGlobal("_G", global)
	basemod := L.RegisterModule("_G", secureBaseFuncs)
	global.RawSetString("ipairs", L.NewClosure(baseIpairs, L.NewFunction(ipairsaux)))
	global.RawSetString("pairs", L.NewClosure(basePairs, L.NewFunction(pairsaux)))
	L.Push(basemod)
	return 1
}

func OpenSecureOs(L *LState) int {
	osmod := L.RegisterModule(OsLibName, secureOsFuncs)
	L.Push(osmod)
	return 1
}

/*
 LICENSE

 The MIT License (MIT)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
