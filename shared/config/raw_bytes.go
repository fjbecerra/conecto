package config

type ConfigBytes []byte

func (c *ConfigBytes) UnmarshalJSON(data []byte) error {
    *c = append((*c)[:0], data...)
    return nil
}