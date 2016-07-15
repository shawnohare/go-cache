# go-cache
Simple Golang caching package.

## Example

```go
var p *redis.Pool
// ... initialize p
cache := &gcredis.Cache{Pool: p, HashKeys: true}
namespace := []string{"myapp", "mysection"}
objectID := "id"
if err := cache.Set(namespace, objectID, 123); err != nil {
  fmt.Println("Error setting value:", err)
}
```
