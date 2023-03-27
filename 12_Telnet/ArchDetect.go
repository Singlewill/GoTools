package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	EI_NIDENT = 16 // Side of e_ident in elf header
	EI_DATA   = 5  // Offset endianness in e_ident
)

//ELF 大小端定义
const (
	EE_NONE   = 0 // No endianness ????
	EE_LITTLE = 1 // Little endian
	EE_BIG    = 2 // Big endian
)

// ELF target machine定义
const (
	EM_NONE  = 0
	EM_M32   = 1
	EM_SPARC = 2
	EM_386   = 3
	EM_68K   = 4 // m68k
	EM_88K   = 5 // m68k
	EM_486   = 6 // x86
	EM_860   = 7 // Unknown
	EM_MIPS  = 8 /* MIPS R3000 (officially, big-endian only) */
	/* Next two are historical and binaries and
	modules of these types will be rejected by
	Linux.  */
	EM_MIPS_RS3_LE = 10 /* MIPS R3000 little-endian */
	EM_MIPS_RS4_BE = 10 /* MIPS R4000 big-endian */

	EM_PARISC      = 15 /* HPPA */
	EM_SPARC32PLUS = 18 /* Sun's "v8plus" */

	EM_PPC          = 20     /* PowerPC */
	EM_PPC64        = 21     /* PowerPC64 */
	EM_SPU          = 23     /* Cell BE SPU */
	EM_ARM          = 40     /* ARM 32 bit */
	EM_SH           = 42     /* SuperH */
	EM_SPARCV9      = 43     /* SPARC v9 64-bit */
	EM_H8_300       = 46     /* Renesas H8/300 */
	EM_IA_64        = 50     /* HP/Intel IA-64 */
	EM_X86_64       = 62     /* AMD x86-64 */
	EM_S390         = 22     /* IBM S/390 */
	EM_CRIS         = 76     /* Axis Communications 32-bit embedded processor */
	EM_M32R         = 88     /* Renesas M32R */
	EM_MN10300      = 89     /* Panasonic/MEI MN10300, AM33 */
	EM_OPENRISC     = 92     /* OpenRISC 32-bit embedded processor */
	EM_BLACKFIN     = 106    /* ADI Blackfin Processor */
	EM_ALTERA_NIOS2 = 113    /* Altera Nios II soft-core processor */
	EM_TI_C6000     = 140    /* TI C6X DSPs */
	EM_AARCH64      = 183    /* ARM 64 bit */
	EM_TILEPRO      = 188    /* Tilera TILEPro */
	EM_MICROBLAZE   = 189    /* Xilinx MicroBlaze */
	EM_TILEGX       = 191    /* Tilera TILE-Gx */
	EM_FRV          = 0x5441 /* Fujitsu FR-V */
	EM_AVR32        = 0x18ad /* Atmel AVR32 */
)

type ElfHdr struct {
	eIdent   [EI_NIDENT]byte
	eType    uint16
	eMachine uint16
	eVersion uint32
}

const ElfHdrSize = 24

func (a *ElfHdr) unmarshal(b []byte, order binary.ByteOrder) error {
	if len(b) < ElfHdrSize {
		return fmt.Errorf("invalid attribute; length too short or too large")
	}
	copy(a.eIdent[:], b[:EI_NIDENT])
	a.eType = order.Uint16(b[EI_NIDENT : EI_NIDENT+2])
	a.eMachine = order.Uint16(b[EI_NIDENT+2 : EI_NIDENT+4])
	a.eVersion = order.Uint32(b[EI_NIDENT+4 : EI_NIDENT+8])
	return nil
}

func ArchDetect(info []byte) (string, error) {

	start := strings.Index(string(info), "ELF")
	//ELF文件头为 "7f 45 4c 46 01 01 01 00 00 00 00 00 00 00 00 00"
	//"ELF在文件头的位置为1
	if start < 1 {
		return "", fmt.Errorf("ELF not found")
	}
	elf := &ElfHdr{}
	if err := elf.unmarshal(info[start-1:], binary.LittleEndian); err != nil {
		return "", err
	}

	//var mystruct *Info = *(**Info)(unsafe.Pointer(&data))
	if elf.eMachine == EM_ARM || elf.eMachine == EM_AARCH64 {
		return "arm", nil
	} else if elf.eMachine == EM_MIPS || elf.eMachine == EM_MIPS_RS3_LE {
		if elf.eIdent[EI_DATA] == EE_LITTLE {
			return "mpsl", nil
		} else {
			return "mps", nil
		}
	} else if elf.eMachine == EM_386 || elf.eMachine == EM_486 || elf.eMachine == EM_860 || elf.eMachine == EM_X86_64 {
		return "x86", nil
	} else if elf.eMachine == EM_SPARC || elf.eMachine == EM_SPARC32PLUS || elf.eMachine == EM_SPARCV9 {
		return "spc", nil
	} else if elf.eMachine == EM_68K || elf.eMachine == EM_88K {
		return "m68k", nil
	} else if elf.eMachine == EM_PPC || elf.eMachine == EM_PPC64 {
		return "ppc", nil
	} else if elf.eMachine == EM_SH {
		return "sh4", nil
	}
	return "", fmt.Errorf("Arch misMatch")

}
