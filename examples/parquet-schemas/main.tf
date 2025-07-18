resource "criblio_parquet_schema" "my_parquet_schema" {
  description = "test parquest"
  group_id    = "default"
  id          = "my_parquet_schema"
  schema      = <<-EOT
{
  "compressed_GZIP": {
    "type": "STRING",
    "compression": "GZIP"
  },
  "compressed_SNAPPY": {
    "type": "STRING",
    "compression": "SNAPPY"
  },
  "compressed_BROTLI": {
    "type": "STRING",
    "compression": "BROTLI"
  },
  "compressed_ZSTD": {
    "type": "STRING",
    "compression": "ZSTD"
  },
  "compressed_LZ4": {
    "type": "STRING",
    "compression": "LZ4"
  },
  "compressed_BOOLEAN": {
    "type": "BOOLEAN",
    "compression": "GZIP"
  },
  "compressed_INT32": {
    "type": "INT",
    "bitWidth": 32,
    "signed": true,
    "compression": "SNAPPY"
  },
  "compressed_INT64": {
    "type": "INT",
    "bitWidth": 64,
    "signed": true,
    "compression": "BROTLI"
  },
  "compressed_FLOAT": {
    "type": "FLOAT",
    "compression": "ZSTD"
  },
  "compressed_DOUBLE": {
    "type": "DOUBLE",
    "compression": "GZIP"
  },
  "compressed_BYTE_ARRAY": {
    "type": "BYTE_ARRAY",
    "compression": "GZIP"
  },
  "compressed_FIXED_LEN_BYTE_ARRAY": {
    "type": "FIXED_LEN_BYTE_ARRAY",
    "typeLength": 10,
    "compression": "GZIP"
  },
  "encoding_PLAIN": {
    "type": "STRING",
    "encoding": "PLAIN"
  },
  "encoding_DICTIONARY": {
    "type": "STRING",
    "encoding": "DICTIONARY"
  },
  "encoding_BYTE_STREAM_SPLIT": {
    "type": "FLOAT",
    "encoding": "BYTE_STREAM_SPLIT"
  },
  "compressed_encoded_STRING": {
    "type": "STRING",
    "compression": "GZIP",
    "encoding": "DICTIONARY"
  },
  "compressed_encoded_INT32": {
    "type": "INT",
    "bitWidth": 32,
    "signed": true,
    "compression": "ZSTD",
    "encoding": "DICTIONARY"
  },
  "compressed_encoded_FLOAT": {
    "type": "FLOAT",
    "compression": "GZIP",
    "encoding": "BYTE_STREAM_SPLIT"
  },
  "compressed_encoded_DOUBLE": {
    "type": "DOUBLE",
    "compression": "GZIP",
    "encoding": "BYTE_STREAM_SPLIT"
  },
  "compressed_encoded_BYTE_ARRAY": {
    "type": "BYTE_ARRAY",
    "compression": "GZIP",
    "encoding": "DICTIONARY"
  },
  "compressed_encoded_FIXED_LEN_BYTE_ARRAY": {
    "type": "FIXED_LEN_BYTE_ARRAY",
    "typeLength": 10,
    "compression": "GZIP",
    "encoding": "PLAIN"
  }
}
EOT
}

output "parquet_schema" {
  value = criblio_parquet_schema.my_parquet_schema
}
