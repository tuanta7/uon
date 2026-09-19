package netplan

func DefaultRoutes(gateway string) []Route {
	if gateway == "" {
		return nil
	}
	return []Route{
		{
			To:  "default",
			Via: gateway,
		},
	}
}
