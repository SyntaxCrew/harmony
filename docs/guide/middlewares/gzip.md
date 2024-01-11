# Gzip

## Usage
``` go
app.Use(middleware.Gzip())
```

## Custom Config
``` go
type GzipConfig struct {
    // Skipper defines a function to skip middleware.
    Skipper Skipper

    // Gzip compression level.
    // Optional. Default value -1.
    // -2 means use Huffman-only compression.
    // -1 means use default compression level.
    // 1 (BestSpeed) to 9 (BestCompression).
    Level int

    // MinLength is the minimum length required for compression.
    // Optional. Default value 0.
    MinLength int
}
```
### Example:
``` go
app.Use(middleware.Gzip(middleware.GzipConfig{
    Level: 6,
}))
```