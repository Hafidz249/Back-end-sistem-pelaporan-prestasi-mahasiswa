-- Table: notifications
-- Untuk menyimpan notifikasi ke user (dosen wali, admin, dll)

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data TEXT,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index untuk performa query
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);

-- Notification types:
-- - achievement_submitted: Prestasi disubmit untuk verifikasi
-- - achievement_verified: Prestasi diverifikasi
-- - achievement_rejected: Prestasi ditolak
-- - achievement_updated: Prestasi diupdate
