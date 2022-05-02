package cmd

// providers is a registry of available ranking providers
var providers map[string]RankingProvider

// RankingProvider defines methods that need to be implemented
// when adding sources of URL rankings
type RankingProvider interface {
	Provider() string
	TopSitesGlobal() ([]string, error)
	TopSitesCountry() (map[string][]string, error)
}

// addRankingProvider registers a RankingProvider
// to the list of available ranking provider
// This is called inside the init() function of the file
// containing the definition of the RankingProvider being added
func addRankingProvider(rp RankingProvider) {
	if providers == nil {
		providers = make(map[string]RankingProvider)
	}

	providers[rp.Provider()] = rp
}
