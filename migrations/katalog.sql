-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Mar 03, 2026 at 09:18 PM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `katalog`
--

-- --------------------------------------------------------

--
-- Table structure for table `product`
--

CREATE TABLE `product` (
  `id` varchar(36) NOT NULL,
  `userid` varchar(36) NOT NULL,
  `namaprod` varchar(255) NOT NULL,
  `kategori` varchar(255) NOT NULL,
  `price` float NOT NULL,
  `stock` int(5) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `product`
--

INSERT INTO `product` (`id`, `userid`, `namaprod`, `kategori`, `price`, `stock`) VALUES
('a905746a-f644-4005-bf6b-ca06e7a652b8', '5da48bc0-c695-4fec-a358-fffb7f61f562', 'Versase', 'Baju', 50.9, 12);

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `id` varchar(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `role` varchar(10) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`id`, `name`, `password`, `email`, `role`) VALUES
('5da48bc0-c695-4fec-a358-fffb7f61f562', 'Rio Jerniko', '$2a$10$buIpt.Y4zd0NJ/Ie9uz0TOd8UO4TlfXynjyUqNAe49mgtXjfHfbz6', 'rio@gmail.com', 'User'),
('b30a37f5-b89a-4190-8651-bdbd1c403c73', 'Axel', '$2a$10$VPQE6gx1wrDCV8THjoh8be2gfPMyM9zkHwT9eVjqTmu9HH/f.wfRe', 'axel@gmail.com', 'User'),
('c1862178-e659-4091-99fd-1e34b66dd965', 'Rexi Leon', '$2a$10$cmdnfP9XoGzgmynaKP7ale85ACI/0xcKUswUeMCscxprKEQ2emF4i', 'rexi@gmail.com', 'Admin');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `product`
--
ALTER TABLE `product`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
