/**
* Copyright 2016 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package gaad

// num_swb_long_window param from tables 4.129 through 4.140
// 960 uses the same tables as 1024, but can't use the last
// value for some of the tables.  This table specifies the
// maximum index into the swb_offset tables each bitrate
// is allowed to use
var num_swb_long_windows = [][]uint8{
	{40, 40, 46, 49, 49, 49, 46, 46, 42, 42, 42, 40, 40}, // 960
	{41, 41, 47, 49, 49, 51, 47, 47, 43, 43, 43, 40, 40}, // 1024
}

// num_swb_short_window param from tables 4.130 through 1.141
var num_swb_short_window = []uint8{
	12, 12, 12, 14, 14, 14, 15, 15, 15, 15, 15, 15,
}

// Table 4.140
var swb_offset_1024_96 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56,
	64, 72, 80, 88, 96, 108, 120, 132, 144, 156, 172, 188, 212, 240,
	276, 320, 384, 448, 512, 576, 640, 704, 768, 832, 896, 960, 1024,
}

// Table 4.141
var swb_offset_128_96 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 32, 40, 48, 64, 92, 128,
}

// Table 4.138
var swb_offset_1024_64 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56,
	64, 72, 80, 88, 100, 112, 124, 140, 156, 172, 192, 216, 240, 268,
	304, 344, 384, 424, 464, 504, 544, 584, 624, 664, 704, 744, 784, 824,
	864, 904, 944, 984, 1024,
}

// Table 4.139
var swb_offset_128_64 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 32, 40, 48, 64, 92, 128,
}

// Table 4.129
var swb_offset_1024_48 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 48, 56, 64, 72,
	80, 88, 96, 108, 120, 132, 144, 160, 176, 196, 216, 240, 264, 292,
	320, 352, 384, 416, 448, 480, 512, 544, 576, 608, 640, 672, 704, 736,
	768, 800, 832, 864, 896, 928, 1024,
}

// Table 4.130
var swb_offset_128_48 = []uint16{
	0, 4, 8, 12, 16, 20, 28, 36, 44, 56, 68, 80, 96, 112, 128,
}

// Table 4.131
var swb_offset_1024_32 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 48, 56, 64, 72,
	80, 88, 96, 108, 120, 132, 144, 160, 176, 196, 216, 240, 264, 292,
	320, 352, 384, 416, 448, 480, 512, 544, 576, 608, 640, 672, 704, 736,
	768, 800, 832, 864, 896, 928, 960, 992, 1024,
}

// Table 4.136
var swb_offset_1024_24 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 52, 60, 68,
	76, 84, 92, 100, 108, 116, 124, 136, 148, 160, 172, 188, 204, 220,
	240, 260, 284, 308, 336, 364, 396, 432, 468, 508, 552, 600, 652, 704,
	768, 832, 896, 960, 1024,
}

// Table 4.137
var swb_offset_128_24 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 28, 36, 44, 52, 64, 76, 92, 108, 128,
}

// Table 4.134
var swb_offset_1024_16 = []uint16{
	0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 100, 112, 124,
	136, 148, 160, 172, 184, 196, 212, 228, 244, 260, 280, 300, 320, 344,
	368, 396, 424, 456, 492, 532, 572, 616, 664, 716, 772, 832, 896, 960, 1024,
}

// Table 4.135
var swb_offset_128_16 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 28, 32, 40, 48, 60, 72, 88, 108, 128,
}

// Table 4.132
var swb_offset_1024_8 = []uint16{
	0, 12, 24, 36, 48, 60, 72, 84, 96, 108, 120, 132, 144, 156, 172,
	188, 204, 220, 236, 252, 268, 288, 308, 328, 348, 372, 396, 420, 448,
	476, 508, 544, 580, 620, 664, 712, 764, 820, 880, 944, 1024,
}

// Table 4.133
var swb_offset_128_8 = []uint16{
	0, 4, 8, 12, 16, 20, 24, 28, 36, 44, 52, 60, 72, 88, 108, 128,
}

var swb_offset_long_window = [][]uint16{
	swb_offset_1024_96, swb_offset_1024_96, // 96000, 88200
	swb_offset_1024_64,                     // 64000
	swb_offset_1024_48, swb_offset_1024_48, // 48000, 441000
	swb_offset_1024_32,                     // 32000
	swb_offset_1024_24, swb_offset_1024_24, // 24000, 22050
	swb_offset_1024_16, swb_offset_1024_16, swb_offset_1024_16, // 16000, 12000, 11025
	swb_offset_1024_8, swb_offset_1024_8, // 8000
}

var swb_offset_short_window = [][]uint16{
	swb_offset_128_96, swb_offset_128_96, // 96000, 88200
	swb_offset_128_64,                                       // 64000
	swb_offset_128_48, swb_offset_128_48, swb_offset_128_48, // 48000, 44100, 32000
	swb_offset_128_24, swb_offset_128_24, // 24000, 22050
	swb_offset_128_16, swb_offset_128_16, swb_offset_128_16, // 16000, 12000, 11025
	swb_offset_128_8, swb_offset_128_8, // 8000
}

////////////////////////////////////////////////////////////////////////////////
// 4.5.2.3.4 - Scalefactor bands and grouping
////////////////////////////////////////////////////////////////////////////////
func window_grouping(info *ics_info, sfi uint8, framelength uint16) {
	idx := 1
	if framelength == 960 {
		idx = 0
	}

	switch info.Window_sequence {
	case ONLY_LONG_SEQUENCE, LONG_START_SEQUENCE, LONG_STOP_SEQUENCE:
		info.num_windows = 1
		info.num_window_groups = 1
		info.window_group_length = make([]uint8, info.num_window_groups)
		info.window_group_length[info.num_window_groups-1] = 1
		info.num_swb = num_swb_long_windows[idx][sfi]

		t := make([]uint16, len(swb_offset_long_window[sfi][:info.num_swb]))
		copy(t, swb_offset_long_window[sfi][:info.num_swb])
		info.sect_sfb_offset = append(info.sect_sfb_offset, t)

		t = make([]uint16, len(swb_offset_long_window[sfi][:info.num_swb]))
		copy(t, swb_offset_long_window[sfi][:info.num_swb])
		info.swb_offset = t

		// Special cases for 960's final values so we don't have to duplicate tables
		info.sect_sfb_offset[0] = append(info.sect_sfb_offset[0], framelength)
		info.swb_offset = append(info.swb_offset, framelength)
	case EIGHT_SHORT_SEQUENCE:
		info.num_windows = 8
		info.num_window_groups = 1
		info.window_group_length = make([]uint8, info.num_window_groups)
		info.window_group_length[info.num_window_groups-1] = 1
		info.num_swb = num_swb_short_window[sfi]

		t := make([]uint16, len(swb_offset_short_window[sfi][:info.num_swb]))
		copy(t, swb_offset_short_window[sfi][:info.num_swb])
		info.swb_offset = t
		info.swb_offset = append(info.swb_offset, framelength/8)

		for i := uint8(0); i < info.num_windows-1; i++ {
			bit := 6 - i
			if ((info.Scale_factor_grouping >> bit) & 0x1) == 0 {
				info.num_window_groups++
				info.window_group_length = append(info.window_group_length, 1)
			} else {
				info.window_group_length[info.num_window_groups-1]++
			}
		}

		info.sect_sfb_offset = make([][]uint16, info.num_window_groups)
		for g := range info.sect_sfb_offset {
			sect_sfb := 0
			offset := uint16(0)

			info.sect_sfb_offset[g] = make([]uint16, info.num_swb)
			for i := range info.sect_sfb_offset[g] {
				width := swb_offset_short_window[sfi][i+1] - swb_offset_short_window[sfi][i]
				// account for special table cases for 120
				if uint8(i+1) == info.num_swb {
					width = framelength/8 - swb_offset_short_window[sfi][i]
				}
				width *= uint16(info.window_group_length[g])
				info.sect_sfb_offset[g][sect_sfb] = offset
				sect_sfb++
				offset += width
			}
			info.sect_sfb_offset[g] = append(info.sect_sfb_offset[g], offset)
		}
	}
}
