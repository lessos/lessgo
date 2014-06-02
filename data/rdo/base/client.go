package base

type Client struct {
    Config  Config
    Base    *Base
    Dialect DialectInterface
}
