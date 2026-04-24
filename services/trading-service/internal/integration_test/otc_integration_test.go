//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- Seed helpers ---

func seedAssetOwnership(t *testing.T, db *gorm.DB, identityID uint, ownerType model.OwnerType, assetID uint, amount float64) *model.AssetOwnership {
	t.Helper()
	o := &model.AssetOwnership{
		UserId:    identityID,
		OwnerType: ownerType,
		AssetID:   assetID,
		Amount:    amount,
		UpdatedAt: time.Now(),
	}
	require.NoError(t, db.Create(o).Error)
	return o
}

func setPublicAmount(t *testing.T, db *gorm.DB, ownershipID uint, publicAmount, reservedAmount float64) {
	t.Helper()
	require.NoError(t, db.Model(&model.AssetOwnership{}).
		Where("asset_ownership_id = ?", ownershipID).
		Updates(map[string]any{
			"public_amount":   publicAmount,
			"reserved_amount": reservedAmount,
		}).Error)
}

func uniqueMIC(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("X%d", uniqueCounter.Add(1))
}

func uniqueTicker(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("T%d", uniqueCounter.Add(1))
}

// --- Publish endpoint tests ---

func TestOTCHandler_PublishAsset_ClientSuccess(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	ex := seedExchange(t, db, uniqueMIC(t))
	listing := seedListing(t, db, uniqueTicker(t), ex.MicCode, model.AssetTypeStock, 100.0)
	ownership := seedAssetOwnership(t, db, 50, model.OwnerTypeClient, listing.AssetID, 20)

	path := fmt.Sprintf("/api/client/50/assets/%d/publish", ownership.AssetOwnershipID)
	rec := performRequest(t, router, http.MethodPatch, path, map[string]any{"amount": 5}, authHeaderForClient(t, 50, 50))
	requireStatus(t, rec, http.StatusNoContent)
}

func TestOTCHandler_PublishAsset_ActuarySuccess(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	ex := seedExchange(t, db, uniqueMIC(t))
	listing := seedListing(t, db, uniqueTicker(t), ex.MicCode, model.AssetTypeStock, 50.0)
	ownership := seedAssetOwnership(t, db, 20, model.OwnerTypeActuary, listing.AssetID, 15)

	path := fmt.Sprintf("/api/actuary/20/assets/%d/publish", ownership.AssetOwnershipID)
	rec := performRequest(t, router, http.MethodPatch, path, map[string]any{"amount": 3}, authHeaderForAgent(t))
	requireStatus(t, rec, http.StatusNoContent)
}

func TestOTCHandler_PublishAsset_Unauthenticated(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	rec := performRequest(t, router, http.MethodPatch, "/api/client/50/assets/1/publish", map[string]any{"amount": 5}, "")
	requireStatus(t, rec, http.StatusUnauthorized)
}

func TestOTCHandler_PublishAsset_WrongOwner(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	ex := seedExchange(t, db, uniqueMIC(t))
	listing := seedListing(t, db, uniqueTicker(t), ex.MicCode, model.AssetTypeStock, 100.0)
	ownership := seedAssetOwnership(t, db, 50, model.OwnerTypeClient, listing.AssetID, 20)

	path := fmt.Sprintf("/api/client/99/assets/%d/publish", ownership.AssetOwnershipID)
	rec := performRequest(t, router, http.MethodPatch, path, map[string]any{"amount": 5}, authHeaderForClient(t, 99, 99))
	requireStatus(t, rec, http.StatusForbidden)
}

func TestOTCHandler_PublishAsset_NotFound(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	rec := performRequest(t, router, http.MethodPatch, "/api/client/50/assets/99999/publish", map[string]any{"amount": 1}, authHeaderForClient(t, 50, 50))
	requireStatus(t, rec, http.StatusNotFound)
}

