/*
 * Copyright(c) 2018 Larry Rau, all rights reserved.
 *
 * LICENSE: see end of file.
 */
package lua

import (
	"bufio"
	"github.com/vmihailenco/msgpack"
	"io/ioutil"
)

/*
 * if a FunctionProto load it.
 * return nil if not detected; no state change
 */
func (ls *LState) loadCheckForFunctionProto(b *bufio.Reader) *LFunction {
	if sbuf, err := b.Peek(4); err == nil {
		if string(sbuf) == dumpSignature {
			b.Discard(4)
			buf, err := ioutil.ReadAll(b)
			if err != nil {
				ls.RaiseError(err.Error())
			}
			var container FunctionProtoContainer
			if err := msgpack.Unmarshal(buf, &container); err != nil {
				ls.RaiseError(err.Error())
			}
			return newLFunctionL(container2proto(&container), ls.currentEnv(), 0)
		}
	}

	return nil
}

type lvalueContainer struct {
	Type   LValueType
	String string
	Number float64
}

type FunctionProtoContainer struct {
	SourceName         string
	LineDefined        int
	LastLineDefined    int
	NumUpvalues        uint8
	NumParameters      uint8
	IsVarArg           uint8
	NumUsedRegisters   uint8
	Code               []uint32
	Constants          []interface{}
	FunctionPrototypes []*FunctionProtoContainer
	DbgSourcePositions []int
	DbgLocals          []*DbgLocalInfo
	DbgCalls           []DbgCall
	DbgUpvalues        []string
	StringConstants    []string
}

const dumpSignature = "\033GoL"

func Proto2Container(proto *FunctionProto) *FunctionProtoContainer {
	return proto2container(proto)
}

func proto2container(proto *FunctionProto) *FunctionProtoContainer {
	protos := []*FunctionProtoContainer{}
	for _, c := range proto.FunctionPrototypes {
		protos = append(protos, proto2container(c))
	}

	constants := []interface{}{}
	for _, c := range proto.Constants {
		constants = append(constants, c)
	}

	stringConstants := []string{}
	for _, s := range proto.stringConstants {
		stringConstants = append(stringConstants, s)
	}

	return &FunctionProtoContainer{
		SourceName:         proto.SourceName,
		LineDefined:        proto.LineDefined,
		LastLineDefined:    proto.LastLineDefined,
		NumUpvalues:        proto.NumUpvalues,
		NumParameters:      proto.NumParameters,
		IsVarArg:           proto.IsVarArg,
		NumUsedRegisters:   proto.NumUsedRegisters,
		Code:               proto.Code,
		Constants:          constants,
		FunctionPrototypes: protos,
		DbgSourcePositions: proto.DbgSourcePositions,
		DbgLocals:          proto.DbgLocals,
		DbgCalls:           proto.DbgCalls,
		DbgUpvalues:        proto.DbgUpvalues,
		StringConstants:    stringConstants,
	}
}

func Container2Proto(container *FunctionProtoContainer) *FunctionProto {
	return container2proto(container)
}

func container2proto(container *FunctionProtoContainer) *FunctionProto {
	protos := []*FunctionProto{}
	for _, c := range container.FunctionPrototypes {
		protos = append(protos, container2proto(c))
	}

	constants := []LValue{}
	for _, c := range container.Constants {
		if s, ok := c.(string); ok {
			constants = append(constants, LString(s))
		} else {
			constants = append(constants, LNumber(c.(float64)))
		}
	}

	stringConstants := []string{}
	for _, s := range container.StringConstants {
		stringConstants = append(stringConstants, s)
	}

	return &FunctionProto{
		SourceName:         container.SourceName,
		LineDefined:        container.LineDefined,
		LastLineDefined:    container.LastLineDefined,
		NumUpvalues:        container.NumUpvalues,
		NumParameters:      container.NumParameters,
		IsVarArg:           container.IsVarArg,
		NumUsedRegisters:   container.NumUsedRegisters,
		Code:               container.Code,
		Constants:          constants,
		FunctionPrototypes: protos,
		DbgSourcePositions: container.DbgSourcePositions,
		DbgLocals:          container.DbgLocals,
		DbgCalls:           container.DbgCalls,
		DbgUpvalues:        container.DbgUpvalues,
		stringConstants:    stringConstants,
	}

}

func (ls *LState) IsStopped() int32 {
	return ls.stop
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
