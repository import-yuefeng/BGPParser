#!/usr/bin/env python3
# based on shawn1m/overture (MIT LICENSE)
# The MIT License (MIT)
# Copyright (c) 2019 import-yuefeng
# Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
# The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import subprocess
import sys

GO_OS_ARCH_LIST = [
    ["darwin", "amd64"],
    ["linux", "386"],
    ["linux", "amd64"],
    ["linux", "arm"],
    ["linux", "arm64"],
    ["linux", "mips", "softfloat"],
    ["linux", "mips", "hardfloat"],
    ["linux", "mipsle", "softfloat"],
    ["linux", "mipsle", "hardfloat"],
    ["linux", "mips64"],
    ["linux", "mips64le"],
    ["freebsd", "386"],
    ["freebsd", "amd64"],
    ["windows", "386"],
    ["windows", "amd64"]
              ]


def go_build_zip():
    subprocess.check_call("GOOS=windows go get -v github.com/import-yuefeng/BGPParser", shell=True)
    for o, a, *p in GO_OS_ARCH_LIST:
        zip_name = "BGPParser-" + o + "-" + a + ("-" + (p[0] if p else "") if p else "")
        binary_name = zip_name + (".exe" if o == "windows" else "")
        version = subprocess.check_output("git describe --tags", shell=True).decode()
        mipsflag = (" GOMIPS=" + (p[0] if p else "") if p else "")
        try:
            subprocess.check_call("GOOS=" + o + " GOARCH=" + a + mipsflag + " CGO_ENABLED=0" + " go build -ldflags \"-s -w " +
                                  "-X main.version=" + version + "\" -o " + binary_name + " main.go", shell=True)
            subprocess.check_call("zip " + zip_name + ".zip " + binary_name, shell=True)
        except subprocess.CalledProcessError:
            print(o + " " + a + " " + (p[0] if p else "") + " failed.")



if __name__ == "__main__":


    if "-build" in sys.argv:
        go_build_zip()
