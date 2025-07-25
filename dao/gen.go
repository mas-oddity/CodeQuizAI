// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q             = new(Query)
	Paper         *paper
	PaperQuestion *paperQuestion
	Question      *question
	TempQuestion  *tempQuestion
	User          *user
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Paper = &Q.Paper
	PaperQuestion = &Q.PaperQuestion
	Question = &Q.Question
	TempQuestion = &Q.TempQuestion
	User = &Q.User
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:            db,
		Paper:         newPaper(db, opts...),
		PaperQuestion: newPaperQuestion(db, opts...),
		Question:      newQuestion(db, opts...),
		TempQuestion:  newTempQuestion(db, opts...),
		User:          newUser(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Paper         paper
	PaperQuestion paperQuestion
	Question      question
	TempQuestion  tempQuestion
	User          user
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:            db,
		Paper:         q.Paper.clone(db),
		PaperQuestion: q.PaperQuestion.clone(db),
		Question:      q.Question.clone(db),
		TempQuestion:  q.TempQuestion.clone(db),
		User:          q.User.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:            db,
		Paper:         q.Paper.replaceDB(db),
		PaperQuestion: q.PaperQuestion.replaceDB(db),
		Question:      q.Question.replaceDB(db),
		TempQuestion:  q.TempQuestion.replaceDB(db),
		User:          q.User.replaceDB(db),
	}
}

type queryCtx struct {
	Paper         IPaperDo
	PaperQuestion IPaperQuestionDo
	Question      IQuestionDo
	TempQuestion  ITempQuestionDo
	User          IUserDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Paper:         q.Paper.WithContext(ctx),
		PaperQuestion: q.PaperQuestion.WithContext(ctx),
		Question:      q.Question.WithContext(ctx),
		TempQuestion:  q.TempQuestion.WithContext(ctx),
		User:          q.User.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
