package handlers

const MAX_UPLOAD_FILE_SIZE int64 = 1 << 30 // 1GB
var channelMapping = make(map[string](chan uint8))

var UPLOAD_DIRECTORY = "./uploads/"
