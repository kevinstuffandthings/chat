module kevinstuffandthings/chat

go 1.17

replace kevinstuffandthings/chat/handshake => ./handshake

require (
	kevinstuffandthings/chat/handshake v0.0.0-00010101000000-000000000000
	kevinstuffandthings/chat/server v0.0.0-00010101000000-000000000000
)

replace kevinstuffandthings/chat/server => ./server
