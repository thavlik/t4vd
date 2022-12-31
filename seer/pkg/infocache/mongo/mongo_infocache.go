package mongo

import (
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	channelRecencyCollection  = "channelrecency"  // tracks how recent the channel cache is
	playlistRecencyCollection = "playlistrecency" // tracks how recent the playlist cache is
	videoRecencyCollection    = "videorecency"    // tracks how recent the playlist cache is
	channelJoinCollection     = "channeljoins"    // tracks which videos are in which channels
	playlistJoinCollection    = "playlistjoins"   // tracks which videos are in which playlists
	cachedVideosCollection    = "cachedvideos"    // cache of video info
	cachedChannelsCollection  = "cachedchannels"  // cache of video info
	cachedPlaylistsCollection = "cachedplaylists" // cache of video info
	channelOriginKey          = "c"
	playlistOriginKey         = "p"
)

type mongoInfoCache struct {
	channelRecencyCollection  *mongo.Collection
	playlistRecencyCollection *mongo.Collection
	videoRecencyCollection    *mongo.Collection
	channelJoinCollection     *mongo.Collection
	playlistJoinCollection    *mongo.Collection
	cachedVideosCollection    *mongo.Collection
	cachedPlaylistsCollection *mongo.Collection
	cachedChannelsCollection  *mongo.Collection
}

func NewMongoInfoCache(db *mongo.Database) infocache.InfoCache {
	return &mongoInfoCache{
		channelRecencyCollection:  db.Collection(channelRecencyCollection),
		playlistRecencyCollection: db.Collection(playlistRecencyCollection),
		videoRecencyCollection:    db.Collection(videoRecencyCollection),
		channelJoinCollection:     db.Collection(channelJoinCollection),
		playlistJoinCollection:    db.Collection(playlistJoinCollection),
		cachedVideosCollection:    db.Collection(cachedVideosCollection),
		cachedPlaylistsCollection: db.Collection(cachedPlaylistsCollection),
		cachedChannelsCollection:  db.Collection(cachedChannelsCollection),
	}
}
