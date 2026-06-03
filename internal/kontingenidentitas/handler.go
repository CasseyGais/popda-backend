package kontingenidentitas

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(c *gin.Context) {
	identitas, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data identitas kontingen",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data identitas kontingen berhasil diambil",
		"data":    identitas,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	identitas, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data identitas kontingen berhasil diambil",
		"data":    identitas,
	})
}

func (h *Handler) GetByKontingenID(c *gin.Context) {
	kontingenID, err := strconv.ParseUint(c.Param("kontingen_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Kontingen ID tidak valid",
		})
		return
	}

	identitas, err := h.service.GetByKontingenID(uint(kontingenID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data identitas kontingen berhasil diambil",
		"data":    identitas,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var request CreateKontingenIdentitasRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	identitas, err := h.service.Create(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Identitas kontingen berhasil dibuat",
		"data":    identitas,
	})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	var request UpdateKontingenIdentitasRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	identitas, err := h.service.Update(uint(id), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Identitas kontingen berhasil diupdate",
		"data":    identitas,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Identitas kontingen berhasil dihapus",
	})
}

func (h *Handler) UpdateKepalaFoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	// Coba terima upload file (multipart/form-data)
	file, fileErr := c.FormFile("foto")
	if fileErr == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
		dst := filepath.Join("uploads", "kepala", filename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal menyimpan file foto kepala",
				"error":   err.Error(),
			})
			return
		}

		fotoPath := "/uploads/kepala/" + filename
		if err := h.service.UpdateKepalaFoto(uint(id), fotoPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Foto kepala berhasil diupdate",
			"foto":    fotoPath,
		})
		return
	}

	// Fallback: terima path string via JSON
	var request struct {
		Foto string `json:"foto" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid — kirim file (multipart) atau JSON {\"foto\": \"path\"}",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.UpdateKepalaFoto(uint(id), request.Foto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Foto kepala berhasil diupdate",
		"foto":    request.Foto,
	})
}

func (h *Handler) UpdatePICFoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	// Coba terima upload file (multipart/form-data)
	file, fileErr := c.FormFile("foto")
	if fileErr == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
		dst := filepath.Join("uploads", "pic", filename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal menyimpan file foto PIC",
				"error":   err.Error(),
			})
			return
		}

		fotoPath := "/uploads/pic/" + filename
		if err := h.service.UpdatePICFoto(uint(id), fotoPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Foto PIC berhasil diupdate",
			"foto":    fotoPath,
		})
		return
	}

	// Fallback: terima path string via JSON
	var request struct {
		Foto string `json:"foto" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid — kirim file (multipart) atau JSON {\"foto\": \"path\"}",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.UpdatePICFoto(uint(id), request.Foto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Foto PIC berhasil diupdate",
		"foto":    request.Foto,
	})
}
