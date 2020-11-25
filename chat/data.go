package chat

type DataFinder interface {
	GetIntID() int
	GetStringID() string
	GetName() string
	GetAvatarURL() string
}

type ClientDataFinder interface {
	DataFinder
}

type RoomDataFinder interface {
	DataFinder
}
