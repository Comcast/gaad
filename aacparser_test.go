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
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestAacLcADTS(t *testing.T) {
	// ID                                       : 102 (0x66)
	// Menu ID                                  : 1 (0x1)
	// Format                                   : AAC
	// Format/Info                              : Advanced Audio Codec
	// Format version                           : Version 2
	// Format profile                           : LC
	// Muxing mode                              : ADTS
	// Codec ID                                 : 15
	// Duration                                 : 6s 570ms
	// Bit rate mode                            : Variable
	// Bit rate                                 : 96.0 Kbps
	// Channel(s)                               : 2 channels
	// Channel positions                        : Front: L R
	// Sampling rate                            : 48.0 KHz
	// Compression mode                         : Lossy
	// Stream size                              : 77.0 KiB (1%)
	// Language                                 : English
	buf, _ := hex.DecodeString("fff94c8021600c21098bfffff7fff21646d20c0a44441020211465018f4dfe6fbb9d46e1bd57a7287afdc88a54ab41c9657781d576e1a4969091168afbb4b99bd76bd5cebeeedeeb35e3364f21dab6a4ade022f200cf5e1988838045824d12d1244fcf9effc766dd7cecd36e01ea20ad7b0368552f274bdf84480892931ac5195b02f25116b63359639b42b4f7571433320b57a13cec7831ddfb559135c9652ae1a396b5be6ebb631162")
	adts, err := ParseADTS(buf)
	// buffer is truncated, make sure we error out
	if err == nil {
		t.Errorf("This buffer is truncated, err must not be nil")
	}
	if adts.MpegVersion != uint8(1) {
		t.Errorf("MpegVersion (%d) must be equal to 1", adts.MpegVersion)
	}
	if adts.ChannelConfiguration != uint8(2) {
		t.Errorf("ChannelConfiguration (%d) must be equal to 2", adts.ChannelConfiguration)
	}
	if adts.SamplingFrequency != uint32(48000) {
		t.Errorf("SamplingFrequency (%d) must be equal to 48000", adts.SamplingFrequency)
	}
}

func TestHeAacAdts(t *testing.T) {
	// ID                                       : 482 (0x1E2)
	// Menu ID                                  : 1 (0x1)
	// Format                                   : AAC
	// Format/Info                              : Advanced Audio Codec
	// Format version                           : Version 4
	// Format profile                           : HE-AAC / LC
	// Muxing mode                              : ADTS
	// Codec ID                                 : 15
	// Duration                                 : 42mn 32s
	// Bit rate mode                            : Variable
	// Bit rate                                 : 62.5 Kbps
	// Channel count                            : 2 channels
	// Channel positions                        : Front: L R
	// Sampling rate                            : 44.1 KHz / 22.05 KHz
	// Compression mode                         : Lossy
	// Delay relative to video                  : -93ms
	// Stream size                              : 19.0 MiB (20%)
	buf, _ := hex.DecodeString("fff15c802b0000210bcbfffffffffff61bd05489222100016a9cc4b4c3e13d1603272651054e3dc1b9f69c338cee6ec068d2105c2e6ffb9eeac00154c1bb82f4b7c24c60cd5d81d3fe87f48da3d877c762e3f15cd60fd4ef3fba5601c86726a2ff5ff6f3a08824dfbcce81ec12620f7710316cc0f6a65cb03f49fd020349058fcd7a8e35cee4f3393c04d28f5eecc8ff2612a416e820125881b224bcbf9e3667d22ec0f63ec6c003")
	adts, err := ParseADTS(buf)
	// buffer is truncated, make sure we error out
	if err == nil {
		t.Errorf("This buffer is truncated, err must not be nil")
	}
	if adts.MpegVersion != uint8(0) {
		t.Errorf("MpegVersion (%d) must be equal to 0", adts.MpegVersion)
	}
	if adts.ChannelConfiguration != uint8(2) {
		t.Errorf("ChannelConfiguration (%d) must be equal to 2", adts.ChannelConfiguration)
	}
}

