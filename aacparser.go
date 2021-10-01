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

	"github.com/Comcast/gaad/bitreader"
)

const MaxBitsLeft = 131072

type ADTS struct {
	Bitrate              uint32
	ChannelConfiguration uint8
	Layer                uint8
	MpegVersion          uint8
	Profile              uint8
	SamplingFrequency    uint32
	VbrMode              bool
	Frame_length         uint16

	reader              *bitreader.BitReader
	aac_frame_length    uint16
	sfi                 uint8
	num_raw_data_blocks uint8
	protection_absent   bool

	Single_channel_elements   []*single_channel_element
	Channel_pair_elements     []*channel_pair_element
	Coupling_channel_elements []*coupling_channel_element
	Lfe_channel_elements      []*lfe_channel_element
	Data_stream_elements      []*data_stream_element
	Program_config_elements   []*program_config_element
	Fill_elements             []*fill_element
}

// Begin Main AAC Element Types
type single_channel_element struct {
	Element_instance_tag uint8
	Channel_stream       *individual_channel_stream
}

type channel_pair_element struct {
	Element_instance_tag uint8

	Common_window bool
	Ics_info      *ics_info
	Ms_used       [][]bool

	Channel_stream1 *individual_channel_stream
	Channel_stream2 *individual_channel_stream
}

type coupling_channel_element struct {
	Element_instance_tag uint8
	Ind_sw_cce_flag      bool
	Num_coupled_elements uint8
	Cc_target_is_cpe     []bool
	Cc_target_tag_select []uint8
	Cc_l                 []bool
	Cc_r                 []bool
	Cc_domain            bool
	Gain_element_sign    bool
	Gain_element_scale   uint8

	Channel_stream *individual_channel_stream

	Common_gain_element_present []bool
	Common_gain_element         []uint8

	DCPM_gain_element [][][]uint8
}

type lfe_channel_element struct {
	Element_instance_tag uint8
	Channel_stream       *individual_channel_stream
}

type data_stream_element struct {
	Element_instance_tag uint8
	Data_byte_align_flag bool
	Count                uint8
	Esc_count            uint8
	Data_stream_byte     [][]uint8
}

type program_config_element struct {
	Element_instance_tag     uint8
	Object_type              uint8
	Sampling_frequency_index uint8

	Num_front_channel_elements uint8
	Num_side_channel_elements  uint8
	Num_back_channel_elements  uint8
	Num_lfe_channel_elements   uint8
	Num_assoc_data_elements    uint8
	Num_valid_cc_elements      uint8

	Mono_mixdown_present     bool
	Mono_mixdown_element_num uint8

	Stereo_mixdown_present     bool
	Stereo_mixdown_element_num uint8

	Matrix_mixdown_idx_present bool
	Matrix_mixdown_idx         uint8
	Pseudo_surround_enable     bool

	Front_element_is_cpe     []bool
	Front_element_tag_select []uint8

	Side_element_is_cpe     []bool
	Side_element_tag_select []uint8

	Back_element_is_cpe     []bool
	Back_element_tag_select []uint8

	Lfe_element_tag_select        []uint8
	Assoc_data_element_tag_select []uint8

	Cc_element_is_ind_sw        []bool
	Valid_cc_element_tag_select []uint8

	Comment_field_bytes uint8
	Comment_field_data  []byte
}

type fill_element struct {
	Count     uint16
	Esc_count uint8

	Extension_payload *extension_payload
}

// End Main AAC element types
// Begin AAC element sub components
type adts_error_check struct {
	Crc_check uint16
}

type adts_header_error_check struct {
	Raw_data_block_position []uint16
	Crc_check               uint16
}

type adts_raw_data_block_error_check struct {
	Crc_check uint16
}

type dynamic_range_info struct {
	Pce_tag_present      bool
	Pce_instance_tag     uint8
	Drc_tag_reserve_bits uint8

	Excluded_chns_present    bool
	Excluded_chns            *excluded_channels
	Drc_bands_present        bool
	Drc_band_incr            uint8
	Drc_interpolation_scheme uint8
	Drc_band_top             []byte

	Prog_ref_level_present       bool
	Prog_ref_level               uint8
	Prog_ref_level_reserved_bits byte

	Dyn_range_sign []uint8
	Dyn_range_cnt  []uint8
}

type excluded_channels struct {
	Exclude_mask             []bool
	Additional_excluded_chns []bool
}

type extension_payload struct {
	Extension_type uint8
	Fill_nibble    uint8
	Fill_byte      []byte

	Data_element_version uint8
	// Yes, this one item is camel case.  It's that way in the spec
	// and all the other data items follow the exact syntax in the
	// spec for easy reference.  Weird.  Just have to deal with it.
	DataElementLengthPart uint8
	Data_element_byte     []byte

	Dynamic_range_info *dynamic_range_info
	Sac_extension_data *sac_extension_data
	Sbr_extension_data *sbr_extension_data

	Other_bits []bool
}

type gain_control_data struct {
	Max_band   uint8
	Alevcode   [][][]uint8
	Aloccode   [][][]uint8
	Adjust_num [][]uint8
}

type individual_channel_stream struct {
	Global_gain uint8

	Ics_info          *ics_info
	Section_data      *section_data
	Scale_factor_data *scale_factor_data

	Pulse_data_present bool
	Pulse_data         *pulse_data

	Tns_data_present bool
	Tns_data         *tns_data

	Gain_control_data_present bool
	Gain_control_data         *gain_control_data

	Spectral_data *spectral_data

	Length_of_reordered_spectral_data uint16
	Length_of_longest_code_word       uint8
	Reordered_spectral_data           *reordered_spectral_data
}

type ics_info struct {
	Window_sequence           uint8
	Window_shape              uint8
	Max_sfb                   uint8
	Scale_factor_grouping     uint8
	Predictor_data_present    bool
	Predictor_reset           bool
	Predictor_reset_group_num uint8
	Prediction_used           []bool

	num_windows         uint8
	num_window_groups   uint8
	window_group_length []uint8
	sect_sfb_offset     [][]uint16
	swb_offset          []uint16
	sfb_cb              [][]uint8
	num_swb             uint8

	Ltp_data_present bool
	Ltp_data         *ltp_data
}

type ltp_data struct {
	Ltp_lag       uint
	Ltp_coef      uint8
	Ltp_long_used []bool
}

type pulse_data struct {
	Number_pulse    uint8
	Pulse_start_sfb uint8
	Pulse_offset    []uint8
	Pulse_amp       []uint8
}

type reordered_spectral_data struct {
	Data []uint8
}

type sac_extension_data struct {
	// For some reason this data element in the spec
	// decided to deviate from the norm and use camel case.
	AncType            uint8
	AncStart           bool
	AncStop            bool
	AncDataSegmentByte []byte
}

type sbr_extension_data struct {
	Bs_sbr_crc_bits uint16
	Bs_header_flag  bool
	Bs_fill_bits    []byte

	Sbr_header *sbr_header
	Sbr_data   *sbr_data

	num_sbr_bits   uint
	num_align_bits uint

	// Derived frequency table parameters
	k0           uint8
	k2           uint8
	f_master     []int
	f_tablehigh  []int
	f_tablelow   []int
	f_tablenoise []int
	M            uint8
	k_x          int
	N_master     uint8
	N_high       uint8
	N_low        uint8
	n            []uint8
	N_Q          uint8
}

type sbr_header struct {
	Bs_amp_res        bool
	Bs_start_freq     uint8
	Bs_stop_freq      uint8
	Bs_xover_band     uint8
	Bs_reserved       uint8
	Bs_header_extra_1 bool
	Bs_header_extra_2 bool

	Bs_freq_scale  uint8
	Bs_alter_scale uint8
	Bs_noise_bands uint8

	Bs_limiter_bands  uint8
	Bs_limiter_gains  uint8
	Bs_interpol_freq  uint8
	Bs_smoothing_mode uint8
}

type sbr_data struct {
	Sbr_single_channel_element       *sbr_single_channel_element
	Sbr_channel_pair_element         *sbr_channel_pair_element
	Sbr_channel_pair_base_element    *sbr_channel_pair_base_element
	Sbr_channel_pair_enhance_element *sbr_channel_pair_enhance_element
}

type sbr_single_channel_element struct {
	Bs_data_extra bool
	Bs_reserved   uint8

	Sbr_grid     *sbr_grid
	Sbr_dtdf     *sbr_dtdf
	Sbr_invf     *sbr_invf
	Sbr_envelope *sbr_envelope
	Sbr_noise    *sbr_noise

	Bs_add_harmonic_flag  bool
	Sbr_sinusoidal_coding *sbr_sinusoidal_coding

	Bs_extended_data  bool
	Bs_extension_size uint8
	Bs_esc_count      uint8

	Bs_extension_id []uint8
	Sbr_extension   []*sbr_extension

	Bs_fill_bits []byte
}

type sbr_channel_pair_element struct {
	Bs_data_extra bool
	Bs_reserved_0 uint8
	Bs_reserved_1 uint8

	Bs_coupling bool
	Sbr_grid    *sbr_grid
	Sbr_dtdf    *sbr_dtdf
	Sbr_invf    *sbr_invf

	Sbr_envelope *sbr_envelope
	Sbr_noise    *sbr_noise

	Bs_add_harmonic_flag  []bool
	Sbr_sinusoidal_coding *sbr_sinusoidal_coding

	Bs_extended_data  bool
	Bs_extension_size uint8
	Bs_esc_count      uint8

	Bs_extension_id []uint8
	Sbr_extension   []*sbr_extension

	Bs_fill_bits []byte
}

type sbr_channel_pair_base_element struct {
	Bs_data_extra bool
	Bs_reserved_0 uint8
	Bs_reserved_1 uint8

	Bs_coupling bool
	Sbr_grid    *sbr_grid
	Sbr_dtdf    *sbr_dtdf
	Sbr_invf    *sbr_invf

	Sbr_envelope *sbr_envelope
	Sbr_noise    *sbr_noise

	Bs_add_harmonic_flag  bool
	Sbr_sinusoidal_coding *sbr_sinusoidal_coding

	Bs_extended_data  bool
	Bs_extension_size uint8
	Bs_esc_count      uint8

	Bs_extension_id []uint8
	Sbr_extension   []*sbr_extension

	Bs_fill_bits []byte
}

type sbr_channel_pair_enhance_element struct {
	Sbr_dtdf     *sbr_dtdf
	Sbr_envelope *sbr_envelope
	Sbr_noise    *sbr_noise

	Bs_add_harmonic_flag  bool
	Sbr_sinusoidal_coding *sbr_sinusoidal_coding
}

type sbr_grid struct {
	Bs_frame_class uint8

	Tmp         uint8 // Yes, this is an official bit field in the spec...
	Bs_freq_res [][]uint8

	Bs_var_bord_0 []uint8
	Bs_var_bord_1 []uint8
	Bs_num_rel_0  []uint8
	Bs_num_rel_1  []uint8
	Bs_pointer    []uint

	bs_num_env    []uint8
	bs_num_noise  []uint8
	bs_rel_bord_0 [][]uint8
	bs_rel_bord_1 [][]uint8
}

type sbr_dtdf struct {
	Bs_df_env   [][]bool
	Bs_df_noise [][]bool
}

type sbr_invf struct {
	Bs_invf_mode [][]uint8
}

type sbr_envelope struct {
	t_huff uint
	f_huff uint

	Bs_env_start_value_balance uint8
	Bs_env_start_value_level   uint8

	Bs_data_env [][][]int
}

type sbr_noise struct {
	t_huff uint
	f_huff uint

	Bs_noise_start_value_balance uint8
	Bs_noise_start_value_level   uint8

	Bs_data_noise [][][]int
}

type sbr_extension struct {
	Bs_fill_bits []byte
}

type sbr_sinusoidal_coding struct {
	Bs_add_harmonic [][]bool
}

type scale_factor_data struct {
	Dcpm_is_position [][]uint8
	Dcpm_noise_nrg   [][]uint16
	Dcpm_sf          [][]uint8

	Sf_concealment      bool
	Rev_global_gain     uint8
	Len_of_rvlc_sf      uint16
	Rvlc_cod_sf         uint8
	Sf_escapes_present  bool
	Len_of_rvlc_escapes uint8
	rvlc_esc_sf         uint8
	Dcpm_noise_last_pos uint16
}

type section_data struct {
	Sect_cb  [][]uint8
	Sect_len uint8

	sect_start [][]uint8
	sect_end   [][]uint16
	num_sec    []uint8
}

type spectral_data struct {
	Hcod           [][]int8
	Quad_sign_bits uint8
	Pair_sign_bits uint8
	Hcod_esc_y     uint32
	Hcod_esc_z     uint32
}

type tns_data struct {
	N_filt        []uint8
	Coef_res      []uint8
	Len           [][]uint8
	Order         [][]uint8
	Direction     [][]bool
	Coef_compress [][]uint8
	Coef          [][][]uint8
}

