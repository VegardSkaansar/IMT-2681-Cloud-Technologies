package igcsdb

// this file have controll over all new datastructures
// in this application

var idForURL = 0

// ------------------------------------------------------------------------------------------------------------------

// this section is the declaration of the the different types

//ServerInfo is a struct with the declaration
type ServerInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// IgcURL datastructure for posting
type IgcURL struct {
	URL string `json:"url"`
	ID  int
}

type igcTrack struct {
	HDate       string  `json:"hdate"`
	Pilot       string  `json:"pilot"`
	Glider      string  `json:"glider"`
	GliderID    string  `json:"gliderid"`
	TrackLength float64 `json:"tracklength"`
}

//---------------------------------------------------------------------------------------------------------------------

// This section is where the databases will be established
// from the types in the section above

// IgcURLStorage represents a way to access the url from different igc files
type IgcURLStorage interface {
	Init()
	Add(i IgcURL) error
	Countl() int
	GetAll() []IgcURL
	GetURL(id int) (IgcURL, bool)
}

// IgcURLDB is the handle to URLS in-memory storage.
type IgcURLDB struct {
	igcUrls map[int]IgcURL
}

// Init initializes the in-memory storage.
func (db *IgcURLDB) Init() {
	db.igcUrls = make(map[int]IgcURL)
}

// Add adds a new url and a new id into the storage
func (db *IgcURLDB) Add(i IgcURL) error {
	db.igcUrls[i.ID] = i
	return nil
}

// Countl returns the current count of the urls in-memory
func (db *IgcURLDB) Countl() int {
	return len(db.igcUrls)
}

//GetAll returns all the urls as slice instead of map
func (db *IgcURLDB) GetAll() []IgcURL {
	all := make([]IgcURL, 0, db.Countl())
	for _, i := range db.igcUrls {
		all = append(all, i)
	}
	return all
}

// GetURL returns one url with information, and returns it, if the id doesnt exist a empty body will be sent
func (db *IgcURLDB) GetURL(id int) (IgcURL, bool) {
	igc, ok := db.igcUrls[id]
	return igc, ok
}

// -----------------------------------------------------------------------------

// GlobalDb is storage for the whole server
var GlobalDb IgcURLStorage
