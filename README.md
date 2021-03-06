# jpath

Package jpath provides a query language for selecting elements from an JSON data.

## Examples

```json
{
    "arranged": "alphabetically",
    "companies": [
        {"name": "apple"},
        {"name": "facebook"},
        {"name": "github"},
        {"name": "google"}
    ]
}
```

| Query                | Result                       |
|----------------------|------------------------------|
| .arranged            | alphabetically               |
| .companies[:].name   | apple facebook github google |
| .companies[1:3].name | facebook github |
| .companies[^2].name  | github |
