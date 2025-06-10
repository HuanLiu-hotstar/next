package main

import (
    "context"

    "fmt"
)

func main() {
    ctx := context.Background()
    ctx = context.WithValue(ctx, "hello", "world")
    fmt.Println(ctx)
    for i := 0; i < 1; i++ {
        ctx := context.WithValue(ctx, i, fmt.Sprintf("%d", i))
        fmt.Println(ctx)
        for j := 0; j < 1; j++ {
            ctx := context.WithValue(ctx, 100+i, fmt.Sprintf("%d", 100-j))
            fmt.Println(ctx)
        }
    }
    m := map[string]string{
        "hello": "world",
    }
    ctx = context.WithValue(ctx, "map", m)
    fmt.Println(ctx)
    updateMap(ctx)
    v := ctx.Value("map").(map[string]string)
    fmt.Println("map after", v)

}
func updateMap(ctx context.Context) {
    v := ctx.Value("map").(map[string]string)
    v["xyz"] = "100"
    fmt.Println("map update", v)
}
