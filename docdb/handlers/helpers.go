package handlers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"platform/services/docdb/models"
	"strings"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/patrickmn/go-cache"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	authz "gitlab.com/youwol/platform/libs/go-libs/middleware"
	utils "gitlab.com/youwol/platform/libs/go-libs/utils"
	zap "go.uber.org/zap"
)

var OwnerIDColumn = "owner_id"
var OwnerNameColumn = "owner_name"
var OwnerKindColumn = "owner_kind"

var sessionMutexes map[string]*sync.Mutex = make(map[string]*sync.Mutex)
var c = cache.New(10*time.Minute, 10*time.Minute)
var scyllaHosts = os.Getenv("SCYLLA_HOSTS")

// GetKeyspaceName transforms the keyspace name so that it is tied to the current environment
func GetKeyspaceName(keyspaceName string) (string, error) {
	if strings.HasPrefix(keyspaceName, "system") {
		return keyspaceName, nil
	}

	forbidden := []string{"dev", "test", "prod", "demo", "staging"}
	for _, f := range forbidden {
		if strings.HasPrefix(keyspaceName, f) {
			return "", fmt.Errorf("Keyspace name cannot start with %s", f)
		}
	}
	return os.Getenv("ENVIRONMENT") + "_" + keyspaceName, nil
}

// GetHandlerContext retrieves the main objects that are necessary for handlers implementation
func GetHandlerContext(keyspaceName string, r *http.Request) (string, context.Context, *gocql.Session, zap.Logger, error) {
	// Common initialization of request context
	ksName, err := GetKeyspaceName(keyspaceName)
	ctx := r.Context()
	logger := utils.ContextLogger(ctx)
	if err != nil {
		return "", ctx, nil, logger, err
	}
	session := getSession(ctx, ksName)

	return ksName, ctx, session, logger, err
}

// getSession instantiates a new CQL session
func getSession(ctx context.Context, Keyspace string) *gocql.Session {
	lock(Keyspace)
	defer unlock(Keyspace)

	session, found := c.Get("session." + Keyspace)
	if !found || session == nil || session.(*gocql.Session).Closed() {
		session, _ = createSession(ctx, Keyspace)
	}
	return session.(*gocql.Session)
}

func lock(Keyspace string) {
	mtx, found := sessionMutexes[Keyspace]
	if !found {
		sessionMutexes[Keyspace] = &sync.Mutex{}
		mtx, _ = sessionMutexes[Keyspace]

	}
	mtx.Lock()
}

func unlock(Keyspace string) {
	mtx, _ := sessionMutexes[Keyspace]
	mtx.Unlock()
}

func createSession(ctx context.Context, Keyspace string) (*gocql.Session, error) {
	// Resolve IP addresses from hosts
	hosts := strings.Split(scyllaHosts, ",")
	IPs := make([]string, 0)
	logger := utils.ContextLogger(ctx)
	for _, h := range hosts {
		addr, err := net.DefaultResolver.LookupIP(ctx, "ip4", h)

		if err == nil {
			for _, a := range addr {
				IPs = append(IPs, a.String())
			}
		} else {
			logger.Error("Could not resolve IP", zap.String("host", h), zap.Error(err))
		}
	}
	cluster := gocql.NewCluster(IPs...)
	// TODO: check this out
	// fallback := gocql.DCAwareRoundRobinPolicy(localDC)
	// cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(localDC)
	cluster.Keyspace = Keyspace
	cluster.Timeout = 5000 * time.Millisecond
	cluster.ConnectTimeout = 5000 * time.Millisecond
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{NumRetries: 5, Min: 100 * time.Millisecond, Max: 5000 * time.Millisecond}
	session, err := cluster.CreateSession()
	if err != nil {
		c.Delete("session." + Keyspace)
		logger.Fatal("Could not create session", zap.Error(err), zap.String("scyllaHost", scyllaHosts))
	} else {
		c.Set("session."+Keyspace, session, cache.DefaultExpiration)
	}

	return session, err
}

// InvalidateSession closes the session after a workspace modification
func InvalidateSession(ctx context.Context, Keyspace string) {
	lock(Keyspace)
	defer unlock(Keyspace)

	_, found := c.Get("session." + Keyspace)
	if found {
		c.Delete("session." + Keyspace)
		createSession(ctx, Keyspace)
	}
}

