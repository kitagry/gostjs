## gostjs

Create JSON Schema from golang struct.

### Installation

```
go get github.com/kitagry/gostjs
```

### Usage

```
Usage of gostjs:
  -output string
        output file name
  -src-path string
        go build src path
  -tag string
        tag name (default "json")
```

```
$ gostjs github.com/kitagry/gostjs/test.Test | gojq
{
  "$schema": "http://json-schema.org/schema#",
  "description": "Test's docs\n",
  "properties": {
    "": {
      "properties": {
        "Name": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "Child2": {
      "properties": {
        "Name": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "Interface": {
      "type": ""
    },
    "array": {
      "contains": {
        "type": "string"
      },
      "type": "array"
    },
    "child": {
      "properties": {
        "Name": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "id": {
      "type": "number"
    },
    "map": {
      "additionalProperties": {
        "type": "string"
      },
      "type": "object"
    },
    "name": {
      "description": "Name is name\n",
      "type": "string"
    }
  },
  "type": "object"
}
```

### License

MIT
