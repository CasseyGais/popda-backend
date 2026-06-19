package middleware

// IsSuperadmin cek apakah role dari JWT adalah superadmin.
// Dipanggil dari handler sebelum resolveKontingenID.
func IsSuperadmin(role string) bool {
	return role == "superadmin"
}
