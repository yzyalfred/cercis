package network

// socket状态
const (
	SS_ACCEPT uint8 = iota		// 接收
	SS_SHUT_DOWN				// 关闭
)

// 连接类型
const (
	CONNECT_CLIENT uint8 = iota		// 客户端的连接
	CONNECT_SERVER			     	// 服务端的连接
)

// 最大连接数
const MAX_CONN_NUM = 3000

// 发送chan的大小
const SEND_CHAN_SIZE = 100

const (
	MSG_LEN_SIZE = 2											// 消息包长度大小
)