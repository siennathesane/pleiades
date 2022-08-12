# Pleiades Protocols

## Rules 

1. Stream must open with state message
2. Stream must send headers before each message
3. Stream must send invalid state if recoverable
4. Stream must close if in unrecoverable state
5. Stream must send state message after each payload message

## Headers

1. All headers must contain 3 fields:
  1. Type
  1. Size
  1. Checksum

The type field must be an enumeration of all types for the specific protocol

The size field must be the size of the next message

The checksum must be a CRC-32 checksum using the IEEE polynomial

## Workflows

