package qr

import "aboutblank/qr-code/gf256"

type ErrorCorrectionLevel int

const (
	EC_L ErrorCorrectionLevel = iota
	EC_M
	EC_Q
	EC_H
)

type blockGroup struct {
	Blocks        int // number of blocks
	DataCodewords int // data codewords per block
}

type ecInfo struct {
	TotalDataCodewords  int
	ECCodewordsPerBlock int
	Group1              blockGroup
	Group2              blockGroup // Blocks == 0 if unused
}

// TODO: Stop generating the generator polynomial every time..
// Just do it once and reuse.
func GenerateErrorCorrectionCodeWords(dataCodeWords []byte, ecInfo ecInfo) []byte {
	genPoly := buildGenerator(ecInfo.ECCodewordsPerBlock)

	// Multiply message by x^n 
	messagePoly := shiftPoly(dataCodeWords, ecInfo.ECCodewordsPerBlock)

	// Division loop
	for i := 0; i <= len(messagePoly)-len(genPoly); i++ {
		lead := messagePoly[i]
		if lead == 0 {
			continue
		}
		for j := range genPoly {
			messagePoly[i+j] ^= gf256.Multiply(genPoly[j], lead)
		}
	}

	// Last n bytes = EC codewords
	ecBytes := messagePoly[len(messagePoly)-ecInfo.ECCodewordsPerBlock:]
	return ecBytes
}

func getEcInfo(version Version, ecLevel ErrorCorrectionLevel) ecInfo {
	return ecTable[version][ecLevel]
}

func (e ecInfo) TotalBlocks() int {
	return e.Group1.Blocks + e.Group2.Blocks
}

func (e ecInfo) TotalECCodewords() int {
	return e.TotalBlocks() * e.ECCodewordsPerBlock
}

func (e ecInfo) TotalCodewords() int {
	return e.TotalDataCodewords + e.TotalECCodewords()
}

func (e ecInfo) TotalDataBits() int {
	return e.TotalDataCodewords * 8
}

func (e ecInfo) TotalRequiredBits() int {
	return e.TotalCodewords() * 8
}

