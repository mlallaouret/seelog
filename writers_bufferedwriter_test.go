// Copyright (c) 2012 - Cloud Instruments Co., Ltd.
// 
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met: 
// 
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer. 
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution. 
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package seelog

import (
	"testing"
)

func TestChunkWriteOnFilling(t *testing.T) {
	writer, _ := newBytesVerifier(t)
	bufferedWriter, err := newBufferedWriter(writer, 1024, 0)

	if err != nil {
		t.Fatalf("Unexpected buffered writer creation error: %s", err.Error())
	}

	bytes := make([]byte, 1000)

	bufferedWriter.Write(bytes)
	writer.ExpectBytes(bytes)
	bufferedWriter.Write(bytes)

	// bufferedWriter writes another chunk not at once but in goroutine (with nondetermined delay)
	// so we wait some time
	writer.MustNotExpectWithDelay(0.1 * 1e8)
}

func TestFlushByTimePeriod(t *testing.T) {
	writer, _ := newBytesVerifier(t)
	bufferedWriter, err := newBufferedWriter(writer, 1024, 10)

	if err != nil {
		t.Fatalf("Unexpected buffered writer creation error: %s", err.Error())
	}

	bytes := []byte("Hello")

	for i := 0; i < 1; i++ {
		writer.ExpectBytes(bytes)
		bufferedWriter.Write(bytes)
		writer.MustNotExpectWithDelay(0.2 * 1e8)
	}
}

func TestBigMessageMustPassMemoryBuffer(t *testing.T) {
	writer, _ := newBytesVerifier(t)
	bufferedWriter, err := newBufferedWriter(writer, 1024, 0)

	if err != nil {
		t.Fatalf("Unexpected buffered writer creation error: %s", err.Error())
	}

	bytes := make([]byte, 5000)

	for i := 0; i < len(bytes); i++ {
		bytes[i] = uint8(i % 255)
	}

	writer.ExpectBytes(bytes)
	bufferedWriter.Write(bytes)
	writer.MustNotExpectWithDelay(0.1 * 1e8)
}
