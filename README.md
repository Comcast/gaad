
# GAAD (Go Advanced Audio Decoder)

Package currently provides AAC parsing capabilities.  This package performs a full parse of AAC-LC and HE-AACv1 bitstreams.  Bitstreams with Parametric Stereo (HE-AACv2) are not yet supported, but AAC bitstream data and SBR data will be extracted. The AAC decode from the parsed data to LPCM (.wav) is not yet implemented.  Please help us expand and test this library!

## AACParser
This package currently supports AAC audio data contained in an ADTS header.  All availabe data is returned in the `adts` struct and can be accessed as nested objects as presented in the AAC specification.  All parameter names shoudl be verbatim from the AAC specification, if you find an issue with this please file a bug or submit a pull request.  

### AAC Types

| Type      | Description     | CODEC        |
| :-------: | --------------- | :----------: |
| AACLC     | AAC             | mp4a.40.2    |
| HEAAC     | AAC + SBR       | mp4a.40.5    |
| HEAACv2   | AAC + SBR + PS  | mp4a.40.29   |

where:
+ SBR = Spectral band replication
+ PS = Parametric Stereo

### VBR vs CBR

VBR (Variable bitrate) and CBR (Constant bitrate) is derived from the bitstream_type attribute in the adif_header section.  It is VBR if bitstream_type is true, and CBR otherwise.

### References

#### ISO/IEC STANDARD 14496-3

```
Title         : Coding of audio-visual objects — Part 3: Audio
File          : ISO_14496-3-4th-Edition.pdf
Edition       : Fourth edition (2009-09-01)
Relevant Sections
    - Page 64   : 1.6.5 Signaling of SBR
    - Page 120  : 1.A.2 AAC Interchange formats
    - Page 489  : 4.4.2 GA bitstream payloads
```
#### Related Books

* Video Demystified (5th Edition) http://my.safaribooksonline.com/book/-/9780750683951

#### Related Links

* http://www.mp4ra.org/codecs.html

