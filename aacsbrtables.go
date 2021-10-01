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

import (
	"fmt"
	"math"
	"sort"
)

var startOffset = [][]int8{
	{-8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7}, // sfi = 8
	{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 9, 11, 13},  // sfi = 7
	{-5, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 9, 11, 13, 16},  // sfi = 6
	{-6, -4, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 9, 11, 13, 16},  // sfi = 5
	{-4, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 9, 11, 13, 16, 20},  // sfi = 4...2
	{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 9, 11, 13, 16, 20, 24},  // sfi < 2
}

// startMin equations from 4.6.18.3.2.1 broken out into a table.  This table
// is indexed by the sfi
var startMin = []uint8{
	7, 7, 10, 11, 12, 16, 16, 17, 24, 32, 35, 48,
}

// Diverging power series broken out into a table for quick lookup.  Each row
// is indexed by sfi
var stopOffset = [][]int8{
	{0, 2, 4, 6, 8, 11, 14, 18, 22, 26, 31, 37, 44, 51},
	{0, 2, 4, 6, 8, 11, 14, 18, 22, 26, 31, 36, 42, 49},
	{0, 2, 4, 6, 8, 11, 14, 17, 21, 25, 29, 34, 39, 44},
	{0, 2, 4, 6, 8, 11, 14, 17, 20, 24, 28, 33, 38, 43},
	{0, 2, 4, 6, 8, 11, 14, 17, 20, 24, 28, 32, 36, 41},
	{0, 2, 4, 6, 8, 10, 12, 14, 17, 20, 23, 26, 29, 32},
	{0, 2, 4, 6, 8, 10, 12, 14, 17, 20, 23, 26, 29, 32},
	{0, 1, 3, 5, 7, 9, 11, 13, 15, 17, 20, 23, 26, 29},
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 14, 16},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, -1, -2, -3, -4, -5, -6, -6, -6, -6, -6, -6, -6, -6},
	{0, -3, -6, -9, -12, -15, -18, -20, -22, -24, -26, -28, -30, -32},
}

// This table fucntions exactly like startMin
var stopMin = []int8{
	13, 15, 20, 21, 23, 32, 32, 35, 48, 64, 70, 96,
}

// Frequency band table derivation is described by ISO-IEC 14496-3 4.6.18.3.2
func derive_sbr_tables(data *sbr_extension_data, sfi uint8, bs_start_freq uint8, bs_stop_freq uint8,
	bs_freq_scale uint8, bs_alter_scale uint8, bs_xover_band uint8) error {
	data.k0 = uint8(qmf_lower_boundary(bs_start_freq, sfi))
	data.k2 = qmf_upper_boundary(bs_stop_freq, sfi, data.k0)

	// Restrictions on k0 and k2 from 4.6.18.3.6
	if sfi < 4 { // 48000+
		if data.k2-data.k0 > 32 {
			return fmt.Errorf("Error: k0 (%d) and k2 (%d) out of range", data.k0, data.k2)
		}
	} else if sfi > 4 { // 32000-
		if data.k2-data.k0 > 48 {
			return fmt.Errorf("Error: k0 (%d) and k2 (%d) out of range", data.k0, data.k2)
		}
	} else { // 44100
		if data.k2-data.k0 > 45 {
			return fmt.Errorf("Error: k0 (%d) and k2 (%d) out of range", data.k0, data.k2)
		}
	}

	if bs_freq_scale == 0 {
		freq_master_fs0(data, data.k0, data.k2, bs_alter_scale)
	} else {
		freq_master(data, data.k0, data.k2, bs_freq_scale, bs_alter_scale)
	}
	if err := freq_derived(data, bs_xover_band, data.k2); err != nil {
		return err
	}

	return nil
}

// k_0
func qmf_lower_boundary(bs_start_freq uint8, sfi uint8) int8 {
	i := 5
	switch sfi {
	case 8:
		i = 0
	case 7:
		i = 1
	case 6:
		i = 2
	case 5:
		i = 3
	case 4, 3, 2:
		i = 4
	case 0, 1:
		i = 5
	}

	return int8(startMin[sfi]) + startOffset[i][bs_start_freq]
}

// k_2
func qmf_upper_boundary(bs_stop_freq uint8, sfi uint8, k0 uint8) uint8 {
	val := minInt(64, int(stopMin[sfi]+stopOffset[sfi][minInt(13, int(bs_stop_freq))]))
	return uint8(val)
}

func freq_master_fs0(data *sbr_extension_data, k0 uint8, k2 uint8, bs_alter_scale uint8) {
	dk := 1
	numBands := 0
	if bs_alter_scale == 0 {
		// 2*INT(k2-k0)/(dk*2) why use division when you have bit shifts?
		numBands = int((((k2 - k0) >> 1) << 1))
	} else {
		// 2*NINT(k2-k0)/(dk*2)
		dk = 2
		numBands = int((((k2 - k0 + 2) >> 2) << 1))
	}

	numBands = minInt(63, numBands)
	k2Achieved := int(k0) + numBands*dk
	k2Diff := int(k2) - k2Achieved

	vDk := make([]int, numBands, numBands)
	for k := 0; k < numBands; k++ {
		vDk[k] = dk
	}

	if k2Diff != 0 {
		var k, incr int
		if k2Diff < 0 {
			incr = 1
			k = 0
		} else {
			incr = -1
			k = numBands - 1
		}

		for k2Diff != 0 {
			vDk[k] -= incr
			k += incr
			k2Diff += incr
		}
	}

	data.f_master = make([]int, numBands)
	data.f_master[0] = int(k0)
	for k := 1; k < numBands; k++ {
		data.f_master[k] = data.f_master[k-1] + vDk[k-1]
	}
	data.N_master = uint8(numBands)
}