// ecTable[version][ecLevel]
var ecTable = [41][4]ecInfo{
	{}, // version 0 unused

	// 1
	{
		{19, 7, blockGroup{1, 19}, blockGroup{}},  // L
		{16, 10, blockGroup{1, 16}, blockGroup{}}, // M
		{13, 13, blockGroup{1, 13}, blockGroup{}}, // Q
		{9, 17, blockGroup{1, 9}, blockGroup{}},   // H
	},
	// 2
	{
		{34, 10, blockGroup{1, 34}, blockGroup{}},
		{28, 16, blockGroup{1, 28}, blockGroup{}},
		{22, 22, blockGroup{1, 22}, blockGroup{}},
		{16, 28, blockGroup{1, 16}, blockGroup{}},
	},
	// 3
	{
		{55, 15, blockGroup{1, 55}, blockGroup{}},
		{44, 26, blockGroup{1, 44}, blockGroup{}},
		{34, 18, blockGroup{2, 17}, blockGroup{}},
		{26, 22, blockGroup{2, 13}, blockGroup{}},
	},
	// 4
	{
		{80, 20, blockGroup{1, 80}, blockGroup{}},
		{64, 18, blockGroup{2, 32}, blockGroup{}},
		{48, 26, blockGroup{2, 24}, blockGroup{}},
		{36, 16, blockGroup{4, 9}, blockGroup{}},
	},
	// 5
	{
		{108, 26, blockGroup{1, 108}, blockGroup{}},
		{86, 24, blockGroup{2, 43}, blockGroup{}},
		{62, 18, blockGroup{2, 15}, blockGroup{2, 16}},
		{46, 22, blockGroup{2, 11}, blockGroup{2, 12}},
	},
	// 6
	{
		{136, 18, blockGroup{2, 68}, blockGroup{}},
		{108, 16, blockGroup{4, 27}, blockGroup{}},
		{76, 24, blockGroup{4, 19}, blockGroup{}},
		{60, 28, blockGroup{4, 15}, blockGroup{}},
	},
	// 7
	{
		{156, 20, blockGroup{2, 78}, blockGroup{}},
		{124, 18, blockGroup{4, 31}, blockGroup{}},
		{88, 18, blockGroup{2, 14}, blockGroup{4, 15}},
		{66, 26, blockGroup{4, 13}, blockGroup{1, 14}},
	},
	// 8
	{
		{194, 24, blockGroup{2, 97}, blockGroup{}},
		{154, 22, blockGroup{2, 38}, blockGroup{2, 39}},
		{110, 22, blockGroup{4, 18}, blockGroup{2, 19}},
		{86, 26, blockGroup{4, 14}, blockGroup{2, 15}},
	},
	// 9
	{
		{232, 30, blockGroup{2, 116}, blockGroup{}},
		{182, 22, blockGroup{3, 36}, blockGroup{2, 37}},
		{132, 20, blockGroup{4, 16}, blockGroup{4, 17}},
		{100, 24, blockGroup{4, 12}, blockGroup{4, 13}},
	},
	// 10
	{
		{274, 18, blockGroup{2, 68}, blockGroup{2, 69}},
		{216, 26, blockGroup{4, 43}, blockGroup{1, 44}},
		{154, 24, blockGroup{6, 19}, blockGroup{2, 20}},
		{122, 28, blockGroup{6, 15}, blockGroup{2, 16}},
	},
	// 11
	{
		{324, 20, blockGroup{4, 81}, blockGroup{}},
		{254, 30, blockGroup{1, 50}, blockGroup{4, 51}},
		{180, 28, blockGroup{4, 22}, blockGroup{4, 23}},
		{140, 24, blockGroup{3, 12}, blockGroup{8, 13}},
	},
	// 12
	{
		{370, 24, blockGroup{2, 92}, blockGroup{2, 93}},
		{290, 22, blockGroup{6, 36}, blockGroup{2, 37}},
		{206, 26, blockGroup{4, 20}, blockGroup{6, 21}},
		{158, 28, blockGroup{7, 14}, blockGroup{4, 15}},
	},
	// 13
	{
		{428, 26, blockGroup{4, 107}, blockGroup{}},
		{334, 22, blockGroup{8, 37}, blockGroup{1, 38}},
		{244, 24, blockGroup{8, 20}, blockGroup{4, 21}},
		{180, 22, blockGroup{12, 11}, blockGroup{4, 12}},
	},
	// 14
	{
		{461, 30, blockGroup{3, 115}, blockGroup{1, 116}},
		{365, 24, blockGroup{4, 40}, blockGroup{5, 41}},
		{261, 20, blockGroup{11, 16}, blockGroup{5, 17}},
		{197, 24, blockGroup{11, 12}, blockGroup{5, 13}},
	},
	// 15
	{
		{523, 22, blockGroup{5, 87}, blockGroup{1, 88}},
		{415, 24, blockGroup{5, 41}, blockGroup{5, 42}},
		{295, 30, blockGroup{5, 24}, blockGroup{7, 25}},
		{223, 24, blockGroup{11, 12}, blockGroup{7, 13}},
	},
	// 16
	{
		{589, 24, blockGroup{5, 98}, blockGroup{1, 99}},
		{453, 28, blockGroup{7, 45}, blockGroup{3, 46}},
		{325, 24, blockGroup{15, 19}, blockGroup{2, 20}},
		{253, 30, blockGroup{3, 15}, blockGroup{13, 16}},
	},
	// 17
	{
		{647, 28, blockGroup{1, 107}, blockGroup{5, 108}},
		{507, 28, blockGroup{10, 46}, blockGroup{1, 47}},
		{367, 28, blockGroup{1, 22}, blockGroup{15, 23}},
		{283, 28, blockGroup{2, 14}, blockGroup{17, 15}},
	},
	// 18
	{
		{721, 30, blockGroup{5, 120}, blockGroup{1, 121}},
		{563, 26, blockGroup{9, 43}, blockGroup{4, 44}},
		{397, 28, blockGroup{17, 22}, blockGroup{1, 23}},
		{313, 28, blockGroup{2, 14}, blockGroup{19, 15}},
	},
	// 19
	{
		{795, 28, blockGroup{3, 113}, blockGroup{4, 114}},
		{627, 26, blockGroup{3, 44}, blockGroup{11, 45}},
		{445, 26, blockGroup{17, 21}, blockGroup{4, 22}},
		{341, 26, blockGroup{9, 13}, blockGroup{16, 14}},
	},
	// 20
	{
		{861, 28, blockGroup{3, 107}, blockGroup{5, 108}},
		{669, 26, blockGroup{3, 41}, blockGroup{13, 42}},
		{485, 30, blockGroup{15, 24}, blockGroup{5, 25}},
		{385, 28, blockGroup{15, 15}, blockGroup{10, 16}},
	},
	// 21
	{
		{932, 28, blockGroup{4, 116}, blockGroup{4, 117}},
		{714, 26, blockGroup{17, 42}, blockGroup{}},
		{512, 28, blockGroup{17, 22}, blockGroup{6, 23}},
		{406, 30, blockGroup{19, 16}, blockGroup{6, 17}},
	},
	// 22
	{
		{1006, 28, blockGroup{2, 111}, blockGroup{7, 112}},
		{782, 28, blockGroup{17, 46}, blockGroup{}},
		{568, 30, blockGroup{7, 24}, blockGroup{16, 25}},
		{442, 24, blockGroup{34, 13}, blockGroup{}},
	},
	// 23
	{
		{1094, 30, blockGroup{4, 121}, blockGroup{5, 122}},
		{860, 28, blockGroup{4, 47}, blockGroup{14, 48}},
		{614, 30, blockGroup{11, 24}, blockGroup{14, 25}},
		{464, 30, blockGroup{16, 15}, blockGroup{14, 16}},
	},
	// 24
	{
		{1174, 30, blockGroup{6, 117}, blockGroup{4, 118}},
		{914, 28, blockGroup{6, 45}, blockGroup{14, 46}},
		{664, 30, blockGroup{11, 24}, blockGroup{16, 25}},
		{514, 30, blockGroup{30, 16}, blockGroup{2, 17}},
	},
	// 25
	{
		{1276, 26, blockGroup{8, 106}, blockGroup{4, 107}},
		{1000, 28, blockGroup{8, 47}, blockGroup{13, 48}},
		{718, 30, blockGroup{7, 24}, blockGroup{22, 25}},
		{538, 30, blockGroup{22, 15}, blockGroup{13, 16}},
	},
	// 26
	{
		{1370, 28, blockGroup{10, 114}, blockGroup{2, 115}},
		{1062, 28, blockGroup{19, 46}, blockGroup{4, 47}},
		{754, 28, blockGroup{28, 22}, blockGroup{6, 23}},
		{596, 30, blockGroup{33, 16}, blockGroup{4, 17}},
	},
	// 27
	{
		{1468, 30, blockGroup{8, 122}, blockGroup{4, 123}},
		{1128, 28, blockGroup{22, 45}, blockGroup{3, 46}},
		{808, 30, blockGroup{8, 23}, blockGroup{26, 24}},
		{628, 30, blockGroup{12, 15}, blockGroup{28, 16}},
	},
	// 28
	{
		{1531, 30, blockGroup{3, 117}, blockGroup{10, 118}},
		{1193, 28, blockGroup{3, 45}, blockGroup{23, 46}},
		{871, 30, blockGroup{4, 24}, blockGroup{31, 25}},
		{661, 30, blockGroup{11, 15}, blockGroup{31, 16}},
	},
	// 29
	{
		{1631, 30, blockGroup{7, 116}, blockGroup{7, 117}},
		{1267, 28, blockGroup{21, 45}, blockGroup{7, 46}},
		{911, 30, blockGroup{1, 23}, blockGroup{37, 24}},
		{701, 30, blockGroup{19, 15}, blockGroup{26, 16}},
	},

	// 30
	{
		{1735, 30, blockGroup{5, 115}, blockGroup{10, 116}},
		{1373, 28, blockGroup{19, 47}, blockGroup{10, 48}},
		{985, 30, blockGroup{15, 24}, blockGroup{25, 25}},
		{745, 30, blockGroup{23, 15}, blockGroup{25, 16}},
	},
	// 31
	{
		{1843, 30, blockGroup{13, 115}, blockGroup{3, 116}},
		{1455, 28, blockGroup{2, 46}, blockGroup{29, 47}},
		{1033, 30, blockGroup{42, 24}, blockGroup{1, 25}},
		{793, 30, blockGroup{23, 15}, blockGroup{28, 16}},
	},
	// 32
	{
		{1955, 30, blockGroup{17, 115}, blockGroup{}},
		{1541, 28, blockGroup{10, 46}, blockGroup{23, 47}},
		{1115, 30, blockGroup{10, 24}, blockGroup{35, 25}},
		{845, 30, blockGroup{19, 15}, blockGroup{35, 16}},
	},
	// 33
	{
		{2071, 30, blockGroup{17, 115}, blockGroup{1, 116}},
		{1631, 28, blockGroup{14, 46}, blockGroup{21, 47}},
		{1171, 30, blockGroup{29, 24}, blockGroup{19, 25}},
		{901, 30, blockGroup{11, 15}, blockGroup{46, 16}},
	},
	// 34
	{
		{2191, 30, blockGroup{13, 115}, blockGroup{6, 116}},
		{1725, 28, blockGroup{14, 46}, blockGroup{23, 47}},
		{1231, 30, blockGroup{44, 24}, blockGroup{7, 25}},
		{961, 30, blockGroup{59, 16}, blockGroup{1, 17}},
	},
	// 35
	{
		{2306, 30, blockGroup{12, 121}, blockGroup{7, 122}},
		{1812, 28, blockGroup{12, 47}, blockGroup{26, 48}},
		{1286, 30, blockGroup{39, 24}, blockGroup{14, 25}},
		{986, 30, blockGroup{22, 15}, blockGroup{41, 16}},
	},
	// 36
	{
		{2434, 30, blockGroup{6, 121}, blockGroup{14, 122}},
		{1914, 28, blockGroup{6, 47}, blockGroup{34, 48}},
		{1354, 30, blockGroup{46, 24}, blockGroup{10, 25}},
		{1054, 30, blockGroup{2, 15}, blockGroup{64, 16}},
	},
	// 37
	{
		{2566, 30, blockGroup{17, 122}, blockGroup{4, 123}},
		{1992, 28, blockGroup{29, 46}, blockGroup{14, 47}},
		{1426, 30, blockGroup{49, 24}, blockGroup{10, 25}},
		{1096, 30, blockGroup{24, 15}, blockGroup{46, 16}},
	},
	// 38
	{
		{2702, 30, blockGroup{4, 122}, blockGroup{18, 123}},
		{2102, 28, blockGroup{13, 46}, blockGroup{32, 47}},
		{1502, 30, blockGroup{48, 24}, blockGroup{14, 25}},
		{1142, 30, blockGroup{42, 15}, blockGroup{32, 16}},
	},
	// 39
	{
		{2812, 30, blockGroup{20, 117}, blockGroup{4, 118}},
		{2216, 28, blockGroup{40, 47}, blockGroup{7, 48}},
		{1582, 30, blockGroup{43, 24}, blockGroup{22, 25}},
		{1222, 30, blockGroup{10, 15}, blockGroup{67, 16}},
	},
	// 40
	{
		{2956, 30, blockGroup{19, 118}, blockGroup{6, 119}},
		{2334, 28, blockGroup{18, 47}, blockGroup{31, 48}},
		{1666, 30, blockGroup{34, 24}, blockGroup{34, 25}},
		{1276, 30, blockGroup{20, 15}, blockGroup{61, 16}},
	},
}