func TestAacLcSBR(t *testing.T) {
	// Complete name                            : /mnt/jitp/columbus_test_assets/d4_HDCC0056300001984003_new_1850/1.ts
	// ID                                       : 257 (0x101)
	// Menu ID                                  : 1 (0x1)
	// Format                                   : AAC
	// Format/Info                              : Advanced Audio Codec
	// Format version                           : Version 4
	// Format profile                           : LC
	// Muxing mode                              : ADTS
	// Codec ID                                 : 15
	// Duration                                 : 25mn 45s
	// Bit rate mode                            : Variable
	// Bit rate                                 : 92.5 Kbps
	// Channel(s)                               : 2 channels
	// Channel positions                        : Front: L R
	// Sampling rate                            : 44.1 KHz
	// Compression mode                         : Lossy
	// Delay relative to video                  : -23ms
	// Stream size                              : 17.0 MiB (4%)
	// Language                                 : English\
	buf, _ := hex.DecodeString("fff1508022c270210a11e5a43e00145aec00340cc3fcc84004bfd39ff91127fb0dff1d00010fe333fcd2097edbbf44847fb907f2c721f84cfece44ff787f4bc8ff581fe27a4bf938ff23827fe0ecfc1e93fc417e8d08fcccf4a90fe255f66087ebbbf3c2001f4ad76001a0661fe642df2c214b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4bc0")
	adts, err := ParseADTS(buf)
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}

	if adts.MpegVersion != uint8(0) {
		t.Errorf("MpegVersion (%d) must be equal to 0", adts.MpegVersion)
	}
	if adts.ChannelConfiguration != uint8(2) {
		t.Errorf("ChannelConfiguration (%d) must be equal to 2", adts.ChannelConfiguration)
	}
	if adts.SamplingFrequency != uint32(44100) {
		t.Errorf("SamplingFrequency (%d) must be equal to 44100", adts.SamplingFrequency)
	}

	if adts.Fill_elements == nil {
		t.Errorf("Fill_elements cannot be nil")
	}
	for _, e := range adts.Fill_elements {
		if e.Extension_payload.Sbr_extension_data != nil {
			t.Errorf("Sbr_extension_data is not nil, AAC-LC annot have an SBR")
		}
	}
}

// TODO: Bring these back once gots gets released
/*func TestSyncWord(t *testing.T) {

	// Packet was captured from Cisco Transcoder (AMC), the ADTS is prefixed with
	// a lot of 0's then 380, before the syncword FFF.  In this case, we need the
	// aacparser to ignore those until the syncword is found
	byteArray, _ := hex.DecodeString(
		"47406834071059e9594f7ed3000001dd021a80800525cf4bbdf7000000000000" +
			"00000000000000000000000000000000000000000000000000000380fff95880" +
			"400000217bd4bd1ac3484440d050a60c058822009d70fa7e9488fd7f6087238a" +
			"92a969489140ee8ae4c49beb7ba25e6f1ee3b57edb1b5daca09561670ca8f426" +
			"a92420c86c6ad53dedde7dfad327686da0ba94e8d26682be54d385715fd152aa" +
			"9c508e238ae7b98b95dd633430add6b2e154dff95d6a17392e17cded")

	pay, _ := packet.Payload(byteArray)
	pes, _ := pes.NewPESHeader(pay)
	adts, err := ParseADTS(pes.Data())
	if err == nil {
		t.Errorf("Truncated buffer, err must not be nil")
	}
	if adts.SamplingFrequency != uint32(24000) {
		t.Errorf("SamplingFrequency (%d) must be 24000", adts.SamplingFrequency)
	}
	if adts.Bitrate != uint32(94208) {
		t.Errorf("Bitrate (%d) must be 94208", adts.Bitrate)
	}
}

func TestAdtsMpegVersion(t *testing.T) {

	// Packet was captured from Cisco Transcoder (AMC), the ADTS.mpegVersion is 0
	// which means it's MPEG4 -- unsupported.  We only want MPEG2.
	byteArray, _ := hex.DecodeString(
		"4740683e071063a530687e8f000001dd021a808005271d2b1ff708498ee76a8a" +
			"e9ad1c46fa54d5e9bf7ad23560770bf67be89ba898e768705da8f670ead7048a" +
			"e684c06f05d01084d31a094cfff7001c614cfff7201c3800c56001bcb8000000" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000380fff958804000" +
			"00211bd4fdad87458332d06451196c0d4f795a522ea92b5552c44210")

	pay, _ := packet.Payload(byteArray)
	pes, _ := pes.NewPESHeader(pay)
	adts, err := ParseADTS(pes.Data())
	if err == nil {
		t.Errorf("Truncated buffer, err must not be nil")
	}
	if adts.MpegVersion != uint8(0) {
		t.Errorf("MpegVersion (%d) must be 0", adts.MpegVersion)
	}
}*/

