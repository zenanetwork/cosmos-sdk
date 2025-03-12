package types

// 이벤트 타입
const (
	EventTypeCreateCheckpoint = "create_checkpoint"
	EventTypeUpdateParams     = "update_params"
	EventTypeNewCheckpoint    = "new_checkpoint"
	EventTypeVerifyCheckpoint = "verify_checkpoint"
)

// 이벤트 속성 키
const (
	AttributeKeyCheckpointNumber     = "checkpoint_number"
	AttributeKeyStartBlock           = "start_block"
	AttributeKeyEndBlock             = "end_block"
	AttributeKeyCheckpointInterval   = "checkpoint_interval"
	AttributeKeyCheckpointBufferSize = "checkpoint_buffer_size"
	AttributeKeyRootHash             = "root_hash"
	AttributeKeyProposer             = "proposer"
	AttributeKeyTimestamp            = "timestamp"
)
