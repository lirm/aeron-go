package logbuffer

/*
struct DataFrameHeaderDefn {
    std::int32_t frameLength;
    std::int8_t version;
    std::int8_t flags;
    std::uint16_t type;
    std::int32_t termOffset;
    std::int32_t sessionId;
    std::int32_t streamId;
    std::int32_t termId;
    std::int64_t reservedValue;
},
*/

var DataFrameHeader = struct {
	FRAME_LENGTH_FIELD_OFFSET   int32
	VERSION_FIELD_OFFSET        int32
	FLAGS_FIELD_OFFSET          int32
	TYPE_FIELD_OFFSET           int32
	TERM_OFFSET_FIELD_OFFSET    int32
	SESSION_ID_FIELD_OFFSET     int32
	STREAM_ID_FIELD_OFFSET      int32
	TERM_ID_FIELD_OFFSET        int32
	RESERVED_VALUE_FIELD_OFFSET int32
	DATA_OFFSET                 int32

	LENGTH int32

	HDR_TYPE_PAD   uint16
	HDR_TYPE_DATA  uint16
	HDR_TYPE_NAK   uint16
	HDR_TYPE_SM    uint16
	HDR_TYPE_ERR   uint16
	HDR_TYPE_SETUP uint16
	HDR_TYPE_EXT   uint16

	CURRENT_VERSION int8
}{
	0,
	4,
	5,
	6,
	8,
	12,
	16,
	20,
	24,
	32,

	32,

	0x00,
	0x01,
	0x02,
	0x03,
	0x04,
	0x05,
	0xFFFF,

	0x0,
}