// AAC Audio frame has an SBR extension, but no SBR header.  We Should successfully
// detect an SBR, but no real SBR data can exist
func TestSbrExists(t *testing.T) {
	buf, err := hex.DecodeString(
		"fff95880400000211bce77fffffffff97b1a92c682b0b9a8ca52128c5441" +
			"109bf194d09c8133a6288b644104aba0878d369fe51d8d338f2f55034e0f" +
			"23f27a9b868841f7d690aa2f1f7b9c4416e27cad2232ca5260df17901808" +
			"862617eff1e518306083b9241418390001136c5001d38bd5fa7a5502b600" +
			"d0216011a881837c0793e71afdf563d008b77ca02fbf3bd462d79e8d5908" +
			"444698f874aeb8c948dbf1ed70adfde5277867a5a4604df893c15fcec989" +
			"72bcc2b1afc700f4de8228d5d4d0ef1ac593a46b7f85efde655949177709" +
			"2a95a8c3997bf684993357f7fb91b190945f0bd16f1804bd14034400c045" +
			"a822323004aee54aae23dc10a8b8024215600331ac6c3d19e71902b5d968" +
			"1dde6ec6af2d433b5091960de8c7ac3afe76fd84ac667d533b4076c1bfdd" +
			"b98515369a1be440b98ed22160093496c2792d4a11776196bcde017c832c" +
			"6382bf5bf825f9c01b1c99c0eba01dca298b55918aebfc4bc2af0396877a" +
			"e8aa96cf552eeec3c938f326075409a925eaa881fbf5a324a425d9c2cb28" +
			"770de2962af40fcb8d97a15e68da79baea105ba46f236a2b0c0133b24815" +
			"eafc1db8c1955c4d5b3dda9422a930da6c2ee75b4a0c56f9addf7846818d" +
			"deaa8f05dffbfa415811347dffe3f416480000062b000de2e00000000000" +
			"000000000000000000000000000000000000000000000000000000000000" +
			"001cfff95880400000211bceff7ffffffff93b2c1e82c9819345247f18ac" +
			"ead580a6254b22aeaac4")
	if err != nil {
		t.Errorf("DecodeString: %s", err)
	}

	adts, err := ParseADTS(buf)

	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if adts.Fill_elements == nil {
		t.Errorf("Fill_elements must not be nil")
	}
	for i := range adts.Fill_elements {
		if adts.Fill_elements[i].Extension_payload == nil {
			t.Errorf("Extension_payload %d must not be nil", i)
		}
	}
}

