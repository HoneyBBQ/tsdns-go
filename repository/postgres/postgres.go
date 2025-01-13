package postgres

import (
	"github.com/honeybbq/tsdns-go/repository/postgres/model"
	"github.com/honeybbq/tsdns-go/repository/postgres/query"
	"github.com/honeybbq/tsdns-go/types"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
	q  *query.Query
}

// NewRepository creates a new PostgreSQL storage implementation
//
// dsn is the PostgreSQL connection string
func NewRepository(dsn string) (types.RecordRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	query.SetDefault(db)

	return &repository{
		db: db,
		q:  query.Q,
	}, nil
}

// MustNewRepository creates a new PostgreSQL storage implementation
//
// dsn is the PostgreSQL connection string
//
// MustNewRepository panics if an error occurs
func MustNewRepository(dsn string) types.RecordRepository {
	repo, err := NewRepository(dsn)
	if err != nil {
		panic(err)
	}
	return repo
}

// Convert model.Record to tsdns.Record
func (p *repository) toRecord(m *model.Record) *types.Record {
	if m == nil {
		return nil
	}

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &types.Record{
		ID:         m.ID,
		InstanceID: m.InstanceID,
		Domain:     m.Domain,
		Target:     m.Target,
		Port:       m.Port,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

// Convert tsdns.Record to model.Record
func (p *repository) toModel(r *types.Record) *model.Record {
	if r == nil {
		return nil
	}

	var deletedAt gorm.DeletedAt
	if r.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{
			Time:  *r.DeletedAt,
			Valid: true,
		}
	}

	return &model.Record{
		ID:         r.ID,
		InstanceID: r.InstanceID,
		Domain:     r.Domain,
		Target:     r.Target,
		Port:       r.Port,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

func (p *repository) Find() ([]*types.Record, error) {
	models, err := p.q.Record.Find()
	if err != nil {
		return nil, err
	}

	records := make([]*types.Record, len(models))
	for i, m := range models {
		records[i] = p.toRecord(m)
	}
	return records, nil
}

// FindByDomain finds a record by domain name
func (p *repository) FindByDomain(domain string) (*types.Record, error) {
	m, err := p.q.Record.Where(p.q.Record.Domain.Eq(domain)).First()
	if err != nil {
		return nil, err
	}
	return p.toRecord(m), nil
}

// Create creates a new DNS record
func (p *repository) Create(record *types.Record) error {
	m := p.toModel(record)
	return p.q.Record.Create(m)
}

// Delete removes a DNS record by domain
func (p *repository) Delete(domain string) error {
	_, err := p.q.Record.Where(p.q.Record.Domain.Eq(domain)).Delete()
	return err
}

// DeleteByInstanceID removes all records for a specific instance
func (p *repository) DeleteByInstanceID(instanceID int64) error {
	_, err := p.q.Record.Where(p.q.Record.InstanceID.Eq(instanceID)).Delete()
	return err
}

// Close closes the storage connection
func (p *repository) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