// This function build metadata dynamically from the table in order to bind
// to entity objects (TODO: use a cache)
func getMeta(ctx context.Context, session *gocql.Session, keyspaceName string, tableName string) (*table.Metadata, map[string]*gocql.ColumnMetadata, []*gocql.ColumnMetadata, []*gocql.ColumnMetadata, error) {
	// generic way of using gocqlx, from our json structure
	// would be faster with a precise data model
	logger := utils.ContextLogger(ctx)
	keyspaceMetadata, err := session.KeyspaceMetadata(keyspaceName)
	if keyspaceMetadata == nil {
		logger.Error("Empty keyspace metadata", zap.Error(err))
		return nil, nil, nil, nil, err
	}
	tmd, found := keyspaceMetadata.Tables[tableName]
	if !found {
		err = errors.New("Table not found in the keyspace")
		return nil, nil, nil, nil, err
	}

	logger.Debug("meta", zap.Any("data", tmd))

	// Unfortunately we need to map manually from CQL to our API
	// GoCQL does not have the same exact structure as the Cassandra API
	columns := make([]string, len(tmd.Columns))
	partitionKey := make([]string, len(tmd.PartitionKey))
	clusteringCols := make([]string, len(tmd.ClusteringColumns))

	var colIndex = 0
	for _, val := range tmd.Columns {
		columns[colIndex] = val.Name
		colIndex++
	}
	for i, val := range tmd.PartitionKey {
		partitionKey[i] = val.Name
	}
	for i, val := range tmd.ClusteringColumns {
		clusteringCols[i] = val.Name
	}

	ret := &table.Metadata{
		Name:    tableName,
		Columns: columns,
		PartKey: partitionKey,
		SortKey: clusteringCols,
	}

	return ret, tmd.Columns, tmd.PartitionKey, tmd.ClusteringColumns, err
}

func buildSelectClause(selectBuilder *qb.SelectBuilder, sel *[]*models.SelectClause) {
	columns := make([]string, len(*sel))
	for i, s := range *sel {
		if len(s.Identifier) > 0 && len(s.Selector) > 0 {
			columns[i] = qb.As(s.Selector, s.Identifier)
		} else if len(s.Selector) > 0 {
			columns[i] = s.Selector
		}
	}
	selectBuilder.Columns(columns...)
}

func buildWhereClause(selectBuilder *qb.SelectBuilder, deleteBuilder *qb.DeleteBuilder, mapValues *qb.M, whereClause *[]*models.QueryRelation) {
	for _, v := range *whereClause {
		cmp := qb.Cmp{}
		switch v.Relation {
		case models.RelationOperatorEq:
			cmp = qb.Eq(*v.Column)
		case models.RelationOperatorLt:
			cmp = qb.Lt(*v.Column)
		case models.RelationOperatorLeq:
			cmp = qb.LtOrEq(*v.Column)
		case models.RelationOperatorGt:
			cmp = qb.Gt(*v.Column)
		case models.RelationOperatorGeq:
			cmp = qb.GtOrEq(*v.Column)
		case models.RelationOperatorIn:
			cmp = qb.In(*v.Column)
		case models.RelationOperatorCnt:
			cmp = qb.Contains(*v.Column)
		case models.RelationOperatorCntKey:
			cmp = qb.ContainsKey(*v.Column)
		case models.RelationOperatorLike:
			cmp = qb.Like(*v.Column)
		}
		if selectBuilder != nil {
			selectBuilder = selectBuilder.Where(cmp)
		}
		if deleteBuilder != nil {
			deleteBuilder = deleteBuilder.Where(cmp)
		}
		(*mapValues)[*v.Column] = v.Term
	}
}

func buildOrderByClause(selectBuilder *qb.SelectBuilder, orderClause *[]*models.Order) {
	for _, v := range *orderClause {
		if *v.Order == "ASC" {
			selectBuilder.OrderBy(*v.Name, qb.ASC)
		} else if *v.Order == "DESC" {
			selectBuilder.OrderBy(*v.Name, qb.DESC)
		}
	}
}

// Contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// GetCredentials gets ownership metadata when uploading documents or making queries
func GetCredentials(ctx context.Context, owner *string, auth *models.Principal) (string, string, string, error) {
	var ownerID string
	var ownerName string
	var ownerKind string
	if owner != nil && len(*owner) > 0 {
		// Check that we are member of this group
		bMatched := false
		for _, authGroup := range auth.MemberOf {
			if authGroup.Path == *owner { // Group name matches exactly
				ownerID = authGroup.ID
				ownerName = authGroup.Path
				ownerKind = "group"
				bMatched = true
			} else if strings.HasPrefix(authGroup.Path, *owner) { // We are using a parent group, then we need to retrieve its information
				grp := authz.GetGroupByPath(ctx, "docdb", *owner) // Get the parent group data
				if grp != nil {
					ownerID = *grp.ID
					ownerName = *grp.Path
					ownerKind = "group"
					bMatched = true
				} else {
					return "", "", "", fmt.Errorf("Error retrieving parent group")
				}
			}
		}
		if !bMatched {
			return "", "", "", fmt.Errorf("You are not member of this group")
		}
	} else {
		// Set default ownership to the current user
		ownerID = *auth.Sub
		ownerName = auth.PreferredUsername
		ownerKind = "user"
	}

	return ownerID, ownerName, ownerKind, nil
}