func TestSbrChannelPairElementParse(t *testing.T) {
	buf, err := base64.StdEncoding.DecodeString(
		"//FYgD+BnCEbQ9uJooYaHEyKrWaAZUmYoqZdSkJUqh/ZoD0ZT3iVbMTgNe8FuI6HSt1vZoAcY13PUkjIqz2PXN82feqhznnAoQwAnAYzxDC3VbYX7x1oyZMJJDh0nozSuFKVaJUjeJVUi01puKLL3LdgXpQyDr0hZa9PAvAtQvxBLVNTXGa6FRWh2VU1r+fA3WhRZP1ChIYijIlCvUJgErbYRTWCgmIghE9cuudbSXxrNF5VNQ3LbkkuRBA6hl73rHJa8C2vUAhslyiC8nspS800OvMc3YzrTHnNlrcYMCXC/Bwha6I5KDVXi3O5UQaDWyinMlymVxOeNUp3FbdbWC4+KfCKGNPbI0aNSYxoo1b0xpUqoWRTZ1p8UU5gqnHFRzrmnThZGNyD8MZa7Y/NYOdthRkyZyD+S1nc0HDm2WoPc7O23AWWNkcLa/kTFekvbxFTk6DcJjxM3IQeqK2gKo5NynkRqIKyFspol79egkg3O8DIf68aWV/C4qjg0r5ODFsMHZQfTQI7ukK2Fh33vrHrN2OlhayZAQJp0gE2DCaInTHGgCpmmag8WmNqWBDZmRmZgmxZobHrmRPKCRtiEkAX8pJblEvdB61vNxNE3WJlvefJV7Y+FLJAXUA3gG8N37ATwNsYoCIAAAAAE/fwe/879/u5fmSNCBAA8w==")
	adts, err := ParseADTS(buf)
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if adts.Fill_elements == nil {
		t.Errorf("Fill_elements cannot be nil")
	}
	if adts.Fill_elements[0].Extension_payload == nil {
		t.Errorf("Extension_payload cannot be nil")
	}
	if int(adts.Fill_elements[0].Extension_payload.Extension_type) != EXT_SBR_DATA {
		t.Errorf("Extension_type (%d) must be of type EXT_SBR_DATA (%d)", adts.Fill_elements[0].Extension_payload.Extension_type, EXT_SBR_DATA)
	}
}

// AAC Audio frame with a full, parsable SBR
func TestSbrParse(t *testing.T) {
	buf, err := base64.StdEncoding.DecodeString(
		"//lYsED//CEblJXelQhhoI7AOFHFc2cMsXe1xEy5V4Brm15v1sz4eh7l7KRzWGdAkVA2OrmHG7LgFU0bV955cbSq3cp2WaD3C6ffadAnTyO0LAp0eHtwGg+300aUm66XufR7VGTEoqJ5sZeCebnxe5OBC1RyIkjiTBsxfxmvQ+Pw07tWjs8/DRtAqLDkM31bHKwkz3XI4eYtj9QzqmKEVyCrRm0keaDTOBlxk0A+Li2lkrjBWyvwnGH2/8xm+DSPD3itzYvkJk8q1UkEpWfqSq+SfDtYuykL63NmMahaF9hsgrTsBOqhQcBQ1xbOlVJUpetZZUVLmRYCcCVxVVEqRcBIJgggLQErUzSqDmUJQ5SUZaoYLQB2ZFKun714ZTLTCoCECqL7Lw3Iz3ybsZMtpiG+RCrAs5DacKMpmNImR9MyprUFyhLKOZZEJzTBll39USy0y2nR1keZfnZGXeB2IGMlnGm+6kKJluBqAek2ov8kqUSRqnhOUhn5djRiGQKMubUzrSEzBprrSR/5/oQNqYUB3EBrjXU9RS0DJ7OwaKERWjsxysubUGKdAH3jeEb9iJ4EIIokxrt/vwAmK9z7/gB9H0AG8yAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP//")
	adts, err := ParseADTS(buf)
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if adts.Fill_elements == nil {
		t.Errorf("Fill_elements cannot be nil")
	}
	if adts.Fill_elements[0].Extension_payload == nil {
		t.Errorf("Extension_payload cannot be nil")
	}
	if int(adts.Fill_elements[0].Extension_payload.Extension_type) != EXT_SBR_DATA {
		t.Errorf("Extension_type (%d) must be of type EXT_SBR_DATA (%d)", adts.Fill_elements[0].Extension_payload.Extension_type, EXT_SBR_DATA)
	}
}

