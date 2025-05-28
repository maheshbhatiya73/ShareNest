# ShareNest - Secure File Sharing

ShareNest is a self-hosted file sharing tool that allows users to upload files with an optional password and get a secure, expiring download link. Files are automatically deleted after 1 hour.

---

## Features

- Upload any file up to 10 GB
- Optional password protection for each file
- Expiring download links (1 hour lifetime)
- Files and metadata stored securely in MongoDB
- Periodic cleanup of expired files
- Built with Next.js (Frontend) and Go (Backend)
- RESTful API design with secure password hashing

---

## Technologies Used

### Frontend
- Next.js (React)
- Tailwind CSS
- Framer Motion
- React Icons (optional)

### Backend
- Go (Golang)
- Gorilla Mux for routing
- MongoDB for storage
- bcrypt for secure password hashing

---

## Folder Structure

