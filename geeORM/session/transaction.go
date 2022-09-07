package session

import "geeorm/log"

func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	defer func() {
		s.tx = nil
	}()
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
		return err
	}
	return
}

func (s *Session) Rollback() (err error) {
	defer func() {
		s.tx = nil
	}()
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
		return err
	}
	return
}