////////////////////////////////////////////////////////////////////////////////
// ID_SYN_ELE (Syntactic Element)
////////////////////////////////////////////////////////////////////////////////
const (
	ID_SCE = 0x00 // "Single Channel Element"
	ID_CPE = 0x01 // "Channel Pair Element"
	ID_CCE = 0x02 // "Coupling Channel Element"
	ID_LFE = 0x03 // "LFE Channel Element"
	ID_DSE = 0x04 // "Data Stream Element"
	ID_PCE = 0x05 // "Program Config Element"
	ID_FIL = 0x06 // "Fill Element"
	ID_END = 0x07 // "End"
)

var SyntacticElement = [...]string{
	"ID_SCE: Single Channel Element",
	"ID_CPE: Channel Pair Element",
	"ID_CCE: Coupling Channel Element",
	"ID_LFE: LFE Channel Element",
	"ID_DSE: Data Stream Element",
	"ID_PCE: Program Config Element",
	"ID_FIL: Fill Element",
	"ID_END: End",
}

////////////////////////////////////////////////////////////////////////////////
// Table 1.17 – Audio Object Types
////////////////////////////////////////////////////////////////////////////////
const (
	AUDIO_OBJECT_TYPE_NULL                uint8 = 0
	AUDIO_OBJECT_TYPE_AAC_MAIN                  = 1
	AUDIO_OBJECT_TYPE_AAC_LC                    = 2
	AUDIO_OBJECT_TYPE_SSR                       = 3
	AUDIO_OBJECT_TYPE_LTP                       = 4
	AUDIO_OBJECT_TYPE_SBR                       = 5
	AUDIO_OBJECT_TYPE_AAC_SCALABLE              = 6
	AUDIO_OBJECT_TYPE_TWINVQ                    = 7
	AUDIO_OBJECT_TYPE_CELP                      = 8
	AUDIO_OBJECT_TYPE_HXVC                      = 9
	AUDIO_OBJECT_TYPE_TTSI                      = 12
	AUDIO_OBJECT_TYPE_MAIN_SYNTHESIS            = 13
	AUDIO_OBJECT_TYPE_WAVETABLE_SYNTHESIS       = 14
	AUDIO_OBJECT_TYPE_GENERAL_MIDI              = 15
	AUDIO_OBJECT_TYPE_ASAE                      = 16
	AUDIO_OBJECT_TYPE_ER                        = 17
	AUDIO_OBJECT_TYPE_ER_AAC_LTP                = 19
	AUDIO_OBJECT_TYPE_ER_AAC_SCALABLE           = 20
	AUDIO_OBJECT_TYPE_ER_TWINVQ                 = 21
	AUDIO_OBJECT_TYPE_ER_BSAC                   = 22
	AUDIO_OBJECT_TYPE_ER_AAC_LD                 = 23
	AUDIO_OBJECT_TYPE_ER_CELP                   = 24
	AUDIO_OBJECT_TYPE_ER_HVXC                   = 25
	AUDIO_OBJECT_TYPE_ER_HILN                   = 26
	AUDIO_OBJECT_TYPE_ER_PARAMETRIC             = 27
	AUDIO_OBJECT_TYPE_SSC                       = 28
	AUDIO_OBJECT_TYPE_PS                        = 29
	AUDIO_OBJECT_TYPE_MPEG_SURROUND             = 30
	AUDIO_OBJECT_TYPE_LAYER_1                   = 32
	AUDIO_OBJECT_TYPE_LAYER_2                   = 33
	AUDIO_OBJECT_TYPE_LAYER_3                   = 34
	AUDIO_OBJECT_TYPE_DST                       = 35
	AUDIO_OBJECT_TYPE_ALS                       = 36
	AUDIO_OBJECT_TYPE_SLS                       = 37
	AUDIO_OBJECT_TYPE_SLS_NON_CORE              = 38
	AUDIO_OBJECT_TYPE_ER_AAC_ELD                = 39
	AUDIO_OBJECT_TYPE_SMR                       = 40
	AUDIO_OBJECT_TYPE_SMR_MAIN                  = 41
	AUDIO_OBJECT_TYPE_USAC_NO_SBR               = 42
	AUDIO_OBJECT_TYPE_SAOC                      = 43
	AUDIO_OBJECT_TYPE_LD_MPEG_SURROUND          = 44
	AUDIO_OBJECT_TYPE_USAC                      = 45
)

var AACProfileType = [...]string{
	"0: Null",
	"1: AAC Main",
	"2: AAC LC (Low Complexity)",
	"3: AAC SSR (Scalable Sample Rate)",
	"4: AAC LTP (Long Term Prediction)",
	"5: SBR (Spectral Band Replication)",
	"6: AAC Scalable",
	"7: TwinVQ",
	"8: CELP (Code Excited Linear Prediction)",
	"9: HXVC (Harmonic Vector eXcitation Coding)",
	"10: Reserved",
	"11: Reserved",
	"12: TTSI (Text-To-Speech Interface)",
	"13: Main Synthesis",
	"14: Wavetable Synthesis",
	"15: General MIDI",
	"16: Algorithmic Synthesis and Audio Effects",
	"17: ER (Error Resilient) AAC LC",
	"18: Reserved",
	"19: ER AAC LTP",
	"20: ER AAC Scalable",
	"21: ER TwinVQ",
	"22: ER BSAC (Bit-Sliced Arithmetic Coding)",
	"23: ER AAC LD (Low Delay)",
	"24: ER CELP",
	"25: ER HVXC",
	"26: ER HILN (Harmonic and Individual Lines plus Noise)",
	"27: ER Parametric",
	"28: SSC (SinuSoidal Coding)",
	"29: PS (Parametric Stereo)",
	"30: MPEG Surround",
	"31: (Escape value)",
	"32: Layer-1",
	"33: Layer-2",
	"34: Layer-3",
	"35: DST (Direct Stream Transfer)",
	"36: ALS (Audio Lossless)",
	"37: SLS (Scalable LosslesS)",
	"38: SLS non-core",
	"39: ER AAC ELD (Enhanced Low Delay)",
	"40: SMR (Symbolic Music Representation) Simple",
	"41: SMR Main",
	"42: USAC (Unified Speech and Audio Coding) (no SBR)",
	"43: SAOC (Spatial Audio Object Coding)",
	"44: LD MPEG Surround",
	"45: USAC",
}

////////////////////////////////////////////////////////////////////////////////
// Table 1.18 – Sampling Frequency Index
////////////////////////////////////////////////////////////////////////////////
var SamplingFrequency = [...]uint32{
	96000,
	88200,
	64000,
	48000,
	44100,
	32000,
	24000,
	22050,
	16000,
	12000,
	11025,
	8000,
	7350,
	0, // RESERVED
	0, // RESERVED
	0, // ESCAPE VALUE
}

