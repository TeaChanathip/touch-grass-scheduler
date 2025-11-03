CREATE EXTENSION pg_cron;

SELECT cron.schedule('daily-pending-upload-cleanup', '0 0 * * *', $$DELETE FROM pending_uploads WHERE expire_at < now()$$);