func freq_master(data *sbr_extension_data, k0 uint8, k2 uint8, bs_freq_scale uint8, bs_alter_scale uint8) {
	twoRegions := 0
	k1 := k2

	if float64(k2)/float64(k0) > 2.2449 {
		twoRegions = 1
		k1 = k0 << 1 // 2*k0
	}

	temp1 := [...]float64{12.0, 10.0, 8.0}
	bands := temp1[bs_freq_scale-1]
	k0_f := float64(k0)
	k1_f := float64(k1)
	k2_f := float64(k2)
	// 2 * NINT( bands * log(k1/k2) / (2*log(2)))
	numBands0 := 2 * aacRound(bands*math.Log10(k1_f/k0_f)/(2.0*math.Log10(2)))

	vDk0 := make([]int, numBands0+1)
	// ew this is ugly...
	for k, _ := range vDk0 {
		// NINT( k0*(k1/k0)^((k+1)/numBands0) ) - NINT( k0*(k1/k0)^(k/numBands0) )
		vDk0[k] = aacRound(k0_f*(math.Pow((k1_f/k0_f), float64(k+1)/float64(numBands0)))) -
			aacRound(k0_f*math.Pow((k1_f/k0_f), float64(k)/float64(numBands0)))
	}
	sort.Sort(sort.IntSlice(vDk0))

	vk0 := make([]int, numBands0+1)
	vk0[0] = int(k0)
	for k := 1; k <= numBands0; k++ {
		vk0[k] = vk0[k-1] + vDk0[k-1]
	}

	if twoRegions == 0 {
		data.N_master = uint8(numBands0)
		data.f_master = vk0[0 : numBands0+1]
		return
	}

	// replaces allocation of temp2 array from the spec
	warp := 1.0
	if bs_alter_scale == 1 {
		warp = 1.3
	}

	// 2 * NINT( bands * log(k2/k1) / (2*log(2)*warp) )
	numBands1 := 2 * aacRound(bands*math.Log10(k2_f/k1_f)/(2.0*math.Log10(2.0)*warp))
	vDk1 := make([]int, numBands1-1)
	for k, _ := range vDk1 {
		// NINT( k1* (k2/k1)^((k+1)/numBands1) ) - NINT( k1* (k2/k1)^(k/numBands1) )
		vDk1[k] = aacRound(k1_f*math.Pow(k2_f/k1_f, float64((k+1))/float64(numBands1))) -
			aacRound(k1_f*math.Pow(k2_f/k1_f, (float64(k)/float64(numBands1))))
	}
	sort.Sort(sort.IntSlice(vDk1))

	// if min(vDk1) < max(vDk0)
	if vDk1[0] < vDk0[cap(vDk0)-1] {
		change := vDk0[cap(vDk0)-1] - vDk1[0]
		if change > (vDk1[numBands1-1] - vDk1[0]<<1) {
			change = (vDk1[numBands1-1] - vDk1[0]<<1)
		}
		vDk1[0] += change
		vDk1[numBands1-1] -= change
	}

	vk1 := make([]int, numBands1)
	vk1[0] = int(k1)
	for k := 1; k < numBands1; k++ {
		vk1[k] = vk1[k-1] + vDk1[k-1]
	}

	data.N_master = uint8(numBands0 + numBands1)
	data.f_master = make([]int, 0)
	data.f_master = append(data.f_master, vk0...)
	data.f_master = append(data.f_master, vk1...)
}

func freq_derived(data *sbr_extension_data, bs_xover_band uint8, k2 uint8) error {
	data.N_high = data.N_master - bs_xover_band
	data.N_low = (data.N_high >> 1) + (data.N_high - (data.N_high>>1)<<1)

	data.n = make([]uint8, 2)
	data.n[0] = data.N_low
	data.n[1] = data.N_high

	index := data.N_high + bs_xover_band + 1

	// check for overflow or upper bounds
	if index < bs_xover_band || len(data.f_master) < int(index) {
		return fmt.Errorf("f_tablehigh invalid index: index (%d) must be between bs_xover_band (%d) and length of f_master (%d)", index, bs_xover_band, len(data.f_master))
	}
	data.f_tablehigh = data.f_master[bs_xover_band:index]
	// check f_tablehigh bounds
	if len(data.f_tablehigh) < int(data.N_high) {
		return fmt.Errorf("N_high index (%d) too high for length of f_tablehigh (%d)", data.N_high, len(data.f_tablehigh))
	}
	data.M = uint8(data.f_tablehigh[data.N_high] - data.f_tablehigh[0])
	data.k_x = data.f_tablehigh[0]

	data.f_tablelow = make([]int, int(data.N_low+1))
	i := 0
	for k, _ := range data.f_tablelow {
		if k != 0 {
			i = 2*k - (int(data.N_high) % 2)
		}
		data.f_tablelow[k] = data.f_tablehigh[i]
	}

	data.N_Q = 0
	bs_noise_bands := data.Sbr_header.Bs_noise_bands
	k2_f := float64(data.k2)
	k_x_f := float64(data.k_x)
	data.N_Q = uint8(maxInt(1, aacRound(float64(bs_noise_bands)*(math.Log10(k2_f/k_x_f)/math.Log10(2)))))

	data.f_tablenoise = make([]int, int(data.N_Q+1))
	i = 0
	for k, _ := range data.f_tablenoise {
		if k != 0 {
			i = i + (int(data.N_low)-i)/(int(data.N_Q)+1-k)
		}
		data.f_tablenoise[k] = data.f_tablelow[i]
	}

	return nil
}

// TODO: Implement freq_limiter table