func TestOTCHandler_PublishAsset_AmountExceedsOwned(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	ex := seedExchange(t, db, uniqueMIC(t))
	listing := seedListing(t, db, uniqueTicker(t), ex.MicCode, model.AssetTypeStock, 100.0)
	ownership := seedAssetOwnership(t, db, 50, model.OwnerTypeClient, listing.AssetID, 10)

	path := fmt.Sprintf("/api/client/50/assets/%d/publish", ownership.AssetOwnershipID)
	rec := performRequest(t, router, http.MethodPatch, path, map[string]any{"amount": 9999}, authHeaderForClient(t, 50, 50))
	requireStatus(t, rec, http.StatusBadRequest)
}

func TestOTCHandler_PublishAsset_UpdatesExistingPublicAmount(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	ex := seedExchange(t, db, uniqueMIC(t))
	listing := seedListing(t, db, uniqueTicker(t), ex.MicCode, model.AssetTypeStock, 100.0)
	ownership := seedAssetOwnership(t, db, 50, model.OwnerTypeClient, listing.AssetID, 20)
	setPublicAmount(t, db, ownership.AssetOwnershipID, 3, 0)

	path := fmt.Sprintf("/api/client/50/assets/%d/publish", ownership.AssetOwnershipID)
	rec := performRequest(t, router, http.MethodPatch, path, map[string]any{"amount": 8}, authHeaderForClient(t, 50, 50))
	requireStatus(t, rec, http.StatusNoContent)

	var updated model.AssetOwnership
	require.NoError(t, db.First(&updated, ownership.AssetOwnershipID).Error)
	require.Equal(t, float64(11), updated.PublicAmount)
}

// --- GetPublicOTCAssets tests ---

func TestOTCHandler_GetPublicOTCAssets_ReturnsList(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	ex := seedExchange(t, db, uniqueMIC(t))
	ticker := uniqueTicker(t)
	listing := seedListing(t, db, ticker, ex.MicCode, model.AssetTypeStock, 120.0)
	ownership := seedAssetOwnership(t, db, 50, model.OwnerTypeClient, listing.AssetID, 10)
	setPublicAmount(t, db, ownership.AssetOwnershipID, 6, 1)

	rec := performRequest(t, router, http.MethodGet, "/api/otc/public?page=1&page_size=10", nil, authHeaderForClient(t, 50, 50))
	requireStatus(t, rec, http.StatusOK)

	body := decodeResponse[map[string]interface{}](t, rec)
	data, ok := body["data"].([]interface{})
	require.True(t, ok)
	require.GreaterOrEqual(t, len(data), 1)

	entry := data[0].(map[string]interface{})
	require.Equal(t, float64(5), entry["available_amount"]) // 6 - 1
	require.NotEmpty(t, entry["ticker"])
	require.NotEmpty(t, entry["name"])
	require.NotEmpty(t, entry["security_type"])
}

func TestOTCHandler_GetPublicOTCAssets_UnpublishedNotIncluded(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	ex := seedExchange(t, db, uniqueMIC(t))
	ticker := uniqueTicker(t)
	listing := seedListing(t, db, ticker, ex.MicCode, model.AssetTypeStock, 100.0)
	ownership := seedAssetOwnership(t, db, 50, model.OwnerTypeClient, listing.AssetID, 10)
	// public_amount = 0 — should not appear
	setPublicAmount(t, db, ownership.AssetOwnershipID, 0, 0)

	rec := performRequest(t, router, http.MethodGet, "/api/otc/public?page=1&page_size=10", nil, authHeaderForClient(t, 50, 50))
	requireStatus(t, rec, http.StatusOK)

	body := decodeResponse[map[string]interface{}](t, rec)
	data := body["data"].([]interface{})
	for _, item := range data {
		entry := item.(map[string]interface{})
		require.NotEqual(t, ticker, entry["ticker"])
	}
}

func TestOTCHandler_GetPublicOTCAssets_Unauthenticated(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	rec := performRequest(t, router, http.MethodGet, "/api/otc/public", nil, "")
	requireStatus(t, rec, http.StatusUnauthorized)
}
