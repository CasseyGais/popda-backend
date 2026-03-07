package transaksi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// TrxKontingenCabor handlers
func (h *Handler) CreateTrxKontingenCabor(c *gin.Context) {
	kontingenID := c.GetUint("kontingenID")

	caborID, err := strconv.ParseUint(c.PostForm("cabor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cabor_id invalid"})
		return
	}

	putra, _ := strconv.Atoi(c.PostForm("putra"))
	putri, _ := strconv.Atoi(c.PostForm("putri"))
	pelatih, _ := strconv.Atoi(c.PostForm("pelatih"))

	err = h.service.CreateTrxKontingenCabor(kontingenID, uint(caborID), putra, putri, pelatih)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil disimpan"})
}

func (h *Handler) GetTrxKontingenCabor(c *gin.Context) {
	kontingenID := c.GetUint("kontingenID")

	trx, err := h.service.GetTrxKontingenCabor(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trx})
}

func (h *Handler) UpdateTrxKontingenCabor(c *gin.Context) {
	kontingenID := c.GetUint("kontingenID")

	caborID, err := strconv.ParseUint(c.PostForm("cabor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cabor_id invalid"})
		return
	}

	putra, _ := strconv.Atoi(c.PostForm("putra"))
	putri, _ := strconv.Atoi(c.PostForm("putri"))
	pelatih, _ := strconv.Atoi(c.PostForm("pelatih"))

	err = h.service.UpdateTrxKontingenCabor(kontingenID, uint(caborID), putra, putri, pelatih)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diupdate"})
}

// TrxKontingenNomor handlers
func (h *Handler) CreateTrxKontingenNomor(c *gin.Context) {
	kontingenID := c.GetUint("kontingenID")

	nomorID, err := strconv.ParseUint(c.PostForm("nomor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nomor_id invalid"})
		return
	}

	err = h.service.CreateTrxKontingenNomor(kontingenID, uint(nomorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nomor pertandingan berhasil ditambahkan"})
}

func (h *Handler) GetTrxKontingenNomor(c *gin.Context) {
	kontingenID := c.GetUint("kontingenID")

	trx, err := h.service.GetTrxKontingenNomor(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trx})
}

func (h *Handler) DeleteTrxKontingenNomor(c *gin.Context) {
	kontingenID := c.GetUint("kontingenID")

	nomorID, err := strconv.ParseUint(c.Param("nomor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nomor_id invalid"})
		return
	}

	err = h.service.DeleteTrxKontingenNomor(kontingenID, uint(nomorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nomor pertandingan berhasil dihapus"})
}

// TrxPendaftaranAtlet handlers
func (h *Handler) CreateTrxPendaftaranAtlet(c *gin.Context) {
	atletID, err := strconv.ParseUint(c.PostForm("atlet_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "atlet_id invalid"})
		return
	}

	nomorID, err := strconv.ParseUint(c.PostForm("nomor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nomor_id invalid"})
		return
	}

	kelasID, _ := strconv.ParseUint(c.PostForm("kelas_id"), 10, 32)

	err = h.service.CreateTrxPendaftaranAtlet(uint(atletID), uint(nomorID), uint(kelasID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pendaftaran atlet berhasil"})
}

func (h *Handler) GetTrxPendaftaranAtlet(c *gin.Context) {
	kontingenID := c.GetUint("kontingenID")

	trx, err := h.service.GetTrxPendaftaranAtlet(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trx})
}

func (h *Handler) UpdateTrxPendaftaranAtlet(c *gin.Context) {
	atletID, err := strconv.ParseUint(c.Param("atlet_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "atlet_id invalid"})
		return
	}

	nomorID, err := strconv.ParseUint(c.Param("nomor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nomor_id invalid"})
		return
	}

	status := c.PostForm("status")

	err = h.service.UpdateTrxPendaftaranAtlet(uint(atletID), uint(nomorID), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pendaftaran berhasil diupdate"})
}
