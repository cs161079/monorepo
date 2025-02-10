package mapper

type Route01Mapper interface {
}

type route01Mapper struct {
}

func NewRouteDetailMapper() Route01Mapper {
	return route01Mapper{}
}