////////////////////////////////////////////////////////////////////////////////
// Table 1.19 – Channel Configuration
////////////////////////////////////////////////////////////////////////////////
var ChannelConfiguration = [...]string{
	" 0: Defined in AOT Specifc Config",
	" 1: 1 channel: front-center",
	" 2: 2 channels: front-left, front-right",
	" 3: 3 channels: front-center, front-left, front-right",
	" 4: 4 channels: front-center, front-left, front-right, back-center",
	" 5: 5 channels: front-center, front-left, front-right, back-left, back-right",
	" 6: 6 channels: front-center, front-left, front-right, back-left, back-right, LFE-channel",
	" 7: 8 channels: front-center, front-left, front-right, side-left, side-right, back-left, back-right, LFE-channel",
	" 8: Reserved",
	" 9: Reserved",
	"10: Reserved",
	"11: Reserved",
	"12: Reserved",
	"13: Reserved",
	"14: Reserved",
	"15: Reserved",
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.114 – Values of the extension_type field
////////////////////////////////////////////////////////////////////////////////
const (
	FIXFIX = 0
	FIXVAR = 1
	VARFIX = 2
	VARVAR = 3
)

////////////////////////////////////////////////////////////////////////////////
// Table 4.121 – Values of the extension_type field
////////////////////////////////////////////////////////////////////////////////
const (
	EXTENSION_ID_PS = 2
)

////////////////////////////////////////////////////////////////////////////////
// Table 4.121 – Values of the extension_type field
////////////////////////////////////////////////////////////////////////////////
const (
	EXT_FILL          = 0x00 // ‘0000’ bitstream payload filler
	EXT_FILL_DATA     = 0x01 // ‘0001’ bitstream payload data as filler
	EXT_DATA_ELEMENT  = 0x02 // ’0010‘ data element
	EXT_DYNAMIC_RANGE = 0x0b // ‘1011’ dynamic range control
	EXT_SAC_DATA      = 0x0c // ‘1100’ MPEG Surround
	EXT_SBR_DATA      = 0x0d // ‘1101’ SBR enhancement
	EXT_SBR_DATA_CRC  = 0x0e // ‘1110’ SBR enhancement with CRC
)

////////////////////////////////////////////////////////////////////////////////
// Table 4.122 – Values of the data_element_version
////////////////////////////////////////////////////////////////////////////////
const (
	ANC_DATA = 0x00 // ‘0000‘ ￼Ancillary data element
)

////////////////////////////////////////////////////////////////////////////////
// AAC WINDOW SEQUENCE
////////////////////////////////////////////////////////////////////////////////
const (
	ONLY_LONG_SEQUENCE   = 0
	LONG_START_SEQUENCE  = 1
	EIGHT_SHORT_SEQUENCE = 2
	LONG_STOP_SEQUENCE   = 3
)

////////////////////////////////////////////////////////////////////////////////
// MPEG VERSION
////////////////////////////////////////////////////////////////////////////////
const (
	MPEG_VERSION_4 = 0
	MPEG_VERSION_2 = 1
)

////////////////////////////////////////////////////////////////////////////////
// AAC WINDOW SEQUENCE
////////////////////////////////////////////////////////////////////////////////
const (
	/* The noiseless coding tool requires these constants (see Table 4.56). */
	ZERO_HCB       = 0
	FIRST_PAIR_HCB = 5
	ESC_HCB        = 11
	QUAD_LEN       = 4
	PAIR_LEN       = 2
	NOISE_HCB      = 13
	INTENSITY_HCB2 = 14
	INTENSITY_HCB  = 15
	ESC_FLAG       = 16
)

////////////////////////////////////////////////////////////////////////////////
// 4.6.7.2 - Long Term Prediction (LTP) definitions
////////////////////////////////////////////////////////////////////////////////
const (
	MAX_LTP_LONG_SFB uint8 = 40
)

var Aac_PRED_SFB_MAX = [...]uint8{
	33, 33, 38, 40, 40, 40, 41, 41, 37, 37, 37, 34, 64, 64, 64, 64,
}

////////////////////////////////////////////////////////////////////////////////
// MAIN PARSE FUNCTION
////////////////////////////////////////////////////////////////////////////////
func ParseADTS(byteArray []byte) (*ADTS, error) {
	adts := &ADTS{}
	adts.reader = bitreader.NewBitReader(byteArray)
	err := adts.adts_frame()
	return adts, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 1.A.5 – Syntax of adts_frame()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) adts_frame() error {

	if adts.reader.HasByteLeft() {
		err := adts.adts_fixed_header()
		if err != nil {
			return err
		}
	}

	if adts.reader.HasByteLeft() {
		adts.adts_variable_header()
	}
	// Frame Length is fixed at 1024 for and ADTS
	adts.Frame_length = 1024
	if adts.num_raw_data_blocks == 0 {
		adts.adts_error_check()
		err := adts.raw_data_block()
		if err != nil {
			return err
		}
	} else {
		adts.adts_header_error_check()
		for i := uint8(0); i <= adts.num_raw_data_blocks; i++ {
			err := adts.raw_data_block()
			if err != nil {
				return err
			}
			adts.adts_raw_data_block_error_check()
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 1.A.6 – Syntax of adts_fixed_header()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) adts_fixed_header() error {
	sync_word_count := 0
	for sync_word_count < 3 && adts.reader.HasBitLeft() { // syncword 0xfff
		val, err := adts.reader.ReadBitsAsUInt8(4)
		if err != nil {
			return err
		}
		if val == 0x0f {
			sync_word_count++
		} else {
			sync_word_count = 0
		}
	}

	if adts.reader.HasBytesLeft(2) {
		adts.MpegVersion, _ = adts.reader.ReadBit()    // mpeg version
		adts.Layer, _ = adts.reader.ReadBitsAsUInt8(2) // layer; always 0
		if adts.Layer != 0 {
			return fmt.Errorf("ADTS Layer (%d) must be 0", adts.Layer)
		}
		adts.protection_absent, _ = adts.reader.ReadBitAsBool() // protection_absent
		adts.Profile, _ = adts.reader.ReadBitsAsUInt8(2)        // profile object
		adts.Profile += uint8(1)                                // profile object
		adts.sfi, _ = adts.reader.ReadBitsAsUInt8(4)            // sampling_frequency_index
		if adts.sfi > 12 {
			return fmt.Errorf("Sampling Frequency Index (%d) out of acceptable range (0-12)", adts.sfi)
		}

		adts.SamplingFrequency = SamplingFrequency[adts.sfi]          // sampling frequency
		adts.reader.SkipBits(1)                                       // private
		adts.ChannelConfiguration, _ = adts.reader.ReadBitsAsUInt8(3) // channel_configuration
		adts.reader.SkipBits(1)                                       // original
		adts.reader.SkipBits(1)                                       // home
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 1.A.7 – Syntax of adts_variable_header()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) adts_variable_header() {
	if adts.reader.HasBytesLeft(4) {
		adts.reader.SkipBits(1)                                      // copyright_id
		adts.reader.SkipBits(1)                                      // copyright_id_start
		adts.aac_frame_length, _ = adts.reader.ReadBitsAsUInt16(13)  // aac_frame_length
		adts_buffer_fullness, _ := adts.reader.ReadBitsAsUInt16(11)  // adts_buffer_fullness
		adts.num_raw_data_blocks, _ = adts.reader.ReadBitsAsUInt8(2) // num_raw_data_blocks

		if adts_buffer_fullness == 0x7ff {
			adts.VbrMode = true
		} else {
			adts.VbrMode = false
		}

		// ADTS is locked at 1024 samples
		adts.Bitrate = adts.SamplingFrequency / 1024
		adts.Bitrate *= uint32(adts.aac_frame_length) * 8
	}
}

////////////////////////////////////////////////////////////////////////////////
// Table 1.A.8 – Syntax of adts_error_check
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) adts_error_check() {
	if !adts.protection_absent {
		adts.reader.SkipBits(16) // crc_check;
	}
}

////////////////////////////////////////////////////////////////////////////////
// Table 1.A.9 – Syntax of adts_header_error_check
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) adts_header_error_check() {
	data := &adts_header_error_check{}

	if !adts.protection_absent {
		data.Raw_data_block_position = make([]uint16, adts.num_raw_data_blocks+1)
		for i := uint8(1); i <= adts.num_raw_data_blocks; i++ {
			data.Raw_data_block_position[i], _ = adts.reader.ReadBitsAsUInt16(16) // raw_data_block_position
		}
		data.Crc_check, _ = adts.reader.ReadBitsAsUInt16(16) // crc_check
	}
}

////////////////////////////////////////////////////////////////////////////////
// Table 1.A.10 – Syntax of adts_raw_data_block_error_check()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) adts_raw_data_block_error_check() {
	if !adts.protection_absent {
		adts.reader.SkipBits(16) // crc_check
	}
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.2 – Syntax of program_config_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) program_config_element() *program_config_element {
	e := &program_config_element{}
	e.Element_instance_tag, _ = adts.reader.ReadBitsAsUInt8(4) // element_instance_tag

	e.Object_type, _ = adts.reader.ReadBitsAsUInt8(2)                // object_type
	e.Sampling_frequency_index, _ = adts.reader.ReadBitsAsUInt8(4)   // sampling_frequency_index
	e.Num_front_channel_elements, _ = adts.reader.ReadBitsAsUInt8(4) // num_front_channel_elements
	e.Num_side_channel_elements, _ = adts.reader.ReadBitsAsUInt8(4)  // num_side_channel_elements
	e.Num_back_channel_elements, _ = adts.reader.ReadBitsAsUInt8(4)  // num_back_channel_elements
	e.Num_lfe_channel_elements, _ = adts.reader.ReadBitsAsUInt8(2)   // num_lfe_channel_elements
	e.Num_assoc_data_elements, _ = adts.reader.ReadBitsAsUInt8(3)    // num_assoc_data_elements;
	e.Num_valid_cc_elements, _ = adts.reader.ReadBitsAsUInt8(4)      // num_valid_cc_elements

	if e.Mono_mixdown_present, _ = adts.reader.ReadBitAsBool(); e.Mono_mixdown_present {
		e.Mono_mixdown_element_num, _ = adts.reader.ReadBitsAsUInt8(4) // mono_mixdown_element_number
	}
	if e.Stereo_mixdown_present, _ = adts.reader.ReadBitAsBool(); e.Stereo_mixdown_present {
		e.Stereo_mixdown_element_num, _ = adts.reader.ReadBitsAsUInt8(4) // stereo_mixdown_element_number
	}
	if e.Matrix_mixdown_idx_present, _ = adts.reader.ReadBitAsBool(); e.Matrix_mixdown_idx_present {
		e.Matrix_mixdown_idx, _ = adts.reader.ReadBitsAsUInt8(2)  // matrix_mixdown_idx
		e.Pseudo_surround_enable, _ = adts.reader.ReadBitAsBool() // pseudo_surround_enable
	}

	e.Front_element_is_cpe = make([]bool, e.Num_front_channel_elements)
	e.Front_element_tag_select = make([]uint8, e.Num_front_channel_elements)
	for i := range e.Front_element_tag_select {
		e.Front_element_is_cpe[i], _ = adts.reader.ReadBitAsBool()        // front_element_is_cpe[i]
		e.Front_element_tag_select[i], _ = adts.reader.ReadBitsAsUInt8(4) // front_element_tag_select[i]
	}

	e.Side_element_is_cpe = make([]bool, e.Num_side_channel_elements)
	e.Side_element_tag_select = make([]uint8, e.Num_side_channel_elements)
	for i := range e.Side_element_tag_select {
		e.Side_element_is_cpe[i], _ = adts.reader.ReadBitAsBool()        // side_element_is_cpe[i]
		e.Side_element_tag_select[i], _ = adts.reader.ReadBitsAsUInt8(4) // side_element_tag_select[i]
	}

	e.Back_element_is_cpe = make([]bool, e.Num_back_channel_elements)
	e.Back_element_tag_select = make([]uint8, e.Num_back_channel_elements)
	for i := range e.Back_element_tag_select {
		e.Back_element_is_cpe[i], _ = adts.reader.ReadBitAsBool()        // back_element_is_cpe[i]
		e.Back_element_tag_select[i], _ = adts.reader.ReadBitsAsUInt8(4) // back_element_tag_select[i]
	}

	e.Lfe_element_tag_select = make([]uint8, e.Num_lfe_channel_elements)
	for i := range e.Lfe_element_tag_select {
		e.Lfe_element_tag_select[i], _ = adts.reader.ReadBitsAsUInt8(4) // lfe_element_tag_select[i]
	}

	e.Assoc_data_element_tag_select = make([]uint8, e.Num_assoc_data_elements)
	for i := range e.Assoc_data_element_tag_select {
		e.Assoc_data_element_tag_select[i], _ = adts.reader.ReadBitsAsUInt8(4) // assoc_data_element_tag_select[i]
	}

	e.Cc_element_is_ind_sw = make([]bool, e.Num_valid_cc_elements)
	e.Valid_cc_element_tag_select = make([]uint8, e.Num_valid_cc_elements)
	for i := range e.Valid_cc_element_tag_select {
		e.Cc_element_is_ind_sw[i], _ = adts.reader.ReadBitAsBool()           // cc_element_is_ind_sw[i]
		e.Valid_cc_element_tag_select[i], _ = adts.reader.ReadBitsAsUInt8(4) // valid_cc_element_tag_select[i]
	}

	adts.reader.ByteAlign()
	e.Comment_field_bytes, _ = adts.reader.ReadBitsAsUInt8(8)                                  // comment_field_bytes
	e.Comment_field_data, _ = adts.reader.ReadBitsToByteArray(uint(e.Comment_field_bytes) * 8) // comment_field_data[i]

	return e
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.3 – Syntax of top level payload for audio object types AAC Main,
//             SSR, LC, and LTP (raw_data_block())
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) raw_data_block() error {
	var err error
	var id_syn_ele uint8 = 0
	var id_syn_ele_Previous uint8

	for id_syn_ele != ID_END {
		id_syn_ele_Previous = id_syn_ele
		id_syn_ele, _ = adts.reader.ReadBits(3)

		switch id_syn_ele {
		case ID_SCE:
			var e *single_channel_element
			e, err = adts.single_channel_element()
			adts.Single_channel_elements = append(adts.Single_channel_elements, e)
		case ID_CPE:
			var e *channel_pair_element
			e, err = adts.channel_pair_element()
			adts.Channel_pair_elements = append(adts.Channel_pair_elements, e)
		case ID_CCE:
			var e *coupling_channel_element
			e, err = adts.coupling_channel_element()
			adts.Coupling_channel_elements = append(adts.Coupling_channel_elements, e)
		case ID_LFE:
			var e *lfe_channel_element
			e, err = adts.lfe_channel_element()
			adts.Lfe_channel_elements = append(adts.Lfe_channel_elements, e)
		case ID_DSE:
			e := adts.data_stream_element()
			adts.Data_stream_elements = append(adts.Data_stream_elements, e)
		case ID_PCE:
			e := adts.program_config_element()
			adts.Program_config_elements = append(adts.Program_config_elements, e)
		case ID_FIL:
			var e *fill_element
			e, err = adts.fill_element(id_syn_ele_Previous)
			adts.Fill_elements = append(adts.Fill_elements, e)
		case ID_END:
		default:
			err = fmt.Errorf("Error: Unsupported id_syn_ele: %d", id_syn_ele)
		}

		if id_syn_ele != ID_END && adts.reader.HasBitLeft() == false {
			err = fmt.Errorf("Error: Buffer empty parsing id_syn_ele %d", id_syn_ele)
		}
		if err != nil {
			return err
		}
	}

	adts.reader.ByteAlign()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.4 – Syntax of single_channel_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) single_channel_element() (*single_channel_element, error) {
	var err error
	e := &single_channel_element{}
	e.Element_instance_tag, _ = adts.reader.ReadBitsAsUInt8(4)                // element_instance_tag
	e.Channel_stream, err = adts.individual_channel_stream(false, false, nil) // individual_channel_stream(0,0)
	return e, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.5 – Syntax of channel_pair_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) channel_pair_element() (*channel_pair_element, error) {
	var err error
	e := &channel_pair_element{}
	e.Element_instance_tag, _ = adts.reader.ReadBitsAsUInt8(4) // element_instance_tag

	e.Common_window, _ = adts.reader.ReadBitAsBool() // common_window
	if e.Common_window {
		e.Ics_info, err = adts.ics_info(e.Common_window)
		if err != nil {
			return e, err
		}
		ms_mask_present, _ := adts.reader.ReadBitsAsUInt8(2) // ms_mask_present
		if ms_mask_present == 3 {
			return e, fmt.Errorf("Error: ms_mask_present (%d) out of range", ms_mask_present)
		}

		if ms_mask_present == 1 {
			e.Ms_used = make([][]bool, e.Ics_info.num_window_groups)
			for g, _ := range e.Ms_used {
				e.Ms_used[g] = make([]bool, e.Ics_info.Max_sfb)
				for sfb, _ := range e.Ms_used[g] {
					e.Ms_used[g][sfb], _ = adts.reader.ReadBitAsBool() // ms_used[g][sfb]
				}
			}
		}
	}

	e.Channel_stream1, err = adts.individual_channel_stream(e.Common_window, false, e.Ics_info)
	if err != nil {
		return e, err
	}

	e.Channel_stream2, err = adts.individual_channel_stream(e.Common_window, false, e.Ics_info)
	return e, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.6 – Syntax of ics_info()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) ics_info(common_window bool) (*ics_info, error) {
	var err error
	info := &ics_info{}
	if ics_reserved_bit, _ := adts.reader.ReadBitsAsUInt8(1); ics_reserved_bit != 0 {
		err = fmt.Errorf("Error: ics_reserved_bit must equal 0")
		return nil, err
	}

	info.Window_sequence, _ = adts.reader.ReadBits(2)     // window_sequence
	info.Window_shape, _ = adts.reader.ReadBitsAsUInt8(1) // window_shape

	if info.Window_sequence == EIGHT_SHORT_SEQUENCE {
		info.Max_sfb, _ = adts.reader.ReadBits(4)
		info.Scale_factor_grouping, _ = adts.reader.ReadBits(7) // scale_factor_grouping

		window_grouping(info, adts.sfi, adts.Frame_length)
		if info.Max_sfb > info.num_swb {
			err = fmt.Errorf("Error: ics_info.Max_sfb (%d) must be less than ics_info.num_swb (%d)",
				info.Max_sfb, info.num_swb)
			return nil, err
		}
	} else {
		window_grouping(info, adts.sfi, adts.Frame_length)

		info.Max_sfb, _ = adts.reader.ReadBits(6)
		if info.Max_sfb > info.num_swb {
			err = fmt.Errorf("Error: ics_info.Max_sfb (%d) must be less than ics_info.num_swb (%d)",
				info.Max_sfb, info.num_swb)
			return nil, err
		}

		info.Predictor_data_present, _ = adts.reader.ReadBitAsBool()
		if info.Predictor_data_present == true {
			if adts.Profile == AUDIO_OBJECT_TYPE_AAC_MAIN {
				info.Predictor_reset, _ = adts.reader.ReadBitAsBool()
				if info.Predictor_reset == true {
					info.Predictor_reset_group_num, _ = adts.reader.ReadBits(5) // predictor_reset_group_number
				}

				PRED_SFB_MAX := minInt(int(info.Max_sfb), int(Aac_PRED_SFB_MAX[adts.sfi]))
				info.Prediction_used = make([]bool, PRED_SFB_MAX)
				for sfb := range info.Prediction_used {
					info.Prediction_used[sfb], _ = adts.reader.ReadBitAsBool()
				}

			} else {
				if info.Ltp_data_present, _ = adts.reader.ReadBitAsBool(); info.Ltp_data_present {
					info.Ltp_data, err = adts.ltp_data(info)
				}
				if common_window {
					if info.Ltp_data_present, _ = adts.reader.ReadBitAsBool(); info.Ltp_data_present {
						info.Ltp_data, err = adts.ltp_data(info)
					}
				}
			}
		}
	}
	return info, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.7 – Syntax of pulse_data()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) pulse_data() *pulse_data {
	data := &pulse_data{}
	data.Number_pulse, _ = adts.reader.ReadBitsAsUInt8(2)    // number_pulse
	data.Pulse_start_sfb, _ = adts.reader.ReadBitsAsUInt8(6) // pulse_start_sfb
	data.Pulse_amp = make([]uint8, data.Number_pulse)
	data.Pulse_offset = make([]uint8, data.Number_pulse)
	for i := range data.Pulse_amp {
		data.Pulse_offset[i], _ = adts.reader.ReadBitsAsUInt8(5) // pulse_offset[i]
		data.Pulse_amp[i], _ = adts.reader.ReadBitsAsUInt8(4)    // pulse_amp[i]
	}

	return data
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.8 – Syntax of coupling_channel_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) coupling_channel_element() (*coupling_channel_element, error) {
	var err error
	e := &coupling_channel_element{}
	e.Element_instance_tag, _ = adts.reader.ReadBitsAsUInt8(4) // element_instance_tag

	e.Ind_sw_cce_flag, _ = adts.reader.ReadBitAsBool()         // ind_sw_cce_flag
	e.Num_coupled_elements, _ = adts.reader.ReadBitsAsUInt8(3) // num_coupled_elements
	num_gain_element_lists := 0

	e.Cc_target_is_cpe = make([]bool, e.Num_coupled_elements)
	e.Cc_target_tag_select = make([]uint8, e.Num_coupled_elements)
	for c, _ := range e.Cc_target_is_cpe {
		num_gain_element_lists++
		e.Cc_target_is_cpe[c], _ = adts.reader.ReadBitAsBool()        // cc_target_is_cpe[c]
		e.Cc_target_tag_select[c], _ = adts.reader.ReadBitsAsUInt8(4) // cc_target_tag_select[c]
		if e.Cc_target_is_cpe[c] {
			if e.Cc_l == nil && e.Cc_r == nil {
				e.Cc_l = make([]bool, e.Num_coupled_elements)
				e.Cc_r = make([]bool, e.Num_coupled_elements)
			}
			e.Cc_l[c], _ = adts.reader.ReadBitAsBool() // cc_l[c]
			e.Cc_r[c], _ = adts.reader.ReadBitAsBool() // cc_r[c]
			if e.Cc_l[c] == true && e.Cc_r[c] == true {
				num_gain_element_lists++
			}
		}
	}

	e.Cc_domain, _ = adts.reader.ReadBitAsBool()
	e.Gain_element_sign, _ = adts.reader.ReadBitAsBool()
	e.Gain_element_scale, _ = adts.reader.ReadBitsAsUInt8(2)

	e.Channel_stream, err = adts.individual_channel_stream(false, false, nil)
	if err != nil {
		return e, err
	}

	e.Common_gain_element_present = make([]bool, num_gain_element_lists)
	e.Common_gain_element = make([]uint8, num_gain_element_lists)
	e.DCPM_gain_element = make([][][]uint8, num_gain_element_lists)
	for c := 1; c < num_gain_element_lists; c++ {
		cge := false
		if e.Ind_sw_cce_flag {
			cge = true
		} else {
			e.Common_gain_element_present[c], _ = adts.reader.ReadBitAsBool() // common_gain_element_present[c]
			cge = e.Common_gain_element_present[c]
		}

		if cge {
			e.Common_gain_element[c], err = hcod_sf(adts.reader) // common_gain_element[c])
		} else {
			info := e.Channel_stream.Ics_info
			e.DCPM_gain_element[c] = make([][]uint8, info.num_window_groups)
			for g := range e.DCPM_gain_element[c] {
				e.DCPM_gain_element[c][g] = make([]uint8, info.Max_sfb)
				for sfb := range e.DCPM_gain_element[c][g] {
					if info.sfb_cb[g][sfb] != ZERO_HCB {
						e.DCPM_gain_element[c][g][sfb], err = hcod_sf(adts.reader) //[dpcm_gain_element[c][g][sfb]]; 1..19
						if err != nil {
							return e, err
						}
					}
				}
			}
		}
	}

	return e, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.9 – Syntax of lfe_channel_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) lfe_channel_element() (*lfe_channel_element, error) {
	var err error
	e := &lfe_channel_element{}
	e.Element_instance_tag, _ = adts.reader.ReadBitsAsUInt8(4) // element_instance_tag

	e.Channel_stream, err = adts.individual_channel_stream(false, false, nil)
	return e, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.10 – Syntax of data_stream_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) data_stream_element() *data_stream_element {
	e := &data_stream_element{}
	e.Element_instance_tag, _ = adts.reader.ReadBitsAsUInt8(4) // element_instance_tag

	e.Data_byte_align_flag, _ = adts.reader.ReadBitAsBool() // data_byte_align_flag
	e.Count, _ = adts.reader.ReadBitsAsUInt8(8)
	if e.Count == 255 {
		e.Esc_count, _ = adts.reader.ReadBitsAsUInt8(8)
		e.Count += e.Esc_count
	}
	if e.Data_byte_align_flag == true {
		adts.reader.ByteAlign()
	}
	e.Data_stream_byte = make([][]uint8, e.Element_instance_tag+1)
	e.Data_stream_byte[e.Element_instance_tag] = make([]uint8, e.Count)
	for i := range e.Data_stream_byte[e.Element_instance_tag] {
		e.Data_stream_byte[e.Element_instance_tag][i], _ = adts.reader.ReadBitsAsUInt8(8) // data_stream_byte[element_instance_tag][i]
	}
	return e
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.11 – Syntax of fill_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) fill_element(id_syn_ele uint8) (*fill_element, error) {
	var err error
	e := &fill_element{}

	e.Count, _ = adts.reader.ReadBitsAsUInt16(4)
	if e.Count == 15 {
		e.Esc_count, _ = adts.reader.ReadBitsAsUInt8(8)
		e.Count += uint16(e.Esc_count) - 1
	}
	for cnt := int(e.Count); cnt > 0; {
		var sub int
		sub, e.Extension_payload, err = adts.extension_payload(cnt, id_syn_ele)
		cnt -= sub
		if err != nil {
			return e, err
		}
	}

	return e, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.12 – Syntax of gain_control_data()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) gain_control_data(info *ics_info) *gain_control_data {
	data := &gain_control_data{}

	data.Max_band, _ = adts.reader.ReadBitsAsUInt8(2)
	data.Adjust_num = make([][]uint8, data.Max_band)
	data.Alevcode = make([][][]uint8, data.Max_band)
	data.Aloccode = make([][][]uint8, data.Max_band)
	switch {
	case info.Window_sequence == ONLY_LONG_SEQUENCE:
		for bd := uint8(1); bd < data.Max_band; bd++ {
			data.Adjust_num[bd] = make([]uint8, 1)
			data.Alevcode[bd] = make([][]uint8, 1)
			data.Aloccode[bd] = make([][]uint8, 1)
			for wd := range data.Adjust_num[bd] {
				data.Adjust_num[bd][wd], _ = adts.reader.ReadBitsAsUInt8(3)
				data.Alevcode[bd][wd] = make([]uint8, data.Adjust_num[bd][wd])
				data.Aloccode[bd][wd] = make([]uint8, data.Adjust_num[bd][wd])
				for ad := range data.Alevcode[bd][wd] {
					data.Alevcode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(4)
					data.Aloccode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(5)
				}
			}
		}
	case info.Window_sequence == LONG_START_SEQUENCE:
		for bd := uint8(1); bd < data.Max_band; bd++ {
			data.Adjust_num[bd] = make([]uint8, 2)
			data.Alevcode[bd] = make([][]uint8, 2)
			data.Aloccode[bd] = make([][]uint8, 2)
			for wd := range data.Adjust_num[bd] {
				data.Adjust_num[bd][wd], _ = adts.reader.ReadBitsAsUInt8(3)
				data.Alevcode[bd][wd] = make([]uint8, data.Adjust_num[bd][wd])
				data.Aloccode[bd][wd] = make([]uint8, data.Adjust_num[bd][wd])
				for ad := range data.Alevcode[bd][wd] {
					data.Alevcode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(4)
					if wd == 0 {
						data.Aloccode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(4)
					} else {
						data.Aloccode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(2)
					}
				}
			}
		}
	case info.Window_sequence == EIGHT_SHORT_SEQUENCE:
		for bd := uint8(1); bd < data.Max_band; bd++ {
			data.Adjust_num[bd] = make([]uint8, 8)
			data.Alevcode[bd] = make([][]uint8, 8)
			data.Aloccode[bd] = make([][]uint8, 8)
			for wd := range data.Adjust_num[bd] {
				data.Adjust_num[bd][wd], _ = adts.reader.ReadBitsAsUInt8(3)
				data.Alevcode[bd][wd] = make([]uint8, data.Adjust_num[bd][wd])
				data.Aloccode[bd][wd] = make([]uint8, data.Adjust_num[bd][wd])
				for ad := range data.Alevcode[bd][wd] {
					data.Alevcode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(4)
					data.Aloccode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(2)
				}
			}
		}
	case info.Window_sequence == LONG_STOP_SEQUENCE:
		for bd := uint8(1); bd < data.Max_band; bd++ {
			data.Adjust_num[bd] = make([]uint8, 2)
			data.Alevcode[bd] = make([][]uint8, 2)
			data.Aloccode[bd] = make([][]uint8, 2)
			for wd := range data.Adjust_num[bd] {
				data.Adjust_num[bd][wd], _ = adts.reader.ReadBitsAsUInt8(3)
				data.Alevcode[bd][wd] = make([]uint8, data.Adjust_num[bd][wd])
				data.Aloccode[bd][wd] = make([]uint8, data.Adjust_num[bd][wd])
				for ad := range data.Alevcode[bd][wd] {
					data.Alevcode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(4)
					if wd == 0 {
						data.Aloccode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(4)
					} else {
						data.Aloccode[bd][wd][ad], _ = adts.reader.ReadBitsAsUInt8(5)
					}
				}
			}
		}
	default:
		return nil
	}

	return data
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.50 – Syntax of individual_channel_stream()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) individual_channel_stream(common_window bool, scale_flag bool, info *ics_info) (*individual_channel_stream, error) {
	var err error
	s := &individual_channel_stream{}
	s.Global_gain, _ = adts.reader.ReadBitsAsUInt8(8) // global_gain

	if !common_window && !scale_flag {
		s.Ics_info, err = adts.ics_info(common_window)
		if err != nil {
			return s, err
		}
	} else {
		s.Ics_info = info
	}

	s.Section_data, err = adts.section_data(s.Ics_info)
	if err != nil {
		return s, err
	}

	s.Scale_factor_data, err = adts.scale_factor_data(s.Ics_info)
	if err != nil {
		return s, err
	}

	if !scale_flag {
		s.Pulse_data_present, _ = adts.reader.ReadBitAsBool() // pulse_data_present
		if s.Pulse_data_present == true {
			s.Pulse_data = adts.pulse_data()
		}

		s.Tns_data_present, _ = adts.reader.ReadBitAsBool() // tns_data_present
		if s.Tns_data_present == true {
			s.Tns_data = adts.tns_data(s.Ics_info)
		}

		s.Gain_control_data_present, _ = adts.reader.ReadBitAsBool() // gain_control_data_present
		if s.Gain_control_data_present == true {
			adts.gain_control_data(s.Ics_info)
		}
	}

	// aacSpectralDataResilienceFlag is hard coded to false in other AAC decoders.
	// We'll follow suit here.
	//if !aacSpectralDataResilienceFlag {
	s.Spectral_data, err = adts.spectral_data(s.Ics_info, s.Section_data)
	/*} else {
		adts.reader.SkipBits(14) // length_of_reordered_spectral_data
		adts.reader.SkipBits(6)  // length_of_longest_codeword
		adts.reordered_spectral_data()
	}*/
	return s, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.51 – Syntax of reordered_spectral_data ()
////////////////////////////////////////////////////////////////////////////////
// func (adts *ADTS) reordered_spectral_data() {
// TODO
/* complex reordering, see tool description of Huffman
codeword reordering (subclause 4.6.16.3) */
// }

////////////////////////////////////////////////////////////////////////////////
// Table 4.52 – Syntax of section_data()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) section_data(info *ics_info) (*section_data, error) {
	var err error
	data := &section_data{}

	// This param is hardcoded to false for all other AAC Decoders I've found.
	// We'll just follow suit for now
	aacSectionDataResilienceFlag := false

	var bits uint
	if info.Window_sequence == EIGHT_SHORT_SEQUENCE {
		bits = 3
	} else {
		bits = 5
	}

	sect_esc_val := uint8((1 << uint8(bits)) - 1)
	data.Sect_cb = make([][]uint8, info.num_window_groups)
	data.sect_start = make([][]uint8, info.num_window_groups)
	data.sect_end = make([][]uint16, info.num_window_groups)
	data.num_sec = make([]uint8, info.num_window_groups)
	info.sfb_cb = make([][]uint8, info.num_window_groups)
	for g := range data.Sect_cb {
		i := uint8(0)
		k := uint8(0)
		for k < info.Max_sfb {
			if aacSectionDataResilienceFlag {
				val, _ := adts.reader.ReadBits(5) // sect_cb[g][i]
				data.Sect_cb[g] = append(data.Sect_cb[g], val)
			} else {
				val, _ := adts.reader.ReadBits(4) // sect_cb[g][i]
				data.Sect_cb[g] = append(data.Sect_cb[g], val)
			}

			sect_len := uint8(0)
			sect_len_incr := uint8(0)

			if !aacSectionDataResilienceFlag ||
				data.Sect_cb[g][i] < 11 ||
				(data.Sect_cb[g][i] > 11 && data.Sect_cb[g][i] < 16) {

				sect_len_incr, _ = adts.reader.ReadBitsAsUInt8(bits)
				for sect_len_incr == sect_esc_val {
					sect_len += sect_len_incr
					sect_len_incr, err = adts.reader.ReadBitsAsUInt8(bits)
					if err != nil {
						// We need to check this err so we don't get stuck in
						// an infinite loop
						return data, err
					}
				}
			} else {
				sect_len_incr = 1
			}

			sect_len += sect_len_incr
			data.sect_start[g] = append(data.sect_start[g], uint8(k))
			data.sect_end[g] = append(data.sect_end[g], uint16(k+sect_len))

			for j := uint8(0); j < sect_len; j++ {
				info.sfb_cb[g] = append(info.sfb_cb[g], data.Sect_cb[g][i])
			}

			k += sect_len
			i++

			upper_bound := info.Max_sfb
			if info.Window_sequence == EIGHT_SHORT_SEQUENCE {
				upper_bound = 8 * 15
			}
			if k > upper_bound || i > upper_bound {
				err = fmt.Errorf("Error: Section Codebook param out of bounds(%d). End (%d), Index (%d)",
					upper_bound, k, i)
				return data, err
			}

		}
		data.num_sec[g] = i
		if k != info.Max_sfb {
			err = fmt.Errorf("Error: Total length (%d) does not equal Max_sfb (%d)", k, info.Max_sfb)
			return data, err
		}
	}
	return data, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.53 – Syntax of scale_factor_data()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) scale_factor_data(info *ics_info) (*scale_factor_data, error) {
	var err error
	data := &scale_factor_data{}
	// This param is hardcoded to false for all other AAC Decoders I've found.
	// We'll just follow suit for now
	//	if !aacSectionDataResilienceFlag {
	noise_pcm_flag := true

	data.Dcpm_is_position = make([][]uint8, info.num_window_groups)
	data.Dcpm_noise_nrg = make([][]uint16, info.num_window_groups)
	data.Dcpm_sf = make([][]uint8, info.num_window_groups)
	for g := uint8(0); g < info.num_window_groups; g++ {
		data.Dcpm_is_position[g] = make([]uint8, info.Max_sfb)
		data.Dcpm_noise_nrg[g] = make([]uint16, info.Max_sfb)
		data.Dcpm_sf[g] = make([]uint8, info.Max_sfb)

		for sfb := range data.Dcpm_sf[g] {
			if info.sfb_cb[g][sfb] != ZERO_HCB {
				if is_intensity(info, g, uint8(sfb)) != 0 {
					data.Dcpm_is_position[g][sfb], err = hcod_sf(adts.reader)
				} else {
					if is_noise(info, g, uint8(sfb)) {
						if noise_pcm_flag {
							noise_pcm_flag = false
							data.Dcpm_noise_nrg[g][sfb], _ = adts.reader.ReadBitsAsUInt16(9)
						} else {
							var val uint8
							val, err = hcod_sf(adts.reader)
							data.Dcpm_noise_nrg[g][sfb] = uint16(val)
						}
					} else {
						data.Dcpm_sf[g][sfb], err = hcod_sf(adts.reader)
					}
					if err != nil {
						return data, err
					}
				}
			} else {
				data.Dcpm_sf[g][sfb] = 0
			}
		}
	}

	// else {} from the spec is not implemented due to hard coded aacSectionDataResilienceFlag
	return data, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.54 – Syntax of tns_data()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) tns_data(info *ics_info) *tns_data {
	data := &tns_data{}

	filt_bits := uint(2)
	len_bits := uint(6)
	order_bits := uint(5)
	if info.Window_sequence == EIGHT_SHORT_SEQUENCE {
		filt_bits = 1
		len_bits = 4
		order_bits = 3
	}

	data.N_filt = make([]uint8, info.num_windows)
	data.Len = make([][]uint8, info.num_windows)
	data.Order = make([][]uint8, info.num_windows)
	data.Direction = make([][]bool, info.num_windows)
	data.Coef_res = make([]uint8, info.num_windows)
	data.Coef_compress = make([][]uint8, info.num_windows)
	data.Coef = make([][][]uint8, info.num_windows)
	for w := range data.N_filt {
		data.N_filt[w], _ = adts.reader.ReadBitsAsUInt8(filt_bits)
		if data.N_filt[w] != 0 {
			data.Coef_res[w], _ = adts.reader.ReadBitsAsUInt8(1)
		}

		data.Len[w] = make([]uint8, data.N_filt[w])
		data.Order[w] = make([]uint8, data.N_filt[w])
		data.Direction[w] = make([]bool, data.N_filt[w])
		data.Coef_compress[w] = make([]uint8, data.N_filt[w])
		data.Coef[w] = make([][]uint8, data.N_filt[w])
		for filt := range data.Len[w] {
			data.Len[w][filt], _ = adts.reader.ReadBitsAsUInt8(len_bits)
			data.Order[w][filt], _ = adts.reader.ReadBitsAsUInt8(order_bits)
			if data.Order[w][filt] != 0 {
				data.Direction[w][filt], _ = adts.reader.ReadBitAsBool()
				data.Coef_compress[w][filt], _ = adts.reader.ReadBitsAsUInt8(1)

				coef_bits := data.Coef_res[w] + 3 - data.Coef_compress[w][filt]
				data.Coef[w][filt] = make([]uint8, data.Order[w][filt])
				for i := range data.Coef[w][filt] {
					data.Coef[w][filt][i], _ = adts.reader.ReadBitsAsUInt8(uint(coef_bits))
				}
			}
		}
	}
	return data
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.55 – Syntax of ltp_data()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) ltp_data(info *ics_info) (*ltp_data, error) {
	data := &ltp_data{}
	if adts.Profile == AUDIO_OBJECT_TYPE_ER_AAC_LD {
		ltp_lag_update, _ := adts.reader.ReadBitAsBool()
		if ltp_lag_update {
			data.Ltp_lag, _ = adts.reader.ReadBitsAsUInt(10)
		} else {
			data.Ltp_lag = info.Ltp_data.Ltp_lag
		}

		if data.Ltp_lag > uint(adts.Frame_length<<1) {
			return data, fmt.Errorf("Error: Ltp_lag (%d) out of range (%d)", data.Ltp_lag, adts.Frame_length<<1)
		}

		data.Ltp_coef, _ = adts.reader.ReadBitsAsUInt8(3)
		data.Ltp_long_used = make([]bool, minInt(int(info.Max_sfb), int(MAX_LTP_LONG_SFB)))
		for sfb := range data.Ltp_long_used {
			data.Ltp_long_used[sfb], _ = adts.reader.ReadBitAsBool()
		}
	} else {
		data.Ltp_lag, _ = adts.reader.ReadBitsAsUInt(11)
		if data.Ltp_lag > uint(adts.Frame_length<<1) {
			return data, fmt.Errorf("Error: Ltp_lag (%d) out of range (%d)", data.Ltp_lag, adts.Frame_length<<1)
		}

		data.Ltp_coef, _ = adts.reader.ReadBitsAsUInt8(3)
		if info.Window_sequence != EIGHT_SHORT_SEQUENCE {
			data.Ltp_long_used = make([]bool, minInt(int(info.Max_sfb), int(MAX_LTP_LONG_SFB)))
			for sfb := range data.Ltp_long_used {
				data.Ltp_long_used[sfb], _ = adts.reader.ReadBitAsBool()
			}
		}
	}

	return data, nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.56 – Syntax of spectral_data()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) spectral_data(info *ics_info, sec_data *section_data) (*spectral_data, error) {
	data := &spectral_data{}

	for g := uint8(0); g < info.num_window_groups; g++ {
		for i := uint8(0); i < sec_data.num_sec[g]; i++ {
			switch sec_data.Sect_cb[g][i] {
			case ZERO_HCB, NOISE_HCB, INTENSITY_HCB, INTENSITY_HCB2:
			default:
				inc := uint16(4)
				if sec_data.Sect_cb[g][i] >= FIRST_PAIR_HCB {
					inc = 2
				}

				start := info.sect_sfb_offset[g][sec_data.sect_start[g][i]]
				end := info.sect_sfb_offset[g][sec_data.sect_end[g][i]]
				for k := start; k < end; k += inc {
					if sec_data.Sect_cb[g][i] != 0 {
						if adts.reader.HasBitLeft() == false {
							return data, fmt.Errorf("Error: Spectral Data parsing ran out of bits")
						}
						val, err := hcod(adts.reader, sec_data.Sect_cb[g][i])
						if err != nil {
							return data, err
						}

						data.Hcod = append(data.Hcod, val)
					}
				}
			}
		}
	}

	return data, nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.57 – Syntax of extension_payload()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) extension_payload(cnt int, id_adts uint8) (int, *extension_payload, error) {
	var err error
	data := &extension_payload{}
	data.Extension_type, _ = adts.reader.ReadBitsAsUInt8(4) // extension_type

	switch data.Extension_type {
	case EXT_DYNAMIC_RANGE:
		cnt, data.Dynamic_range_info = adts.dynamic_range_info()
		return cnt, data, err
	case EXT_SAC_DATA:
		cnt, data.Sac_extension_data = adts.sac_extension_data(cnt)
		return cnt, data, err
	case EXT_SBR_DATA:
		cnt, data.Sbr_extension_data, err = adts.sbr_extension_data(cnt, id_adts, false)
		return cnt, data, err
	case EXT_SBR_DATA_CRC:
		cnt, data.Sbr_extension_data, err = adts.sbr_extension_data(cnt, id_adts, true)
		return cnt, data, err
	case EXT_FILL_DATA:
		data.Fill_nibble, _ = adts.reader.ReadBitsAsUInt8(4) // fill_nibble - must be ‘0000’
		if int(adts.reader.BitsLeft()) > cnt {
			data.Fill_byte, _ = adts.reader.ReadBitsToByteArray(uint(8 * (cnt - 1)))
		}
	case EXT_DATA_ELEMENT:
		data.Data_element_version, _ = adts.reader.ReadBitsAsUInt8(4) // data_element_version
		switch data.Data_element_version {
		case ANC_DATA:
			dataElementLength := uint(0)
			for {
				data.DataElementLengthPart, _ = adts.reader.ReadBitsAsUInt8(8) // dataElementLengthPart
				dataElementLength += uint(data.DataElementLengthPart)
				if data.DataElementLengthPart != 255 {
					break
				}
			}
			data.Data_element_byte, _ = adts.reader.ReadBytes(dataElementLength)
		}
	case EXT_FILL:
		fallthrough
	default:
		adts.reader.SkipBits(uint(8*(cnt-1) + 4))
	}

	return cnt, data, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.58 – Syntax of dynamic_range_info()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) dynamic_range_info() (int, *dynamic_range_info) {
	info := &dynamic_range_info{}

	n := 1
	drc_num_bands := uint8(1)
	info.Pce_tag_present, _ = adts.reader.ReadBitAsBool() // pce_tag_present
	if info.Pce_tag_present {
		info.Pce_instance_tag, _ = adts.reader.ReadBitsAsUInt8(4)     // pce_instance_tag
		info.Drc_tag_reserve_bits, _ = adts.reader.ReadBitsAsUInt8(4) // drc_tag_reserved_bits
	}

	info.Excluded_chns_present, _ = adts.reader.ReadBitAsBool() // excluded_chns_present
	if info.Excluded_chns_present {
		exclude := 0
		exclude, info.Excluded_chns = adts.excluded_channels()
		n += exclude
	}

	info.Drc_bands_present, _ = adts.reader.ReadBitAsBool() // drc_bands_present
	if info.Drc_bands_present {
		info.Drc_band_incr, _ = adts.reader.ReadBitsAsUInt8(4)            // drc_band_incr
		info.Drc_interpolation_scheme, _ = adts.reader.ReadBitsAsUInt8(4) // drc_interpolation_scheme

		n++
		drc_num_bands += info.Drc_band_incr
		info.Drc_band_top, _ = adts.reader.ReadBytes(uint(drc_num_bands))
	}

	info.Prog_ref_level_present, _ = adts.reader.ReadBitAsBool() // prog_ref_level_present
	if info.Prog_ref_level_present {
		info.Prog_ref_level, _ = adts.reader.ReadBitsAsUInt8(7)      // prog_ref_level
		info.Prog_ref_level_reserved_bits, _ = adts.reader.ReadBit() // prog_ref_level_reserved_bitsBits(1)
		n++
	}

	info.Dyn_range_sign = make([]uint8, drc_num_bands)
	info.Dyn_range_cnt = make([]uint8, drc_num_bands)
	for i := range info.Dyn_range_sign {
		info.Dyn_range_sign[i], _ = adts.reader.ReadBitsAsUInt8(1) // dyn_rng_sgn[i]
		info.Dyn_range_cnt[i], _ = adts.reader.ReadBitsAsUInt8(7)  // dyn_rng_ctl[i]
		n++
	}

	return n, info
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.59 – Syntax of excluded_channels()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) excluded_channels() (int, *excluded_channels) {
	data := &excluded_channels{}
	n := 0
	num_excl_chan := 7
	data.Exclude_mask = make([]bool, 7)
	for i := range data.Exclude_mask {
		data.Exclude_mask[i], _ = adts.reader.ReadBitAsBool()
	}

	n++

	data.Additional_excluded_chns = make([]bool, 0)
	additional_excluded_chn, _ := adts.reader.ReadBitAsBool()
	data.Additional_excluded_chns = append(data.Additional_excluded_chns, additional_excluded_chn)
	for data.Additional_excluded_chns[n-1] {
		for i := num_excl_chan; i < num_excl_chan+7; i++ {
			mask, _ := adts.reader.ReadBitAsBool()
			data.Exclude_mask = append(data.Exclude_mask, mask)
		}
		n++
		num_excl_chan += 7

		additional_excluded_chn, _ := adts.reader.ReadBitAsBool()
		data.Additional_excluded_chns = append(data.Additional_excluded_chns, additional_excluded_chn)
	}

	return n, data
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.60 – Syntax of ms_data()
////////////////////////////////////////////////////////////////////////////////
// func (adts *ADTS) ms_data() {
//	for g := 0; g < num_window_groups; g++ {
//		for sfb := last_max_sfb_ms; sfb < max_sfb; sfb++ {
//			reader.SkipBits(1) // ms_used[g][sfb]
//		}
//	}
// }

////////////////////////////////////////////////////////////////////////////////
// Table 4.61 – Syntax of sac_extension_data()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sac_extension_data(cnt int) (int, *sac_extension_data) {
	data := &sac_extension_data{}

	data.AncType, _ = adts.reader.ReadBitsAsUInt8(2)                  // ancType
	data.AncStart, _ = adts.reader.ReadBitAsBool()                    // ancStart
	data.AncStop, _ = adts.reader.ReadBitAsBool()                     // ancStop
	data.AncDataSegmentByte, _ = adts.reader.ReadBytes(uint(cnt - 1)) // ancDataSegmentByte[i]
	return cnt, data
}

//////////////////////////////////////////////////////////////////////////////////
// Table 4.62 – Syntax of sbr_extension_data()
//////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_extension_data(cnt int, id_aac uint8, crc_flag bool) (int, *sbr_extension_data, error) {
	var err error
	data := &sbr_extension_data{}
	num_sbr_bits := uint(0)

	if crc_flag {
		data.Bs_sbr_crc_bits, _ = adts.reader.ReadBitsAsUInt16(10) // bs_sbr_crc_bits
		num_sbr_bits += 10
	}

	// Both FFMPEG and FAAD2 assume sbr_layer != SBR_STEREO_ENHANCE, we'll do the same
	// if sbr_layer != SBR_STEREO_ENHANCE
	num_sbr_bits++
	var data_bits uint
	if data.Bs_header_flag, _ = adts.reader.ReadBitAsBool(); data.Bs_header_flag {
		data_bits, data.Sbr_header = adts.sbr_header()
		num_sbr_bits += data_bits
	}

	if data.Sbr_header != nil {
		// Sampling freq for the sbr is twice the stream sample rate (4.6.18.2.5).  This means the sbr
		// sampling freq index moves down by 3, except for 0, 1, 2 and 12
		sfi := int(adts.sfi) - 3
		if sfi < 0 {
			sfi = 0
		} else if sfi > 8 {
			sfi = 8
		}

		err = derive_sbr_tables(data, uint8(sfi), data.Sbr_header.Bs_start_freq, data.Sbr_header.Bs_stop_freq,
			data.Sbr_header.Bs_freq_scale, data.Sbr_header.Bs_alter_scale, data.Sbr_header.Bs_xover_band)
		if err != nil {
			return 0, data, err
		}

		data_bits, data.Sbr_data, err = adts.sbr_data(data, id_aac, data.Sbr_header.Bs_amp_res)
		num_sbr_bits += data_bits
	}

	num_align_bits := (8*uint(cnt) - 4 - num_sbr_bits)

	if 8*uint(cnt) < (4 + num_sbr_bits) {
		return 0, data, fmt.Errorf("sbr extension payload malformed")
	}

	data.Bs_fill_bits, _ = adts.reader.ReadBitsToByteArray(num_align_bits)

	return int(num_sbr_bits+num_align_bits+4) / 8, data, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.63 – Syntax of sbr_header()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_header() (uint, *sbr_header) {
	data := &sbr_header{}
	start_bits := adts.reader.BitsLeft()

	data.Bs_amp_res, _ = adts.reader.ReadBitAsBool()
	data.Bs_start_freq, _ = adts.reader.ReadBitsAsUInt8(4)
	data.Bs_stop_freq, _ = adts.reader.ReadBitsAsUInt8(4)
	data.Bs_xover_band, _ = adts.reader.ReadBitsAsUInt8(3)
	data.Bs_reserved, _ = adts.reader.ReadBitsAsUInt8(2)

	data.Bs_header_extra_1, _ = adts.reader.ReadBitAsBool()
	data.Bs_header_extra_2, _ = adts.reader.ReadBitAsBool()
	if data.Bs_header_extra_1 {
		data.Bs_freq_scale, _ = adts.reader.ReadBitsAsUInt8(2)
		data.Bs_alter_scale, _ = adts.reader.ReadBitsAsUInt8(1)
		data.Bs_noise_bands, _ = adts.reader.ReadBitsAsUInt8(2)
	} else {
		// Defaults defined in 4.2.8.1: 4.105-4.107
		data.Bs_freq_scale = 2
		data.Bs_alter_scale = 1
		data.Bs_noise_bands = 2
	}

	if data.Bs_header_extra_2 {
		data.Bs_limiter_bands, _ = adts.reader.ReadBitsAsUInt8(2)
		data.Bs_limiter_gains, _ = adts.reader.ReadBitsAsUInt8(2)
		data.Bs_interpol_freq, _ = adts.reader.ReadBitsAsUInt8(1)
		data.Bs_smoothing_mode, _ = adts.reader.ReadBitsAsUInt8(1)
	} else {
		// Defaults defined in 4.2.8.1: 4.108-4.111
		data.Bs_limiter_bands = 2
		data.Bs_limiter_gains = 2
		data.Bs_interpol_freq = 1
		data.Bs_smoothing_mode = 1
	}

	return (start_bits - adts.reader.BitsLeft()), data
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.64 – Syntax of sbr_data()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_data(ext_data *sbr_extension_data, id_aac uint8, bs_amp_res bool) (uint, *sbr_data, error) {
	var err error
	data := &sbr_data{}
	data_bits := adts.reader.BitsLeft()
	// FFMPEG and FAAD2 both blindly assume sbr_layer = SBR_NOT_SCALABLE, I guess we will too
	// switch sbr_layer
	//		case SBR_NOT_SCALABLE:
	switch id_aac {
	case ID_SCE:
		data.Sbr_single_channel_element, err = adts.sbr_single_channel_element(ext_data, bs_amp_res)
	case ID_CPE:
		data.Sbr_channel_pair_element, err = adts.sbr_channel_pair_element(ext_data, bs_amp_res)
	}
	return (data_bits - adts.reader.BitsLeft()), data, err
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.65 – Syntax of sbr_single_channel_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_single_channel_element(ext_data *sbr_extension_data, bs_amp_res bool) (*sbr_single_channel_element, error) {
	e := &sbr_single_channel_element{}

	if e.Bs_data_extra, _ = adts.reader.ReadBitAsBool(); e.Bs_data_extra { // bs_data_extra
		e.Bs_reserved, _ = adts.reader.ReadBitsAsUInt8(4) // bs_reserved
	}

	e.Sbr_grid = &sbr_grid{
		Bs_var_bord_0: make([]uint8, 2),
		Bs_var_bord_1: make([]uint8, 2),
		Bs_num_rel_0:  make([]uint8, 2),
		Bs_num_rel_1:  make([]uint8, 2),
		bs_rel_bord_0: make([][]uint8, 2),
		bs_rel_bord_1: make([][]uint8, 2),
	}
	e.Sbr_dtdf = &sbr_dtdf{}
	e.Sbr_invf = &sbr_invf{}
	e.Sbr_envelope = &sbr_envelope{}
	e.Sbr_noise = &sbr_noise{}
	if err := adts.sbr_grid(0, e.Sbr_grid, ext_data.Sbr_header); err != nil {
		return e, err
	}

	adts.sbr_dtdf(0, e.Sbr_dtdf, e.Sbr_grid)
	adts.sbr_invf(0, e.Sbr_invf, ext_data)
	adts.sbr_envelope(0, false, bs_amp_res, e.Sbr_envelope, ext_data, e.Sbr_grid, e.Sbr_dtdf)
	adts.sbr_noise(0, false, e.Sbr_noise, ext_data, e.Sbr_grid, e.Sbr_dtdf)
	if e.Bs_add_harmonic_flag, _ = adts.reader.ReadBitAsBool(); e.Bs_add_harmonic_flag { // bs_add_harmonic_flag
		e.Sbr_sinusoidal_coding = &sbr_sinusoidal_coding{}
		adts.sbr_sinusoidal_coding(0, e.Sbr_sinusoidal_coding, ext_data)
	}

	if e.Bs_extended_data, _ = adts.reader.ReadBitAsBool(); e.Bs_extended_data { // bs_extended_data
		e.Bs_extension_size, _ = adts.reader.ReadBitsAsUInt8(4)
		var cnt uint
		if cnt = uint(e.Bs_extension_size); cnt == 15 {
			e.Bs_esc_count, _ = adts.reader.ReadBitsAsUInt8(8)
			cnt += uint(e.Bs_esc_count)
		}

		num_bits_left := cnt * 8
		if num_bits_left > MaxBitsLeft {
			return e, fmt.Errorf("Too many bits left, check bitstream continuity")
		}

		e.Bs_extension_id = make([]uint8, 0)
		e.Sbr_extension = make([]*sbr_extension, 0)
		for i := 0; num_bits_left > 7; i++ {
			ext_id, _ := adts.reader.ReadBitsAsUInt8(2)
			e.Bs_extension_id = append(e.Bs_extension_id, ext_id)
			if e.Bs_extension_id[i] == EXTENSION_ID_PS {
				num_bits_left -= 2

				bits_read, ext, err := adts.sbr_extension(e.Bs_extension_id[i], num_bits_left)
				if err != nil {
					return e, err
				}
				e.Sbr_extension = append(e.Sbr_extension, ext)

				num_bits_left -= bits_read
				if num_bits_left < 0 {
					return e, fmt.Errorf("Error: SBR parsing overran available bits")
				}
			}
		}

		e.Bs_fill_bits, _ = adts.reader.ReadBitsToByteArray(num_bits_left)
	}

	return e, nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.66 – Syntax of sbr_channel_pair_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_channel_pair_element(ext_data *sbr_extension_data, bs_amp_res bool) (*sbr_channel_pair_element, error) {
	e := &sbr_channel_pair_element{}

	if e.Bs_data_extra, _ = adts.reader.ReadBitAsBool(); e.Bs_data_extra { // bs_data_extra
		e.Bs_reserved_0, _ = adts.reader.ReadBitsAsUInt8(4) // bs_reserved
		e.Bs_reserved_1, _ = adts.reader.ReadBitsAsUInt8(4) // bs_reserved
	}

	e.Sbr_grid = &sbr_grid{
		Bs_var_bord_0: make([]uint8, 2),
		Bs_var_bord_1: make([]uint8, 2),
		Bs_num_rel_0:  make([]uint8, 2),
		Bs_num_rel_1:  make([]uint8, 2),
		bs_rel_bord_0: make([][]uint8, 2),
		bs_rel_bord_1: make([][]uint8, 2),
	}
	e.Sbr_dtdf = &sbr_dtdf{}
	e.Sbr_invf = &sbr_invf{}
	e.Sbr_envelope = &sbr_envelope{}
	e.Sbr_noise = &sbr_noise{}
	if e.Bs_coupling, _ = adts.reader.ReadBitAsBool(); e.Bs_coupling { // bs_coupling
		if err := adts.sbr_grid(0, e.Sbr_grid, ext_data.Sbr_header); err != nil {
			return e, err
		}
		// The spec isn't clear about this, but because this is a coupled channel we
		// need to copy the grid and inverting filter from ch0 to ch1
		grid_copy(e.Sbr_grid)

		adts.sbr_dtdf(0, e.Sbr_dtdf, e.Sbr_grid)
		adts.sbr_dtdf(1, e.Sbr_dtdf, e.Sbr_grid)
		adts.sbr_invf(0, e.Sbr_invf, ext_data)

		adts.sbr_envelope(0, true, bs_amp_res, e.Sbr_envelope, ext_data, e.Sbr_grid, e.Sbr_dtdf)
		adts.sbr_noise(0, true, e.Sbr_noise, ext_data, e.Sbr_grid, e.Sbr_dtdf)
		adts.sbr_envelope(1, true, bs_amp_res, e.Sbr_envelope, ext_data, e.Sbr_grid, e.Sbr_dtdf)
		adts.sbr_noise(1, true, e.Sbr_noise, ext_data, e.Sbr_grid, e.Sbr_dtdf)

	} else {
		if err := adts.sbr_grid(0, e.Sbr_grid, ext_data.Sbr_header); err != nil {
			return e, err
		}
		if err := adts.sbr_grid(1, e.Sbr_grid, ext_data.Sbr_header); err != nil {
			return e, err
		}
		adts.sbr_dtdf(0, e.Sbr_dtdf, e.Sbr_grid)
		adts.sbr_dtdf(1, e.Sbr_dtdf, e.Sbr_grid)
		adts.sbr_invf(0, e.Sbr_invf, ext_data)
		adts.sbr_invf(1, e.Sbr_invf, ext_data)

		adts.sbr_envelope(0, false, bs_amp_res, e.Sbr_envelope, ext_data, e.Sbr_grid, e.Sbr_dtdf)
		adts.sbr_envelope(1, false, bs_amp_res, e.Sbr_envelope, ext_data, e.Sbr_grid, e.Sbr_dtdf)
		adts.sbr_noise(0, false, e.Sbr_noise, ext_data, e.Sbr_grid, e.Sbr_dtdf)
		adts.sbr_noise(1, false, e.Sbr_noise, ext_data, e.Sbr_grid, e.Sbr_dtdf)
	}

	flag, _ := adts.reader.ReadBitAsBool()
	e.Bs_add_harmonic_flag = append(e.Bs_add_harmonic_flag, flag)
	e.Sbr_sinusoidal_coding = &sbr_sinusoidal_coding{}
	if e.Bs_add_harmonic_flag[0] { // bs_add_harmonic_flag
		adts.sbr_sinusoidal_coding(0, e.Sbr_sinusoidal_coding, ext_data)
	} else {
		// We need to fill in the 0 channel in the struct to it indexes correctly
		e.Sbr_sinusoidal_coding.Bs_add_harmonic = make([][]bool, 1)
	}

	flag, _ = adts.reader.ReadBitAsBool()
	e.Bs_add_harmonic_flag = append(e.Bs_add_harmonic_flag, flag)
	if e.Bs_add_harmonic_flag[1] { // bs_add_harmonic_flag
		adts.sbr_sinusoidal_coding(1, e.Sbr_sinusoidal_coding, ext_data)
	}

	if e.Bs_extended_data, _ = adts.reader.ReadBitAsBool(); e.Bs_extended_data { // bs_extended_data
		e.Bs_extension_size, _ = adts.reader.ReadBitsAsUInt8(4)
		var cnt uint
		if cnt = uint(e.Bs_extension_size); cnt == 15 {
			e.Bs_esc_count, _ = adts.reader.ReadBitsAsUInt8(8)
			cnt += uint(e.Bs_esc_count)
		}

		// TODO: could we be a bit more graceful about this block?
		num_bits_left := cnt * 8

		if num_bits_left > MaxBitsLeft {
			return e, fmt.Errorf("Too many bits left, check bitstream continuity")
		}

		// Extentions are currently unsupported
		adts.reader.SkipBits(num_bits_left)
	}
	return e, nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.67 – Syntax of sbr_channel_pair_base_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_channel_pair_base_element(bs_amp_res bool, ext_data *sbr_extension_data) (*sbr_channel_pair_base_element, error) {
	e := &sbr_channel_pair_base_element{}

	if e.Bs_data_extra, _ = adts.reader.ReadBitAsBool(); e.Bs_data_extra { // bs_data_extra
		e.Bs_reserved_0, _ = adts.reader.ReadBitsAsUInt8(4) // bs_reserved
		e.Bs_reserved_1, _ = adts.reader.ReadBitsAsUInt8(4) // bs_reserved
	}

	e.Bs_coupling, _ = adts.reader.ReadBitAsBool()

	e.Sbr_grid = &sbr_grid{
		Bs_var_bord_0: make([]uint8, 2),
		Bs_var_bord_1: make([]uint8, 2),
		Bs_num_rel_0:  make([]uint8, 2),
		Bs_num_rel_1:  make([]uint8, 2),
		bs_rel_bord_0: make([][]uint8, 2),
		bs_rel_bord_1: make([][]uint8, 2),
	}
	e.Sbr_dtdf = &sbr_dtdf{}
	e.Sbr_invf = &sbr_invf{}
	e.Sbr_envelope = &sbr_envelope{}
	e.Sbr_noise = &sbr_noise{}
	if err := adts.sbr_grid(0, e.Sbr_grid, ext_data.Sbr_header); err != nil {
		return e, err
	}

	adts.sbr_dtdf(0, e.Sbr_dtdf, e.Sbr_grid)
	adts.sbr_invf(0, e.Sbr_invf, ext_data)
	adts.sbr_envelope(0, true, bs_amp_res, e.Sbr_envelope, ext_data, e.Sbr_grid, e.Sbr_dtdf)
	adts.sbr_noise(0, true, e.Sbr_noise, ext_data, e.Sbr_grid, e.Sbr_dtdf)

	if e.Bs_add_harmonic_flag, _ = adts.reader.ReadBitAsBool(); e.Bs_add_harmonic_flag { // bs_add_harmonic_flag
		e.Sbr_sinusoidal_coding = &sbr_sinusoidal_coding{}
		adts.sbr_sinusoidal_coding(0, e.Sbr_sinusoidal_coding, ext_data)
	}

	if e.Bs_extended_data, _ = adts.reader.ReadBitAsBool(); e.Bs_extended_data { // bs_extended_data
		e.Bs_extension_size, _ = adts.reader.ReadBitsAsUInt8(4)
		var cnt uint
		if cnt = uint(e.Bs_extension_size); cnt == 15 {
			e.Bs_esc_count, _ = adts.reader.ReadBitsAsUInt8(8)
			cnt += uint(e.Bs_esc_count)
		}

		// TODO: could we be a bit more graceful about this block?
		num_bits_left := cnt * 8
		if num_bits_left > MaxBitsLeft {
			return e, fmt.Errorf("Too many bits left, check bitstream continuity")
		}

		e.Bs_extension_id = make([]uint8, 0)
		e.Sbr_extension = make([]*sbr_extension, 0)
		for i := 0; num_bits_left > 7; i++ {
			ext_id, _ := adts.reader.ReadBitsAsUInt8(2)
			e.Bs_extension_id = append(e.Bs_extension_id, ext_id)
			if e.Bs_extension_id[i] == EXTENSION_ID_PS {
				num_bits_left -= 2

				bits_read, ext, err := adts.sbr_extension(e.Bs_extension_id[i], num_bits_left)
				if err != nil {
					return e, err
				}
				e.Sbr_extension = append(e.Sbr_extension, ext)

				num_bits_left -= bits_read
				if num_bits_left < 0 {
					return e, fmt.Errorf("Error: SBR parsing overran available bits")
				}
			}
		}

		e.Bs_fill_bits, _ = adts.reader.ReadBitsToByteArray(num_bits_left)
	}

	return e, nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.68 – Syntax of sbr_channel_pair_enhance_element()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_channel_pair_enhance_element(bs_amp_res bool) *sbr_channel_pair_enhance_element {
	// This is not implemented in either FFMPEG or FAAD2.  As is it would lack the grid element for the
	// sbr_dtdf decoding and the spec doesn't specify default values.  We'll leave it commented out
	//
	/*e := &sbr_channel_pair_enhance_element{}

	e.Sbr_dtdf = adts.sbr_dtdf(1, ???)
	e.Sbr_envelope = adts.sbr_envelope(1, 1, bs_amp_res, nil, e.Sbr_dtdf)
	e.Sbr_noise = adts.sbr_noise(1, 1)

	if e.Bs_add_harmonic_flag, _ = adts.reader.ReadBitAsBool(); e.Bs_add_harmonic_flag {
		e.Sbr_sinusoidal_coding = adts.sbr_sinusoidal_coding(1)
	}*/

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.69 – Syntax of sbr_grid()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_grid(ch uint, data *sbr_grid, header *sbr_header) error {
	data.Bs_frame_class, _ = adts.reader.ReadBitsAsUInt8(2)
	switch data.Bs_frame_class {
	case FIXFIX:
		data.Tmp, _ = adts.reader.ReadBitsAsUInt8(2)
		data.bs_num_env = append(data.bs_num_env, 1<<data.Tmp)

		if data.bs_num_env[ch] == 1 {
			header.Bs_amp_res = false
		}

		data.Bs_freq_res = append(data.Bs_freq_res, make([]uint8, data.bs_num_env[ch]))
		data.Bs_freq_res[ch][0], _ = adts.reader.ReadBitsAsUInt8(1)
		for env := uint8(1); env < data.bs_num_env[ch]; env++ {
			data.Bs_freq_res[ch][env] = data.Bs_freq_res[ch][0]
		}

		// Initialize this to 0 for this channel to keep the chan index correct
		data.Bs_pointer = append(data.Bs_pointer, 0)
	case FIXVAR:
		tmp, _ := adts.reader.ReadBitsAsUInt8(2)
		data.Bs_var_bord_1[ch] = tmp
		tmp, _ = adts.reader.ReadBitsAsUInt8(2)
		data.Bs_num_rel_1[ch] = tmp
		data.bs_num_env = append(data.bs_num_env, data.Bs_num_rel_1[ch]+1)

		data.bs_rel_bord_1[ch] = make([]uint8, data.bs_num_env[ch]-1)
		for rel := range data.bs_rel_bord_1[ch] {
			data.Tmp, _ = adts.reader.ReadBitsAsUInt8(2)
			data.bs_rel_bord_1[ch][rel] = 2*data.Tmp + 1
		}

		ptr_bits := ceil_log2(data.bs_num_env[ch] + 1)
		ptr, _ := adts.reader.ReadBitsAsUInt(uint(ptr_bits))
		data.Bs_pointer = append(data.Bs_pointer, ptr)

		data.Bs_freq_res = append(data.Bs_freq_res, make([]uint8, data.bs_num_env[ch]))
		for env := range data.Bs_freq_res[ch] {
			data.Bs_freq_res[ch][data.bs_num_env[ch]-1-uint8(env)], _ = adts.reader.ReadBitsAsUInt8(1)
		}
	case VARFIX:
		tmp, _ := adts.reader.ReadBitsAsUInt8(2)
		data.Bs_var_bord_0[ch] = tmp
		tmp, _ = adts.reader.ReadBitsAsUInt8(2)
		data.Bs_num_rel_0[ch] = tmp
		data.bs_num_env = append(data.bs_num_env, data.Bs_num_rel_0[ch]+1)

		data.bs_rel_bord_0[ch] = make([]uint8, data.bs_num_env[ch]-1)
		for rel := range data.bs_rel_bord_0[ch] {
			data.Tmp, _ = adts.reader.ReadBitsAsUInt8(2)
			data.bs_rel_bord_0[ch][rel] = 2*data.Tmp + 2
		}

		ptr_bits := ceil_log2(data.bs_num_env[ch] + 1)
		ptr, _ := adts.reader.ReadBitsAsUInt(uint(ptr_bits))
		data.Bs_pointer = append(data.Bs_pointer, ptr)

		data.Bs_freq_res = append(data.Bs_freq_res, make([]uint8, data.bs_num_env[ch]))
		for env := range data.Bs_freq_res[ch] {
			data.Bs_freq_res[ch][env], _ = adts.reader.ReadBitsAsUInt8(1)
		}
	case VARVAR:
		tmp, _ := adts.reader.ReadBitsAsUInt8(2)
		data.Bs_var_bord_0[ch] = tmp
		tmp, _ = adts.reader.ReadBitsAsUInt8(2)
		data.Bs_var_bord_1[ch] = tmp
		tmp, _ = adts.reader.ReadBitsAsUInt8(2)
		data.Bs_num_rel_0[ch] = tmp
		tmp, _ = adts.reader.ReadBitsAsUInt8(2)
		data.Bs_num_rel_1[ch] = tmp

		data.bs_num_env = append(
			data.bs_num_env,
			uint8(minInt(5, int(data.Bs_num_rel_0[ch]+data.Bs_num_rel_1[ch]+1))),
		)

		// This is how the spec spells this section out.  We could make Bs_num_rel into
		// a 2d array and use one loop, but for the sake of keeping variable names consistent
		// with the spec we'll follow their two loop method
		data.bs_rel_bord_0[ch] = make([]uint8, data.Bs_num_rel_0[ch])
		for rel := range data.bs_rel_bord_0[ch] {
			data.Tmp, _ = adts.reader.ReadBitsAsUInt8(2)
			data.bs_rel_bord_0[ch][rel] = data.Tmp*2 + 2
		}

		data.bs_rel_bord_1[ch] = make([]uint8, data.Bs_num_rel_1[ch])
		for rel := range data.bs_rel_bord_1[ch] {
			data.Tmp, _ = adts.reader.ReadBitsAsUInt8(2)
			data.bs_rel_bord_1[ch][rel] = data.Tmp*2 + 2
		}

		ptr_bits := ceil_log2(data.bs_num_env[ch] + 1)
		ptr, _ := adts.reader.ReadBitsAsUInt(uint(ptr_bits))
		data.Bs_pointer = append(data.Bs_pointer, ptr)

		data.Bs_freq_res = append(data.Bs_freq_res, make([]uint8, data.bs_num_env[ch]))
		for env := range data.Bs_freq_res[ch] {
			data.Bs_freq_res[ch][env], _ = adts.reader.ReadBitsAsUInt8(1)
		}
	}

	if data.bs_num_env[ch] > 1 {
		data.bs_num_noise = append(data.bs_num_noise, 2)
	} else if data.bs_num_env[ch] == 1 {
		data.bs_num_noise = append(data.bs_num_noise, 1)
	} else {
		return fmt.Errorf("Error: bs_num_env[%d] (%d) is out of range", ch, data.bs_num_env[ch])
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.70 – Syntax of sbr_dtdf()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_dtdf(ch uint8, data *sbr_dtdf, grid *sbr_grid) {
	data.Bs_df_env = append(data.Bs_df_env, make([]bool, grid.bs_num_env[ch]))
	for env := range data.Bs_df_env[ch] {
		data.Bs_df_env[ch][env], _ = adts.reader.ReadBitAsBool() // bs_df_env[ch][env]
	}

	data.Bs_df_noise = append(data.Bs_df_noise, make([]bool, grid.bs_num_noise[ch]))
	for noise := range data.Bs_df_noise[ch] {
		data.Bs_df_noise[ch][noise], _ = adts.reader.ReadBitAsBool() // bs_df_noise[ch][noise]
	}
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.71 – Syntax of sbr_invf()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_invf(ch uint, data *sbr_invf, ext_data *sbr_extension_data) {
	data.Bs_invf_mode = append(data.Bs_invf_mode, make([]uint8, ext_data.N_Q))
	for n := range data.Bs_invf_mode[ch] {
		data.Bs_invf_mode[ch][n], _ = adts.reader.ReadBitsAsUInt8(2) // bs_invf_mode[ch][n]
	}
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.72 – Syntax of sbr_envelope()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_envelope(
	ch uint, bs_coupling bool, bs_amp_res bool,
	e *sbr_envelope, ext_data *sbr_extension_data, grid *sbr_grid, dtdf *sbr_dtdf,
) {
	var t_huff [][]int8
	var f_huff [][]int8

	amp_res := bs_amp_res
	if grid.bs_num_env[ch] == 1 && grid.Bs_frame_class == FIXFIX {
		amp_res = false
	}

	if bs_coupling && ch == 1 {
		if amp_res {
			t_huff = t_huffman_env_bal_3_0dB
			f_huff = f_huffman_env_bal_3_0dB
		} else {
			t_huff = t_huffman_env_bal_1_5dB
			f_huff = f_huffman_env_bal_1_5dB
		}
	} else {
		if amp_res {
			t_huff = t_huffman_env_3_0dB
			f_huff = f_huffman_env_3_0dB
		} else {
			t_huff = t_huffman_env_1_5dB
			f_huff = f_huffman_env_1_5dB
		}
	}

	e.Bs_data_env = append(e.Bs_data_env, make([][]int, grid.bs_num_env[ch]))
	for env := range e.Bs_data_env[ch] {
		if dtdf.Bs_df_env[ch][env] == false {
			num_bands := ext_data.n[grid.Bs_freq_res[ch][env]]
			e.Bs_data_env[ch][env] = make([]int, num_bands)
			if bs_coupling && ch == 1 {
				if amp_res {
					e.Bs_env_start_value_balance, _ = adts.reader.ReadBitsAsUInt8(5) // bs_env_start_value_balance
					e.Bs_data_env[ch][env][0] = int(e.Bs_env_start_value_balance)
				} else {
					e.Bs_env_start_value_balance, _ = adts.reader.ReadBitsAsUInt8(6) // bs_env_start_value_balance
					e.Bs_data_env[ch][env][0] = int(e.Bs_env_start_value_balance)
				}
			} else {
				if amp_res {
					e.Bs_env_start_value_level, _ = adts.reader.ReadBitsAsUInt8(6) // bs_env_start_value_level
					e.Bs_data_env[ch][env][0] = int(e.Bs_env_start_value_level)
				} else {
					e.Bs_env_start_value_level, _ = adts.reader.ReadBitsAsUInt8(7) // bs_env_start_value_level
					e.Bs_data_env[ch][env][0] = int(e.Bs_env_start_value_level)
				}
			}

			for band := uint8(1); band < num_bands; band++ {
				e.Bs_data_env[ch][env][band] = sbr_huff_dec(adts.reader, f_huff)
			}
		} else {
			num_bands := ext_data.n[grid.Bs_freq_res[ch][env]]
			e.Bs_data_env[ch][env] = make([]int, num_bands)
			for band := range e.Bs_data_env[ch][env] {
				e.Bs_data_env[ch][env][band] = sbr_huff_dec(adts.reader, t_huff)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.73 – Syntax of sbr_noise()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_noise(
	ch uint, bs_coupling bool,
	data *sbr_noise, ext_data *sbr_extension_data, grid *sbr_grid, dtdf *sbr_dtdf,
) {
	var t_huff [][]int8
	var f_huff [][]int8

	if bs_coupling && ch == 1 {
		t_huff = t_huffman_noise_bal_3_0dB
		f_huff = f_huffman_env_bal_3_0dB
	} else {
		t_huff = t_huffman_noise_3_0dB
		f_huff = f_huffman_env_3_0dB
	}

	data.Bs_data_noise = append(data.Bs_data_noise, make([][]int, grid.bs_num_noise[ch]))
	for noise := range data.Bs_data_noise[ch] {
		data.Bs_data_noise[ch][noise] = make([]int, ext_data.N_Q)
		if dtdf.Bs_df_noise[ch][noise] == false {
			if bs_coupling && ch == 1 {
				data.Bs_noise_start_value_balance, _ = adts.reader.ReadBitsAsUInt8(5)
				data.Bs_data_noise[ch][noise][0] = int(data.Bs_noise_start_value_balance)
			} else {
				data.Bs_noise_start_value_level, _ = adts.reader.ReadBitsAsUInt8(5)
				data.Bs_data_noise[ch][noise][0] = int(data.Bs_noise_start_value_level)
			}
			for band := uint8(1); band < ext_data.N_Q; band++ {
				data.Bs_data_noise[ch][noise][band] = sbr_huff_dec(adts.reader, f_huff)
			}
		} else {
			for band := range data.Bs_data_noise[ch][noise] {
				data.Bs_data_noise[ch][noise][band] = sbr_huff_dec(adts.reader, t_huff)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// Table 4.74 – Syntax of sbr_sinusoidal_coding()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_sinusoidal_coding(ch uint8, data *sbr_sinusoidal_coding, ext_data *sbr_extension_data) {
	data.Bs_add_harmonic = append(data.Bs_add_harmonic, make([]bool, ext_data.N_high))
	for n := range data.Bs_add_harmonic[ch] {
		data.Bs_add_harmonic[ch][n], _ = adts.reader.ReadBitAsBool()
	}
}

////////////////////////////////////////////////////////////////////////////////
// Table 8.A.1 – Syntax of sbr_extension()
////////////////////////////////////////////////////////////////////////////////
func (adts *ADTS) sbr_extension(bs_extension_id uint8, num_bits_left uint) (uint, *sbr_extension, error) {
	data := &sbr_extension{}
	switch bs_extension_id {
	case EXTENSION_ID_PS:
		// he-AAC v2 - currently unsupported
		// TODO
		// num_bits_left -= ps_data()
		// returning num_bits_left tells the caller that all bits have been read
		return num_bits_left, nil, fmt.Errorf("bs_extension_id of 2 (EXTENSION_ID_PS) unsupported")
	default:
		//data.Bs_fill_bits, _ = adts.reader.ReadBitsToByteArray(num_bits_left)
	}
	// returning num_bits_left tells the caller that all bits have been read
	return num_bits_left, data, nil
}

func is_intensity(info *ics_info, group, sfb uint8) int {
	if info.sfb_cb[group][sfb] == INTENSITY_HCB {
		return 1
	} else if info.sfb_cb[group][sfb] == INTENSITY_HCB2 {
		return -1
	}
	return 0
}

func is_noise(info *ics_info, group, sfb uint8) bool {
	if info.sfb_cb[group][sfb] == NOISE_HCB {
		return true
	}
	return false
}

// Implements `ceil(log( bs_num_ch_env[ch] + 1 ) / log(2))` for sbr_grid decoding
func ceil_log2(val uint8) uint8 {
	log2 := [...]uint8{0, 0, 1, 2, 2, 3, 3, 3, 3, 4}
	if val < 10 && val >= 0 {
		return log2[val]
	}

	return 0
}

// Copy function for coupled channels
func grid_copy(grid *sbr_grid) {
	grid.Bs_freq_res = append(grid.Bs_freq_res, grid.Bs_freq_res[0])
	grid.Bs_pointer = append(grid.Bs_pointer, grid.Bs_pointer[0])
	grid.bs_num_env = append(grid.bs_num_env, grid.bs_num_env[0])
	grid.bs_num_noise = append(grid.bs_num_noise, grid.bs_num_noise[0])

	// We initialize these in each function before we populate the initial grid, so they're
	// treated differently here as well
	grid.Bs_num_rel_0[1] = grid.Bs_num_rel_0[0]
	grid.Bs_num_rel_1[1] = grid.Bs_num_rel_1[0]
	grid.Bs_var_bord_0[1] = grid.Bs_var_bord_0[0]
	grid.Bs_var_bord_1[1] = grid.Bs_var_bord_1[0]
	grid.bs_rel_bord_0[1] = grid.bs_rel_bord_0[0]
	grid.bs_rel_bord_1[1] = grid.bs_rel_bord_1[0]
}

func (adts *ADTS) Debug() string {
	return fmt.Sprintf("MPEGVer: %d, Layer: %d, Profile: %s, SampFreq: %d, ChannelConfig: %s, Bitrate: %d",
		adts.MpegVersion,
		adts.Layer,
		AACProfileType[adts.Profile],
		adts.SamplingFrequency,
		ChannelConfiguration[adts.ChannelConfiguration],
		adts.Bitrate)
}
