-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Waktu pembuatan: 08 Mar 2026 pada 06.07
-- Versi server: 10.4.32-MariaDB
-- Versi PHP: 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `popda_2026`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `kontingen`
--

CREATE TABLE `kontingen` (
  `id` bigint(20) NOT NULL,
  `territory_id` bigint(20) NOT NULL,
  `nama_kontingen` varchar(150) DEFAULT NULL,
  `tahap1_status` enum('DRAFT','SUBMITTED') DEFAULT 'DRAFT',
  `tahap1_submitted_at` datetime DEFAULT NULL,
  `tahap2_status` enum('DRAFT','SUBMITTED') DEFAULT 'DRAFT',
  `tahap2_submitted_at` datetime DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `kontingen`
--

INSERT INTO `kontingen` (`id`, `territory_id`, `nama_kontingen`, `tahap1_status`, `tahap1_submitted_at`, `tahap2_status`, `tahap2_submitted_at`, `created_at`, `updated_at`) VALUES
(2, 2, 'Kontingen 2', '', NULL, '', NULL, '2026-02-13 04:33:26', '2026-02-13 04:33:26');

-- --------------------------------------------------------

--
-- Struktur dari tabel `kontingen_identitas`
--

CREATE TABLE `kontingen_identitas` (
  `id` bigint(20) NOT NULL,
  `kontingen_id` bigint(20) NOT NULL,
  `kepala_nama` varchar(150) DEFAULT NULL,
  `kepala_jabatan` varchar(150) DEFAULT NULL,
  `kepala_nip` varchar(50) DEFAULT NULL,
  `kepala_telepon` varchar(30) DEFAULT NULL,
  `kepala_foto` varchar(255) DEFAULT NULL,
  `pic_nama` varchar(150) DEFAULT NULL,
  `pic_jabatan` varchar(150) DEFAULT NULL,
  `pic_telepon` varchar(30) DEFAULT NULL,
  `pic_foto` varchar(255) DEFAULT NULL,
  `alamat` text DEFAULT NULL,
  `email_instansi` varchar(150) DEFAULT NULL,
  `phone_instansi` varchar(30) DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `kontingen_identitas`
--

INSERT INTO `kontingen_identitas` (`id`, `kontingen_id`, `kepala_nama`, `kepala_jabatan`, `kepala_nip`, `kepala_telepon`, `kepala_foto`, `pic_nama`, `pic_jabatan`, `pic_telepon`, `pic_foto`, `alamat`, `email_instansi`, `phone_instansi`, `updated_at`) VALUES
(14, 2, 'Bejo Sutisna', 'Ketua Kontingen', '1234567890', '+628979819142', '/uploads/kepala/20260213115553_gambar 1.PNG', '', '', '', '', '', 'casseyvienyard@gmail.com', '+628979819142', '2026-02-13 04:55:53');

-- --------------------------------------------------------

--
-- Struktur dari tabel `master_atlet`
--

CREATE TABLE `master_atlet` (
  `id` bigint(20) NOT NULL,
  `kontingen_id` bigint(20) NOT NULL,
  `sekolah_id` bigint(20) NOT NULL,
  `nisn` varchar(50) DEFAULT NULL,
  `nama` varchar(150) DEFAULT NULL,
  `jenis_kelamin` enum('PUTRA','PUTRI') DEFAULT NULL,
  `tanggal_lahir` date DEFAULT NULL,
  `kelas` varchar(20) DEFAULT NULL,
  `tinggi` int(11) DEFAULT NULL,
  `berat` decimal(5,2) DEFAULT NULL,
  `foto` varchar(255) DEFAULT NULL,
  `status_verifikasi` enum('PENDING','VALID','DITOLAK') DEFAULT 'PENDING',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `master_cabor`
--

CREATE TABLE `master_cabor` (
  `id` int(11) NOT NULL,
  `nama` varchar(150) NOT NULL,
  `max_putra` int(11) DEFAULT 0,
  `max_putri` int(11) DEFAULT 0,
  `max_pelatih` int(11) DEFAULT 0,
  `is_active` tinyint(1) DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `master_cabor`
--

INSERT INTO `master_cabor` (`id`, `nama`, `max_putra`, `max_putri`, `max_pelatih`, `is_active`, `created_at`) VALUES
(1, 'Atletik', 26, 26, 4, 1, '2026-02-11 08:25:14'),
(2, 'Angkat Besi', 5, 5, 3, 1, '2026-02-11 08:25:14'),
(3, 'Bola Basket', 16, 16, 6, 1, '2026-02-11 08:25:14'),
(4, 'Bola Voli Indoor', 12, 12, 4, 1, '2026-02-11 08:25:14'),
(5, 'Bola Voli Pasir', 2, 2, 2, 1, '2026-02-11 08:25:14'),
(6, 'Bulutangkis', 5, 5, 2, 1, '2026-02-11 08:25:14'),
(7, 'Catur', 5, 5, 3, 1, '2026-02-11 08:25:14'),
(8, 'Dayung', 20, 20, 3, 1, '2026-02-11 08:25:14'),
(9, 'Gulat', 8, 3, 2, 1, '2026-02-11 08:25:14'),
(10, 'Hockey', 12, 12, 2, 1, '2026-02-11 08:25:14'),
(11, 'Judo', 5, 5, 2, 1, '2026-02-11 08:25:14'),
(12, 'Karate', 10, 9, 4, 1, '2026-02-11 08:25:14'),
(13, 'Kempo', 11, 11, 2, 1, '2026-02-11 08:25:14'),
(14, 'Menembak', 8, 8, 4, 1, '2026-02-11 08:25:14'),
(15, 'Panahan', 16, 16, 4, 1, '2026-02-11 08:25:14'),
(16, 'Panjat Tebing', 7, 7, 2, 1, '2026-02-11 08:25:14'),
(17, 'Pencak Silat', 15, 14, 4, 1, '2026-02-11 08:25:14'),
(18, 'Renang', 10, 10, 2, 1, '2026-02-11 08:25:14'),
(19, 'Sepak Bola', 18, 0, 3, 1, '2026-02-11 08:25:14'),
(20, 'Sepak Takraw', 5, 5, 2, 1, '2026-02-11 08:25:14'),
(21, 'Senam', 5, 8, 2, 1, '2026-02-11 08:25:14'),
(22, 'Tae Kwon Do', 12, 12, 4, 1, '2026-02-11 08:25:14'),
(23, 'Tenis Lapangan', 4, 4, 2, 1, '2026-02-11 08:25:14'),
(24, 'Tenis Meja', 4, 4, 2, 1, '2026-02-11 08:25:14'),
(25, 'Tinju', 8, 6, 4, 1, '2026-02-11 08:25:14'),
(26, 'Wushu', 5, 4, 4, 1, '2026-02-11 08:25:14');

-- --------------------------------------------------------

--
-- Struktur dari tabel `master_nomor`
--

CREATE TABLE `master_nomor` (
  `id` bigint(20) NOT NULL,
  `cabor_id` int(11) NOT NULL,
  `nama` varchar(255) NOT NULL,
  `jenis_kelamin` enum('PUTRA','PUTRI','CAMPURAN') NOT NULL,
  `tipe` enum('INDIVIDU','BEREGU') DEFAULT 'INDIVIDU',
  `is_active` tinyint(1) DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `master_nomor`
--

INSERT INTO `master_nomor` (`id`, `cabor_id`, `nama`, `jenis_kelamin`, `tipe`, `is_active`, `created_at`) VALUES
(1, 1, '100 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(2, 1, '200 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(3, 1, '400 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(4, 1, '800 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(5, 1, '1500 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(6, 1, '5000 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(7, 1, 'Lompat Jauh', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(8, 1, '4x100 M Estafet', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:25:32'),
(9, 1, '100 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(10, 1, '200 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(11, 1, '400 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(12, 1, '800 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(13, 1, '1500 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(14, 1, '5000 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(15, 1, 'Lompat Jauh', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:25:32'),
(16, 1, '4x100 M Estafet', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:25:32'),
(17, 19, 'Sepak Bola Putra', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:26:01'),
(18, 24, 'Tunggal', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:26:28'),
(19, 24, 'Ganda', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:26:28'),
(20, 24, 'Beregu', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:26:28'),
(21, 24, 'Tunggal', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:26:28'),
(22, 24, 'Ganda', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:26:28'),
(23, 24, 'Beregu', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:26:28'),
(24, 24, 'Ganda Campuran', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:26:28'),
(25, 1, '100 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(26, 1, '200 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(27, 1, '400 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(28, 1, '800 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(29, 1, '1500 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(30, 1, '5000 M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(31, 1, 'Lompat Jauh', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(32, 1, 'Lompat Tinggi', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(33, 1, 'Tolak Peluru', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(34, 1, 'Lempar Lembing', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(35, 1, 'Lempar Cakram', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(36, 1, '4x100 M', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(37, 1, '4x400 M', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(38, 1, '100 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(39, 1, '200 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(40, 1, '400 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(41, 1, '800 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(42, 1, '1500 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(43, 1, '5000 M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(44, 1, 'Lompat Jauh', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(45, 1, 'Lompat Tinggi', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(46, 1, 'Tolak Peluru', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(47, 1, 'Lempar Lembing', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(48, 1, 'Lempar Cakram', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(49, 1, '4x100 M', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(50, 1, '4x400 M', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(51, 2, '56 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(52, 2, '60 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(53, 2, '65 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(54, 2, '71 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(55, 2, '+71 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(56, 2, '44 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(57, 2, '48 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(58, 2, '53 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(59, 2, '58 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(60, 2, '+58 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(61, 3, 'Basket 5x5', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(62, 3, 'Basket 3x3', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(63, 3, 'Basket 5x5', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(64, 3, 'Basket 3x3', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(65, 4, 'Voli Indoor', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(66, 4, 'Voli Indoor', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(67, 5, 'Voli Pasir', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(68, 5, 'Voli Pasir', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(69, 6, 'Tunggal', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(70, 6, 'Ganda', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(71, 6, 'Beregu', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(72, 6, 'Tunggal', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(73, 6, 'Ganda', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(74, 6, 'Beregu', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(75, 6, 'Ganda Campuran', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(76, 7, 'Cepat Beregu', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(77, 7, 'Kilat Beregu', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(78, 7, 'Standar Beregu', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(79, 7, 'Cepat Beregu', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(80, 7, 'Kilat Beregu', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(81, 7, 'Standar Beregu', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(82, 7, 'Mix Cepat', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(83, 7, 'Mix Kilat', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(84, 7, 'Mix Standar', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(85, 8, 'K1 200M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(86, 8, 'K1 500M', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(87, 8, 'Dragon Boat', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(88, 8, 'K1 200M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(89, 8, 'K1 500M', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(90, 8, 'Dragon Boat', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(91, 8, 'Dragon Boat Mix', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(92, 9, 'Gaya Bebas', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(93, 9, 'Gaya Bebas', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(94, 10, 'Hockey', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(95, 10, 'Hockey', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(96, 11, '50 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(97, 11, '55 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(98, 11, '60 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(99, 11, '52 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(100, 11, '57 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(101, 11, 'Beregu Campuran', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(102, 12, 'Kata', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(103, 12, 'Kumite', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(104, 12, 'Kata', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(105, 12, 'Kumite', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(106, 12, 'Kumite Beregu', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(107, 12, 'Kumite Beregu', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(108, 13, 'Randori', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(109, 13, 'Randori', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(110, 13, 'Embu Pasangan', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(111, 14, 'Air Pistol', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(112, 14, 'Air Rifle', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(113, 14, 'Air Pistol', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(114, 14, 'Air Rifle', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(115, 14, 'Mixed Team', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(116, 15, 'Recurve', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(117, 15, 'Compound', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(118, 15, 'Recurve', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(119, 15, 'Compound', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(120, 15, 'Mix Team', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(121, 16, 'Speed', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(122, 16, 'Lead', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(123, 16, 'Speed', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(124, 16, 'Lead', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(125, 16, 'Speed Mix', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(126, 17, 'Tunggal', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(127, 17, 'Ganda', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(128, 17, 'Regu', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(129, 17, 'Tunggal', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(130, 17, 'Ganda', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(131, 17, 'Regu', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(132, 18, '50 M Bebas', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(133, 18, '100 M Bebas', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(134, 18, '50 M Bebas', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(135, 18, '100 M Bebas', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(136, 19, 'Sepak Bola', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(137, 20, 'Inter Regu', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(138, 20, 'Double Event', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(139, 20, 'Inter Regu', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(140, 20, 'Double Event', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(141, 21, 'Artistik', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(142, 21, 'Artistik', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(143, 21, 'Ritmik', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(144, 22, 'Poomsae', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(145, 22, 'Kyorugi', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(146, 22, 'Poomsae', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(147, 22, 'Kyorugi', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(148, 23, 'Tunggal', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(149, 23, 'Ganda', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(150, 23, 'Tunggal', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(151, 23, 'Ganda', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(152, 23, 'Ganda Campuran', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(153, 24, 'Tunggal', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(154, 24, 'Ganda', 'PUTRA', 'BEREGU', 1, '2026-02-11 08:38:00'),
(155, 24, 'Tunggal', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(156, 24, 'Ganda', 'PUTRI', 'BEREGU', 1, '2026-02-11 08:38:00'),
(157, 24, 'Ganda Campuran', 'CAMPURAN', 'BEREGU', 1, '2026-02-11 08:38:00'),
(158, 25, '45 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(159, 25, '48 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(160, 25, '50 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(161, 25, '52 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(162, 26, '48 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(163, 26, '52 Kg', 'PUTRA', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(164, 26, '48 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00'),
(165, 26, '52 Kg', 'PUTRI', 'INDIVIDU', 1, '2026-02-11 08:38:00');

-- --------------------------------------------------------

--
-- Struktur dari tabel `master_official`
--

CREATE TABLE `master_official` (
  `id` bigint(20) NOT NULL,
  `kontingen_id` bigint(20) NOT NULL,
  `nama` varchar(150) DEFAULT NULL,
  `jabatan` varchar(100) DEFAULT NULL,
  `no_hp` varchar(30) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `master_pelatih`
--

CREATE TABLE `master_pelatih` (
  `id` bigint(20) NOT NULL,
  `kontingen_id` bigint(20) NOT NULL,
  `nama` varchar(150) DEFAULT NULL,
  `no_hp` varchar(30) DEFAULT NULL,
  `sertifikat` varchar(255) DEFAULT NULL,
  `foto` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `master_sekolah`
--

CREATE TABLE `master_sekolah` (
  `id` bigint(20) NOT NULL,
  `nama` varchar(200) DEFAULT NULL,
  `npsn` varchar(20) DEFAULT NULL,
  `alamat` text DEFAULT NULL,
  `kabupaten` varchar(150) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `permissions`
--

CREATE TABLE `permissions` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `description` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `permissions`
--

INSERT INTO `permissions` (`id`, `name`, `description`) VALUES
(1, 'dashboard.view', 'Melihat dashboard statistik'),
(2, 'master.cabor.view', 'Melihat data master cabor'),
(3, 'master.cabor.manage', 'Menambah / mengubah / menghapus master cabor'),
(4, 'master.nomor.view', 'Melihat data master nomor'),
(5, 'master.nomor.manage', 'Menambah / mengubah / menghapus master nomor'),
(6, 'master.sekolah.view', 'Melihat data sekolah'),
(7, 'master.sekolah.manage', 'Mengelola data sekolah'),
(8, 'master.atlet.view', 'Melihat data atlet'),
(9, 'master.atlet.manage', 'Mengelola data atlet'),
(10, 'master.pelatih.view', 'Melihat data pelatih'),
(11, 'master.pelatih.manage', 'Mengelola data pelatih'),
(12, 'master.official.view', 'Melihat data official'),
(13, 'master.official.manage', 'Mengelola data official'),
(14, 'tahap1.view', 'Melihat data Tahap I'),
(15, 'tahap1.submit', 'Submit Tahap I'),
(16, 'tahap2.view', 'Melihat data Tahap II'),
(17, 'tahap2.submit', 'Submit Tahap II'),
(18, 'tahap3.view', 'Melihat data Tahap III'),
(19, 'tahap3.submit', 'Submit Tahap III'),
(20, 'verifikasi.view', 'Melihat data untuk verifikasi'),
(21, 'verifikasi.approve', 'Approve atau tolak data'),
(22, 'kontingen.view', 'Melihat data kontingen'),
(23, 'kontingen.manage', 'Mengelola data kontingen'),
(24, 'user.view', 'Melihat data user'),
(25, 'user.manage', 'Mengelola data user'),
(26, 'role.view', 'Melihat data role'),
(27, 'role.manage', 'Mengelola data role');

-- --------------------------------------------------------

--
-- Struktur dari tabel `roles`
--

CREATE TABLE `roles` (
  `id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `description` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `roles`
--

INSERT INTO `roles` (`id`, `name`, `description`) VALUES
(1, 'SUPERADMIN', 'Admin pusat Dispora Provinsi'),
(2, 'ADMIN', 'Admin masing-masing Kabupaten/Kota'),
(3, 'STAFF_LAPANGAN', 'Petugas lapangan / verifikator');

-- --------------------------------------------------------

--
-- Struktur dari tabel `role_permissions`
--

CREATE TABLE `role_permissions` (
  `role_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `territories`
--

CREATE TABLE `territories` (
  `id` bigint(20) NOT NULL,
  `name` varchar(100) NOT NULL,
  `type` enum('PROVINSI','KABUPATEN','KOTA') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `territories`
--

INSERT INTO `territories` (`id`, `name`, `type`) VALUES
(1, 'Provinsi Banten', 'PROVINSI'),
(2, 'Kabupaten Tangerang', 'KABUPATEN'),
(3, 'Kabupaten Serang', 'KABUPATEN'),
(4, 'Kabupaten Lebak', 'KABUPATEN'),
(5, 'Kabupaten Pandeglang', 'KABUPATEN'),
(6, 'Kota Tangerang', 'KOTA'),
(7, 'Kota Tangerang Selatan', 'KOTA'),
(8, 'Kota Serang', 'KOTA'),
(9, 'Kota Cilegon', 'KOTA');

-- --------------------------------------------------------

--
-- Struktur dari tabel `trx_kontingen_cabor`
--

CREATE TABLE `trx_kontingen_cabor` (
  `id` bigint(20) NOT NULL,
  `kontingen_id` bigint(20) NOT NULL,
  `cabor_id` int(11) NOT NULL,
  `putra` int(11) DEFAULT 0,
  `putri` int(11) DEFAULT 0,
  `pelatih` int(11) DEFAULT 0,
  `total_atlet` int(11) DEFAULT 0,
  `total_personel` int(11) DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Trigger `trx_kontingen_cabor`
--
DELIMITER $$
CREATE TRIGGER `before_insert_trx_kontingen_cabor` BEFORE INSERT ON `trx_kontingen_cabor` FOR EACH ROW BEGIN
    DECLARE max_putra INT;
    DECLARE max_putri INT;
    DECLARE max_pelatih INT;

    SELECT max_putra, max_putri, max_pelatih
    INTO max_putra, max_putri, max_pelatih
    FROM master_cabor
    WHERE id = NEW.cabor_id;

    IF NEW.putra > max_putra THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Jumlah atlet putra melebihi kuota';
    END IF;

    IF NEW.putri > max_putri THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Jumlah atlet putri melebihi kuota';
    END IF;

    IF NEW.pelatih > max_pelatih THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Jumlah pelatih melebihi kuota';
    END IF;

END
$$
DELIMITER ;

-- --------------------------------------------------------

--
-- Struktur dari tabel `trx_kontingen_nomor`
--

CREATE TABLE `trx_kontingen_nomor` (
  `id` bigint(20) NOT NULL,
  `kontingen_id` bigint(20) NOT NULL,
  `nomor_id` bigint(20) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `trx_pendaftaran_atlet`
--

CREATE TABLE `trx_pendaftaran_atlet` (
  `id` bigint(20) NOT NULL,
  `atlet_id` bigint(20) NOT NULL,
  `nomor_id` bigint(20) NOT NULL,
  `kelas_id` bigint(20) DEFAULT NULL,
  `status` enum('PENDING','VERIFIED','REJECTED') DEFAULT 'PENDING',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Trigger `trx_pendaftaran_atlet`
--
DELIMITER $$
CREATE TRIGGER `before_insert_trx_pendaftaran_atlet` BEFORE INSERT ON `trx_pendaftaran_atlet` FOR EACH ROW BEGIN
    DECLARE kontingenId BIGINT;
    DECLARE total INT;

    SELECT kontingen_id INTO kontingenId
    FROM master_atlet
    WHERE id = NEW.atlet_id;

    SELECT COUNT(*) INTO total
    FROM trx_kontingen_nomor
    WHERE kontingen_id = kontingenId
      AND nomor_id = NEW.nomor_id;

    IF total = 0 THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Nomor belum dipilih pada Tahap II';
    END IF;

END
$$
DELIMITER ;
DELIMITER $$
CREATE TRIGGER `trg_validasi_berat_atlet` BEFORE INSERT ON `trx_pendaftaran_atlet` FOR EACH ROW BEGIN
    DECLARE v_berat DECIMAL(5,2);
    DECLARE v_min DECIMAL(5,2);
    DECLARE v_max DECIMAL(5,2);

    -- Jika kelas NULL (tidak pakai berat), skip validasi
    IF NEW.kelas_id IS NOT NULL THEN

        -- Ambil berat atlet
        SELECT berat INTO v_berat
        FROM master_atlet
        WHERE id = NEW.atlet_id;

        -- Ambil batas berat kelas
        SELECT berat_min, berat_max
        INTO v_min, v_max
        FROM master_kelas
        WHERE id = NEW.kelas_id;

        -- Validasi bawah
        IF v_min IS NOT NULL AND v_berat < v_min THEN
            SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'Berat atlet di bawah batas kelas!';
        END IF;

        -- Validasi atas
        IF v_max IS NOT NULL AND v_berat > v_max THEN
            SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'Berat atlet melebihi batas kelas!';
        END IF;

    END IF;

END
$$
DELIMITER ;

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `name` varchar(150) NOT NULL,
  `email` varchar(150) NOT NULL,
  `password` varchar(255) NOT NULL,
  `is_active` tinyint(1) DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `avatar` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `name`, `email`, `password`, `is_active`, `created_at`, `avatar`) VALUES
(1, 'Super Admin Dispora', 'superadmin@popda.id', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 1, '2026-02-11 07:58:38', NULL),
(2, 'Admin Kab Tangerang', 'admin.kabtangerang@popda.id', '871241dbf332bd665f337dbbe036fd560ee1af5731fd29a9d0fc449c44548d4a', 1, '2026-02-11 07:58:38', '/avatar/tangerangkab.png'),
(3, 'Admin Kab Serang', 'admin.kabserang@popda.id', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 1, '2026-02-11 07:58:38', NULL),
(4, 'Admin Kab Lebak', 'admin.kablebak@popda.id', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 1, '2026-02-11 07:58:38', NULL),
(5, 'Admin Kab Pandeglang', 'admin.kabpandeglang@popda.id', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 1, '2026-02-11 07:58:38', NULL),
(6, 'Admin Kota Tangerang', 'admin.kotatangerang@popda.id', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 1, '2026-02-11 07:58:38', NULL),
(7, 'Admin Tangsel', 'admin.tangsel@popda.id', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 1, '2026-02-11 07:58:38', NULL),
(8, 'Admin Kota Serang', 'admin.kotaserang@popda.id', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 1, '2026-02-11 07:58:38', NULL),
(9, 'Admin Kota Cilegon', 'admin.kotacilegon@popda.id', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 1, '2026-02-11 07:58:38', NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `user_roles`
--

CREATE TABLE `user_roles` (
  `user_id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `user_roles`
--

INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES
(1, 1),
(2, 2),
(3, 2),
(4, 2),
(5, 2),
(6, 2),
(7, 2),
(8, 2),
(9, 2);

-- --------------------------------------------------------

--
-- Struktur dari tabel `user_territories`
--

CREATE TABLE `user_territories` (
  `user_id` int(11) NOT NULL,
  `territory_id` bigint(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `user_territories`
--

INSERT INTO `user_territories` (`user_id`, `territory_id`) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5),
(6, 6),
(7, 7),
(8, 8),
(9, 9);

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `kontingen`
--
ALTER TABLE `kontingen`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_kontingen_territory` (`territory_id`);

--
-- Indeks untuk tabel `kontingen_identitas`
--
ALTER TABLE `kontingen_identitas`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `kontingen_id` (`kontingen_id`);

--
-- Indeks untuk tabel `master_atlet`
--
ALTER TABLE `master_atlet`
  ADD PRIMARY KEY (`id`),
  ADD KEY `kontingen_id` (`kontingen_id`),
  ADD KEY `sekolah_id` (`sekolah_id`);

--
-- Indeks untuk tabel `master_cabor`
--
ALTER TABLE `master_cabor`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `master_nomor`
--
ALTER TABLE `master_nomor`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_cabor` (`cabor_id`),
  ADD KEY `idx_kelamin` (`jenis_kelamin`),
  ADD KEY `idx_active` (`is_active`);

--
-- Indeks untuk tabel `master_official`
--
ALTER TABLE `master_official`
  ADD PRIMARY KEY (`id`),
  ADD KEY `kontingen_id` (`kontingen_id`);

--
-- Indeks untuk tabel `master_pelatih`
--
ALTER TABLE `master_pelatih`
  ADD PRIMARY KEY (`id`),
  ADD KEY `kontingen_id` (`kontingen_id`);

--
-- Indeks untuk tabel `master_sekolah`
--
ALTER TABLE `master_sekolah`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `permissions`
--
ALTER TABLE `permissions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `role_permissions`
--
ALTER TABLE `role_permissions`
  ADD PRIMARY KEY (`role_id`,`permission_id`),
  ADD KEY `permission_id` (`permission_id`);

--
-- Indeks untuk tabel `territories`
--
ALTER TABLE `territories`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `trx_kontingen_cabor`
--
ALTER TABLE `trx_kontingen_cabor`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_kontingen_cabor` (`kontingen_id`,`cabor_id`),
  ADD KEY `cabor_id` (`cabor_id`);

--
-- Indeks untuk tabel `trx_kontingen_nomor`
--
ALTER TABLE `trx_kontingen_nomor`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_kontingen_nomor` (`kontingen_id`,`nomor_id`),
  ADD KEY `nomor_id` (`nomor_id`);

--
-- Indeks untuk tabel `trx_pendaftaran_atlet`
--
ALTER TABLE `trx_pendaftaran_atlet`
  ADD PRIMARY KEY (`id`),
  ADD KEY `atlet_id` (`atlet_id`),
  ADD KEY `nomor_id` (`nomor_id`),
  ADD KEY `kelas_id` (`kelas_id`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`);

--
-- Indeks untuk tabel `user_roles`
--
ALTER TABLE `user_roles`
  ADD PRIMARY KEY (`user_id`,`role_id`),
  ADD KEY `role_id` (`role_id`);

--
-- Indeks untuk tabel `user_territories`
--
ALTER TABLE `user_territories`
  ADD PRIMARY KEY (`user_id`,`territory_id`),
  ADD KEY `fk_user_territory` (`territory_id`);

--
-- AUTO_INCREMENT untuk tabel yang dibuang
--

--
-- AUTO_INCREMENT untuk tabel `kontingen`
--
ALTER TABLE `kontingen`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT untuk tabel `kontingen_identitas`
--
ALTER TABLE `kontingen_identitas`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;

--
-- AUTO_INCREMENT untuk tabel `master_atlet`
--
ALTER TABLE `master_atlet`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `master_cabor`
--
ALTER TABLE `master_cabor`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=30;

--
-- AUTO_INCREMENT untuk tabel `master_nomor`
--
ALTER TABLE `master_nomor`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=166;

--
-- AUTO_INCREMENT untuk tabel `master_official`
--
ALTER TABLE `master_official`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `master_pelatih`
--
ALTER TABLE `master_pelatih`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `master_sekolah`
--
ALTER TABLE `master_sekolah`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `permissions`
--
ALTER TABLE `permissions`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=28;

--
-- AUTO_INCREMENT untuk tabel `roles`
--
ALTER TABLE `roles`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT untuk tabel `territories`
--
ALTER TABLE `territories`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

--
-- AUTO_INCREMENT untuk tabel `trx_kontingen_cabor`
--
ALTER TABLE `trx_kontingen_cabor`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `trx_kontingen_nomor`
--
ALTER TABLE `trx_kontingen_nomor`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `trx_pendaftaran_atlet`
--
ALTER TABLE `trx_pendaftaran_atlet`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `kontingen`
--
ALTER TABLE `kontingen`
  ADD CONSTRAINT `fk_kontingen_territory` FOREIGN KEY (`territory_id`) REFERENCES `territories` (`id`) ON UPDATE CASCADE;

--
-- Ketidakleluasaan untuk tabel `kontingen_identitas`
--
ALTER TABLE `kontingen_identitas`
  ADD CONSTRAINT `fk_identitas_kontingen` FOREIGN KEY (`kontingen_id`) REFERENCES `kontingen` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ketidakleluasaan untuk tabel `master_atlet`
--
ALTER TABLE `master_atlet`
  ADD CONSTRAINT `master_atlet_ibfk_1` FOREIGN KEY (`kontingen_id`) REFERENCES `kontingen` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `master_atlet_ibfk_2` FOREIGN KEY (`sekolah_id`) REFERENCES `master_sekolah` (`id`);

--
-- Ketidakleluasaan untuk tabel `master_nomor`
--
ALTER TABLE `master_nomor`
  ADD CONSTRAINT `master_nomor_ibfk_1` FOREIGN KEY (`cabor_id`) REFERENCES `master_cabor` (`id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `master_official`
--
ALTER TABLE `master_official`
  ADD CONSTRAINT `master_official_ibfk_1` FOREIGN KEY (`kontingen_id`) REFERENCES `kontingen` (`id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `master_pelatih`
--
ALTER TABLE `master_pelatih`
  ADD CONSTRAINT `master_pelatih_ibfk_1` FOREIGN KEY (`kontingen_id`) REFERENCES `kontingen` (`id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `role_permissions`
--
ALTER TABLE `role_permissions`
  ADD CONSTRAINT `role_permissions_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `role_permissions_ibfk_2` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `trx_kontingen_cabor`
--
ALTER TABLE `trx_kontingen_cabor`
  ADD CONSTRAINT `trx_kontingen_cabor_ibfk_1` FOREIGN KEY (`kontingen_id`) REFERENCES `kontingen` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `trx_kontingen_cabor_ibfk_2` FOREIGN KEY (`cabor_id`) REFERENCES `master_cabor` (`id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `trx_kontingen_nomor`
--
ALTER TABLE `trx_kontingen_nomor`
  ADD CONSTRAINT `trx_kontingen_nomor_ibfk_1` FOREIGN KEY (`kontingen_id`) REFERENCES `kontingen` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `trx_kontingen_nomor_ibfk_2` FOREIGN KEY (`nomor_id`) REFERENCES `master_nomor` (`id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `trx_pendaftaran_atlet`
--
ALTER TABLE `trx_pendaftaran_atlet`
  ADD CONSTRAINT `trx_pendaftaran_atlet_ibfk_1` FOREIGN KEY (`atlet_id`) REFERENCES `master_atlet` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `trx_pendaftaran_atlet_ibfk_2` FOREIGN KEY (`nomor_id`) REFERENCES `master_nomor` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `trx_pendaftaran_atlet_ibfk_3` FOREIGN KEY (`kelas_id`) REFERENCES `master_kelas` (`id`) ON DELETE SET NULL;

--
-- Ketidakleluasaan untuk tabel `user_roles`
--
ALTER TABLE `user_roles`
  ADD CONSTRAINT `user_roles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `user_roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `user_territories`
--
ALTER TABLE `user_territories`
  ADD CONSTRAINT `fk_user_territory` FOREIGN KEY (`territory_id`) REFERENCES `territories` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `user_territories_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
