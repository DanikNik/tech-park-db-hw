package db

import (
	"github.com/jackc/pgx"
	"tech-park-db-hw/internal/pkg/models"
)

func DoVote(threadSlugOrId string, vote models.Vote) (*models.Thread, error) {
	threadData, err := GetThread(threadSlugOrId)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, err
	}

	_, err = Exec(CreateVoteQuery, vote.Nickname, threadData.Id, vote.Voice)
	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case foreignKeyError:
				return nil, ErrNotFound
			}
		}
		//return nil, err
		panic(err)
	}

	row := dbObj.QueryRow(threadGetVotesCountQuery, threadData.Id)
	row.Scan(&threadData.Votes)
	return threadData, nil
}
