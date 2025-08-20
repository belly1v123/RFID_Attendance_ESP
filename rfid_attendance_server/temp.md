## middleware

```go 
package mt_middlewares

import (
	"context"
	"strings"

	"errors"
	"fmt"

	"log"
	"net/http"
	"regexp"
	"sync"

	common_responses "github.com/DipakShrestha-ADS/rms_go_api/common/responses"
	"github.com/DipakShrestha-ADS/rms_go_api/database/publicquery"
	"github.com/DipakShrestha-ADS/rms_go_api/internal/auth"
	mt_constants "github.com/DipakShrestha-ADS/rms_go_api/modules/multi-tenant/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type contextKey string

// const (
// 	TenantSchemaNameContextKey      contextKey = "xittooRmsTenantSchemaName"
// 	TenantAssetsDirectoryContextKey contextKey = "xittooRmsTenantAssetsDirectory"
// 	TenantQrCodeUrlContextKey       contextKey = "xittooRmsTenantQrCodeUrl"
// 	DBContextKey                    contextKey = "xittooRmsDB"
// 	RoleContextKey                  contextKey = "xittooRmsRole"
// )

type TenantData struct {
	Schema          string
	AssetsDirectory string
	QRCodeURL       string
}
type TenantCache struct {
	mu    sync.RWMutex
	cache map[string]TenantData // tenantID -> TenantData
}

func NewTenantCache() *TenantCache {
	return &TenantCache{
		cache: make(map[string]TenantData),
	}
}

func (tc *TenantCache) Get(tenantID string) (TenantData, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	data, ok := tc.cache[tenantID]
	return data, ok
}

func (tc *TenantCache) Set(tenantID string, data TenantData) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.cache[tenantID] = data
}

func SetTenantMiddleware(db *gorm.DB, cache *TenantCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			urlPath := r.URL.Path
			// referer := r.Header.Get("Referer")
			fmt.Println("url path tm: ", urlPath, " || host : ", r.Host, " || url -host : ", r.URL.Host, " || host name : ", r.URL.Hostname(), " || port : ", r.URL.Port(), " ||  origin : ", r.Header.Get("Origin"), " || referer : ", r.Header.Get("Referer"), " || schema : ", r.URL.Scheme)
			// skip for static files
			if strings.Contains(urlPath, "/asdvip") || strings.Contains(urlPath, "/rms-promotion-video") {
				// need to verify the url link to assets here
				ctx = context.WithValue(ctx, mt_constants.TenantSchemaNameContextKey, ".")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// if strings.Contains(urlPath, "/rms-promotion-video"){
			// 	// need to verify the url link to assets here
			// 	ctx = context.WithValue(ctx, mt_constants.TenantSchemaNameContextKey, ".")
			// 	next.ServeHTTP(w, r.WithContext(ctx))
			// 	return
			// }
			role := r.Header.Get("X-Role")
			tenantID := r.Header.Get("X-Tenant-ID")
			schemaName := r.Header.Get("Ads_sn")
			log.Println("role: ", role, " | tenant id: ", tenantID)
			tenantSchemaName := ""
			tenantAssetsDir := ""
			tenantQrCodeUrl := ""

			if role == "Super Admin" && tenantID == "" {
				// accessing only public schema for superadmin
				ctx = context.WithValue(ctx, mt_constants.TenantSchemaNameContextKey, "public")
				// if isSuperAdmin := auth.CheckForSuperAdminUserFromToken(r, ctx); !isSuperAdmin {
				// 	common_responses.ERROR(w, http.StatusUnauthorized, errors.New("su : unauthorized"))
				// 	return
				// }
				ctx = context.WithValue(ctx, mt_constants.DBContextKey, db)
				ctx = context.WithValue(ctx, mt_constants.RoleContextKey, "Super Admin")
				ctx = context.WithValue(ctx, mt_constants.TenantAssetsDirectoryContextKey, "")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else if role == "Super Admin" && tenantID != "" {
				ctx = context.WithValue(ctx, mt_constants.TenantSchemaNameContextKey, "public")
				ctx = context.WithValue(ctx, mt_constants.DBContextKey, db)
				// accessing the selected tenant info by superadmin
				if isSuperAdmin := auth.CheckForSuperAdminUserFromToken(r, ctx); !isSuperAdmin {
					common_responses.ERROR(w, http.StatusUnauthorized, errors.New("su 1 : unauthorized"))
					return
				}
				teanntUUID, err := uuid.Parse(tenantID)
				if err != nil {
					common_responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("invalid tenant id"))
					return
				}
				// tenantSchemaName, tenantAssetsDir, tenantQrCodeUrl, exists := cache.GetSchema(tenantID)
				tenantData, exists := cache.Get(tenantID)
				log.Println("tenant data inside middleware accessed by superamdin : ", tenantData, " | exists: ", exists)
				tenantSchemaName = tenantData.Schema
				tenantAssetsDir = tenantData.AssetsDirectory
				tenantQrCodeUrl = tenantData.QRCodeURL
				log.Println("tenant schema name accessed by superadmin: ", tenantSchemaName, " | tenant assets dir: ", tenantAssetsDir, " | tenant qr code url: ", tenantQrCodeUrl, " | exists: ", exists)
				if !exists {
					if tenantID != "" {
						q := publicquery.Use(db).RmsTenant
						q.UseTable(fmt.Sprintf("%s.%s", "public", q.TableName()))
						tenantInfo, err := q.WithContext(ctx).Where(q.ID.Eq(teanntUUID)).Where(q.Status).First()

						if err != nil {
							common_responses.ERROR(w, http.StatusNotFound, fmt.Errorf("su: tenant not found: %v", err))
							return
						}
						sanitizedTenantSchemaName, err := sanitizeTenantID(tenantInfo.TenantSchemaName)
						if err != nil {
							common_responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("su1 : invalid tenant schema name"))
							return
						}
						log.Println("tenant info: ", tenantInfo)
						tenantSchemaName = sanitizedTenantSchemaName
						tenantAssetsDir = tenantInfo.TenantAssetsDirectory
						tenantQrCodeUrl = tenantInfo.QrCodeUrl
						// caching the tenant info to avoid further db queries
						cache.Set(tenantID, TenantData{
							Schema:          tenantSchemaName,
							AssetsDirectory: tenantAssetsDir,
							QRCodeURL:       tenantQrCodeUrl,
						})
					}
				}

				ctx = context.WithValue(ctx, mt_constants.TenantSchemaNameContextKey, tenantSchemaName)
				ctx = context.WithValue(ctx, mt_constants.TenantAssetsDirectoryContextKey, tenantAssetsDir)
				ctx = context.WithValue(ctx, mt_constants.TenantQrCodeUrlContextKey, tenantQrCodeUrl)
				ctx = context.WithValue(ctx, mt_constants.DBContextKey, db)
				ctx = context.WithValue(ctx, mt_constants.RoleContextKey, role)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			if tenantID == "" {
				// handling the case for qr order
				if strings.Contains(urlPath, "/auth/verify-qr-code") || strings.Contains(urlPath, "/get-food-category-with-menus") || strings.Contains(urlPath, "/filter-food-menu") || strings.Contains(urlPath, "/filter-food-sub-menu") || strings.Contains(urlPath, "/create-multiple-qr-order") {
					canPass := true
					if(strings.Contains(urlPath, "/create-multiple-qr-order")){
						query := r.URL.Query()
						qt := query.Get("qt")
						if qt == "" {
							canPass = false
						}
					}
					if role == "ADS@Qr_OrDer" && canPass {
						if schemaName == "" {
							schemaName = "vip"
						}
						q := publicquery.Use(db).RmsTenant
						q.UseTable(fmt.Sprintf("%s.%s", "public", q.TableName()))
						tenantInfo, err := q.WithContext(ctx).Where(q.TenantSchemaName.Eq(schemaName)).Where(q.Status).First()
						if err != nil {
							common_responses.ERROR(w, http.StatusNotFound, fmt.Errorf("qr: tenant not found: %v", err))
							return
						}
						ctx = context.WithValue(ctx, mt_constants.TenantSchemaNameContextKey, schemaName)
						ctx = context.WithValue(ctx, mt_constants.TenantAssetsDirectoryContextKey, tenantInfo.TenantAssetsDirectory)
						ctx = context.WithValue(ctx, mt_constants.TenantQrCodeUrlContextKey, tenantInfo.QrCodeUrl)
						ctx = context.WithValue(ctx, mt_constants.DBContextKey, db)
						ctx = context.WithValue(ctx, mt_constants.RoleContextKey, role)
					}else{
						common_responses.ERROR(w, http.StatusForbidden, errors.New("qr : invalid tenant id"))
						return
					}
				} else {
					// fallback value for production vip
					// tenantSchemaName = "vip"
					// tenantAssetsDir = ""
					// tenantQrCodeUrl = "https://xittoorms.com/qr-order-v2-verify"
					// ctx = context.WithValue(ctx, TenantSchemaNameContextKey, tenantSchemaName)
					// ctx = context.WithValue(ctx, TenantAssetsDirectoryContextKey, tenantAssetsDir)
					// ctx = context.WithValue(ctx, TenantQrCodeUrlContextKey, tenantQrCodeUrl)
					// removed the fallback because the production to vip updated to tenantId
					common_responses.ERROR(w, http.StatusForbidden, errors.New("em : invalid tenant id"))
					return
				}

			} else {
				teanntUUID, err := uuid.Parse(tenantID)
				if err != nil {
					common_responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("pa : invalid tenant id"))
					return
				}
				tenantData, exists := cache.Get(tenantID)
				log.Println("tenant data inside middleware : ", tenantData, " | exists: ", exists)
				tenantSchemaName = tenantData.Schema
				tenantAssetsDir = tenantData.AssetsDirectory
				tenantQrCodeUrl = tenantData.QRCodeURL
				log.Println("tenant schema name: ", tenantSchemaName, " | tenant assets dir: ", tenantAssetsDir, " | tenant qr code url: ", tenantQrCodeUrl, " | exists: ", exists)
				if !exists {
					if tenantID != "" {

						ctx = context.WithValue(ctx, mt_constants.TenantSchemaNameContextKey, "public")
						q := publicquery.Use(db).RmsTenant
						q.UseTable(fmt.Sprintf("%s.%s", "public", q.TableName()))
						tenantInfo, err := q.WithContext(ctx).Where(q.ID.Eq(teanntUUID)).Where(q.Status).First()
						if err != nil {
							common_responses.ERROR(w, http.StatusNotFound, fmt.Errorf("tenant not found: %v", err))
							return
						}
						sanitizedTenantSchemaName, err := sanitizeTenantID(tenantInfo.TenantSchemaName)
						if err != nil {
							common_responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("san : invalid tenant schema name"))
							return
						}
						log.Println("tenant info: ", tenantInfo)
						tenantSchemaName = sanitizedTenantSchemaName
						tenantAssetsDir = tenantInfo.TenantAssetsDirectory
						tenantQrCodeUrl = tenantInfo.QrCodeUrl
						// caching the tenant info to avoid further db queries
						cache.Set(tenantID, TenantData{
							Schema:          tenantSchemaName,
							AssetsDirectory: tenantAssetsDir,
							QRCodeURL:       tenantQrCodeUrl,
						})
					}
				}

				ctx = context.WithValue(ctx, mt_constants.TenantSchemaNameContextKey, tenantSchemaName)
				ctx = context.WithValue(ctx, mt_constants.TenantAssetsDirectoryContextKey, tenantAssetsDir)
				ctx = context.WithValue(ctx, mt_constants.TenantQrCodeUrlContextKey, tenantQrCodeUrl)
			}

			ctx = context.WithValue(ctx, mt_constants.DBContextKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func sanitizeTenantID(tenantID string) (string, error) {
	// Allow only alphanumeric characters and underscores
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, tenantID); !matched {
		return "", fmt.Errorf("invalid tenant ID")
	}
	return tenantID, nil
}

func GetTenantSchemaName(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(mt_constants.TenantSchemaNameContextKey).(string)
	return tenantID, ok
}
func GetTenantAssetsDirectory(ctx context.Context) (string, bool) {
	assetsDir, ok := ctx.Value(mt_constants.TenantAssetsDirectoryContextKey).(string)
	if assetsDir == "" {
		return "static", ok
	}
	return "static/" + assetsDir, ok
}
func GetTenantAssetsDirectoryReplace(ctx context.Context) string {
	assetsDir, ok := ctx.Value(mt_constants.TenantAssetsDirectoryContextKey).(string)
	if !ok {
		assetsDir = ""
	}
	return "asdvip/" + assetsDir
}
func GetTenantQrCodeUrl(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(mt_constants.TenantQrCodeUrlContextKey).(string)
	return tenantID, ok
}

// GetQualifiedTableName returns the schema-qualified table name based on the context
func GetQualifiedTableName(ctx context.Context, table string) (string, bool) {
	schema, ok := GetTenantSchemaName(ctx)
	if !ok {
		return table, false
	}
	return fmt.Sprintf("%s.%s", schema, table), true
}

// GetDB retrieves tenant-specific DB from context
func GetDB(ctx context.Context) (*gorm.DB, bool) {
	baseDB, ok := ctx.Value(mt_constants.DBContextKey).(*gorm.DB)
	if !ok {
		return nil, false
	}
	return baseDB, true
}

// func GetDB(ctx context.Context) (*gorm.DB, bool) {
// 	mu := sync.RWMutex{}
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	schema, ok := ctx.Value(TenantSchemaNameContextKey).(string)
// 	if !ok {
// 		return nil, false
// 	}
// 	log.Println("schema Name inside get db: ", schema)
// 	baseDB, ok := ctx.Value(DBContextKey).(*gorm.DB)
// 	if !ok {
// 		return nil, false
// 	}
// 	db := baseDB.Session(&gorm.Session{SkipDefaultTransaction: true, NewDB: true, PrepareStmt: false})
// 	db.Set("gorm:prepare_stmt", false)
// 	txt := db.Exec(fmt.Sprintf("SET search_path TO %s", schema))
// 	if err := txt.Error; err != nil {
// 		log.Printf("Failed to set search_path for schema %s: %v", schema, err)
// 		return nil, false
// 	}
// 	var searchPath string
// 	err := txt.Raw("SHOW search_path").Scan(&searchPath)
// 	log.Printf("DB search_path: %s", searchPath)
// 	if err.Error != nil || searchPath != schema {
// 		fmt.Printf("failed to verify search_path for schema %s, got %s: %v", schema, searchPath, err.Error)
// 		return nil, false
// 	} else {
// 		fmt.Println("matched search path")
// 	}
// 	log.Printf("Created new connection for schema %s with search_path: %s", schema, searchPath)
// 	return txt, true
// }

// GetPublicDB retrieves a DB connection for the public schema
func GetPublicDB(ctx context.Context) (*gorm.DB, bool) {
	baseDB, ok := ctx.Value(mt_constants.DBContextKey).(*gorm.DB)
	if !ok {
		return nil, false
	}
	db := baseDB.Session(&gorm.Session{SkipDefaultTransaction: true, NewDB: true})
	err := db.Exec(fmt.Sprintf("SET search_path TO %s", "public")).Error
	if err != nil {
		log.Printf("Failed to get public DB: %v", err)
		return nil, false
	}
	return db, true
}

func IsSuperadmin(ctx context.Context) bool {
	role, _ := ctx.Value(mt_constants.RoleContextKey).(string)
	return role == "Super Admin"
}

```



## to generate gorm functions
```go
package main

import (
	"log"

	"github.com/DipakShrestha-ADS/rms_go_api/api/models"
	// tb_models "github.com/DipakShrestha-ADS/rms_go_api/modules/table_booking/models"
	"gorm.io/gen"
)


const (
	PUBLIC_OUTPUT_DIR = "./database/publicquery"
	OUTPUT_DIR = "./database/rmsquery"
)
func main(){
	log.Println("hello this is gen")
	allOtherModels := models.GetAllModels()
	// allOtherModels = append(allOtherModels, tb_models.TableBooking{})
	GenereateQuery( allOtherModels, OUTPUT_DIR)
	GenereateQuery( models.GetPublicModels(), PUBLIC_OUTPUT_DIR)
}

func GenereateQuery( allModels []interface{}, outputPath string) {
	g := gen.NewGenerator(gen.Config{
		OutPath:       outputPath,
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		FieldNullable: true,
	})

	// g.UseDB(db)
	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(allModels...)

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {},models.GetAllModels()...)

	// Generate the code
	g.Execute()
}
```