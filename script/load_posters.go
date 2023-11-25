package script

type LoadPostersCmd struct{}

// LoadPosters queries Crunchyroll for a complete list of series, downloads the smallest size "posters" for each one,
// then converts them to base64 encoded strings the Tidbyt can use.
//
// This is a separate process since these aren't expected to change very often, and downloading loads of photos all the time
// would put unnecessary strain on Crunchyroll's servers.
func (cmd LoadPostersCmd) LoadPosters() {
	// imgconv
}
