# Bitcask Simple Key-Value Store database in Golang
https://riak.com/assets/bitcask-intro.pdf

## Goals
1. DB Write - store data in file
2. DB Read - be able to read data
3. Handle file compaction
4. Handle segment merging
5. frequently writing hashmap into disk for faster reloading when database restarts

## Socket communication protocol

### Operations

1. get key
2. put key=value
3. delete key

format:
`[OPCODE][KEY SIZE][KEY][VALUE SIZE][VALUE]`
`[1][1][2][any size max 65,535 bytes]`

### OPCODE

We could represent op codes with 1 byte

1. GET -> 0
2. PUT -> 1
3. DELETE -> 2

### KEY SIZE

We would represent our key size with 1 byte giving us a maximum of key size of 255 bytes

### VALUE SIZE

We would represent our data length with 2 bytes giving us a maximum of approximately 64KiB of single data write

1 Byte = 8 bits
2 Byte = 16 bits
(2^16 - 1) = 0 - 65,535

### VALUE

we can read this based on the data length

Note: we can improve this protocol later
