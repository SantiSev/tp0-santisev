package protocol

const HEADER = "\x01"
const EOF = "\xFF"
const SUCCESS_HEADER_SIZE = 1
const WINNER_HEADER_SIZE = 1
const SUCCESS_HEADER = "\x02"
const SUCCESS_MESSAGE_SIZE = 4
const WINNERS_HEADER = "\x03"
const MAX_BATCH_SIZE = 8192 // 8 kB
const WINNER_COUNT_SIZE = 2