// AAC Audio frame with EIGHT_SHORT_SEQUENCE
func TestEightShortSequence(t *testing.T) {
	buf, err := base64.StdEncoding.DecodeString(
		"//FYQC3AOAE2v+oTcrBNds4KmC34SXcuXJI78bpLJJjK3ElySygIWU7DiWttnUA6dGYOwhq8LCU2rjXxxe7LpiAHOb1NzwOwQog5lWi+c5WtG8WtmBaKhodeTPZfIWe9L1wVLMlnU7pUuNlysIHrnYyC3FsZxBb5gCkQMOi1UcJ/t32juqK2nUZli+uLL/P6MRvfVytWTOMZ3SJUwp0Ii+nw0L55xE1WcRERKKsxc6Qj2zz3JTdzZlYKKCgDSAAABBfRe173cWSt4FgGRuFKBAziNCf626dvrMZIgDLJkz0QDm1JwdWdZtrRdgr6q6c5IpGM7qekCaMx2SUf1hxKfX2p8P63evN0A+E1M4AUgExBq5UxcZP/WDjafvjZeem1R3Y9ke3uzwaGqQZDSma22qFYJApP73OsIBIjk6WU3Y7CBgITOe+hpkOaZ65KnAd4foEAJK9+3PsOPgdmRZfLySssSyVX+L1emNzzYAAH//FYQCuAIAFO16CNdDsFCMVBsFBEgkRcSSEhEJIkACz6gkmNiYhmTsJzYTVy+WWjfODzgFMqE5tYKBkRkBgohZURlF7zzYdV1k602q7mAJCIVAKgnCnbeW7EaNkmwh2bqaEwWQDKZF/pnO+B4b5Po3lYegtuVHbWQHDZcPZnxE0m/Y6Iq3OrmJbR26tl/YJGgwoTMwDrLXjWufzZkYCkUi86I0bJWfMbbVe2hZjDlYyI5EoFtfLZ43d+Jb7t11DzIfOweMuNF8AGUJbMCo6vKAqDpBzTMlQySBUUApgYVHaqaabVmkHKy4FIrAiejnUl7le8Lakm7fp7rXqAhRjY3EUBL+wE7OgWVQapk2qiNr44RsLWOXVDQXkqr9TXQUdp9G58x7AHmcg8FJhhbMiC4cq55d1WL1ob3iegwAK+R9/jftT8bT/yXDhx8AsyLZUNwQssO3Xut7LNYAHA//FYQCiAHAFMF6CNFCsFDMIiMJBqNgAQSJIkSQkiQSJEWscESAPEsYOTC3Z9GXOZa2ILetFdq5+0lrUjKUUdTphAxGK2GfVltodfN4Nf5pg8Zkm7Jl0VFkD7EUr5sF+PhjWZSjkF2ulTCc6da4G46Jb3Q9pk97Aw8IPKxTKH151kBtBM5ni+dbZaF3Ld4d20J4dZhBV/ivcVucP8Y6tTMp1mohtAFWE=")
	if err != nil {
		t.Errorf("DecodeString: %s", err)
	}
	adts, err := ParseADTS(buf)
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if adts.Fill_elements == nil {
		t.Errorf("Fill_elements cannot be nil")
	}
	for _, e := range adts.Fill_elements {
		if e.Extension_payload == nil {
			t.Errorf("Extension_payload cannot be nil")
		}
		if int(e.Extension_payload.Extension_type) != EXT_SBR_DATA {
			t.Errorf("Extension_type (%d) must be of type EXT_SBR_DATA (%d)", e.Extension_payload.Extension_type, EXT_SBR_DATA)
		}
	}
}

// Fuzz testing courtesy of gy741
func TestWindowGroupingFuzz(t *testing.T) {

	var crashers = []string{
		"\xff\xf10\x850",
	}

	for _, f := range crashers {
		ParseADTS([]byte(f))
	}
}

func TestADTSHeaderFuzz(t *testing.T) {

	var crashers = []string{
		"\xff\xf00000\x010",
	}

	for _, f := range crashers {
		ParseADTS([]byte(f))
	}
}
