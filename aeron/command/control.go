package command

const (
	// Clients to Media Driver

	/** Add Publication */
	ADD_PUBLICATION int32 = 0x01
	/** Remove Publication */
	REMOVE_PUBLICATION int32 = 0x02
	/** Add Subscriber */
	ADD_SUBSCRIPTION int32 = 0x04
	/** Remove Subscriber */
	REMOVE_SUBSCRIPTION int32 = 0x05
	/** Keepalive from Client */
	CLIENT_KEEPALIVE int32 = 0x06
)
