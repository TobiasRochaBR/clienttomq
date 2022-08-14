package main

const (
	// Publish Message
	MsgPublishTcpReq = 1001
	MsgPublishTcpAck = 1002
	// Unacknowledge Message
	MsgNAckTcpReq = 1003
	MsgNAckTcpAck = 1004
	// Acknowledge Message
	MsgAckTcpReq = 1005
	MsgAckTcpAck = 1006
	// Reject Message
	MsgRejectTcpReq = 1007
	MsgRejectTcpAck = 1008
	// Next Message
	MsgNextTcpReq = 1009
	MsgNextTcpAck = 1010
	// Create Channel
	ChannelCreateTcpReq = 1011
	ChannelCreateTcpAck = 1012
	// Join Channel
	ChannelJoinTcpReq = 1013
	ChannelJoinTcpAck = 1014
	// Register Consumer
	ConsumerRegisterTcpReq = 1015
	ConsumerRegisterTcpAck = 1016
	// Register Publisher
	ConsumerPublisherTcpReq = 1017
	ConsumerPublisherTcpAck = 1018
	// distribute
	MsgDistributeTcpReq = 1019
	MsgDistributeTcpAck = 1020
	// LOGIN
	// LOGOFF
)
