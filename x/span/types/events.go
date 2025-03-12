package types

// 이벤트 타입
const (
	EventTypeCreateSpan   = "create_span"
	EventTypeUpdateParams = "update_params"
	EventTypeNewSpan      = "new_span"
	EventTypeSpanEnd      = "span_end"
)

// 이벤트 속성 키
const (
	AttributeKeySpanID          = "span_id"
	AttributeKeyStartBlock      = "start_block"
	AttributeKeyEndBlock        = "end_block"
	AttributeKeySpanLength      = "span_length"
	AttributeKeyActiveSpanCount = "active_span_count"
	AttributeKeyChainID         = "chain_id"
	AttributeKeyValidatorCount  = "validator_count"
	AttributeKeyProducerCount   = "producer_count"
)
