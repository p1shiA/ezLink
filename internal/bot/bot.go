package bot

import (
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
	"github.com/p1shiA/ezLink/config"
)

func Run(cfg *config.Config) (*gotgproto.Client, error) {
	bot, err := gotgproto.NewClient(
		cfg.ApiID,
		cfg.ApiHash,
		gotgproto.ClientTypeBot(cfg.Token),
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqlSession(sqlite.Open("ezlink.db")),
		},
	)

	if err != nil {
		return nil, err
	}

	return bot, nil
}