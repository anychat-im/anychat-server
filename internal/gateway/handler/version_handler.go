package handler

import (
	"net/http"
	"strconv"

	versionpb "github.com/anychat/server/api/proto/version"
	"github.com/gin-gonic/gin"
)

type VersionHandler struct {
	clientManager interface {
		Version() versionpb.VersionServiceClient
	}
}

func NewVersionHandler(cm interface {
	Version() versionpb.VersionServiceClient
}) *VersionHandler {
	return &VersionHandler{clientManager: cm}
}

func (h *VersionHandler) CheckVersion(c *gin.Context) {
	platformValue := c.Query("platform")
	version := c.Query("version")
	buildNumber := c.DefaultQuery("build_number", "0")

	platform, err := strconv.Atoi(platformValue)
	if platformValue == "" || version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "platform and version are required",
		})
		return
	}
	if err != nil || platform < int(versionpb.Platform_PLATFORM_IOS) || platform > int(versionpb.Platform_PLATFORM_H5) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "platform must be one of 1,2,3,4,5",
		})
		return
	}

	// Convert buildNumber to int32
	var bn int32
	for _, ch := range buildNumber {
		if ch >= '0' && ch <= '9' {
			bn = bn*10 + int32(ch-'0')
		}
	}

	resp, err := h.clientManager.Version().CheckVersion(c.Request.Context(), &versionpb.CheckVersionRequest{
		Platform:    versionpb.Platform(platform),
		Version:     version,
		BuildNumber: bn,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"has_update":          resp.HasUpdate,
			"latest_version":      resp.LatestVersion,
			"latest_build_number": resp.LatestBuildNumber,
			"force_update":        resp.ForceUpdate,
			"min_version":         resp.MinVersion,
			"min_build_number":    resp.MinBuildNumber,
			"update_info":         resp.UpdateInfo,
		},
	})
}

func (h *VersionHandler) GetLatestVersion(c *gin.Context) {
	platformValue := c.Query("platform")
	releaseTypeValue := c.DefaultQuery("release_type", "1")

	platform, err := strconv.Atoi(platformValue)
	if platformValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "platform is required",
		})
		return
	}
	if err != nil || platform < int(versionpb.Platform_PLATFORM_IOS) || platform > int(versionpb.Platform_PLATFORM_H5) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "platform must be one of 1,2,3,4,5",
		})
		return
	}
	releaseType, err := strconv.Atoi(releaseTypeValue)
	if err != nil || releaseType < int(versionpb.ReleaseType_RELEASE_TYPE_STABLE) || releaseType > int(versionpb.ReleaseType_RELEASE_TYPE_ALPHA) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "release_type must be one of 1,2,3",
		})
		return
	}

	resp, err := h.clientManager.Version().GetLatestVersion(c.Request.Context(), &versionpb.GetLatestVersionRequest{
		Platform:    versionpb.Platform(platform),
		ReleaseType: versionpb.ReleaseType(releaseType),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"version": resp.Version,
		},
	})
}

func (h *VersionHandler) ListVersions(c *gin.Context) {
	platformValue := c.Query("platform")
	releaseTypeValue := c.DefaultQuery("release_type", "0")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	platform, err := strconv.Atoi(platformValue)
	if platformValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "platform is required",
		})
		return
	}
	if err != nil || platform < int(versionpb.Platform_PLATFORM_IOS) || platform > int(versionpb.Platform_PLATFORM_H5) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "platform must be one of 1,2,3,4,5",
		})
		return
	}
	releaseType, err := strconv.Atoi(releaseTypeValue)
	if err != nil || releaseType < int(versionpb.ReleaseType_RELEASE_TYPE_UNSPECIFIED) || releaseType > int(versionpb.ReleaseType_RELEASE_TYPE_ALPHA) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "release_type must be one of 0,1,2,3",
		})
		return
	}

	var pageInt, pageSizeInt int32
	for _, ch := range page {
		if ch >= '0' && ch <= '9' {
			pageInt = pageInt*10 + int32(ch-'0')
		}
	}
	for _, ch := range pageSize {
		if ch >= '0' && ch <= '9' {
			pageSizeInt = pageSizeInt*10 + int32(ch-'0')
		}
	}

	resp, err := h.clientManager.Version().ListVersions(c.Request.Context(), &versionpb.ListVersionsRequest{
		Platform:    versionpb.Platform(platform),
		Page:        pageInt,
		PageSize:    pageSizeInt,
		ReleaseType: versionpb.ReleaseType(releaseType),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"total":    resp.Total,
			"versions": resp.Versions,
		},
	})
}

func (h *VersionHandler) ReportVersion(c *gin.Context) {
	var req struct {
		Platform    int32  `json:"platform"`
		Version     string `json:"version"`
		BuildNumber int32  `json:"build_number"`
		DeviceID    string `json:"device_id"`
		OsVersion   string `json:"os_version"`
		SdkVersion  string `json:"sdk_version"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "invalid request body",
		})
		return
	}

	if req.Platform < int32(versionpb.Platform_PLATFORM_IOS) || req.Platform > int32(versionpb.Platform_PLATFORM_H5) || req.Version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "platform must be 1..5 and version is required",
		})
		return
	}

	_, err := h.clientManager.Version().ReportVersion(c.Request.Context(), &versionpb.ReportVersionRequest{
		Platform:    versionpb.Platform(req.Platform),
		Version:     req.Version,
		BuildNumber: req.BuildNumber,
		DeviceId:    req.DeviceID,
		OsVersion:   req.OsVersion,
		SdkVersion:  req.SdkVersion,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
